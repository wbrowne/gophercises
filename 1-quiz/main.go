package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	CsvFilename      = "problems.csv"
	MaxProblems      = 100
	TimeoutInSeconds = 30
)

func (p Problem) String() string {
	return fmt.Sprintf("Problem: %s? | Answer: %s", p.question, p.answer)
}

type Problem struct {
	question string
	answer   string
}

func main() {
	problemsFilePtr := flag.String("f", CsvFilename, "a CSV file of problems")
	timeoutPtr := flag.Int("t", TimeoutInSeconds, "a problem timeout value")
	flag.Parse()

	problems, err := getProblemsFromCsv(*problemsFilePtr)

	if err != nil {
		log.Fatal(err)
	}

	if len(problems) > MaxProblems {
		log.Printf("That's too many problems! Try less than %d\n", MaxProblems)
		return
	}

	log.Println("Hit enter to start")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

	var numCorrect = 0
	for _, problem := range problems {
		log.Println(problem.question)

		input := make(chan string)
		go getUserInput(input)

		select {
		case userInput, _ := <-input:
			if userInput == problem.answer {
				numCorrect++
			}
		case <-time.After(time.Duration(*timeoutPtr) * time.Second):
			log.Println("Not quick enough!")
			return
		}
	}

	log.Printf("You answered %d/%d correctly", numCorrect, len(problems))
}

func getUserInput(input chan string) {
	reader := bufio.NewReader(os.Stdin)
	result, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	input <- strings.TrimRight(result, "\n")
}

func getProblemsFromCsv(filename string) ([]Problem, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	rows, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	var problems []Problem
	for _, row := range rows {
		problems = append(problems, Problem{question: row[0], answer: row[1]})
	}

	return problems, nil
}
