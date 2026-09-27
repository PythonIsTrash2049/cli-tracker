package main

import (
	"fmt"
	"os"
)

type Task struct {
	Id int `json:"id"`
	Description string `json:"description"`
	Status string `json:"status"`
	CreatedAt time.Date `json:"createdAt"`
	UpdateAt time.Date `json:"updateAt"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal(fmt.Errorf("need more arguments"))
		return
	}

	switch os.Args[1] {
	default:
		log.Fatal(fmt.Errorf("invalid arguments"))
		return
	}
}
