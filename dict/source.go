package dict

import (
	"fmt"

	"github.com/basebytes/component/database/rdb"
)

func NewNormal[T NormalDict](name string) *Normal[T] {
	return &Normal[T]{name: name}
}

type Normal[T NormalDict] struct {
	name   string
	values []Dict
}

func (d *Normal[T]) Name() string {
	return d.name
}

func (d *Normal[T]) Values() []Dict {
	return d.values
}

func (d *Normal[T]) Load(dbName string, _ map[string]any) (err error) {
	var (
		t       T
		results []T
		ins, ok = rdb.GetConnection(dbName)
	)
	if ok {
		if err = ins.FindByCondition(t, &results).Error; err == nil {
			d.values = make([]Dict, 0, len(results))
			for _, result := range results {
				d.values = append(d.values, result.Trans()...)
			}
		}
	} else {
		err = fmt.Errorf("rdb[%s] instance not found", dbName)
	}
	return
}
