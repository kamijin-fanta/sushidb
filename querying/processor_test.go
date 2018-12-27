package querying

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func JsonToInterface(t *testing.T, str string) interface{} {
	var res interface{}
	err := json.Unmarshal([]byte(str), &res)
	assert.Nil(t, err)

	return res
}

func TestFilterRowSingleEq(t *testing.T) {
	singleEq := QueryProcessor{
		Query: QueryAstRoot{
			Lower: 0,
			Upper: 0,
			Sort:  "desc",
			Filters: []FilterExpr{
				{
					Type:         "eq",
					Path:         "$.id",
					Value:        "target-id",
					ChildrenExpr: nil,
				},
			},
		},
	}

	condition, err := singleEq.FilterRow(JsonToInterface(t, `{
		"id": "target-id"
	}`))
	assert.Nil(t, err)
	assert.True(t, condition)

	condition, err = singleEq.FilterRow(JsonToInterface(t, `{
		"id": "not-found"
	}`))
	assert.Nil(t, err)
	assert.False(t, condition)
}

func TestFilterRowNumberEq(t *testing.T) {
	numberEq := QueryProcessor{
		Query: QueryAstRoot{
			Lower: 0,
			Upper: 0,
			Sort:  "desc",
			Filters: []FilterExpr{
				{
					Type:         "eq",
					Path:         "$.id",
					Value:        123,
					ChildrenExpr: nil,
				},
			},
		},
	}

	condition, err := numberEq.FilterRow(JsonToInterface(t, `{
		"id": 123
	}`))
	assert.Nil(t, err)
	assert.True(t, condition)

	condition, err = numberEq.FilterRow(JsonToInterface(t, `{
		"id": 1234
	}`))
	assert.Nil(t, err)
	assert.False(t, condition)

	floatEq := QueryProcessor{
		Query: QueryAstRoot{
			Lower: 0,
			Upper: 0,
			Sort:  "desc",
			Filters: []FilterExpr{
				{
					Type:         "eq",
					Path:         "$.id",
					Value:        123.456,
					ChildrenExpr: nil,
				},
			},
		},
	}

	condition, err = floatEq.FilterRow(JsonToInterface(t, `{
		"id": 123.456
	}`))
	assert.Nil(t, err)
	assert.True(t, condition)

	condition, err = floatEq.FilterRow(JsonToInterface(t, `{
		"id": 123
	}`))
	assert.Nil(t, err)
	assert.False(t, condition)
}

func TestFilterRowOr(t *testing.T) {
	orExpr := QueryProcessor{
		Query: QueryAstRoot{
			Lower: 0,
			Upper: 0,
			Sort:  "desc",
			Filters: []FilterExpr{
				{
					Type: "or",
					ChildrenExpr: []FilterExpr{
						{
							Type:  "eq",
							Path:  "$.id",
							Value: "hogehoge",
						},
						{
							Type:  "eq",
							Path:  "$.id",
							Value: "test-data",
						},
					},
				},
			},
		},
	}

	condition, err := orExpr.FilterRow(JsonToInterface(t, `{
		"id": "test-data"
	}`))
	assert.Nil(t, err)
	assert.True(t, condition)

	condition, err = orExpr.FilterRow(JsonToInterface(t, `{
		"id": "not-found"
	}`))
	assert.Nil(t, err)
	assert.False(t, condition)
}
