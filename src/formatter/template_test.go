package formatter

import (
	"bytes"
	"encoding/base64"
	"github.com/stretchrcom/testify/assert"
	"testing"
)

func TestTemplateSubstr(t *testing.T) {
	templ, err := NewTemplate("test").Parse(`{{"0123456789" | substr 2 4}}`)
	if err != nil {
		t.Fatalf("unexpect error: %s", err)
	}
	buf := bytes.NewBuffer(nil)
	err = templ.Execute(buf, nil)
	if err != nil {
		t.Fatalf("unexpect error: %s", err)
	}
	assert.Equal(t, buf.String(), "2345", "should equal")
}

func TestTemplateAppend(t *testing.T) {
	templ, err := NewTemplate("test").Parse(`{{append "1234567890" "abcd"}}`)
	if err != nil {
		t.Fatalf("unexpect error: %s", err)
	}
	buf := bytes.NewBuffer(nil)
	err = templ.Execute(buf, nil)
	if err != nil {
		t.Fatalf("unexpect error: %s", err)
	}
	assert.Equal(t, buf.String(), "1234567890abcd", "should equal")
}

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

func TestTemplateFor(t *testing.T) {
	templ, err := NewTemplate("test").Parse(`{{range for .}}{{if not .First}}{{if not .Last}}, {{else}} and {{end}}{{end}}{{.Index}} - {{.V}}{{end}}`)
	if err != nil {
		t.Fatalf("unexpect error: %s", err)
	}
	buf := bytes.NewBuffer(nil)
	err = templ.Execute(buf, []string{"a", "b", "c", "d"})
	assert.Equal(t, err, nil)
	assert.Equal(t, buf.String(), "1 - a, 2 - b, 3 - c and 4 - d")
}

func TestTemplateLimit(t *testing.T) {
	templ, err := NewTemplate("test").Parse(`{{"测试文字测试文字" | limit 4}}`)
	if err != nil {
		t.Fatalf("unexpect error: %s", err)
	}
	buf := bytes.NewBuffer(nil)
	err = templ.Execute(buf, nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, buf.String(), "测...")
}

func TestTemplatePlural(t *testing.T) {
	templ, err := NewTemplate("test").Parse(`{{plural "is" "are" 1}}{{plural "is" "are" 2}}`)
	if err != nil {
		t.Fatalf("unexpect error: %s", err)
	}
	buf := bytes.NewBuffer(nil)
	err = templ.Execute(buf, nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, buf.String(), "isare")
}

func TestTemplateSub(t *testing.T) {
	templ, err := NewTemplate("test").Parse(`{{sub . "a"}} {{sub . "b"}} {{sub . "a" "b"}} {{sub . "c" "b"}}`)
	if err != nil {
		t.Fatalf("unexpect error: %s", err)
	}
	_, err = templ.New("a").Parse(`aaa`)
	if err != nil {
		t.Fatalf("unexpect error: %s", err)
	}
	_, err = templ.New("b").Parse(`bbb`)
	if err != nil {
		t.Fatalf("unexpect error: %s", err)
	}
	buf := bytes.NewBuffer(nil)
	err = templ.Execute(buf, nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, buf.String(), "aaa bbb aaa bbb")
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
	assert.Equal(t, buf.String(), "1234 abccc\n")

	buf.Reset()
	l.Execute(buf, "zh_CN", "test.template", nil)
	assert.Equal(t, buf.String(), "abcd\n")

	buf.Reset()
	l.Execute(buf, "en_CN", "test.template", nil)
	assert.Equal(t, buf.String(), "1234 abccc\n")
}

func TestLocalTemplateExist(t *testing.T) {
	l, err := NewLocalTemplate("./template_test", "en_US")
	assert.Equal(t, err, nil)
	assert.Equal(t, l.defaultLang, "en_US")

	assert.Equal(t, l.IsExist("en_US", "test.template"), true)
	assert.Equal(t, l.IsExist("en_CN", "test.template"), true)
	assert.Equal(t, l.IsExist("en_US", "nonexist.template"), false)
	assert.Equal(t, l.IsExist("en_CN", "nonexist.template"), false)
}
