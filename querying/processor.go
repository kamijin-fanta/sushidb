package querying

import (
	"errors"
	"github.com/oliveagle/jsonpath"
)

type QueryProcessor struct {
	Query QueryAstRoot
}

func New(queryData []byte) (*QueryProcessor, error) {
	parsedQuery, err := QueryParser(queryData)
	if err != nil {
		return nil, err
	}
	processor := QueryProcessor{
		Query: *parsedQuery,
	}
	return &processor, nil
}

func (p *QueryProcessor) FilterRow(row interface{}) (bool, error) {
	return EvaluationFilter(p.Query.Filters, row, true)
}

func EvaluationFilter(filters []FilterExpr, row interface{}, mustAll bool) (bool, error) {
	condition := false

	if len(filters) == 0 {
		return true, nil
	}

	for _, expr := range filters {
		switch expr.Type {
		case "eq":
			res, err := jsonpath.JsonPathLookup(row, expr.Path)
			if err != nil {
				return false, err
			}
			switch v := res.(type) {
			case float64:
				switch exprValue := expr.Value.(type) {
				case float64:
					condition = condition || floatEquals(v, exprValue)
					break
				case int:
					condition = condition || floatEquals(v, float64(exprValue))
					break
				default:
					break
				}

				break
			default:
				condition = condition || res == expr.Value
				break
			}
			break
		case "and":
			res, err := EvaluationFilter(*expr.ChildrenExpr, row, true)
			if err != nil {
				return false, nil
			}
			condition = condition || res
			break
		case "or":
			res, err := EvaluationFilter(*expr.ChildrenExpr, row, false)
			if err != nil {
				return false, nil
			}
			condition = condition || res
			break
		default:
			return false, errors.New("undefined expression type '" + expr.Type + "'")
		}

		if mustAll && !condition {
			return false, nil
		}
	}

	return condition, nil
}

var EPSILON = 0.00000001

func floatEquals(a, b float64) bool {
	return (a-b) < EPSILON && (b-a) < EPSILON
}
