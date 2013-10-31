package main

import (
	"log"
	"archive/zip"
)

type Exporter interface {
	Export(e *Experiment) ([]byte, error)
	ContentType() string
}

type FileData struct {
	Name, Content string
}

func addToArchive(zw *zip.Writer, fd FileData) {
	f, err := zw.Create(fd.Name)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write([]byte(fd.Content))
	if err != nil {
		log.Fatal(err)
	}
}
