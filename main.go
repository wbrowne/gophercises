package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	CsvFilename      = "problems.csv"
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
	problemsFilePtr := flag.String("f", CsvFilename, "a CSV file of problems with format <question>,<answer>")
	timeoutPtr := flag.Int("t", TimeoutInSeconds, "quiz timeout value")
	shufflePtr := flag.Bool("s", true, "shuffle order of problems")
	flag.Parse()

	problems, err := getProblemsFromCsv(*problemsFilePtr)

	if err != nil {
		log.Fatal(err)
	}

	if *shufflePtr {
		problems = shuffle(problems)
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

func shuffle(p []Problem) []Problem {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(p), func(i, j int) { p[i], p[j] = p[j], p[i] })

	return p
}

func getUserInput(input chan string) {
	reader := bufio.NewReader(os.Stdin)
	result, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	input <- normalizeString(result)
}

func normalizeString(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	return s
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
		problems = append(problems, Problem{question: row[0], answer: normalizeString(row[1])})
	}

	return problems, nil
}
