package main

import (
	"fmt"

	"github.com/abulhanifah/classroom"
)

func main() {
	fmt.Println("Initial commit")
	db := classroom.Connect()
	defer db.Close()
}
