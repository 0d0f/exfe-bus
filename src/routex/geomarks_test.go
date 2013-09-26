package routex

import (
	"github.com/googollee/go-assert"
	"github.com/googollee/go-pubsub"
	"model"
	"routex/model"
	"testing"
	"time"
)

type FakeConversion struct{}

func (c *FakeConversion) EarthToMars(lat, lng float64) (float64, float64) {
	return lat, lng
}

func (c *FakeConversion) MarsToEarth(lat, lng float64) (float64, float64) {
	return lat, lng
}

type FakeGeomarkRepo struct {
	geomarks map[string]rmodel.Geomark
}

func NewFakeGeomarkRepo() *FakeGeomarkRepo {
	geomarks := make(map[string]rmodel.Geomark)
	geomarks["location.1234"] = rmodel.Geomark{
		Id:   "location.1234",
		Type: "location",
		Tags: []string{DestinationTag},
	}
	geomarks["route.2234"] = rmodel.Geomark{
		Id:   "route.2234",
		Type: "route",
		Tags: []string{},
	}
	geomarks["location.3234"] = rmodel.Geomark{
		Id:   "location.3234",
		Type: "location",
		Tags: []string{},
	}
	return &FakeGeomarkRepo{
		geomarks: geomarks,
	}
}

func (r *FakeGeomarkRepo) Set(crossId int64, mark rmodel.Geomark) error {
	r.geomarks[mark.Id] = mark
	return nil
}

func (r *FakeGeomarkRepo) Get(crossId int64) ([]rmodel.Geomark, error) {
	var ret []rmodel.Geomark
	for _, v := range r.geomarks {
		ret = append(ret, v)
	}
	return ret, nil
}

func (r *FakeGeomarkRepo) Delete(crossId int64, type_, id, by string) error {
	delete(r.geomarks, id)
	return nil
}

func TestGeomarksCheck(t *testing.T) {
	routex := new(RouteMap)
	routex.geomarksRepo = NewFakeGeomarkRepo()
	routex.conversion = new(FakeConversion)
	routex.pubsub = pubsub.New(10)
	cross := model.Cross{
		ID: 789,
		By: model.Identity{
			ID:         567,
			ExternalID: "external",
			Provider:   "provider",
			UserID:     765,
		},
		Place: &model.Place{
			Title: "cross place",
			Lng:   "1.0",
			Lat:   "2.0",
		},
	}
	c := make(chan interface{}, 100)
	routex.pubsub.Subscribe(routex.publicName(int64(cross.ID)), c)
	marks, err := routex.getGeomarks_(cross, false)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(marks), 4)
	var xplace rmodel.Geomark
	var dest rmodel.Geomark
	for _, mark := range marks {
		if mark.HasTag(XPlaceTag) {
			xplace = mark
		}
		if mark.HasTag(DestinationTag) {
			dest = mark
		}
	}
	assert.Equal(t, xplace.HasTag(XPlaceTag), true)
	assert.Equal(t, xplace.HasTag(DestinationTag), false)
	assert.Equal(t, xplace.Id, "xplace.789")
	assert.Equal(t, dest.Id, "location.1234")
	dest.Tags = append(dest.Tags, XPlaceTag)
	dest.Action = "update"
	routex.checkGeomarks(cross, dest)
	select {
	case m := <-c:
		if mark, ok := m.(rmodel.Geomark); ok {
			assert.Equal(t, mark.Id, "location.1234")
			assert.Equal(t, mark.Action, "delete")
		} else {
			t.Error("should be geomark")
		}
	case <-time.After(time.Second):
		t.Fatal("wait too long")
	}
	select {
	case m := <-c:
		if mark, ok := m.(rmodel.Geomark); ok {
			assert.Equal(t, mark.Id, "xplace.789")
			assert.Equal(t, mark.Action, "update")
			assert.Equal(t, mark.HasTag(XPlaceTag), true)
			assert.Equal(t, mark.HasTag(DestinationTag), true)
			xplace = mark
		} else {
			t.Error("should be geomark")
		}
	case <-time.After(time.Second):
		t.Fatal("wait too long")
	}
	select {
	case m := <-c:
		t.Fatal("no more update", m)
	default:
	}
	dest.Id = "location.1234"
	dest.Tags = []string{DestinationTag}
	dest.Action = "update"
	routex.checkGeomarks(cross, dest)
	select {
	case m := <-c:
		if mark, ok := m.(rmodel.Geomark); ok {
			assert.Equal(t, mark.Id, "xplace.789")
			assert.Equal(t, mark.Action, "update")
			assert.Equal(t, mark.HasTag(DestinationTag), false)
		} else {
			t.Error("should be geomark")
		}
	case <-time.After(time.Second):
		t.Fatal("wait too long")
	}
	select {
	case m := <-c:
		if mark, ok := m.(rmodel.Geomark); ok {
			assert.Equal(t, mark.Id, "location.1234")
			assert.Equal(t, mark.Action, "update")
			assert.Equal(t, mark.HasTag(DestinationTag), true)
		} else {
			t.Error("should be geomark")
		}
	case <-time.After(time.Second):
		t.Fatal("wait too long")
	}
	select {
	case <-c:
		t.Fatal("no more update")
	default:
	}
	dest.Tags = []string{DestinationTag}
	dest.Action = "delete"
	routex.checkGeomarks(cross, dest)
	select {
	case m := <-c:
		if mark, ok := m.(rmodel.Geomark); ok {
			assert.Equal(t, mark.Id, "location.1234")
			assert.Equal(t, mark.Action, "delete")
		} else {
			t.Error("should be geomark")
		}
	case <-time.After(time.Second):
		t.Fatal("wait too long")
	}
	select {
	case m := <-c:
		if mark, ok := m.(rmodel.Geomark); ok {
			assert.Equal(t, mark.Id, "xplace.789")
			assert.Equal(t, mark.Action, "update")
			assert.Equal(t, mark.HasTag(DestinationTag), true)
		} else {
			t.Error("should be geomark")
		}
	case <-time.After(time.Second):
		t.Fatal("wait too long")
	}
	select {
	case m := <-c:
		t.Fatal("no more update", m)
	default:
	}
}
