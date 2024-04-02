package collector

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Items map[string]float64

type CollectorFunc func() (Items, error)

type Collector struct {
	Name     string
	FuncName CollectorFunc
	Data     []Items
	Limit    int
	Enabled  bool
}

type Collectors []Collector

func (c *Collector) CollectData(doneCh <-chan struct{}) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-doneCh:
			fmt.Printf("%s done.", c.Name)
			return
		case <-ticker.C:
			//log.Println(c.Data)
			data, err := c.FuncName()
			if err != nil {
				log.Printf("error collecting data %s: %s", c.Name, err)
			} else {
				c.Data = append(c.Data, data)
				if len(c.Data) > c.Limit {
					c.Data = c.Data[len(c.Data)-c.Limit:]
				}
			}
		}
	}
}

func (c *Collector) Average(m int) (Items, error) {
	if m == 0 {
		return Items{}, fmt.Errorf("error counting average from empty slice")
	}
	firstItem := c.Data[0]
	keys := make([]string, 0, len(firstItem))
	for k := range firstItem {
		keys = append(keys, k)
	}
	averages := make(Items)
	items := c.Data[len(c.Data)-m:]
	for _, item := range items {
		for _, k := range keys {
			averages[k] += item[k]
		}
	}
	for k, _ := range averages {
		averages[k] = averages[k] / float64(m)
	}
	return averages, nil
}

func RunTopBatch() (string, error) {
	cmd := exec.Command("top", "-b", "-n", "1")
	var buf bytes.Buffer
	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return buf.String(), nil
}

func ParseStrToFloat(value string) (float64, error) {
	result, err := strconv.ParseFloat(strings.ReplaceAll(strings.TrimSpace(value), ",", "."), 64)
	var floatVal float64
	if err != nil {
		return floatVal, fmt.Errorf("error parsing %s: %w", value, err)
	}
	return result, nil
}
