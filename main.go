package main

import (
	"fmt"

	"github.com/abulhanifah/classroom/helper"
)

func main() {
	fmt.Println("Initial commit")
	db := helper.Connect()
	defer db.Close()
}
