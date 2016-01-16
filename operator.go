package main

import "strconv"

type operator interface {
	max() int
	min() int
	size() int
	getNext(int) int
	insertIntoLine(int) string
	isMax(int) bool
}

type deepOperator struct {
	Base string
	Min  int
	Max  int
}

func (o deepOperator) max() int {
	return o.Max
}

func (o deepOperator) min() int {
	return o.Min
}

func (o deepOperator) size() int {
	return (o.Max - o.Min) + 1
}

func (o deepOperator) getNext(i int) int {
	return (o.Min + i) % (o.Max + 1)
}

func (o deepOperator) insertIntoLine(i int) string {
	return deepOperatorRegex.ReplaceAllString(o.Base, strconv.Itoa(o.getNext(i)))
}

func (o deepOperator) isMax(i int) bool {
	return i == o.size()
}

type flatOperator deepOperator

func (f flatOperator) max() int {
	return f.Max
}

func (f flatOperator) min() int {
	return f.Min
}

func (f flatOperator) size() int {
	return (f.Max - f.Min) + 1
}

func (f flatOperator) getNext(i int) int {
	return i
}

func (f flatOperator) insertIntoLine(i int) string {
	return flatOperatorRegex.ReplaceAllString(f.Base, strconv.Itoa(i))
}

func (f flatOperator) isMax(i int) bool {
	return true
}

func parseLine(line string) (operator, bool) {
	flatMatch := flatOperatorRegex.FindAllStringSubmatch(line, 1)
	deepMatch := deepOperatorRegex.FindAllStringSubmatch(line, 1)

	if len(flatMatch) > 0 {
		min, _ := strconv.Atoi(flatMatch[0][1])
		max, _ := strconv.Atoi(flatMatch[0][2])
		return flatOperator{line, min, max}, true
	}

	if len(deepMatch) > 0 {
		min, _ := strconv.Atoi(deepMatch[0][1])
		max, _ := strconv.Atoi(deepMatch[0][2])
		return deepOperator{line, min, max}, true
	}

	return deepOperator{}, false
}
