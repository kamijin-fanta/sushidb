package querying

import (
	"encoding/json"
	"math"
)

type QueryAstRoot struct {
	Lower   int64        `json:"lower"` // nanosecond
	Upper   int64        `json:"upper"` // nanosecond
	Sort    string       `json:"sort"`  // asc or desc
	Limit   int          `json:"limit"` // asc or desc
	Filters []FilterExpr `json:"filters"`
}
type FilterExpr struct {
	Type         string        `json:"type"`
	Path         string        `json:"path"`
	Value        interface{}   `json:"value"`
	ChildrenExpr *[]FilterExpr `json:"children"`
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
	return &query, nil
}
