package rdb

import (
	"fmt"
	"strings"

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

func Or(conditions ...Condition) Condition {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Or(newDryRun().Scopes(conditions...))
		return db
	}
}

func And(conditions ...Condition) Condition {
	return func(db *gorm.DB) *gorm.DB {
		if len(conditions) > 0 {
			var db1 *gorm.DB
			for _, cond := range conditions {
				db1 = cond(db)
			}
			db = db.Where(db1)
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

func Group(fields string) Condition {
	return func(db *gorm.DB) *gorm.DB {
		if fields != "" {
			return db.Group(fields)
		}
		return db
	}
}

func OrderBy(field string, order ...OrderType) Condition {
	oderBy := DESC
	if len(order) > 0 && strings.ToUpper(order[0]) == ASC {
		oderBy = ASC
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(fmt.Sprintf("%s %s", field, oderBy))
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

func JsonContains(field string, value any) Condition {
	return func(db *gorm.DB) *gorm.DB {
		if field != "" && value != nil {
			db = db.Where(fmt.Sprintf("JSON_CONTAINS(%s,JSON_ARRAY(?))", field), value)
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

func TimeScope(field string, start, end *types.Time) Condition {
	return func(db *gorm.DB) *gorm.DB {
		if start == nil && end == nil {
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
