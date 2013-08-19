package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Location struct {
	Lat float64
	Lng float64
}

const allPoints = 8640

// echo '[{"lat":x.xxx,"lng":y.yy},...]' | tutorial 15:10 19:00
func main() {
	if len(os.Args) != 3 {
		fmt.Println(`echo '[{"lat":x.xxx,"lng":y.yy},...]' | tutorial 15:10 19:00`)
		return
	}

	startStr := os.Args[1]
	startTime, err := time.Parse("15:04", startStr)
	if err != nil {
		fmt.Println("invalid start time %s: %s", startStr, err)
		return
	}
	start := startTime.Hour()*60*60 + startTime.Minute()*60

	endStr := os.Args[2]
	endTime, err := time.Parse("15:04", endStr)
	if err != nil {
		fmt.Println("invalid start time %s: %s", endStr, err)
		return
	}

	end := endTime.Hour()*60*60 + endTime.Minute()*60

	var locations []Location
	decoder := json.NewDecoder(os.Stdin)
	err = decoder.Decode(&locations)
	if err != nil {
		panic(err)
	}

	totalInput := len(locations)
	totalOutput := (end - start) / 10
	var ret []Location
	for i, o := 0, 0; i < totalInput; i++ {
		if float64(i)/float64(totalInput) < float64(o)/float64(totalOutput) {
			continue
		}
		ret = append(ret, locations[i])
		o++
	}

	for _, l := range ret {
		fmt.Printf("{\"offset\":%d, \"lat\":%.7f, \"lng\":%.7f, \"acc\":10},\n", offset, l.Lat, l.Lng)
		offset += interval
	}
}
