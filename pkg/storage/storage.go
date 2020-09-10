package storage

type Storager interface {
	Configure() (Storager, error)
	CreateDB(name, schema string) error
	DBExists(name string) (bool, error)
	Insert(name string, body interface{}) (interface{}, error)
}