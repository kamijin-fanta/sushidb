package querying

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQueryParser(t *testing.T) {
	queryStr := `
		{
          "lower": 10000,
          "sort": "desc",
		  "filters": [
			{
			  "type": "eq",
			  "path": ".moduleid",
			  "value": "hoge"
			},
			{
			  "type": "or",
			  "children": [
			  	{
				  "type": "gte",
				  "path": ".tmp",
				  "value": 10
				},
				{
				  "type": "lt",
				  "path": ".tmp",
				  "value": 30
				}
			  ]
			}
		  ]
		}
	`

	expect := QueryAstRoot{
		Lower: 10000,
		Upper: 9223372036854775807,
		Limit: 1000,
		Sort:  "desc",
		Filters: []FilterExpr{
			{
				Type:         "eq",
				Path:         ".moduleid",
				Value:        "hoge",
				ChildrenExpr: nil,
			},
			{
				Type:  "or",
				Path:  "",
				Value: nil,
				ChildrenExpr: []FilterExpr{
					{
						Type:         "gte",
						Path:         ".tmp",
						Value:        float64(10),
						ChildrenExpr: nil,
					},
					{
						Type:         "lt",
						Path:         ".tmp",
						Value:        float64(30),
						ChildrenExpr: nil,
					},
				},
			},
		},
	}
	query, err := QueryParser([]byte(queryStr))
	assert.Nil(t, err)
	assert.Equal(t, expect, *query)
}
