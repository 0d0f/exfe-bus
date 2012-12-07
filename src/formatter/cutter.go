package formatter

import (
	"bytes"
	"fmt"
	"strings"
)

func isNoBreakRune(r rune) bool {
	const noBreakRune = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~:/?#[]@!$&'(%)*+,;=\""
	for _, c := range noBreakRune {
		if r == c {
			return true
		}
	}
	return false
}

func parseLexical(text string) ([]string, error) {
	ret := make([]string, 0)
	buf := bytes.NewBuffer(nil)
	breakSave := func() {
		if buf.Len() > 0 {
			ret = append(ret, buf.String())
			buf.Reset()
		}
	}
	control := false
	nobreak := false
	for _, r := range text {
		if r == '\\' && !control {
			control = true
			continue
		}
		if control {
			control = false
			switch r {
			case '(':
				breakSave()
				nobreak = true
				continue
			case ')':
				breakSave()
				nobreak = false
				continue
			}
		}
		if nobreak || isNoBreakRune(r) {
			buf.WriteRune(r)
			continue
		}
		breakSave()
		buf.WriteRune(r)
	}
	breakSave()
	return ret, nil
}

type LengthFunc func(string) int

type Cutter struct {
	origin string
	nodes  []string
	len    LengthFunc
}

func CutterParse(text string, f LengthFunc) (*Cutter, error) {
	nodes, err := parseLexical(text)
	if err != nil {
		return nil, err
	}
	return &Cutter{
		origin: strings.Join(nodes, ""),
		nodes:  nodes,
		len:    f,
	}, nil
}

const emptyRunes = " \t\n\r"

func (c *Cutter) Limit(max int) []string {
	if c.len(c.origin) <= max {
		return []string{c.origin}
	}

	ret := make([]string, 0)
	buf := bytes.NewBuffer(nil)
	max = max - 6
	insert := func() {
		str := strings.Trim(buf.String(), emptyRunes)
		if str == "" {
			return
		}
		ret = append(ret, str)
		buf.Reset()
	}

	for _, s := range c.nodes {
		str := buf.String()
		if c.len(str+s) > max {
			insert()
			buf.WriteString(s)
			continue
		}
		buf.WriteString(s)
	}
	insert()

	for i, _ := range ret {
		if i != len(ret)-1 {
			ret[i] = fmt.Sprintf("%s (%d/%d)", ret[i], i+1, len(ret))
		} else {
			ret[i] = fmt.Sprintf("%s (%d/%d)", ret[i], i+1, len(ret))
		}
	}
	return ret
}
