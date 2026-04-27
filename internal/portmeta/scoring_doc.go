// Package portmeta — scoring sub-feature
//
// scoring.go provides ScoreLevel, a five-bucket normalisation (Negligible →
// Critical) of the raw 0–100 composite score produced by CompositeScoreFor.
//
// scoring_annotator.go attaches the raw score and its bucket label to each
// PortInfo's Meta map under the keys "score_raw" and "score_level", and
// exposes FilterByMinScore for pipeline filtering.
//
// Typical pipeline usage:
//
//	annotate := portmeta.NewScoringAnnotator()
//	enriched := annotate(ports)
//	high    := portmeta.FilterByMinScore(enriched, portmeta.ScoreHigh)
package portmeta
