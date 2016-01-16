package main

import "strconv"

type operator struct {
	Base string
	Min  int
	Max  int
}

func (o operator) size() int {
	return (o.Max - o.Min) + 1
}

func (o operator) getNext(i int) int {
	return (o.Min + i) % (o.Max + 1)
}

func (o operator) insertIntoLine(i int) string {
	return flatOperator.ReplaceAllString(o.Base, strconv.Itoa(o.getNext(i)))
}

func parseLine(line string) (operator, bool) {
	match := flatOperator.FindAllStringSubmatch(line, 1)
	if len(match) == 0 {
		return operator{}, false
	}
	min, _ := strconv.Atoi(match[0][1])
	max, _ := strconv.Atoi(match[0][2])
	return operator{line, min, max}, true
}
