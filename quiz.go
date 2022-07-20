package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func getExamples(filename string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Couldn't open the file :", filename)
		return false, err
	}
	defer file.Close()
	reader := csv.NewReader(file)

}

func main() {

}
