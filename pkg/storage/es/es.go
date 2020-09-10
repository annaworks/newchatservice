package es

import (
	"context"
	"errors"

	"github.com/annaworks/surubot/pkg/storage"
	"github.com/olivere/elastic/v7"
)

const (
	DocType = "_doc"
)

type es struct {
	client *elastic.Client
	host string
}

func New(host string) *es {
	return &es{
		host: host,
		client: nil,
	}
}

func (e *es) Configure() (storage.Storager, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(e.host),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)
	if err != nil {
		return e, err
	}

	e.client = client
	return e, nil
}

func (e es) CreateDB(name, schema string) error {
	ctx := context.Background()

	res, err := e.client.CreateIndex(name).BodyString(schema).Do(ctx)
	if err != nil {
		return err
	}

	if !res.Acknowledged {
		return errors.New("Registering was not acknowledged")
	}
	
	return nil
}

func (e es) DBExists(name string) (bool, error) {
	ctx := context.Background()

	exists, err := e.client.IndexExists(name).Do(ctx)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (e es) Insert(name string, body interface{}) (interface{}, error) {
	ctx := context.Background()
	resp, err := e.client.Index().
								Type(DocType).
								Index(name).
								BodyJson(body).
								Do(ctx)
	if err != nil {
		return "", err
	}

	return resp, nil
}