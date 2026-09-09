package rdb

import (
	"fmt"

	"gorm.io/gorm"
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

func GTClause(field string, value any) clause.Expression {
	return clause.Gt(eqClause(field, value))
}

func LTEClause(field string, value any) clause.Expression {
	return clause.Lte(eqClause(field, value))
}

func LTClause(field string, value any) clause.Expression {
	return clause.Lt(eqClause(field, value))
}

func eqClause(field string, value any) clause.Eq {
	return clause.Eq{Column: field, Value: value}
}

func RangeClause(field string, op OpType, value any) (exp clause.Expression) {
	switch op {
	case LT:
		exp = LTClause(field, value)
	case GT:
		exp = GTClause(field, value)
	case LTE:
		exp = LTEClause(field, value)
	case GTE:
		exp = GTEClause(field, value)
	default:
		return
	}
	return
}

func RangeExpr(field string, op OpType, value any) (exp clause.Expr) {
	switch op {
	case LT, LTE, GT, GTE:
		exp = gorm.Expr(fmt.Sprintf("%s %s ?", field, op), value)
	default:
		return
	}
	return
}
