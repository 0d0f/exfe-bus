package formatter

import (
	"bytes"
	"encoding/base64"
	"github.com/stretchrcom/testify/assert"
	"testing"
)

func TestTemplateColumn(t *testing.T) {
	templ, err := NewTemplate("test").Parse(`{{"1234567890" | column 5 "\n"}}`)
	if err != nil {
		t.Fatalf("unexpect error: %s", err)
	}
	buf := bytes.NewBuffer(nil)
	err = templ.Execute(buf, nil)
	if err != nil {
		t.Fatalf("unexpect error: %s", err)
	}
	assert.Equal(t, buf.String(), "12345\n67890", "should equal")
}

func TestTemplateEqual(t *testing.T) {
	templ, err := NewTemplate("test").Parse(`{{equal 1 1}}{{"a" | equal "1"}}`)
	if err != nil {
		t.Fatalf("unexpect error: %s", err)
	}
	buf := bytes.NewBuffer(nil)
	err = templ.Execute(buf, nil)
	if err != nil {
		t.Fatalf("unexpect error: %s", err)
	}
	assert.Equal(t, buf.String(), "truefalse", "should equal")
}

func TestTemplateReplace(t *testing.T) {
	templ, err := NewTemplate("test").Parse(`{{"12345" | replace "12" "ab"}}`)
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

func TestTemplateLimit(t *testing.T) {
	templ, err := NewTemplate("test").Parse(`{{"12345678" | limit 3}}`)
	if err != nil {
		t.Fatalf("unexpect error: %s", err)
	}
	buf := bytes.NewBuffer(nil)
	err = templ.Execute(buf, nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, buf.String(), "12…")
}

func TestLocalTemplate(t *testing.T) {
	l, err := NewLocalTemplate("./template_test", "en_US")
	assert.Equal(t, err, nil)
	assert.Equal(t, l.defaultLang, "en_US")

	_, ok := l.templates["en_US"]
	assert.Equal(t, ok, true)

	_, ok = l.templates[".should_ignore"]
	assert.Equal(t, ok, false)

	buf := bytes.NewBuffer(nil)
	l.Execute(buf, "en_US", "test.template", nil)
	assert.Equal(t, buf.String(), "1234\n")

	buf.Reset()
	l.Execute(buf, "zh_CN", "test.template", nil)
	assert.Equal(t, buf.String(), "abcd\n")

	buf.Reset()
	l.Execute(buf, "en_CN", "test.template", nil)
	assert.Equal(t, buf.String(), "1234\n")
}
