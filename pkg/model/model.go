package model

import "encoding/xml"

type MxFile struct {
	XMLName  xml.Name `xml:"mxfile"`
	Host     string   `xml:"host,attr"`
	Modified string   `xml:"modified,attr"`
	Agent    string   `xml:"agent,attr"`
	Etag     string   `xml:"etag,attr"`
	Version  string   `xml:"version,attr"`
	Type     string   `xml:"type,attr"`
	Diagram  struct {
		Value string `xml:",chardata"`
		Name  string `xml:"name,attr"`
		Id    string `xml:"id,attr"`
	} `xml:"diagram"`
}