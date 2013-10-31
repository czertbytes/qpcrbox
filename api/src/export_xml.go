package main

import (
	"encoding/xml"
)

type XMLExport struct {}

type XMLExportExperiment struct {
	XMLName 			xml.Name 						`xml:"experiment"`
	Detectors			[]XMLExportDetector				`xml:"detectors>detector"`
	EndogenousControls 	[]XMLExportEndogenousControl	`xml:"endogenous-controls>endogenous-control"`
}

type XMLExportDetector struct {
	XMLName 	xml.Name				`xml:"detector"`
	Name		string					`xml:"name,attr"`
	TargetGenes	[]XMLExportTargetGene	`xml:"target-genes>target-gene"`
}

type XMLExportTargetGene struct {
	XMLName 	xml.Name 	`xml:"target-gene"`
	Name 		string		`xml:"name,attr"`
	RawValues	[]string	`xml:"raw-values>raw-value"`
	Values 		[]float64	`xml:"values>value"`
	Mean		float64 	`xml:"mean"`
	StdDev		float64		`xml:"stddev"`
	DCt			float64		`xml:"dct"`
	DdCt		float64		`xml:"ddct"`
	DdCtErr		float64		`xml:"ddcterr"`
	RQ			float64		`xml:"rq"`
	RQErr		float64		`xml:"rqerr"`
}

type XMLExportEndogenousControl struct {
	XMLName 					xml.Name								`xml:"endogenous-control"`
	Name						string									`xml:"name,attr"`
	EndogenousControlDetectors 	[]XMLExportEndogenousControlDetector 	`xml:"detectors>detector"`
	Values						[]float64								`xml:"values>value"`
	Mean						float64 								`xml:"mean"`
	StdDev						float64									`xml:"stddev"`
}

type XMLExportEndogenousControlDetector struct {
	XMLName 	xml.Name	`xml:"endogenous-control"`
	Name		string		`xml:"name,attr"`
	RawValues	[]string	`xml:"raw-values>raw-value"`
}

func (export *XMLExport) Export(e *Experiment) ([]byte, error) {
	var endogenousControls  []XMLExportEndogenousControl
	for endogenousControlName, endogenousControl := range e.EndogenousControls {
		var endogenousControlDetectors = []XMLExportEndogenousControlDetector{}
		for detectorName, rawValues := range endogenousControl.Detectors {
			endogenousControlDetectors = append(endogenousControlDetectors, XMLExportEndogenousControlDetector{Name: detectorName, RawValues: rawValues})
		}

		endogenousControls = append(endogenousControls, XMLExportEndogenousControl{Name: endogenousControlName, EndogenousControlDetectors: endogenousControlDetectors, Values: endogenousControl.Values, Mean: endogenousControl.Mean, StdDev: endogenousControl.StdDev})
	}

	var detectors []XMLExportDetector
	for detectorName, detector := range e.Detectors {
		var targetGenes = []XMLExportTargetGene{}
		for targetGeneName, targetGene := range detector {
			targetGenes = append(targetGenes, XMLExportTargetGene{Name: targetGeneName, RawValues: targetGene.RawValues, Values: targetGene.Values, Mean: targetGene.Mean, StdDev: targetGene.StdDev, DCt: targetGene.DCt, DdCt: targetGene.DdCt, DdCtErr: targetGene.DdCtErr, RQ: targetGene.RQ, RQErr: targetGene.RQErr})
		}

		detectors = append(detectors, XMLExportDetector{Name: detectorName, TargetGenes: targetGenes})
	}

	experiment := XMLExportExperiment{Detectors: detectors, EndogenousControls: endogenousControls}

	xmlContent, err := xml.MarshalIndent(experiment, " ", "  ")
	if err != nil {
		return []byte{}, err
	}

	return xmlContent, nil
}

func (export *XMLExport) ContentType() string {
	return "application/xml"
}

