package storage

type Storager interface {
	Configure() (Storager, error)
	CreateDB(name, schema string) error
}