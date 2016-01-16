package main

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type phaseResult struct {
	Result      result
	Score       float32
	Passed      bool
	Combination []int
}

type phase struct {
	Name            string   `json:"name"`
	Top             int      `json:"top"`
	Rounds          int      `json:"rounds"`
	Threshold       int      `json:"threshold"`
	Static          bool     `json:"static"`
	OponentPath     string   `json:"oponent"`
	OponentTypes    []string `json:"oponentTypes"`
	OponentWarriors []string `json:"-"`
	Bestscore       float32  `json:"-"`
	Total           int      `json:"-"`
	Passed          int      `json:"-"`
}

type phaseResultSorter []phaseResult

func (a phaseResultSorter) Len() int {
	return len(a)
}

func (a phaseResultSorter) Swap(i int, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a phaseResultSorter) Less(i int, j int) bool {
	return a[i].Score > a[j].Score
}

func (a *phaseResultSorter) decode(b []byte) {
	buffer := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buffer)
	dec.Decode(&a)
}

func (r *phaseResult) GetID() string {
	out := ""
	for _, v := range r.Combination {
		out += strconv.Itoa(v) + ","
	}
	return out[:len(out)-1]
}

func (r *phaseResult) decode(b []byte) {
	buffer := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buffer)
	dec.Decode(&r)
}

func (p *phase) GetOponentFileName() string {
	i, j := strings.LastIndex(p.OponentPath, "/"), strings.LastIndex(p.OponentPath, path.Ext(p.OponentPath))
	return p.OponentPath[i+1 : j]
}

func (p *phase) GetFailed() int {
	return p.Total - p.Passed
}

func (p *phase) loadOponents(suite string) {
	if p.OponentPath != "" {
		b, _ := ioutil.ReadFile(p.OponentPath)
		p.OponentWarriors = append(p.OponentWarriors, fixWarrior(b))
	}
	for _, v := range p.OponentTypes {
		files, _ := filepath.Glob(suite + "/" + v + "/*.red")
		for _, w := range files {
			b, _ := ioutil.ReadFile(w)
			p.OponentWarriors = append(p.OponentWarriors, fixWarrior(b))
		}
	}
}

func (p *phase) pass(w string, is []int) phaseResult {
	p.Total++

	aw := 0
	ae := 0
	al := 0

	for _, o := range p.OponentWarriors {
		r := fight(w, o, p.Rounds)
		aw += r.Win
		ae += r.Equal
		al += r.Lose
	}

	score := ((float32(aw) * 3.0) + float32(ae)) * (100.0 / (float32(p.Rounds) * float32(len(p.OponentWarriors))))

	if p.Static {
		if score >= float32(p.Threshold) {
			p.Passed++
			return phaseResult{result{aw, al, ae}, score, true, is}
		}
		return phaseResult{result{aw, al, ae}, score, false, is}
	}

	if p.Bestscore >= score {
		if score >= (p.Bestscore * (float32(p.Threshold) / 100.0)) {
			p.Passed++
			return phaseResult{result{aw, al, ae}, score, true, is}
		}
		return phaseResult{result{aw, al, ae}, score, false, is}
	}

	p.Passed++
	p.Bestscore = score
	return phaseResult{result{aw, al, ae}, score, true, is}
}
