package main

import "bytes"

type XLSXExport struct {}

func (export *XLSXExport) Export(e *Experiment) ([]byte, error) {
	content := new(bytes.Buffer)

	return content.Bytes(), nil
}

func (export *XLSXExport) ContentType() string {
	return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
}

