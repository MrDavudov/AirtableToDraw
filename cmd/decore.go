package main

import (
	"encoding/base64"
	"encoding/xml"
	"os"
	"bytes"
	"compress/flate"
	"fmt"
	"io/ioutil"
	"net/url"
)

type MxFile struct {
	XMLName xml.Name `xml:"mxfile"`
	Diagram string   `xml:"diagram"`
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