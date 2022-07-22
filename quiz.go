package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var timer = time.NewTimer(time.Second * 30)
var duration = 30

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

func getExamples(filename string, shuffle bool) ([][]string, error) {
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
	if shuffle == true {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(equationsAndAnswers), func(i, j int) {
			equationsAndAnswers[i], equationsAndAnswers[j] = equationsAndAnswers[j], equationsAndAnswers[i]
		})
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

func Quiz(list [][]string, duration int) int64 {
	var counter int64 = 0
	var timerSet = false

	input := make(chan int)
	timer = time.NewTimer(time.Duration(duration) * time.Second)
	go getInput(input)
	for i, question := range list {
		fmt.Print("What is the answer for ", question[0], "?\n")
		ans, err := quest(question[1], input, &timerSet)
		if err != nil {
			fmt.Println("Time has passed, next question")
		} else if err != nil && i != len(list)-1 {
			fmt.Println("Time has passed")
		} else if i != len(list)-1 {
			fmt.Println("Next question:")
		}
		counter += ans
		if i == len(list)-1 {
			break
		}
	}
	return counter
}

func quest(ans string, input <-chan int, timerSet *bool) (int64, error) {
	for {
		select {
		case <-timer.C:
			*timerSet = false
			timer.Reset(time.Second * time.Duration(duration))
			return 0, fmt.Errorf("Time out!")

		case answer := <-input:
			isRight := 0
			strAns := strconv.Itoa(answer)
			if strAns == ans {
				isRight = 1
			} else {
				isRight = 0
			}

			if *timerSet {
				if !timer.Stop() {
					<-timer.C
				}
				*timerSet = false
			}

			if !*timerSet {
				*timerSet = true
				timer.Reset(time.Second * time.Duration(duration))
			}
			return int64(isRight), nil

		}
	}
}
func main() {
	var filename string
	var shuffle bool

	flag.StringVar(&filename, "file", "problems.csv", "file with questions and answers")
	flag.IntVar(&duration, "dur", 30, "number of seconds before timeout")
	flag.BoolVar(&shuffle, "shuffle", false, "shuffling the questions")
	flag.Parse()

	list, err := getExamples(filename, shuffle)
	if err != nil {
		log.Fatalf("Problems with file: %s", filename)
	}

	fmt.Println("Press the Enter key to start a quiz")
	fmt.Scanln()

	trueAnswers := Quiz(list, duration)

	fmt.Println("You have", trueAnswers, "correct answers from", len(list))
	if trueAnswers < int64(len(list)/2) {
		fmt.Println("Try harder next time!")
	} else if trueAnswers < 3*int64(len(list)/4) {
		fmt.Println("Not so bad!")
	} else {
		fmt.Println("Great work!")
	}
}
