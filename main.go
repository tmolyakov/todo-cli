package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tmolyakov/todo-cli/internal/storage/jsonfile"
	"github.com/tmolyakov/todo-cli/internal/todo"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	repo := jsonfile.New("tasks.json")
	service := todo.NewService(repo)

	for {
		fmt.Println("\n1. Todo list")
		fmt.Println("2. Add todo item")
		fmt.Println("3. Show todo item")
		fmt.Println("4. Remove todo item")
		fmt.Println("5. Exit")
		fmt.Println("> ")

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Println(os.Stderr, err)
			}
			return
		}

		switch strings.TrimSpace(scanner.Text()) {
		case "1":
			list, err := service.List()
			if err != nil {
				fmt.Println(os.Stderr, err)
				continue
			}
			if len(list) == 0 {
				fmt.Println("no tasks")
			}
			for _, t := range list {
				fmt.Printf("%d. %s [%s, %s]\n", t.ID, t.Text, t.Status, t.Priority)
			}

		case "2":
			fmt.Println("2")
		}
	}
}
