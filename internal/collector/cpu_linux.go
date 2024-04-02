package collector

import (
	"fmt"
	"log"
	"regexp"
)

//type CPUStat struct {
//	UserMode   float64
//	SystemMode float64
//	Idle       float64
//}
//
//func CpuStat() (CPUStat, error) {
//	topBatch, err := RunTopBatch()
//	if err != nil {
//		log.Println(err)
//	}
//	cpuRegex := regexp.MustCompile(`%Cpu\(s\):\s+(.+)\sus,\s+(.+)\ssy,\s+(.+)\sni,\s+(.+)\sid`)
//	matches := cpuRegex.FindStringSubmatch(topBatch)
//	if len(matches) < 5 {
//		return CPUStat{}, fmt.Errorf("unable to parse cpu stat values")
//	}
//	var cpuStat CPUStat
//	cpuStat.UserMode, err = ParseStrToFloat(matches[1])
//	if err != nil {
//		return CPUStat{}, err
//	}
//	cpuStat.SystemMode, err = ParseStrToFloat(matches[2])
//	if err != nil {
//		return CPUStat{}, err
//	}
//	cpuStat.Idle, err = ParseStrToFloat(matches[4])
//	if err != nil {
//		return CPUStat{}, err
//	}
//	return cpuStat, nil
//}

func CpuStat() (Items, error) {
	topBatch, err := RunTopBatch()
	if err != nil {
		log.Println(err)
	}
	cpuRegex := regexp.MustCompile(`%Cpu\(s\):\s+(.+)\sus,\s+(.+)\ssy,\s+(.+)\sni,\s+(.+)\sid`)
	matches := cpuRegex.FindStringSubmatch(topBatch)
	if len(matches) < 5 {
		return Items{}, fmt.Errorf("unable to parse cpu stat values")
	}
	cpuStat := make(Items)
	cpuStat["user"], err = ParseStrToFloat(matches[1])
	if err != nil {
		return Items{}, err
	}
	cpuStat["system"], err = ParseStrToFloat(matches[2])
	if err != nil {
		return Items{}, err
	}
	cpuStat["idle"], err = ParseStrToFloat(matches[4])
	if err != nil {
		return Items{}, err
	}
	return cpuStat, nil
}
