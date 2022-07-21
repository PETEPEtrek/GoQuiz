package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

var filename = "problems.csv"

func stripString(string []byte) []byte {
	n := 0
	for _, b := range string {
		if b == '(' || b == ')' ||
			('0' <= b && b <= '9') ||
			b == '+' ||
			b == '-' ||
			b == '*' ||
			b == '^' ||
			b == '%' ||
			b == '/' {
			string[n] = b
			n++
		}
	}
	return string[:n]
}

func getExamples(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	if err != nil {
		errorLog.Fatal(err)
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 2
	equationsAndAnswers := make([][]string, 0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		record[0] = string(stripString([]byte(record[0])))
		record[1] = string(stripString([]byte(record[1])))
		equationsAndAnswers = append(equationsAndAnswers, record)
	}
	return equationsAndAnswers, nil
}

func Quiz(list [][]string) int64 {
	var counter int64 = 0
	for i, question := range list {
		fmt.Print("What is an answer for ", question[0], "?\n")
		answer := 0
		for _, err := fmt.Fscan(os.Stdin, &answer); err != nil; {
			fmt.Println("Please, write the answer one more time")
		}
		strAns := strconv.Itoa(answer)
		if strAns == question[1] {
			counter++
		}
		if i == len(list)-1 {
			break
		}
		fmt.Println("Next question:")
	}
	return counter
}

func main() {
	args := os.Args
	if len(args) > 1 {
		filename = args[1]
	}
	list, err := getExamples(filename)
	if err != nil {
		log.Fatalf("Problems with file: %s", filename)
	}
	trueAnswers := Quiz(list)
	fmt.Println("You have", trueAnswers, "correct answers from", len(list))
	if trueAnswers < int64(len(list)/2) {
		fmt.Println("Try harder next time!")
	} else if trueAnswers < 3*int64(len(list)/4) {
		fmt.Println("Not so bad!")
	} else {
		fmt.Println("Great work!")
	}
}
