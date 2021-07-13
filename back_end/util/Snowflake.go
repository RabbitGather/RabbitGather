package util

import (
	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func init() {
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		panic(err.Error())
	}
}
func Snowflake() snowflake.ID {
	// Create a new Node with a Node number of 1
	//node, err := snowflake.NewNode(1)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//// Generate a snowflake id.
	return node.Generate()

	// Print out the id in a few different ways.
	//fmt.Printf("Int64  id: %d\n", id)
	//fmt.Printf("String id: %s\n", id)
	//fmt.Printf("Base2  id: %s\n", id.Base2())
	//fmt.Printf("Base64 id: %s\n", id.Base64())
	//
	//// Print out the id's timestamp
	//fmt.Printf("id Time  : %d\n", id.Time())
	//
	//// Print out the id's node number
	//fmt.Printf("id Node  : %d\n", id.Node())
	//
	//// Print out the id's sequence number
	//fmt.Printf("id Step  : %d\n", id.Step())
	//
	//// Generate and print, all in one.
	//fmt.Printf("id       : %d\n", node.Generate().Int64())
}
