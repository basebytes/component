package dict

import "sort"

func NewItems() *Items {
	return &Items{}
}

type Items []Item

func (e *Items) Len() int {
	return len(*e)
}

func (e *Items) Swap(i, j int) {
	(*e)[i], (*e)[j] = (*e)[j], (*e)[i]
}

func (e *Items) Less(i, j int) bool {
	return (*e)[i].Seq() < (*e)[j].Seq() || (*e)[i].Key() < (*e)[j].Key()
}

func (e *Items) Append(r Item) {
	if r == nil {
		return
	}
	if !e.Exist(r.Key()) {
		*e = append(*e, r)
		sort.Sort(e)
	}
}

func (e *Items) Remove(key string) (deleted Item) {
	if idx := e.find(key); idx == 0 {
		*e = (*e)[1:]
	} else if idx == e.Len()-1 {
		*e = (*e)[:idx]
	} else if idx > 0 {
		deleted = (*e)[idx]
		newEnums := make([]Item, e.Len()-1)
		copy(newEnums[:idx], (*e)[:idx])
		copy(newEnums[idx:], (*e)[idx+1:])
		*e = newEnums
	}
	return deleted
}

func (e *Items) Get(key string) (Item, bool) {
	if idx := e.find(key); idx >= 0 {
		return (*e)[idx], true
	}
	return nil, false
}

func (e *Items) Exist(key string) bool {
	return e.find(key) >= 0
}

func (e *Items) find(key string) int {
	for i, item := range *e {
		if item.Key() == key {
			return i
		}
	}
	return -1
}

type Item interface {
	Key() string
	Value() string
	Seq() int
}
