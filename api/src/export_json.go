package main

import (
	"log"
	"encoding/json"
)

type JSONExport struct {}

func (export *JSONExport) Export(e *Experiment) ([]byte, error) {
	content, err := json.Marshal(e)
	if err != nil {
		log.Printf("[export|json] experiment marshalling failed! Error: '%s'\n", err)
		return []byte{}, err
	}

	return content, nil
}

func (export *JSONExport) ContentType() string {
	return "application/json"
}
