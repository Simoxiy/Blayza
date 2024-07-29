package main

import "github.com/ysnbhb/group-tracker/serve"

func main() {
	art := serve.NewPort(":8080")
	art.Start()
}
