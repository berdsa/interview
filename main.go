package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
)

type Numbers struct {
	A int `json:"a"`
	B int `json:"b"`
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <filename> <numRoutines>")
		os.Exit(1)
	}

	filename := os.Args[1]
	//filename := "numbers.json"
	numRoutines, err := strconv.Atoi(os.Args[2])
	//numRoutines := 2
	if err != nil {
		fmt.Println("Error: numRoutines must be an integer")
		os.Exit(1)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	var numbers []Numbers
	err = json.Unmarshal(data, &numbers)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		os.Exit(1)
	}

	sumChan := make(chan int, numRoutines)
	var wg sync.WaitGroup

	chunkSize := len(numbers) / numRoutines
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sum := 0
			for j := i * chunkSize; j < (i+1)*chunkSize; j++ {
				sum += numbers[j].A + numbers[j].B
			}
			fmt.Println(sum)
			sumChan <- sum
		}(i)
	}

	go func() {
		wg.Wait()
		close(sumChan)
	}()

	totalSum := 0
	for sum := range sumChan {
		totalSum += sum
	}

	fmt.Println("Total sum:", totalSum)
}
