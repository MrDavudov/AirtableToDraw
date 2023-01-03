package main

import (
	"fmt"

	"github.com/studio-b12/gowebdav"
)

func main() {
	root := "https://cloud.05.ru"
	username := "daudov.r"
	password := "f41i4zvh"
	path := "Marketplace"

	// подключение
	client := gowebdav.NewClient(root, username, password)
	client.Connect()
	
	fmt.Println(client.Connect())

	// создать папку
	err := client.Mkdir("folder", 0644)
	if err != nil {
		fmt.Println(err)
	}

	// зайти в папку Marketplace
	files, err := client.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(files)

	for _, file := range files {
		//notice that [file] has os.FileInfo type
		fmt.Println(file.Name())
	}

	info, err := client.Stat(path)
	if err != nil {
		fmt.Println(err)
	}
	//notice that [info] has os.FileInfo type
	fmt.Println(info)
}