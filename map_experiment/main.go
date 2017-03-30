package main

import (
	"fmt"
)

type key int
type key2 int
const val1 key = 0
const val2 key2 = 0
const val3 int = 0
const val4 int = 0

func main() {

	aMap := make(map[interface{}]string)

	aMap[val1] = "wow"
	aMap[val2] = "wool"
	aMap[val3] = "woollp"
	aMap[val4] = "wyuiyus"
	fmt.Println(aMap)

	fmt.Println(len(aMap))
}
