package report

import "time"

// A Card sums up the whole report for a single repo.
// It contains a set of Scores, a final Grade and some other meta data.
type Card struct {
	Checks               []Score   `json:"checks"`
	Average              float64   `json:"average"`
	Grade                Grade     `json:"grade"`
	Files                int       `json:"files"`
	Issues               int       `json:"issues"`
	Repo                 string    `json:"repo"`
	ResolvedRepo         string    `json:"resolvedRepo"`
	LastRefresh          time.Time `json:"last_refresh"`
	LastRefreshFormatted string    `json:"formatted_last_refresh"`
	LastRefreshHumanized string    `json:"humanized_last_refresh"`
}

// A Score represents a single executed check.
type Score struct {
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	FileSummaries []FileSummary `json:"file_summaries"`
	Weight        float64       `json:"weight"`
	Percentage    float64       `json:"percentage"`
	Error         string        `json:"error"`
}

// ByWeight implements sorting for checks by weight descending
type ByWeight []Score

func (a ByWeight) Len() int           { return len(a) }
func (a ByWeight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByWeight) Less(i, j int) bool { return a[i].Weight > a[j].Weight }
