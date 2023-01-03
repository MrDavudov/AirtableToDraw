package main

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"encoding/xml"
	// "fmt"
	"io"
	"net/url"
	"os"
)

type MxFile struct {
	XMLName xml.Name `xml:"mxfile"`
	Diagram string   `xml:"diagram"`
}

func Decode(file string) (string, error) {
	b, _ := os.ReadFile(file)

	var p MxFile
	if err := xml.Unmarshal(b, &p); err != nil {
		return "", err
	}

	// decode base64
	rawDecodedText, err := base64.StdEncoding.DecodeString(p.Diagram)
	if err != nil {
		return "", err
	}

	// deflate
	enflated, err := io.ReadAll(flate.NewReader(bytes.NewReader(rawDecodedText)))
	if err != nil {
		return "", err
	}
	
	// url decode
	query, err := url.QueryUnescape(string(enflated))
	if err != nil {
		return "", err
	}
	// fmt.Println(query)
	
	f, err := os.Create("copy_" + file)
	if err != nil{
		return "", err 
	}
	defer f.Close()
	f.WriteString(query)
	 
	return query, nil
}