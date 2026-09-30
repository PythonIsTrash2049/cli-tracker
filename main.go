package main

import (
	"fmt"
	"io"
	"os"
	"time"
	"log"
	"encoding/json"
	"strconv"
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
	case "list":
		if len(os.Args) != 2 {
			log.Fatal("invalid arguments")
			return
		}

		err := list()
		if err != nil {
			log.Fatal(err)
			return
		}

	case "update":
		if len(os.Args) != 4 {
			log.Fatal("invalid arguments")
			return
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
			return
		}

		err = update(id, os.Args[3])
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Println("task update successfully")
	default:
		log.Fatal(fmt.Errorf("invalid arguments"))
		return
	}
}

func list() error {
	file, err := os.OpenFile("tasks.jsonl", os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	count, err := countTasks(file)
	if err != nil {
		return err
	}
	defer file.Close()

	if count == 0 {
		fmt.Println("You have no tasks!")
		return nil
	}
	fmt.Printf("You have %d tasks\n", count)
	file.Seek(0, 3)

	dec := json.NewDecoder(file)

	for {
		var task Task
		err = dec.Decode(&task)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		fmt.Printf("Task: %d\n    Description: %s\n    Status: %s\n    Created at: %v\n    Last update at: %v\n", task.Id, task.Description, task.Status, task.CreatedAt.Format("02.01.2006 15:04:05"), task.UpdateAt.Format("02.01.2006 15:04:05"))
	}

	return nil
}

func update(id int, text string) error{
	file, err := os.OpenFile("tasks.jsonl", os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	temp, err := os.CreateTemp("", "temp-*")
	if err != nil {
		return err
	}
	defer os.Remove(temp.Name())
	defer temp.Close()

	dec := json.NewDecoder(file)
	enc := json.NewEncoder(temp)

	updated := false
	for {
		var task Task
		err = dec.Decode(&task)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if task.Id == id {
			updated = true
			task.Description = text
			task.UpdateAt = time.Now()
		}

		err = enc.Encode(task)
		if err != nil {
			return err
		}
	}

	if !updated {
		return fmt.Errorf("task not found")
	}

	file.Seek(0, 3)
	temp.Seek(0, 3)
	file.Truncate(0)
	
	if _, err = io.Copy(file, temp); err != nil {
		return err
	}

	return nil
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
