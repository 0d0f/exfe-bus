package textcutter

import (
	"testing"
)

func equalArray(a, b []string) (bool, int) {
	if len(a) != len(b) {
		return false, -1
	}
	for i := range a {
		if a[i] != b[i] {
			return false, i
		}
	}
	return true, 0
}

func TestParseLex(t *testing.T) {
	str := `\(googol:\)\\ abcddef \(("some\\title" http://exfe.com/abcdefg)\)\(123 456\)`
	nodes, err := parseLexical(str)
	t.Logf("output: %v", nodes)
	if err != nil {
		t.Fatalf("parse lexical error: %s", err)
	}
	rawString := []string{`googol:`, `\`, ` abcddef`, ` `, `("some\title" http://exfe.com/abcdefg)`, `123 456`}

	if ok, i := equalArray(nodes, rawString); !ok {
		if i < 0 {
			t.Fatalf("length not same: %v", nodes)
		} else {
			t.Errorf("%d node got: '%s', expect: '%s'", i, nodes[i], rawString[i])
		}
	}
}

func normalStringLen(s string) int {
	return len(s)
}

func TestCutter(t *testing.T) {
	str := `\(googol:\)\\ 12345678901234567890\\ \(("1234567890123" http://exfe.com/abcdefg)\)`
	cutter, err := Parse(str, normalStringLen)
	if err != nil {
		t.Fatalf("parse string error: %s", err)
	}

	{
		got, expect := cutter.Limit(73), []string{`googol:\ 12345678901234567890\ ("1234567890123" http://exfe.com/abcdefg)`}
		if ok, i := equalArray(got, expect); !ok {
			if i < 0 {
				t.Fatalf("length not same: %v", got)
			} else {
				t.Errorf("%d node got: '%s', expect: '%s'", i, got[i], expect[i])
			}
		}
	}

	{
		got, expect := cutter.Limit(72), []string{`googol:\ 12345678901234567890\ ("1234567890123" http://exfe.com/abcdefg)`}
		if ok, i := equalArray(got, expect); !ok {
			if i < 0 {
				t.Fatalf("length not same: %v", got)
			} else {
				t.Errorf("%d node got: '%s', expect: '%s'", i, got[i], expect[i])
			}
		}
	}

	{
		got, expect := cutter.Limit(71), []string{`(1/2)googol:\ 12345678901234567890\`, `(2/2)("1234567890123" http://exfe.com/abcdefg)`}
		if ok, i := equalArray(got, expect); !ok {
			if i < 0 {
				t.Fatalf("length not same: %v", got)
			} else {
				t.Errorf("%d node got: '%s', expect: '%s'", i, got[i], expect[i])
			}
		}
	}

	{
		got, expect := cutter.Limit(70), []string{`(1/2)googol:\ 12345678901234567890\`, `(2/2)("1234567890123" http://exfe.com/abcdefg)`}
		if ok, i := equalArray(got, expect); !ok {
			if i < 0 {
				t.Fatalf("length not same: %v", got)
			} else {
				t.Errorf("%d node got: '%s', expect: '%s'", i, got[i], expect[i])
			}
		}
	}

	{
		got, expect := cutter.Limit(35), []string{`(1/2)googol:\ 12345678901234567890\`, `(2/2)("1234567890123" http://exfe.com/abcdefg)`}
		if ok, i := equalArray(got, expect); !ok {
			if i < 0 {
				t.Fatalf("length not same: %v", got)
			} else {
				t.Errorf("%d node got: '%s', expect: '%s'", i, got[i], expect[i])
			}
		}
	}

	{
		got, expect := cutter.Limit(34), []string{`(1/3)googol:\ 12345678901234567890`, `(2/3)\`, `(3/3)("1234567890123" http://exfe.com/abcdefg)`}
		if ok, i := equalArray(got, expect); !ok {
			if i < 0 {
				t.Fatalf("length not same: %v", got)
			} else {
				t.Errorf("%d node got: '%s', expect: '%s'", i, got[i], expect[i])
			}
		}
	}

	{
		got, expect := cutter.Limit(33), []string{`(1/3)googol:\`, `(2/3)12345678901234567890\`, `(3/3)("1234567890123" http://exfe.com/abcdefg)`}
		if ok, i := equalArray(got, expect); !ok {
			if i < 0 {
				t.Fatalf("length not same: %v", got)
			} else {
				t.Errorf("%d node got: '%s', expect: '%s'", i, got[i], expect[i])
			}
		}
	}
}
