package main

import (
	"fmt"

	"gitlab.com/abulhanifah/classroom/configs"
	"gitlab.com/abulhanifah/classroom/helpers"
	"gitlab.com/abulhanifah/classroom/routes"
)

func main() {
	fmt.Println("Initial commit")
	db := helpers.Connect()
	defer db.Close()

	r := routes.Init(db)
	err := r.Start(":" + configs.Get("APP_PORT").String())
	if err != nil {
		r.Logger.Fatal(err)
	}
}
