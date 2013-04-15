package here

import (
	"fmt"
	"sort"
	"testing"
	"time"
)

func TestHere(t *testing.T) {
	var results = []string{
		"123 {123 }",
		"123 {123 1234 }",
		"123 {123 1234 1235 }",
		"123 {123 1234 1235 1236 }",
	}
	here := New(0.0001, 1, time.Second)
	go here.Serve()

	go func() {
		for group := range here.UpdateChannel() {
			fmt.Println(stringGroup(&group))
			if len(results) == 0 {
				if group.Name != "" {
					t.Errorf("should recevie empty, got: %s", stringGroup(&group))
				}
				continue
			}
			if results[0] != stringGroup(&group) {
				t.Errorf("should received: %s, got: %s", results[0], stringGroup(&group))
			}
			results = results[1:]
		}
	}()

	fmt.Println("add 123")
	here.Add(&Data{
		Token:     "123",
		Latitude:  "13.4576787",
		Longitude: "14.4324325",
		Accuracy:  "10",
	})
	fmt.Println("add 1234")
	here.Add(&Data{
		Token:     "1234",
		Latitude:  "13.457677",
		Longitude: "14.432435",
		Accuracy:  "10",
	})
	time.Sleep(time.Second / 2)
	fmt.Println("add 1235")
	here.Add(&Data{
		Token:     "1235",
		Latitude:  "13.457677",
		Longitude: "14.432425",
		Accuracy:  "10",
		Traits:    []string{"abc"},
	})
	fmt.Println("add 1236")
	here.Add(&Data{
		Token:     "1236",
		Latitude:  "133.457677",
		Longitude: "142.432425",
		Accuracy:  "10",
		Traits:    []string{"abc"},
	})
	time.Sleep(time.Second * 2)

	if len(results) != 0 {
		fmt.Errorf("should receive %v", results)
	}
}

func stringGroup(group *Group) string {
	if group == nil {
		return ""
	}
	var data []string
	for k := range group.Data {
		data = append(data, k)
	}
	sort.Strings(data)
	ret := group.Name + " {"
	for _, u := range data {
		ret += u + " "
	}
	ret += "}"
	return ret
}
