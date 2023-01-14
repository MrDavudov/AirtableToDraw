package main

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
)

type MxFile struct {
	XMLName xml.Name 	`xml:"mxfile"`
	Host	string	 	`xml:"host,attr"`
	Modified	string	`xml:"modified,attr"`
	Agent	string		`xml:"agent,attr"`
	Etag	string		`xml:"etag,attr"`
	Version	string		`xml:"version,attr"`
	Type	string		`xml:"type,attr"`
	Diagram struct{
		Value	string		`xml:",chardata"`
		Name	string		`xml:"name,attr"`
		Id		string		`xml:"id,attr"`
	}   	`xml:"diagram"`
}

func Decode(file string) (string, error) {
	read, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err)
	}

	var p MxFile
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