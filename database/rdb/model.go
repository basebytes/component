package rdb

import (
	"fmt"
	"sync/atomic"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type Data interface {
	TableName() string
}

type Instance struct {
	name  string
	debug atomic.Bool
	db    *gorm.DB
	cfg   *Config
}

func NewInstance(name string, cfg *Config) (ins *Instance, err error) {
	if cfg == nil {
		err = fmt.Errorf("database [%s] config not found", name)
		return
	}
	ins = &Instance{name: name, cfg: cfg}
	if ins.db, err = gorm.Open(cfg.Dial(), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}); err == nil {
		ins.debug.Store(false)
		sqlDB, _ := ins.db.DB()
		sqlDB.SetMaxIdleConns(ins.cfg.MaxIdleCons)
		sqlDB.SetMaxOpenConns(ins.cfg.MaxOpenCons)
	} else {
		err = fmt.Errorf("create database[%s] instance failed :%s", name, err)
	}
	return
}

func (ins *Instance) Name() string {
	return ins.name
}

func (ins *Instance) DBName() string {
	return ins.cfg.DataBase
}

func (ins *Instance) EnableDebug() {
	ins.debug.Store(true)
}

func (ins *Instance) DisableDebug() {
	ins.debug.Store(false)
}

func (ins *Instance) DB() *gorm.DB {
	if ins.debug.Load() {
		return ins.db.Debug()
	} else {
		return ins.db
	}
}

func (ins *Instance) Create(table any) *gorm.DB {
	return ins.DB().Create(table)
}

func (ins *Instance) FirstOrCreate(data Data, condition ...Condition) *gorm.DB {
	return ins.DB().Model(data).Scopes(condition...).FirstOrCreate(data)
}

func (ins *Instance) CreateIgnoreConflicts(conflicts []string, values any) *gorm.DB {
	columns := make([]clause.Column, 0, len(conflicts))
	for _, column := range conflicts {
		columns = append(columns, clause.Column{Name: column})
	}
	return ins.DB().Clauses(clause.OnConflict{
		Columns:   columns,
		DoNothing: true,
	}).Create(values)
}

func (ins *Instance) AssociationUpdatesNotEmpty(table Data, name string, values []Data) error {
	v := make([]any, 0, len(values))
	for _, value := range values {
		v = append(v, value)
	}
	return ins.DB().Model(table).Updates(table).Association(name).Append(v...)
}

func (ins *Instance) UpdatesNotEmpty(table Data) *gorm.DB {
	return ins.DB().Updates(table).First(table)
}

func (ins *Instance) UpdatesWithCondition(table Data, values any, condition ...Condition) *gorm.DB {
	return ins.DB().Model(table).Scopes(condition...).Updates(values)
}

func (ins *Instance) UpdatesByCondition(table Data, condition ...Condition) *gorm.DB {
	return ins.DB().Model(table).Scopes(condition...).Updates(table)
}

func (ins *Instance) UpdateColumn(table Data, column string, value any) *gorm.DB {
	return ins.DB().Model(table).UpdateColumn(column, value)
}

func (ins *Instance) UpdateColumns(table Data, values any) *gorm.DB {
	return ins.DB().Model(table).UpdateColumns(values)
}

func (ins *Instance) UpdateColumnsById(table Data, columns ...string) *gorm.DB {
	return ins.DB().Model(table).Select(columns).Updates(table)
}

func (ins *Instance) Upsert(conflicts, updates []string, values any) *gorm.DB {
	columns := make([]clause.Column, 0, len(conflicts))
	for _, column := range conflicts {
		columns = append(columns, clause.Column{Name: column})
	}
	return ins.DB().Clauses(clause.OnConflict{
		Columns:   columns,
		DoUpdates: clause.AssignmentColumns(updates),
	}).Create(values)
}

func (ins *Instance) FindById(table Data, id any) *gorm.DB {
	return ins.DB().First(table, id)
}

func (ins *Instance) FindByCondition(table Data, result any) *gorm.DB {
	return ins.DB().Where(table).Find(result)
}

func (ins *Instance) FindFirstByCondition(table Data) *gorm.DB {
	return ins.DB().Where(table).First(table)
}

func (ins *Instance) GetData(table Data, result any, conditions ...Condition) *gorm.DB {
	return ins.DB().Model(table).Where(table).Scopes(conditions...).Find(result)
}

func (ins *Instance) SubQuery(table Data, conditions ...Condition) *gorm.DB {
	return ins.DB().Model(table).Where(table).Scopes(conditions...)
}

func (ins *Instance) DeleteByCondition(table Data, condition ...Condition) *gorm.DB {
	return ins.DB().Model(table).Scopes(condition...).Delete(table)
}

func (ins *Instance) Count(table Data, condition ...Condition) (count int64, err error) {
	err = ins.DB().Model(table).Scopes(condition...).Count(&count).Error
	return
}

func (ins *Instance) PageQuery(table Data, result any, page Condition, conditions ...Condition) (count int64, err error) {
	err = ins.DB().Model(table).Where(table).Scopes(conditions...).Count(&count).Scopes(page).Find(result).Error
	return
}

func (ins *Instance) Raw(sql string, args ...any) *gorm.DB {
	return ins.DB().Raw(sql, args...)
}

func (ins *Instance) Transaction(fc func(tx *gorm.DB) error) error {
	return ins.DB().Transaction(fc)
}

func (ins *Instance) OrClause(condition ...Condition) *gorm.DB {
	return ins.DB().Scopes(condition...)
}

func newDryRun() *gorm.DB {
	return &gorm.DB{
		Config: &gorm.Config{
			DryRun: true,
		},
		Statement: &gorm.Statement{
			Clauses: map[string]clause.Clause{},
		},
	}
}

type Condition = func(db *gorm.DB) *gorm.DB
type OrderType = string
type OpType = string
type FuzzyType int8

const (
	ASC  OrderType = "ASC"
	DESC OrderType = "DESC"
)

const (
	LT  OpType = "<"
	GT  OpType = ">"
	LTE OpType = "<="
	GTE OpType = ">="
)

const (
	FuzzyTypeLeft  FuzzyType = 1
	FuzzyTypeRight FuzzyType = 2
	FuzzyTypeBoth  FuzzyType = 4
	fuzzyTypeMsk             = FuzzyTypeLeft | FuzzyTypeRight | FuzzyTypeBoth
	FuzzySymbol              = "%"
)
