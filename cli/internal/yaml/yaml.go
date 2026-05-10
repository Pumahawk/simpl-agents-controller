package yaml

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

type Obj struct {
	lines []string
}

func (o *Obj) UpdateAttribute(value string, path ...string) (bool, error) {
	deep := 0
	for i, line := range o.lines {
		if deep < len(path) {
			s := strings.Repeat(" ", deep*2)
			rx := "^" + s + path[deep] + ":"
			r, err := regexp.Compile(rx)
			if err != nil {
				return false, err
			}
			if r.MatchString(line) {
				deep++
				if deep == len(path) {
					loc := regexp.MustCompile(rx).FindStringIndex(line)
					if loc == nil {
						return false, fmt.Errorf("unexpected not match regex=%q, line=%q", rx, line)
					}
					o.lines[i] = fmt.Sprintf("%s %q", string(line[loc[0]:loc[1]]), value)
					return true, nil
				}
			} else {
				if deep > 0 {
					s := strings.Repeat(" ", (deep-1)*2)
					if !regexp.MustCompile("^" + s).MatchString(line) {
						deep--
					}
				}
			}
		}
	}
	return false, nil
}

func (o *Obj) Bytes() []byte {
	bf := &bytes.Buffer{}
	for _, line := range o.lines {
		fmt.Fprintf(bf, "%s\n", line)
	}
	return bf.Bytes()
}

func NewObj(content []byte) *Obj {
	var lines []string
	bt := bytes.NewBuffer(content)
	sc := bufio.NewScanner(bt)
	for sc.Scan() {
		lines = append(lines, string(sc.Bytes()))
	}
	return &Obj{
		lines: lines,
	}
}
