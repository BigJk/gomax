package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"strconv"
	"strings"
)

type warrior struct {
	Lines         []string
	Operators     map[int]operator
	OperatorLines []int
}

func (w *warrior) decode(b []byte) {
	buffer := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buffer)
	dec.Decode(&w)
}

func (w *warrior) push(is []int) {
	wcode := w.construct(is)
	if currentConfig.EocTest {
		if !eocTest(wcode, 1) {
			//log.Println("EocTest failed!", is)
			return
		}
	}
	for i, p := range currentConfig.Phases {
		res := currentConfig.Phases[i].pass(wcode, is)
		pushChange("phase"+strconv.Itoa(i), intsToString(is), encodeInterface(res))
		pushWarriorToTop(i, res)
		log.Println(p.Name, "result", res.Combination, res.Score)
		if !res.Passed {
			break
		}
	}
}

func (w *warrior) size() int {
	s := w.getOperator(0).size() + 1
	for i := 1; i < w.operatorSize(); i++ {
		s *= w.getOperator(i).size() + 1
	}
	return s
}

func (w *warrior) operatorSize() int {
	return len(w.Operators)
}

func (w *warrior) getOperator(i int) operator {
	return w.Operators[w.OperatorLines[i]]
}

func (w *warrior) construct(is []int) string {
	var buffer bytes.Buffer
	oi := 0
	for i := 0; i < len(w.Lines); i++ {
		if o, ok := w.Operators[i]; ok {
			buffer.WriteString(o.insertIntoLine(is[oi]) + "\n")
			oi++
		} else {
			buffer.WriteString(w.Lines[i] + "\n")
		}
	}
	return buffer.String()
}

func parseWarrior(code string) warrior {
	lines := strings.Split(code, "\n")
	operators := make(map[int]operator)
	var operatorLines []int
	for i := 0; i < len(lines); i++ {
		o, t := parseLine(lines[i])
		if t {
			operators[i] = o
			operatorLines = append(operatorLines, i)
		}
	}
	return warrior{lines, operators, operatorLines}
}
