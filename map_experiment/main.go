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

	fmt.Println(BinaryStringFromByte(1))
	fmt.Println(BinaryStringFromByte(255))
	fmt.Println(BinaryStringFromByte(16))
	fmt.Println(BinaryStringFromByte(128))
}

func BinaryStringFromByte(b byte) string {
	return fmt.Sprintf("%08b", b) //different from monolith equivalent because go byte is unsigned
}
