package main

import (
	"fmt"
	"io"
	"os"
	"time"
	"log"
	"encoding/json"
)

type Task struct {
	Id int `json:"id"`
	Description string `json:"description"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdateAt time.Time `json:"updateAt"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal(fmt.Errorf("need more arguments"))
		return
	}

	switch os.Args[1] {
	case "add":
		err := addTask(os.Args[2:])
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Println("tasks saved successfully")
	default:
		log.Fatal(fmt.Errorf("invalid arguments"))
		return
	}
}

func countTasks(data *os.File) (int, error){
	dec := json.NewDecoder(data)

	count := 0
	for {
		var v Task
		err := dec.Decode(&v)
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
		count++
	}

	return count, nil
}

func addTask(texts []string) error {
	file, err := os.OpenFile("tasks.jsonl", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	numberTasks, err := countTasks(file)
	if err != nil {
		return err
	}

	for i := range texts {
		date := time.Now()
		task := Task{
			Id: numberTasks + 1 + i,
			Description: texts[i],
			Status: "todo",
			CreatedAt: date,
			UpdateAt: date,
		}

		enc := json.NewEncoder(file)
		if err := enc.Encode(task); err != nil {
			return err
		}
	}

	return nil
}
