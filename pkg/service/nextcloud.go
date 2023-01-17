package service

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/MrDavudov/AirtableToDraw/pkg/model"
	"github.com/studio-b12/gowebdav"
)

type NCService struct{
	client	*gowebdav.Client
}

func NewNextcloudService(n *gowebdav.Client) *NCService {
	return &NCService{
		client: n,
	}
}

const path = "Marketplace"

func (s *NCService) Get() ([]string, error) {
	// подключение
	// client := gowebdav.NewClient(root, username, password)

	// зайти в папку Marketplace и показать все файлы
	files, err := s.client.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}

	var FileName []string
	for _, file := range files {
		FileName = append(FileName, file.Name())
		// fmt.Println(file.Name())
	}

	return FileName, nil
}

func Decode(file string) (string, error) {
	read, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err)
	}

	var p model.MxFile
	if err := xml.Unmarshal(read, &p); err != nil {
		fmt.Println(err)
	}

	// decode base64
	rawDecodedText, err := base64.StdEncoding.DecodeString(p.Diagram.Value)
	if err != nil {
		fmt.Println(err)
	}
	// deflate
	enflated, err := io.ReadAll(flate.NewReader(bytes.NewReader(rawDecodedText)))
	if err != nil {
		fmt.Println(err)
	}
	// decode url
	query, err := url.QueryUnescape(string(enflated))
	if err != nil {
		fmt.Println(err)
	}

	out := strings.Replace(query, "ЗАДАТЬ СРОК ХОЛДИРОВАНИЯ МЕТОДА ОПЛАТЫ", "wolf", -1)
	
	// запись в временный xml
	f, err := os.Create("copy_" + file)  // file[:len(file)-3]
	if err != nil{
		fmt.Println(err)
	}
	defer f.Close()
	f.WriteString(out)

	return "", nil
}