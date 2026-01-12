package dict

import (
	"fmt"

	"github.com/basebytes/component/database/rdb"
)

type Source interface {
	Name() string
	Values() []Dict
	Load(dbName string, params map[string]any) (err error)
}

type DictModifier interface {
	NormalDict
	IgnoreUpdate() bool
}

type NormalDict interface {
	rdb.Data
	Trans() []Dict
}

type Dict interface {
	rdb.Data
	Enum() *Enum
	GetCategory() string
	GetKey() string
	GetMappingKey() string
	GetStatus() int
	UpdateFlag() byte
}

func NewDefaultEnum(key, value string, category Category) *Enum {
	return &Enum{Key: key, Value: value, category: category}
}

func NewEnum(key, value string, seq, status int, category Category) *Enum {
	return &Enum{Key: key, Value: value, Seq: seq, Status: status, category: category}
}

type Enum struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Seq      int    `json:"seq,omitempty"`
	Status   int    `json:"status,omitempty"`
	Children *Enums `json:"children,omitempty"`
	category Category
}

func (e *Enum) Unique() string {
	return fmt.Sprintf("%s/%s/%d", e.category, e.Key, e.Status)
}

func (e *Enum) update(newValue *Enum, updateFlag byte) {
	if newValue == nil {
		return
	}
	if updateFlag &= UpdateFlagMask; updateFlag == 0 {
		return
	}
	if updateFlag&UpdateFlagValue == UpdateFlagValue {
		e.Value = newValue.Value
	}
	if updateFlag&UpdateFlagSeq == UpdateFlagSeq {
		e.Seq = newValue.Seq
	}
	if updateFlag&UpdateFlagStatus == UpdateFlagStatus {
		e.Status = newValue.Status
	}
}

func (e *Enum) SetChildren(children *Enums) *Enum {
	e.Children = children
	return e
}

func (e *Enum) AppendChild(enum *Enum) {
	if enum == nil {
		return
	}
	if e.Children == nil {
		e.Children = NewEnums(10)
	}
	e.Children.Append(enum)
}

func NewEnums(size int) *Enums {
	e := Enums(make([]*Enum, 0, size))
	return &e
}

type Enums []*Enum

func (e *Enums) Append(ne *Enum) {
	if ne == nil {
		return
	}
	if e.find(ne.Unique()) < 0 {
		*e = append(*e, ne)
	}
}

func (e *Enums) Remove(key string) {
	if idx := e.find(key); idx == 0 {
		*e = (*e)[1:]
	} else if idx == len(*e)-1 {
		*e = (*e)[:idx]
	} else if idx > 0 {
		_enums := Enums(make([]*Enum, 0, len(*e)-1))
		copy(_enums[:idx], (*e)[:idx])
		copy(_enums[idx:], (*e)[idx+1:])
		*e = _enums
	}
}
func (e *Enums) Len() int {
	return len(*e)
}

func (e *Enums) find(key string) (idx int) {
	idx = -1
	for i, enum := range *e {
		if enum.Unique() == key {
			idx = i
			break
		}
	}
	return
}

type Category = string

const (
	UpdateFlagValue   byte = 1
	UpdateFlagSeq     byte = 2
	UpdateFlagStatus  byte = 4
	UpdateFlagMapping byte = 8
	UpdateFlagMask         = UpdateFlagValue | UpdateFlagSeq | UpdateFlagStatus
)
