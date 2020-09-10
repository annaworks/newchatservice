package storage

type Storager interface {
	Configure() (Storager, error)
	CreateIndex(name, schema string) error
	IndexExists(name string) (bool, error)
	Insert(name string, body interface{}) (interface{}, error)
}