package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mehanizm/airtable"
)

func Airtable() ([]interface{}, error) {
	// получаем 
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// подключения airtable
	client := airtable.NewClient(os.Getenv("API"))

	// получаем таблицу
	table := client.GetTable("appTh2xBTnix4oi1d", "Features")

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
	var obj []interface{}
	for _, tableRecord := range records.Records {
		obj = append(obj, tableRecord.Fields["CODE"])
		// fmt.Println(tableRecord.Fields["CODE"])
	}

	return obj, nil
}