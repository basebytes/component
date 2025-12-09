package rdb

import (
	"fmt"

	"gorm.io/gorm/clause"
)

func EqualClause(field string, value any) clause.Expression {
	return clause.Eq{
		Column: field,
		Value:  value,
	}
}

func LikeClause(field string, ft FuzzyType, value string) clause.Expression {
	var left, right string
	if ft&fuzzyTypeMsk != FuzzyTypeRight {
		left = FuzzySymbol
	}
	if ft&fuzzyTypeMsk != FuzzyTypeLeft {
		right = FuzzySymbol
	}
	return clause.Like{
		Column: field,
		Value:  fmt.Sprintf("%s%s%s", left, value, right),
	}
}

func GTEClause(field string, value any) clause.Expression {
	return clause.Gte(eqClause(field, value))
}

func LTEClause(field string, value any) clause.Expression {
	return clause.Lte(eqClause(field, value))
}

func eqClause(field string, value any) clause.Eq {
	return clause.Eq{Column: field, Value: value}
}

func JoinClause(field string) clause.Expression {
	return clause.Join{
		Type:       clause.LeftJoin,
		Table:      clause.Table{},
		ON:         clause.Where{},
		Using:      nil,
		Expression: nil,
	}
}
