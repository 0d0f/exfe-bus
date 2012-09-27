package textcutter

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
	nodes []string
	len   LengthFunc
}

func Parse(text string, f LengthFunc) (*Cutter, error) {
	nodes, err := parseLexical(text)
	if err != nil {
		return nil, err
	}
	return &Cutter{
		nodes: nodes,
		len:   f,
	}, nil
}

const emptyRunes = " \t\n\r"

func (c *Cutter) Limit(max int) []string {
	lens := make([]int, len(c.nodes))
	length := 0
	for i, s := range c.nodes {
		lens[i] = c.len(s)
		length += lens[i]
	}
	if length <= max {
		return []string{strings.Join(c.nodes, "")}
	}

	ret := make([]string, 0)
	buf := bytes.NewBuffer(nil)
	max = max - 5
	insert := func() {
		str := strings.Trim(buf.String(), emptyRunes)
		if str == "" {
			return
		}
		ret = append(ret, str)
		buf.Reset()
	}

	for i, s := range c.nodes {
		if length+lens[i] > max {
			insert()
			buf.WriteString(s)
			length = lens[i]
			continue
		}
		buf.WriteString(s)
		length += lens[i]
	}
	insert()

	for i, _ := range ret {
		ret[i] = fmt.Sprintf("(%d/%d)%s", i+1, len(ret), ret[i])
	}
	return ret
}
