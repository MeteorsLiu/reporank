package main

import "testing"

func TestRepoStats(t *testing.T) {
	t.Log(New("llcppg", "goplus").Score())
	t.Log(New("linux", "torvalds").Score())

}
