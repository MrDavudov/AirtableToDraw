package main

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/url"

	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/mehanizm/airtable"
)

func main() {
	// Airtable()

    Decode()
}

type MxFile struct {
    XMLName xml.Name `xml:"mxfile"`
	Diagram string	 `xml:"diagram"`
}

func Decode() {
	b, _ := os.ReadFile("draw.drawio.xml")

    var p MxFile
    if err := xml.Unmarshal(b, &p); err != nil {
        panic(err)
    }

	// decode base64
	rawDecodedText, err := base64.StdEncoding.DecodeString(p.Diagram)
	if err != nil {
		panic(err)
	}

	// deflate
	enflated, _ := ioutil.ReadAll(flate.NewReader(bytes.NewReader(rawDecodedText)))
	
	// url decode
	query, err := url.QueryUnescape(string(enflated))
	if err != nil {
		panic(err)
	}
	fmt.Println(query)
}

func Airtable() {
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