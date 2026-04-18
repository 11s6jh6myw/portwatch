package sampler_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/sampler"
	"github.com/user/portwatch/internal/scanner"
)

func TestMerge_DeduplicatesPorts(t *testing.T) {
	a := []scanner.PortInfo{{Port: 80}, {Port: 443}}
	b := []scanner.PortInfo{{Port: 443}, {Port: 8080}}
	result := sampler.Merge(a, b)
	if len(result) != 3 {
		t.Fatalf("expected 3 unique ports, got %d", len(result))
	}
}

func TestMerge_EmptyInputs(t *testing.T) {
	result := sampler.Merge()
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d", len(result))
	}
}

func TestAggregator_FansInResults(t *testing.T) {
	p1 := startTCP(t)
	p2 := startTCP(t)

	s1 := scanner.NewTCPScanner("127.0.0.1", []int{p1}, 100*time.Millisecond)
	s2 := scanner.NewTCPScanner("127.0.0.1", []int{p2}, 100*time.Millisecond)

	smp1 := sampler.New(s1, 50*time.Millisecond, 0)
	smp2 := sampler.New(s2, 50*time.Millisecond, 0)

	agg := sampler.NewAggregator(smp1, smp2)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go agg.Run(ctx)

	received := 0
	for received < 2 {
		select {
		case r, ok := <-agg.Out():
			if !ok {
				t.Fatal("channel closed early")
			}
			if r.Err == nil && len(r.Ports) > 0 {
				received++
			}
		case <-ctx.Done():
			t.Fatalf("timed out, only received %d results", received)
		}
	}
}

func TestAggregator_ClosesOutOnStop(t *testing.T) {
	s := scanner.NewTCPScanner("127.0.0.1", []int{29999}, 100*time.Millisecond)
	smp := sampler.New(s, 500*time.Millisecond, 0)
	agg := sampler.NewAggregator(smp)

	ctx, cancel := context.WithCancel(context.Background())
	go agg.Run(ctx)
	cancel()

	select {
	case <-agg.Out():
	case <-time.After(time.Second):
		t.Fatal("out channel not closed after cancel")
	}
}
