package main

type DetectorMap map[string]DetectorTargetGeneMap
type DetectorTargetGeneMap map[string]DetectorTargetGene
type EndoTargetGeneMap map[string]EndoTargetGene
type StringArrayMap map[string][]string

type Experiment struct {
	Detectors DetectorMap
	EndogenousControls EndoTargetGeneMap
}

type DetectorTargetGene struct {
	RawValues                                   []string
	Values                                      []float64
	Mean, StdDev, DCt, DdCt, DdCtErr, RQ, RQErr float64
}

type EndoTargetGene struct {
	Detectors   	StringArrayMap
	Values       	[]float64
	Mean, StdDev 	float64
}

type ExperimentComputer interface {
	Compute() (*Experiment, error)
}


type AB7300 struct {
	Content, Mock string
}
