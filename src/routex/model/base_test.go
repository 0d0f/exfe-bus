package model

import (
	"encoding/json"
	"fmt"
	"github.com/googollee/go-assert"
	"testing"
)

type GeoConversionTest struct {
	data map[string][2]int
}

func NewGeoConversionTest() *GeoConversionTest {
	return &GeoConversionTest{
		data: map[string][2]int{
			"1.00,2.00": [2]int{3, -1},
			"2.00,1.00": [2]int{-3, 2},
		},
	}
}

func (g GeoConversionTest) EarthToMars(lat, lng float64) (float64, float64) {
	key := fmt.Sprintf("%.2f,%.2f", lat, lng)
	offset, ok := g.data[key]
	if !ok {
		return lat, lng
	}
	return lat + float64(offset[0]), lng + float64(offset[1])
}

func (g GeoConversionTest) MarsToEarth(lat, lng float64) (float64, float64) {
	key := fmt.Sprintf("%.2f,%.2f", lat, lng)
	offset, ok := g.data[key]
	if !ok {
		return lat, lng
	}
	return lat - float64(offset[0]), lng - float64(offset[1])
}

func TestSimpleLocationJson(t *testing.T) {
	type Test struct {
		t      int64
		gps    [3]float64
		target string
	}
	var tests = []Test{
		{1, [3]float64{1, 2, 3}, `{"t":1,"gps":[1.0000000,2.0000000,3]}`},
		{2, [3]float64{1.00000001, 2.12345678, 3.234}, `{"t":2,"gps":[1.0000000,2.1234568,3]}`},
	}
	for i, test := range tests {
		l := SimpleLocation{test.t, test.gps}
		j, err := json.Marshal(l)
		assert.MustEqual(t, err, nil, "test %d", i)
		assert.Equal(t, string(j), test.target, "test %d", i)
		j, err = json.Marshal(&l)
		assert.MustEqual(t, err, nil, "test %d", i)
		assert.Equal(t, string(j), test.target, "test %d", i)
	}
}

func TestSimpleLocationConversion(t *testing.T) {
	conv := NewGeoConversionTest()
	type Test struct {
		gps   [3]float64
		mars  [3]float64
		earth [3]float64
	}
	var tests = []Test{
		{[3]float64{1, 2, 3}, [3]float64{4, 1, 3}, [3]float64{-2, 3, 3}},
		{[3]float64{1.001, 2.002, 3}, [3]float64{4.001, 1.002, 3}, [3]float64{-1.999, 3.002, 3}},
		{[3]float64{2, 1, 3}, [3]float64{-1, 3, 3}, [3]float64{5, -1, 3}},
		{[3]float64{2.002, 1.001, 3}, [3]float64{-0.998, 3.001, 3}, [3]float64{5.002, -0.999, 3}},
		{[3]float64{1, 1, 3}, [3]float64{1, 1, 3}, [3]float64{1, 1, 3}},
	}
	toString := func(d [3]float64) string {
		ret := ""
		for _, f := range d {
			ret += fmt.Sprintf("%.3f,", f)
		}
		return ret
	}
	for i, test := range tests {
		l := SimpleLocation{0, test.gps}
		mars := l
		mars.ToMars(conv)
		earth := l
		earth.ToEarth(conv)
		assert.Equal(t, toString(l.GPS), toString(test.gps), "test %d", i)
		assert.Equal(t, toString(mars.GPS), toString(test.mars), "test %d", i)
		assert.Equal(t, toString(earth.GPS), toString(test.earth), "test %d", i)
	}
}

func TestGeomarkHasTag(t *testing.T) {
	type Test struct {
		tags []string
		tag  string
		ok   bool
	}
	var tests = []Test{
		{[]string{"a", "b", "c"}, "d", false},
		{[]string{"a", "b", "c"}, "a", true},
		{[]string{}, "d", false},
	}
	for i, test := range tests {
		mark := Geomark{}
		mark.Tags = test.tags
		assert.Equal(t, mark.HasTag(test.tag), test.ok, "test %d", i)
	}
}

func TestGeomarkRemoveTag(t *testing.T) {
	type Test struct {
		tags []string
		tag  string
		ok   bool
	}
	var tests = []Test{
		{[]string{"a", "b", "c"}, "d", false},
		{[]string{"a", "b", "c"}, "a", true},
		{[]string{}, "d", false},
	}
	for i, test := range tests {
		mark := Geomark{}
		mark.Tags = test.tags
		assert.Equal(t, mark.RemoveTag(test.tag), test.ok, "test %d", i)
		assert.Equal(t, mark.HasTag(test.tag), false, "test %d", i)
	}
}

func TestGeomarkConversion(t *testing.T) {
	conv := NewGeoConversionTest()
	type Test struct {
		lat       float64
		lng       float64
		positions [][3]float64
		mars      string
		earth     string
	}
	var tests = []Test{
		{1, 2, nil, `{"type":"location","lat":4,"lng":1}`, `{"type":"location","lat":-2,"lng":3}`},
		{2, 1, nil, `{"type":"location","lat":-1,"lng":3}`, `{"type":"location","lat":5,"lng":-1}`},
		{0, 0, [][3]float64{{1.001, 2.002, 3}, {2.002, 1.001, 3}},
			"{\"type\":\"route\",\"positions\":[{\"t\":0,\"gps\":[4.0010000,1.0020000,3]},{\"t\":0,\"gps\":[-0.9980000,3.0010000,3]}]}",
			"{\"type\":\"route\",\"positions\":[{\"t\":0,\"gps\":[-1.9990000,3.0020000,3]},{\"t\":0,\"gps\":[5.0020000,-0.9990000,3]}]}",
		},
	}
	for i, test := range tests {
		mark := Geomark{
			Type:      "location",
			Latitude:  test.lat,
			Longitude: test.lng,
		}
		for _, p := range test.positions {
			mark.Type = "route"
			l := SimpleLocation{
				GPS: p,
			}
			mark.Positions = append(mark.Positions, l)
		}
		mars := mark
		earth := mark
		mars.ToMars(conv)
		earth.ToEarth(conv)
		m, err := json.Marshal(mars)
		assert.MustEqual(t, err, nil, "test %d", i)
		e, err := json.Marshal(earth)
		assert.MustEqual(t, err, nil, "test %d", i)
		assert.Equal(t, string(m), test.mars, "test %d", i)
		assert.Equal(t, string(e), test.earth, "test %d", i)
	}
}
