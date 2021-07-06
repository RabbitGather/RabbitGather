package main

import (
	"fmt"
	"path/filepath"
)
func main() {
	//theuuid := uuid.New()
	//a := [16]byte(uuid.New())
	p ,_:=filepath.Abs("ssl/crt/meowalien_com.crt")
	fmt.Println("ssl/crt/meowalien_com.crt  --  ",p) // F:\GoTest\GoTest\master.exe <nil>
//fmt.Println(string(a[:]))
//fmt.Println(string(uuid.New().NodeID()))
}