package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
    "github.com/mehanizm/airtable"
)

func main() {
	err := godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
	}

    client := airtable.NewClient(os.Getenv("API"))

    table := client.GetTable("appTh2xBTnix4oi1d", "Features")
    fmt.Println(table)

    records, err := table.GetRecords().
	FromView("Table view").
	ReturnFields("CODE" ,"Name").
	InStringFormat("Europe/Moscow", "ru").
	Do()
    if err != nil {
        panic(err)
    }

	fmt.Println("Best Public Domain Books: ")

	for _, tableRecord := range records.Records {
        fmt.Println(tableRecord.Fields["CODE"])
    }
}