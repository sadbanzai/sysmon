package collector

import (
	"fmt"
	"log"
	"regexp"
)

//type LoadAvg struct {
//	Avg1  float64
//	Avg5  float64
//	Avg15 float64
//}
//
//func LoadAverage() (LoadAvg, error) {
//	topBatch, err := RunTopBatch()
//	if err != nil {
//		log.Println(err)
//	}
//	loadAvgRegex := regexp.MustCompile(`load average: (.+), (.+), (.+)`)
//	matches := loadAvgRegex.FindStringSubmatch(topBatch)
//	if len(matches) < 4 {
//		return LoadAvg{}, fmt.Errorf("unable to parse load average values")
//	}
//	var loadAvg LoadAvg
//	loadAvg.Avg1, err = ParseStrToFloat(matches[1])
//	if err != nil {
//		return LoadAvg{}, err
//	}
//	loadAvg.Avg5, err = ParseStrToFloat(matches[2])
//	if err != nil {
//		return LoadAvg{}, err
//	}
//	loadAvg.Avg15, err = ParseStrToFloat(matches[3])
//	if err != nil {
//		return LoadAvg{}, err
//	}
//	return loadAvg, nil
//}

func LoadAverage() (Items, error) {
	topBatch, err := RunTopBatch()
	if err != nil {
		log.Println(err)
	}
	loadAvgRegex := regexp.MustCompile(`load average: (.+), (.+), (.+)`)
	matches := loadAvgRegex.FindStringSubmatch(topBatch)
	if len(matches) < 4 {
		return Items{}, fmt.Errorf("unable to parse load average values")
	}
	loadAvg := make(Items)
	loadAvg["avg1"], err = ParseStrToFloat(matches[1])
	if err != nil {
		return Items{}, err
	}
	loadAvg["avg5"], err = ParseStrToFloat(matches[2])
	if err != nil {
		return Items{}, err
	}
	loadAvg["avg15"], err = ParseStrToFloat(matches[3])
	if err != nil {
		return Items{}, err
	}
	return loadAvg, nil
}
