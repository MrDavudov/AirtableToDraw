package service

import (
	"fmt"

	"github.com/mehanizm/airtable"
)

type ATService struct {
	client	*airtable.Client
}

func NewAirtableService(a *airtable.Client) *ATService {
	return &ATService{
		client: a,
	}
}

func (s *ATService) Get() ([]string, error) {
	// подключения airtable
	// client := airtable.NewClient(os.Getenv("API"))

	// получаем таблицу
	table := s.client.GetTable("appTh2xBTnix4oi1d", "Features")

	// получаем данные из таблиц
	records, err := table.GetRecords().
		FromView("Table view").	// название таблицы
		ReturnFields("CODE", "Name"). // имена столбцов
		InStringFormat("Europe/Moscow", "ru").
		Do()
	if err != nil {
		return nil, err
	}

	// сохраняем данные
	var obj []string
	for _, tableRecord := range records.Records {
		v := fmt.Sprintf("%v", tableRecord.Fields["CODE"])
		obj = append(obj, v)
		// fmt.Println(tableRecord.Fields["CODE"])
	}

	return obj, nil
}