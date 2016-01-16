package main

import (
	"bytes"
	"encoding/gob"
	"strconv"
)

func encodeInterface(i interface{}) []byte {
	buffer := new(bytes.Buffer)
	enc := gob.NewEncoder(buffer)
	enc.Encode(i)
	return buffer.Bytes()
}

func decodeInts(b []byte) []int {
	buffer := bytes.NewBuffer(b)
	var i []int
	dec := gob.NewDecoder(buffer)
	dec.Decode(&i)
	return i
}

func decodeBool(b []byte) bool {
	buffer := bytes.NewBuffer(b)
	var i bool
	dec := gob.NewDecoder(buffer)
	dec.Decode(&i)
	return i
}

func intsToString(i []int) string {
	out := ""
	for _, v := range i {
		out += strconv.Itoa(v)
	}
	return out
}
