package formater

import (
	"bytes"
	"encoding/base64"
	"github.com/stretchrcom/testify/assert"
	"testing"
)

func TestTemplateReplace(t *testing.T) {
	templ, err := NewTemplate("test").Parse(`{{replace "12345" "12" "ab"}}`)
	if err != nil {
		t.Fatalf("unexpect error: %s", err)
	}
	buf := bytes.NewBuffer(nil)
	err = templ.Execute(buf, nil)
	if err != nil {
		t.Fatalf("unexpect error: %s", err)
	}
	assert.Equal(t, buf.String(), "ab345", "should equal")
}

func TestTemplateBase64(t *testing.T) {
	{
		templ, err := NewTemplate("test").Parse(`{{base64 "12345"}}`)
		if err != nil {
			t.Fatalf("unexpect error: %s", err)
		}
		buf := bytes.NewBuffer(nil)
		err = templ.Execute(buf, nil)
		if err != nil {
			t.Fatalf("unexpect error: %s", err)
		}
		assert.Equal(t, buf.String(), base64.StdEncoding.EncodeToString([]byte("12345")), "should equal")
	}

	{
		templ, err := NewTemplate("test").Parse(`{{base64url "12345"}}`)
		if err != nil {
			t.Fatalf("unexpect error: %s", err)
		}
		buf := bytes.NewBuffer(nil)
		err = templ.Execute(buf, nil)
		if err != nil {
			t.Fatalf("unexpect error: %s", err)
		}
		assert.Equal(t, buf.String(), base64.URLEncoding.EncodeToString([]byte("12345")), "should equal")
	}
}

func TestTemplateSub(t *testing.T) {
	templ, err := NewTemplate("a").Parse(`abcd{{sub "b" .}}`)
	if err != nil {
		t.Fatalf("unexpect error: %s", err)
	}
	_, err = templ.New("b").Parse(`1234`)
	if err != nil {
		t.Fatalf("unexpect error: %s", err)
	}
	buf := bytes.NewBuffer(nil)
	err = templ.Execute(buf, nil)
	if err != nil {
		t.Fatalf("unexpect error: %s", err)
	}
	assert.Equal(t, buf.String(), "abcd1234", "should equal")
}

func TestTemplateFor(t *testing.T) {
	templ, err := NewTemplate("test").Parse(`{{range for .}}{{.Index}} - {{.V}}{{if not .Last}}, {{end}}{{end}}`)
	if err != nil {
		t.Fatalf("unexpect error: %s", err)
	}
	buf := bytes.NewBuffer(nil)
	err = templ.Execute(buf, []string{"a", "b", "c", "d"})
	assert.Equal(t, err, nil)
	assert.Equal(t, buf.String(), "1 - a, 2 - b, 3 - c, 4 - d")
}
