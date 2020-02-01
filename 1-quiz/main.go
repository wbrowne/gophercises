package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	CsvFilename = "problems.csv"
	MaxProblems = 100
)

func (p Problem) String() string {
	return fmt.Sprintf("Problem: %s? | Answer: %s", p.question, p.answer)
}

type Problem struct {
	question string
	answer   string
}

func main() {
	problems, err := getProblemsFromCsv(CsvFilename)

	if len(problems) > MaxProblems {
		log.Printf("That's too many problems! Try less than %d\n", MaxProblems)
		return
	}

	if err != nil {
		log.Println("Error:", err)
	}

	reader := bufio.NewReader(os.Stdin)
	var numCorrect = 0
	for _, problem := range problems {
		log.Print(problem.question)

		userInput := getUserInput(reader)
		if userInput == problem.answer {
			numCorrect++
		}
	}

	log.Printf("You answered %d/%d correctly", numCorrect, len(problems))
}

func getUserInput(reader *bufio.Reader) string {
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		return strings.TrimRight(input, "\n")
	}
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
