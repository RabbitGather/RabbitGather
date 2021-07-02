package main

import (
	"fmt"
	"github.com/google/uuid"
)
func main() {
	//theuuid := uuid.New()
	a := [16]byte(uuid.New())
fmt.Println(string(a[:]))
fmt.Println(string(uuid.New().NodeID()))
}