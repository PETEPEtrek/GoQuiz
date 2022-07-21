package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

var filename = "problems.csv"
var timer = time.NewTicker(time.Second * 30)

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

func getInput(input chan int) {
	for {
		answer := 0
		for _, err := fmt.Fscan(os.Stdin, &answer); err != nil; {
			fmt.Println("Please, write the answer one more time")
		}
		input <- answer
	}
}

func Quiz(list [][]string) int64 {
	var counter int64 = 0
	input := make(chan int)
	tim := time.NewTimer(time.Second * 2)
	go getInput(input)
	for i, question := range list {
		fmt.Print("What is an answer for ", question[0], "?\n")
		ans, err := quest(question[1], tim.C, input)
		if err != nil {
			fmt.Println("Time has passed, next question")
		}
		counter += ans
		if i == len(list)-1 {
			break
		}
	}
	return counter
}

func quest(ans string, timer <-chan time.Time, input <-chan int) (int64, error) {
	for {
		select {
		case <-timer:

			return 0, fmt.Errorf("Time out!")
		case answer := <-input:
			isRight := 0
			strAns := strconv.Itoa(answer)
			if strAns == ans {
				isRight = 1
			} else {
				isRight = 0
			}
			fmt.Println("Next question:")
			return int64(isRight), nil
		}
	}
}
func main() {
	args := os.Args
	if len(args) > 1 && args[1] == "true" {
		fmt.Print("Please, write a new name of csv file: ")
		_, err := fmt.Fscan(os.Stdin, &filename)
		if err != nil {
			log.Fatalf("Can't get the filename")
		}
	}
	if len(args) > 2 {
		num, err := strconv.Atoi(args[2])
		if err != nil {
			log.Fatalf("Not a number: %s", args[2])
		}
		timer = time.NewTicker(time.Second * time.Duration(num))
	}

	list, err := getExamples(filename)
	if err != nil {
		log.Fatalf("Problems with file: %s", filename)
	}

	fmt.Println("Press the Enter key to start a quiz")
	fmt.Scanln()
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
