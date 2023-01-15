package service

import (
	"github.com/mehanizm/airtable"
	"github.com/studio-b12/gowebdav"
)

type Airtable interface {
	Get() ([]string, error)
}

type NextCloud interface {
	Get() ([]string, error)
}

type Service struct {
	Airtable
	NextCloud
}

func New(a *airtable.Client, n *gowebdav.Client) *Service {
	return &Service{
		Airtable:  NewAirtableService(a),
		NextCloud: NewNextcloudService(n),
	}
}