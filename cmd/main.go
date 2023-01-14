package main

import (
	"fmt"

	"github.com/studio-b12/gowebdav"
)

var fileName []string

func main() {
	root := "https://cloud.05.ru/remote.php/dav/files/16cd10dc-1afa-103d-80cb-4777a0293e48/"
	username := "daudov.r"
	password := "f41i4zvh"
	path := "Marketplace"

	// подключение
	client := gowebdav.NewClient(root, username, password)

	// зайти в папку Marketplace и показать все файлы
	files, err := client.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		fileName = append(fileName, file.Name())
		fmt.Println(file.Name())
	}

	// Airtable()
	Decode("Задать срок холдирования метода оплаты.drawio.xml")
	// DecodePath("Задать срок холдирования метода оплаты.drawio.xml")
	// object[name_ucd=""]
}