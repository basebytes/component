package dict

type Source interface {
	Load(*SourceConfig) error
	Name() string
	Add(any)
	Update(any)
	Remove(any)
	Get(any) any
	ExistKey(string, string) bool
	AllDictItems() map[string]*Items
	DictItems(string) *Items
	DictItem(string, string) Item
	AllKVs() map[string]map[string]string
	GetKV(string) map[string]string
	AllReversedKVs() map[string]map[string]string
	GetReversedKV(string) map[string]string
}

type EmptySource struct {
}

func (c *EmptySource) Load(*SourceConfig) error                     { return nil }
func (c *EmptySource) Name() string                                 { return "empty" }
func (c *EmptySource) Add(any)                                      {}
func (c *EmptySource) Update(any)                                   {}
func (c *EmptySource) Remove(any)                                   {}
func (c *EmptySource) Get(any) any                                  { return nil }
func (c *EmptySource) ExistKey(string, string) bool                 { return false }
func (c *EmptySource) AllDictItems() map[string]*Items              { return emptyItems }
func (c *EmptySource) DictItems(string) *Items                      { return emptyItem }
func (c *EmptySource) DictItem(string, string) Item                 { return nil }
func (c *EmptySource) AllKVs() map[string]map[string]string         { return emptyMaps }
func (c *EmptySource) GetKV(string) map[string]string               { return emptyMap }
func (c *EmptySource) AllReversedKVs() map[string]map[string]string { return emptyMaps }
func (c *EmptySource) GetReversedKV(string) map[string]string       { return emptyMap }

var (
	emptyMap   = make(map[string]string)
	emptyMaps  = make(map[string]map[string]string)
	emptyItems = make(map[string]*Items)
	emptyItem  = NewItems()
)
