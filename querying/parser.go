package querying

import (
	"encoding/json"
	"math"
)

type QueryAstRoot struct {
	Lower     int64        `json:"lower"`    // nanosecond
	Upper     int64        `json:"upper"`    // nanosecond
	Sort      string       `json:"sort"`     // asc or desc
	Limit     int          `json:"limit"`    // limit count
	MaxSkip   int          `json:"max_skip"` // limit of skip count
	Cursor    int64        `json:"cursor"`   // cursor bound
	Filters   []FilterExpr `json:"filters"`
	MetricIDs []string     `json:"metric_ids"`
}
type FilterExpr struct {
	Type         string       `json:"type"`
	Path         string       `json:"path"`
	Value        interface{}  `json:"value"`
	ChildrenExpr []FilterExpr `json:"children"`
}

func QueryParser(data []byte) (*QueryAstRoot, error) {
	var query QueryAstRoot
	err := json.Unmarshal(data, &query)
	if err != nil {
		return nil, err
	}
	if query.Upper == 0 {
		query.Upper = math.MaxInt64
	}
	if query.Limit == 0 {
		query.Limit = 1000
	}
	if query.MaxSkip == 0 {
		query.MaxSkip = 1000
	}
	return &query, nil
}
