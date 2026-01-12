package dict

import (
	"fmt"

	"github.com/basebytes/types"
)

type BizDict struct {
	Id         int64       `mapstructure:"id" gorm:"column:id" json:"id,omitempty"`
	Category   Category    `mapstructure:"category" gorm:"column:category" json:"category,omitempty"`
	Key        string      `mapstructure:"key" gorm:"column:key" json:"key,omitempty"`
	Value      string      `mapstructure:"value" gorm:"column:value" json:"value,omitempty"`
	MappingKey string      `mapstructure:"mapping_key" gorm:"column:mapping_key" json:"mappingKey,omitempty"`
	Seq        *int        `mapstructure:"seq" gorm:"column:seq" json:"seq,omitempty"`
	Status     *int        `mapstructure:"status" gorm:"column:status;default:0" json:"status,omitempty"`
	CreateTime *types.Time `mapstructure:"create_time" gorm:"column:create_time;default:CURRENT_TIMESTAMP();<-:create" json:"createTime,omitempty"`
	UpdateTime *types.Time `mapstructure:"update_time" gorm:"column:update_time;default:CURRENT_TIMESTAMP()" json:"updateTime,omitempty"`
}

func (d *BizDict) TableName() string {
	return "biz_dict"
}

func (d *BizDict) Enum() *Enum {
	return NewEnum(d.Key, d.Value, d.GetSeq(), d.GetStatus(), d.Category)
}

func (d *BizDict) GetCategory() string {
	return d.Category
}

func (d *BizDict) GetKey() string {
	return d.Key
}

func (d *BizDict) GetMappingKey() string {
	return d.MappingKey
}

func (d *BizDict) GetStatus() (status int) {
	if d.Status != nil && *d.Status != StatusEnable {
		status = StatusDisable
	}
	return
}

func (d *BizDict) UpdateFlag() (flag byte) {
	if d.Value != "" {
		flag |= UpdateFlagValue
	}
	if d.Seq != nil {
		flag |= UpdateFlagSeq
	}
	if d.Status != nil {
		flag |= UpdateFlagStatus
	}
	return
}

func (d *BizDict) Trans() []Dict {
	return []Dict{d}
}

func (d *BizDict) Unique() string {
	return fmt.Sprintf("%s/%s/%d", d.Category, d.Key, d.GetStatus())
}

func (d *BizDict) GetSeq() (seq int) {
	if d.Seq != nil {
		seq = *d.Seq
	}
	return
}

func (d *BizDict) updateMigration(newValue *BizDict) (update bool) {
	if newValue != nil {
		newValue.Id = d.Id
		newValue.Category = ""
		newValue.Key = ""
		newValue.MappingKey = ""
		if d.Value == newValue.Value {
			newValue.Value = ""
		} else {
			d.Value = newValue.Value
			update = true
		}
		if d.GetSeq() == newValue.GetSeq() {
			newValue.Seq = nil
		} else {
			d.Seq = newValue.Seq
			update = true
		}
		if d.GetStatus() == newValue.GetStatus() {
			newValue.Status = nil
		} else {
			d.Status = newValue.Status
			update = true
		}
	}
	return
}

func (d *BizDict) clear() *BizDict {
	d.Id = 0
	d.MappingKey = ""
	d.Status = nil
	if d.GetSeq() == 0 {
		d.Seq = nil
	}
	return d
}

const (
	StatusEnable  int = 0
	StatusDisable int = 1
)
