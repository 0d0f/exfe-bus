package here

import (
	"fmt"
	"sort"
	"testing"
	"time"
)

func TestHere(t *testing.T) {
	var results = map[string][]string{
		"123":  []string{"123 {123 }", "123 {123 1234 }", "123 {123 1234 1235 }", "123 {123 1234 1235 1236 }", "123 {1235 1236 }"},
		"1234": []string{"123 {123 1234 }", "123 {123 1234 1235 }", "123 {123 1234 1235 1236 }", "123 {1235 1236 }"},
		"1235": []string{"123 {123 1234 1235 }", "123 {123 1234 1235 1236 }", ""},
		"1236": []string{"123 {123 1234 1235 1236 }", ""},
	}
	here := New(0.0001, 1, time.Second)
	go func() {
		for id := range here.UpdateChannel() {
			group := stringGroup(here.UserInGroup(id))
			fmt.Printf("user update: %+v, group:%s\n", id, group)
			if group != results[id][0] {
				t.Errorf("user %s should get %s", id, group)
			}
			results[id] = results[id][1:]
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

	for k, v := range results {
		if len(v) != 0 {
			fmt.Errorf("user %s should receive %v", k, v)
		}
	}
}

func stringGroup(group *Group) string {
	if group == nil {
		return ""
	}
	var users []string
	for k := range group.Users {
		users = append(users, k)
	}
	sort.Strings(users)
	ret := group.Name + " {"
	for _, u := range users {
		ret += u + " "
	}
	ret += "}"
	return ret
}
