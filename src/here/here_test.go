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
			fmt.Printf("user update: %+v, group:%+v\n", id, here.UserInGroup(id))
		}
	}()

	fmt.Println("add 123")
	here.Add(User{
		Id:        "123",
		Name:      "123",
		Latitude:  13.4576787,
		Longitude: 14.4324325,
		Accuracy:  10,
	})
	fmt.Println("add 1234")
	here.Add(User{
		Id:        "1234",
		Name:      "1234",
		Latitude:  13.457677,
		Longitude: 14.432435,
		Accuracy:  10,
	})
	time.Sleep(time.Second / 2)
	fmt.Println("add 1235")
	here.Add(User{
		Id:        "1235",
		Name:      "1235",
		Latitude:  13.457677,
		Longitude: 14.432425,
		Accuracy:  10,
		Traits:    []string{"abc"},
	})
	fmt.Println("add 1236")
	here.Add(User{
		Id:        "1236",
		Name:      "1236",
		Latitude:  133.457677,
		Longitude: 142.432425,
		Accuracy:  10,
		Traits:    []string{"abc"},
	})
	time.Sleep(time.Second * 2)
	// t.Errorf("show")
}
