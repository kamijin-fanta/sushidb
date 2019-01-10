package querying

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"strings"
)

type QueryAstRoot struct {
	Lower      int64        `json:"lower"`    // nanosecond
	Upper      int64        `json:"upper"`    // nanosecond
	Sort       string       `json:"sort"`     // asc or desc
	Limit      int          `json:"limit"`    // limit count
	MaxSkip    int          `json:"max_skip"` // limit of skip count
	Cursor     string       `json:"cursor"`   // cursor bound
	Filters    []FilterExpr `json:"filters"`
	MetricKeys []string     `json:"metric_keys"`
}
type FilterExpr struct {
	Type         string       `json:"type"`
	Path         string       `json:"path"`
	Value        interface{}  `json:"value"`
	ChildrenExpr []FilterExpr `json:"children"`
}

func (q *QueryAstRoot) ParseCursor() (timestamp int64, skipKeys int, err error) {
	if len(q.Cursor) == 0 {
		return
	}

	split := strings.Split(q.Cursor, ",")
	if len(split) != 2 {
		err = errors.New("cannot parse cursor")
		return
	}

	timestamp, err = strconv.ParseInt(split[0], 10, 0)
	if err != nil {
		return
	}
	skipKeys, err = strconv.Atoi(split[1])
	if err != nil {
		return
	}
	return
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
