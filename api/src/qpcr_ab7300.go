package main

import (
	"math"
	"strconv"
	"strings"
	"log"
	"errors"
)

func (md *AB7300) Compute() (*Experiment, error) {
	e := &Experiment{}
	e.Detectors = make(DetectorMap)
	e.EndogenousControls = make(EndoTargetGeneMap)

	if isContentValid(md.Content) {
		section := 1
		lines := strings.Split(md.Content, "\n")
		for _, line := range lines {
			if len(line) == 0 {
				section++
			} else {
				switch section {
				case 11:
					e.parseRow(&line)
				}
			}
		}

		e.computeTargetGenes(md.Mock)

		return e, nil
	}

	return e, errors.New("[ab7300] content is not valid!")
}

func isContentValid(content string) bool {
	return (strings.Contains(content, "Applied Biosystems 7300 Real-Time PCR System") && strings.Contains(content, "SDS v1.4"))
}

func (e *Experiment) computeTargetGenes(mockName string) {
	endoControlMock := e.EndogenousControls[mockName]

	e.computeMocks(mockName, endoControlMock)

	for detectorName, detector := range e.Detectors {
		for targetGeneName, targetGene := range detector {
			if targetGeneName != mockName {
				targetGeneMock := e.Detectors[detectorName][mockName]

				targetGene.mergeRawValues()

				targetGene.Mean, targetGene.StdDev = meanAndStdDev(targetGene.Values)
				targetGene.DCt = targetGene.Mean - endoControlMock.Mean
				targetGene.DdCt = targetGene.DCt - targetGeneMock.DCt
				targetGene.DdCtErr = math.Sqrt((2 * (endoControlMock.StdDev * endoControlMock.StdDev)) + (2 * (targetGene.StdDev * targetGene.StdDev)))
				targetGene.RQ = math.Pow(2, (-1) * targetGene.DdCt)
				targetGene.RQErr = math.Sqrt(targetGene.RQ * targetGene.RQ * math.Ln2 * math.Ln2 * targetGene.DdCtErr * targetGene.DdCtErr)

				e.Detectors[detectorName][targetGeneName] = targetGene
			}
		}
	}
}

func (e *Experiment) computeMocks(mockName string, endoControlMock EndoTargetGene) {
	for detectorName, detector := range e.Detectors {
		if _, found := detector[mockName]; found == true {
			targetGene := detector[mockName]

			targetGene.mergeRawValues()

			targetGene.Mean, targetGene.StdDev = meanAndStdDev(targetGene.Values)
			targetGene.DCt = targetGene.Mean - endoControlMock.Mean
			targetGene.DdCt = 0.0
			targetGene.DdCtErr = math.Sqrt((2 * (endoControlMock.StdDev * endoControlMock.StdDev)) + (2 * (targetGene.StdDev * targetGene.StdDev)))
			targetGene.RQ = 1.0
			targetGene.RQErr = math.Sqrt(math.Ln2 * math.Ln2)

			e.Detectors[detectorName][mockName] = targetGene
		} else {
			log.Printf("[ab7300] mock for detector '%s' not found!\n", detectorName)
		}
	}
}

func (e *Experiment) parseRow(line *string) {
	rowValues := strings.Split(*line, ",")
	if len(rowValues) == 22 {
		name, detector, task, value := rowValues[3], rowValues[4], rowValues[5], rowValues[6]
		switch task {
		case "ENDO":
			e.addEndogenousControlTargetGeneValue(name, detector, value)
		case "Target":
			e.addDetectorTargetGeneValue(name, detector, value)
		default:
			log.Printf("[ab7300] ignoring unknown task type '%s'!\n", task)
		}
	} else {
		log.Printf("[ab7300] line '%s' is not valid ab7300 line!\n", line)
	}
}

func (e *Experiment) addEndogenousControlTargetGeneValue(name, detector, value string) {
	e.createEndogenousControlTargetGene(name)

	e.EndogenousControls[name].Detectors[detector] = append(e.EndogenousControls[name].Detectors[detector], value)

	e.updateEndogenousControlTargetGeneValues(name, value)
}

func (e *Experiment) createEndogenousControlTargetGene(name string) {
	if _, found := e.EndogenousControls[name]; found == false {
		e.EndogenousControls[name] = EndoTargetGene{Detectors: make(StringArrayMap)}
	}
}

func (e *Experiment) updateEndogenousControlTargetGeneValues(name, value string) {
	if v, err := strconv.ParseFloat(value, 64); err == nil {
		endoTargetGene := e.EndogenousControls[name]

		endoTargetGene.Values = append(endoTargetGene.Values, v)
		endoTargetGene.Mean, endoTargetGene.StdDev = meanAndStdDev(endoTargetGene.Values)

		e.EndogenousControls[name] = endoTargetGene
	}
}

func (e *Experiment) addDetectorTargetGeneValue(name, detector, value string) {
	e.createDetectorTargetGene(name, detector)

	targetGene := e.Detectors[detector][name]
	targetGene.RawValues = append(targetGene.RawValues, value)

	e.Detectors[detector][name] = targetGene
}

func (e *Experiment) createDetectorTargetGene(name, detector string) {
	if _, found := e.Detectors[detector]; found == false {
		e.Detectors[detector] = make(DetectorTargetGeneMap)
	}

	if _, found := e.Detectors[detector][name]; found == false {
		e.Detectors[detector][name] = DetectorTargetGene{}
	}
}

func (tg *DetectorTargetGene) mergeRawValues() {
	for _, value := range tg.RawValues {
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			tg.Values = append(tg.Values, v)
		}
	}

	if len(tg.Values) < 1 {
		tg.Values = append(tg.Values, 38.0)
	}
}

func meanAndStdDev(values []float64) (float64, float64) {
	sum := 0.0
	count := len(values)
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(count)

	sum = 0.0
	var x float64
	for _, v := range values {
		x = float64(v) - mean
		sum += x * x
	}
	stdDev := math.Sqrt(sum / float64(count-1))

	if math.IsNaN(stdDev) {
		stdDev = 0.0
	}

	return mean, stdDev
}

