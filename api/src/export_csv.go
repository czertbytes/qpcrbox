package main

import (
	"fmt"
	"bytes"
)

type CSVExport struct {}

func (export *CSVExport) Export(e *Experiment) ([]byte, error) {
	var content bytes.Buffer

	content.WriteString("name,mean,stddev\n")
	for endogenousControlName, endogenousControl := range e.EndogenousControls {
		content.WriteString(fmt.Sprintf("%s,%f,%f\n", endogenousControlName, endogenousControl.Mean, endogenousControl.StdDev))
	}

	content.WriteString("\ndetector,name,mean,stddev,dct,ddct,ddcterr,rq,rqerr\n")
	for detectorName, detector := range e.Detectors {
		for targetGeneName, targetGene := range detector {
			content.WriteString(fmt.Sprintf("%s,%s,%f,%f,%f,%f,%f,%f,%f\n", detectorName, targetGeneName, targetGene.Mean, targetGene.StdDev, targetGene.DCt, targetGene.DdCt, targetGene.DdCtErr, targetGene.RQ, targetGene.RQErr))
		}
	}

	return content.Bytes(), nil
}

func (export *CSVExport) ContentType() string {
	return "text/csv"
}
