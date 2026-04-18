package sampler

import (
	"context"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Aggregator fans-in results from multiple Samplers and deduplicates ports
// seen within a single aggregation cycle.
type Aggregator struct {
	samplers []*Sampler
	out      chan Result
}

// NewAggregator creates an Aggregator over the provided samplers.
func NewAggregator(samplers ...*Sampler) *Aggregator {
	return &Aggregator{
		samplers: samplers,
		out:      make(chan Result, 8),
	}
}

// Out returns the merged result channel.
func (a *Aggregator) Out() <-chan Result { return a.out }

// Run starts all samplers and fans their results into Out().
// Blocks until ctx is cancelled.
func (a *Aggregator) Run(ctx context.Context) {
	var wg sync.WaitGroup
	for _, s := range a.samplers {
		wg.Add(1)
		go func(s *Sampler) {
			defer wg.Done()
			s.Run(ctx)
		}(s)
		wg.Add(1)
		go func(s *Sampler) {
			defer wg.Done()
			for r := range s.Results() {
				select {
				case a.out <- r:
				case <-ctx.Done():
					return
				}
			}
		}(s)
	}
	wg.Wait()
	close(a.out)
}

// Merge combines multiple port lists, deduplicating by port number.
func Merge(lists ...[]scanner.PortInfo) []scanner.PortInfo {
	seen := make(map[int]scanner.PortInfo)
	for _, list := range lists {
		for _, p := range list {
			if _, ok := seen[p.Port]; !ok {
				seen[p.Port] = p
			}
		}
	}
	out := make([]scanner.PortInfo, 0, len(seen))
	for _, p := range seen {
		out = append(out, p)
	}
	return out
}
