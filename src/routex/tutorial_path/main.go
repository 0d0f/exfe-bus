package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"net/textproto"
	"os"
	"time"
)

func Distance(latA, lngA, latB, lngB float64) float64 {
	x := math.Cos(latA*math.Pi/180) * math.Cos(latB*math.Pi/180) * math.Cos((lngA-lngB)*math.Pi/180)
	y := math.Sin(latA*math.Pi/180) * math.Sin(latB*math.Pi/180)
	s := x + y
	if s > 1 {
		s = 1
	}
	if s < -1 {
		s = -1
	}
	alpha := math.Acos(s)
	distance := alpha * 6371000
	return distance
}

type Location struct {
	Offset int
	Lat    float64
	Lng    float64
}

const interval = 10

var buf = bytes.NewBuffer(nil)

func main() {
	if len(os.Args) != 2 {
		fmt.Println(`tutorial [file]`)
		return
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("invalid file:", os.Args[1])
		return
	}
	r := textproto.NewReader(bufio.NewReader(f))

	buf.WriteString("[")
	var start, end int
	var offsetLat, offsetLng float64
	var data []Location
	for l, err := r.ReadLine(); err == nil; l, err = r.ReadLine() {
		l = textproto.TrimString(l)
		if len(l) == 0 {
			continue
		}
		if l[0] == '/' {
			if data != nil {
				parseGps(start, end, data)
			}
			offsetLat, offsetLng = 0, 0
			var startStr, endStr, comment string
			if _, err := fmt.Sscanf(l, "// %s %s %s GPS +offset %f,%f", &startStr, &endStr, &comment, &offsetLat, &offsetLng); err != nil {
				if _, err := fmt.Sscanf(l, "// %s %s", &startStr, &endStr); err != nil {
					panic(err)
				}
			}
			startTime, err := time.Parse("15:04", startStr)
			if err != nil {
				panic(err)
			}
			start = startTime.Hour()*60*60 + startTime.Minute()*60
			endTime, err := time.Parse("15:04", endStr)
			if err != nil {
				panic(err)
			}
			end = endTime.Hour()*60*60 + endTime.Minute()*60
			if end < start {
				if end == 0 {
					end = 24 * 60 * 60
				} else {
					panic(fmt.Sprintf("%s not small than %s", startStr, endStr))
				}
			}
			data = nil
			continue
		}
		var lat, lng float64
		if _, err := fmt.Sscanf(l, "%f,%f", &lat, &lng); err != nil {
			panic(err)
		}
		lat, lng = lat+offsetLat, lng+offsetLng
		data = append(data, Location{0, lat, lng})
	}
	if data != nil {
		parseGps(start, end, data)
	}
	buf.Truncate(buf.Len() - 2)
	buf.WriteString("]")
	fmt.Println(buf.String())
}

func parseGps(start, end int, data []Location) {
	distances := make([]int, len(data)-1)
	sum := 0
	total := (end - start) / 10
	for i, n := 0, len(data)-1; i < n; i++ {
		distances[i] = int(Distance(data[i].Lat, data[i].Lng, data[i+1].Lat, data[i+1].Lng) * 100)
		sum += distances[i]
	}
	var points []Location
	offset := 0
	for i, d := range distances {
		length := total * d / sum
		total -= length
		sum -= d
		intervalLat := (data[i+1].Lat - data[i].Lat) / float64(length)
		intervalLng := (data[i+1].Lng - data[i].Lng) / float64(length)
		for j := 0; j < length; j++ {
			points = append(points, Location{start + offset, data[i].Lat + intervalLat*float64(j), data[i].Lng + intervalLng*float64(j)})
			offset += 10
		}
	}
	if len(points) == 0 {
		for i := 0; i < total; i++ {
			points = append(points, Location{start + offset, data[0].Lat, data[0].Lng})
			offset += 10
		}
	}
	for _, p := range points {
		buf.WriteString(fmt.Sprintf("{\"offset\":%d,\"acc\":10,\"lat\":%.7f,\"lng\":%.7f},\n", p.Offset, p.Lat, p.Lng))
	}
}
