package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"github.com/mesameen/micro-app/metadata/pkg/model"
	"github.com/mesameen/micro-app/src/api/gen"
	"google.golang.org/protobuf/proto"
)

var metadata = &model.Metadata{
	ID:          "1",
	Title:       "title",
	Description: "description",
	Director:    "director",
}

var genMetadata = &gen.Metadata{
	Id:          "1",
	Title:       "title",
	Description: "description",
	Director:    "director",
}

func main() {
	jsonBytes, err := serializeToJSON(metadata)
	if err != nil {
		panic(err)
	}

	xmlBytes, err := serializeToXML(metadata)
	if err != nil {
		panic(err)
	}

	protoBytes, err := serializeToProto(genMetadata)
	if err != nil {
		panic(err)
	}

	fmt.Printf("JSON size: \t%dB\n", len(jsonBytes))
	fmt.Printf("XML size: \t%dB\n", len(xmlBytes))
	fmt.Printf("PROTO size: \t%dB\n", len(protoBytes))
}

func serializeToJSON(m *model.Metadata) ([]byte, error) {
	return json.Marshal(m)
}

func serializeToXML(m *model.Metadata) ([]byte, error) {
	return xml.Marshal(m)
}

func serializeToProto(m *gen.Metadata) ([]byte, error) {
	return proto.Marshal(m)
}
