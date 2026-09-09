package rdb

import (
	"fmt"
	"strings"
	"time"

	"github.com/basebytes/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Select(fields ...string) Condition {
	return func(db *gorm.DB) *gorm.DB {
		if len(fields) > 0 {
			db = db.Select(fields)
		}
		return db
	}
}

func Equal(field string, value any) Condition {
	return func(db *gorm.DB) *gorm.DB {
		if field != "" {
			db = db.Where(fmt.Sprintf("%s = ?", field), value)
		}
		return db
	}
}

func NotEqual(field string, value any) Condition {
	return func(db *gorm.DB) *gorm.DB {
		if field != "" {
			db = db.Where(fmt.Sprintf("%s != ?", field), value)
		}
		return db
	}
}

func In(field string, values ...any) Condition {
	return func(db *gorm.DB) *gorm.DB {
		if field != "" {
			db = db.Where(fmt.Sprintf("%s in ?", field), values)
		}
		return db
	}
}

func NotIn(field string, values ...any) Condition {
	return func(db *gorm.DB) *gorm.DB {
		if field != "" {
			db = db.Where(fmt.Sprintf("%s not in ?", field), values)
		}
		return db
	}
}

func Range(field string, op OpType, value any) Condition {
	return func(db *gorm.DB) *gorm.DB {
		if field == "" || value == nil {
			return db
		}
		switch op {
		case LT, GT, LTE, GTE:
			db = db.Where(fmt.Sprintf("%s %s?", field, op), value)
		}
		return db
	}
}

func Like(field string, ft FuzzyType, value string) Condition {
	var left, right string
	if ft&fuzzyTypeMsk != FuzzyTypeRight {
		left = FuzzySymbol
	}
	if ft&fuzzyTypeMsk != FuzzyTypeLeft {
		right = FuzzySymbol
	}
	return func(db *gorm.DB) *gorm.DB {
		if field != "" && value != "" {
			db = db.Where(fmt.Sprintf("%s like ?", field), fmt.Sprintf("%s%s%s", left, value, right))
		}
		return db
	}
}

func TimeScope(field string, start, end *types.Time) Condition {
	return func(db *gorm.DB) *gorm.DB {
		if field == "" || (start == nil && end == nil) {
			return db
		}
		if start != nil {
			db = db.Where(fmt.Sprintf("%s >=?", field), start.String())
		}
		if end != nil {
			db = db.Where(fmt.Sprintf("%s <=?", field), end.String())
		}
		return db
	}
}

func DateToTimeScope(field string, start, end *types.Date) Condition {
	return func(db *gorm.DB) *gorm.DB {
		if field == "" || (start == nil && end == nil) {
			return db
		}
		if start != nil {
			db = db.Where(fmt.Sprintf("%s >=?", field), start.Time.Format(time.DateTime))
		}
		if end != nil {
			end.Time = end.AddDate(0, 0, 1)
			db = db.Where(fmt.Sprintf("%s <?", field), end.Time.Format(time.DateTime))
		}
		return db
	}
}

func DateScope(field string, start, end *types.Date) Condition {
	return func(db *gorm.DB) *gorm.DB {
		if field == "" || (start == nil && end == nil) {
			return db
		}
		if start != nil {
			db = db.Where(fmt.Sprintf("%s >=?", field), start.String())
		}
		if end != nil {
			db = db.Where(fmt.Sprintf("%s <=?", field), end.String())
		}
		return db
	}
}

func DateTimeStrScope(field, start, end string) Condition {
	return func(db *gorm.DB) *gorm.DB {
		if field == "" || (start == "" && end == "") {
			return db
		}
		if start != "" {
			db = db.Where(fmt.Sprintf("%s >= ?", field), start)
		}
		if end != "" {
			db = db.Where(fmt.Sprintf("%s <= ?", field), end)
		}
		return db
	}
}

func IntScope(field string, start, end *int64) Condition {
	return func(db *gorm.DB) *gorm.DB {
		if field == "" || (start == nil && end == nil) {
			return db
		}
		if start != nil {
			db = db.Where(fmt.Sprintf("%s >=?", field), *start)
		}
		if end != nil {
			db = db.Where(fmt.Sprintf("%s <=?", field), *end)
		}
		return db
	}
}

func JsonArrayContains(field string, value any) Condition {
	return func(db *gorm.DB) *gorm.DB {
		if field != "" && value != nil {
			db = db.Where(fmt.Sprintf("JSON_CONTAINS(%s,JSON_ARRAY(?))", field), value)
		}
		return db
	}
}

func JsonArrayNotContains(field string, value any) Condition {
	return func(db *gorm.DB) *gorm.DB {
		if field != "" && value != nil {
			db = db.Where(fmt.Sprintf("NOT JSON_CONTAINS(%s,JSON_ARRAY(?))", field), value)
		}
		return db
	}
}

func JsonSearch(field string, value any) Condition {
	return func(db *gorm.DB) *gorm.DB {
		if field != "" && value != nil {
			db = db.Where(fmt.Sprintf("JSON_SEARCH(%s,'all',?)", field), value)
		}
		return db
	}
}

func JsonArrayAppendExpr(field string, value any) clause.Expr {
	return gorm.Expr(fmt.Sprintf("JSON_ARRAY_APPEND(COALESCE(%s,'[]'),'$',?)", field), value)
}

func JsonArrayRemoveExpr(field string, value any) clause.Expr {
	return gorm.Expr(fmt.Sprintf("JSON_REMOVE(%s,JSON_UNQUOTE(JSON_SEARCH(%s,'one',?)))", field, field), value)
}

func Expr(expr string, args ...any) clause.Expr {
	return gorm.Expr(expr, args...)
}

func Or(conditions ...Condition) Condition {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Or(newDryRun().Scopes(conditions...))
		return db
	}
}

func And(conditions ...Condition) Condition {
	return func(db *gorm.DB) *gorm.DB {
		if len(conditions) > 0 {
			db = db.Where(newDryRun().Scopes(conditions...))
		}
		return db
	}
}

func Preload(name string, conditions ...Condition) Condition {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Preload(name, func(db *gorm.DB) *gorm.DB {
			if len(conditions) > 0 {
				for _, cond := range conditions {
					db = cond(db)
				}
			}
			return db
		})
		return db
	}
}

func Join(name string, clauses ...clause.Expression) Condition {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Joins(name)
		if len(clauses) > 0 {
			db = db.Clauses(clauses...)
		}
		return db
	}
}

func Association(name string) Condition {
	return func(db *gorm.DB) *gorm.DB {
		if name != "" {
			db.Association(name)
		}
		return db
	}
}

func Page(offset, limit int) Condition {
	if offset < 0 {
		offset = 0
	}
	return func(db *gorm.DB) *gorm.DB {
		if limit < 0 {
			return db
		}
		return db.Limit(limit).Offset(offset)
	}
}

func Group(fields string, having ...clause.Expression) Condition {
	return func(db *gorm.DB) *gorm.DB {
		if fields != "" {
			db = db.Group(fields)
		}
		if len(having) > 0 {
			args := make([]any, len(having))
			for i, c := range having {
				args[i] = c
			}
			db.Having(args[0], args[1:]...)
		}
		return db
	}
}

func OrderBy(fields ...string) Condition {
	orders := make([]string, 0, len(fields))
	for _, field := range fields {
		subs := strings.Split(field, " ")
		order := DESC
		if len(subs) > 1 && strings.ToUpper(subs[1]) == ASC {
			order = ASC
		}
		orders = append(orders, fmt.Sprintf("%s %s", subs[0], order))
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(strings.Join(orders, ", "))
	}
}
