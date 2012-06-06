package main

import (
	"testing"
)

func equal(t *testing.T, expect, got string) {
	if expect != got {
		t.Errorf("expect: %s, but got: %s", expect, got)
	}
}

func TestModel(t *testing.T) {
	Init("tester")

	user := User{
		Id_str: "13145123512",
		Screen_name: "tester",
	}
	direct := DirectMessage{
		Created_at: "Wed Jun 07 06:52:56 +0000 2012",
		Text: "fdafdsfdsa #b3",
		Sender: user,
	}
	tweet := &Tweet{
		Created_at: "Wed Jun 06 06:52:56 +0000 2012",
		Text: "@tester #a4 jkljkljflda",
		User: &user,
		Direct_message: &direct,
	}

	equal(t, "fdafdsfdsa #b3", tweet.text())

	tweet.Direct_message = nil
	equal(t, "#a4 jkljkljflda", tweet.text())

	hash, post := tweet.parse()
	equal(t, "a4", hash)
	equal(t, "jkljkljflda", post)

	tweet.Direct_message = &direct
	hash, post = tweet.parse()
	equal(t, "b3", hash)
	equal(t, "fdafdsfdsa", post)
}
