package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Step struct {
	Distance struct {
		Text  string
		Value int64
	}
	Duration struct {
		Text  string
		Value int64
	}
	EndLocation   Location `json:"end_location"`
	StartLocation Location `json:"start_location"`
}

type Map struct {
	Routes []struct {
		Bounds struct {
			Northeast Location
			Southwest Location
		}
		Copyrights string
		Legs       []struct {
			Step
			EndAddress   string
			StartAddress string
			Steps        []Step
		}
	}
	Status string
}

type Location struct {
	Lat float64
	Lng float64
}

const allPoints = 8640

// stdin google map api route path:
//   http://maps.googleapis.com/maps/api/directions/json?origin=Tiananmen+Square&destination=Tiananmen+Square&waypoints=shanghai&sensor=false
// stdout location array
func main() {
	var ret Map
	decoder := json.NewDecoder(os.Stdin)
	err := decoder.Decode(&ret)
	if err != nil {
		panic(err)
	}
	var total = int64(allPoints)
	var steps []Step
	for _, route := range ret.Routes {
		for _, leg := range route.Legs {
			steps = append(steps, leg.Steps...)
		}
	}
	sum := int64(0)
	points := make([]int64, len(steps))
	for _, step := range steps {
		sum += step.Distance.Value
	}
	for i, step := range steps {
		points[i] = total * step.Distance.Value / sum
		if points[i] == 0 {
			points[i] = 1
		}
		sum -= step.Distance.Value
		total -= points[i]
	}
	sum = 0
	for _, p := range points {
		sum += p
	}
	var locations []Location
	total -= 1
	for i, step := range steps {
		ls := make([]Location, points[i])
		latStep := (step.EndLocation.Lat - step.StartLocation.Lat) / float64(points[i])
		lngStep := (step.EndLocation.Lng - step.StartLocation.Lng) / float64(points[i])
		for j := int64(0); j < points[i]; j++ {
			ls[j] = Location{
				Lat: float64(j)*latStep + step.StartLocation.Lat,
				Lng: float64(j)*lngStep + step.StartLocation.Lng,
			}
		}
		locations = append(locations, ls...)
	}
	offset := 0
	interval := 24 * 60 * 60 / allPoints
	for _, l := range locations {
		fmt.Printf("{\"offset\":%d, \"lat\":%.7f, \"lng\":%.7f, \"acc\":10},\n", offset, l.Lat, l.Lng)
		offset += interval
	}
}
