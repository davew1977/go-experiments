package main

import "testing"

func TestIsYes(t *testing.T) {
	assertYes(t, "yes", "y", "Y", "YES", "yes")
	assertNo(t, "no", "NO", "pig", "n", "N")
}
func assertYes(t *testing.T, args ...string) {
	assert(t, true, args...)
}
func assertNo(t *testing.T, args ...string) {
	assert(t, false, args...)
}
func assert(t *testing.T, expected bool, args ...string) {
	for _,s := range args	 {
		if(isYes(s) != expected) {
			t.Fatal(s);
		}
	}
}

