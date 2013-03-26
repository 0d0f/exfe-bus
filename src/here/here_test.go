package here

import (
	"fmt"
	"testing"
	"time"
)

func TestHere(t *testing.T) {
	here := New(0.0001, 1, time.Second)
	go func() {
		for id := range here.UpdateChannel() {
			fmt.Printf("group update: %+v\n", here.GetGroup(id))
		}
	}()

	here.Add(User{
		Id:        "123",
		Name:      "123",
		Latitude:  13.4576787,
		Longitude: 14.4324325,
		Accuracy:  10,
	})
	here.Add(User{
		Id:        "1234",
		Name:      "1234",
		Latitude:  13.457677,
		Longitude: 14.432435,
		Accuracy:  10,
	})
	time.Sleep(time.Second / 2)
	here.Add(User{
		Id:        "1235",
		Name:      "1235",
		Latitude:  13.457677,
		Longitude: 14.432425,
		Accuracy:  10,
		Traits:    []string{"abc"},
	})
	here.Add(User{
		Id:        "1236",
		Name:      "1236",
		Latitude:  133.457677,
		Longitude: 142.432425,
		Accuracy:  10,
		Traits:    []string{"abc"},
	})
	time.Sleep(time.Second * 2)
	t.Errorf("show")
}
