package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"time"

	"github.com/wcharczuk/go-chart"
)

type Values struct {
	Name     string  `json:"name"`
	TempMin  float32 `json:"tempMin"`
	TempMax  float32 `json:"tempMax"`
	Interval int     `json:"interval"`
	Values   []Value `json:"values"`
}

type Value struct {
	MessageId    int     `json:"messageId"`
	Temperature  float32 `json:"temperature"`
	EnqueuedTime string  `json:"enqueuedTime"`
}

type reading struct {
	hour       int
	normal     float32
	outofrange float32
}

func main() {

	jsonFile, err := os.Open("data.json")
	if err != nil {
		log.Fatalln("File not found")
	}

	byteValue, _ := io.ReadAll(jsonFile)
	var v Values
	json.Unmarshal(byteValue, &v)

	tempMap := make(map[int][]float32)
	for _, value := range v.Values {
		t, err := time.Parse("2006-01-02 15:04:05", value.EnqueuedTime)
		if err != nil {
			log.Fatalln(err)
		}
		h := t.Hour()
		tempMap[h] = append(tempMap[h], value.Temperature)
	}
	var normal, outofrange float32
	var readings []reading
	for x, y := range tempMap {
		normal, outofrange = 0.0, 0.0
		for _, b := range y {
			if b >= v.TempMin && b <= v.TempMax {
				normal++
			} else {
				outofrange++
			}
		}
		read := reading{x, normal, outofrange}
		readings = append(readings, read)
	}
	sort.Slice(readings, func(i, j int) bool {
		return readings[i].hour < readings[j].hour
	})
	printTable(readings)
	printChart(readings)
}

func printChart(r []reading) {
	var bars []chart.StackedBar
	for _, val := range r {
		msg := fmt.Sprintf("Hours %d", val.hour)
		bar := chart.StackedBar{
			Name: msg,
			Values: []chart.Value{
				{Value: float64(val.normal), Label: "Green"},
				{Value: float64(val.outofrange), Label: "Red"},
			},
		}
		bars = append(bars, bar)
	}
	sbc := chart.StackedBarChart{Title: "IOT Sensor Temperature Details",
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Height: 512,
		Width:  2000,
		Bars:   bars,
	}
	f, _ := os.Create("sensorOutput.png")
	sbc.Render(chart.PNG, f)
}
func printTable(r []reading) {
	fmt.Printf("Hours\tTotal\tNormal\tOut of Range\tPercent\n")
	for _, val := range r {
		total := val.normal + val.outofrange
		percent := val.outofrange / total * 100
		fmt.Printf("%v\t%v\t%v\t%5v\t\t%5.1f\n", val.hour, total, val.normal, val.outofrange, percent)
	}
}
