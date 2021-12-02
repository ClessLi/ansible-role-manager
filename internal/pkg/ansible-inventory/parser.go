package ansible_inventory

import (
	"bytes"
	"errors"
	"strings"
)

type Parser interface {
	Parse(data []byte) (Group, error)
	Dump(g Group) ([]byte, error)
}

type parser struct {
}

func NewParser() Parser {
	p := new(parser)
	return Parser(p)
}

func (p parser) Parse(data []byte) (Group, error) {
	if len(data) == 0 {
		return nil, errors.New("data is null")
	}

	g := newGroup()
	peek := 0
	isFirstGroup := true

	for idx := 0; idx < len(data); idx++ {
		if data[idx] == '\n' || idx == len(data)-1 {
			line := string(data[peek:idx])
			if idx == len(data)-1 {
				line = string(data[peek:])
			}
			peek = idx + 1
			line = strings.TrimSpace(line)
			if h := ParseHost(line); h != nil {
				_ = g.AddHost(h)
			} else if line[0] == '[' && line[len(line)-1] == ']' {
				if !isFirstGroup {
					break
				}
				err := g.setName(line[1 : len(line)-1])
				if err != nil {
					return nil, err
				}
				isFirstGroup = false
			}
		}
	}
	return g, nil
}

func (p parser) Dump(g Group) ([]byte, error) {
	if g == nil || g.GetName() == "" {
		return nil, errors.New("nil Group input")
	}
	buff := bytes.NewBuffer(make([]byte, 0))
	buff.WriteString("[" + g.GetName() + "]\n")
	for _, h := range g.GetHosts() {
		buff.WriteString(h.GetIPString() + "\n")
	}
	return buff.Bytes(), nil
}
