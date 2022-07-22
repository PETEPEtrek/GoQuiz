package main

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestGetting(t *testing.T) {
	list, err := getExamples("harder_problems.csv", false)
	if err != nil {
		t.Errorf("Problems with file %s", "harder_problems.csv")
	}
	list_ans := [][]string{{"5-5", "0"}, {"1+1", "2"}, {"2*2", "4"}}
	if len(list) != len(list_ans) {
		t.Errorf("Questions are not the same(len)")
	}
	for i, elem := range list {
		if elem[0] != list_ans[i][0] || elem[1] != list_ans[i][1] {
			t.Errorf("Questions or answers are not the same")
		}
	}
}
func TestQuiz(t *testing.T) {
	in, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()

	_, err = io.WriteString(in, "0\n"+"2\n"+"4\n")
	if err != nil {
		t.Fatal(err)
	}
	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		t.Fatal(err)
	}
	list_ans := [][]string{{"5-5", "0"}, {"1+1", "2"}, {"2*2", "4"}}
	correct_answers := Quiz(list_ans, 30, in)
	if correct_answers != 3 {
		t.Errorf("Wrong number of correct answers")
	}
}

func TestQuiz2(t *testing.T) {
	in, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()

	_, err = io.WriteString(in, "1\n"+"2\n"+"4\n")
	if err != nil {
		t.Fatal(err)
	}
	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		t.Fatal(err)
	}
	list_ans := [][]string{{"5-5", "0"}, {"1+1", "2"}, {"2*2", "4"}}
	correct_answers := Quiz(list_ans, 30, in)
	if correct_answers != 2 {
		t.Errorf("Wrong number of correct answers")
	}
}

func TestStripping(t *testing.T) {
	list := stripString([]byte{'a', ' ', '2', ' ', '+', '2'})
	list_ans := []byte{'2', '+', '2'}

	if len(list) != len(list_ans) {
		t.Errorf("Bytes are not the same(len)")
	}
	for i, elem := range list {
		if elem != list_ans[i] {
			t.Errorf("Bytes are not the same")
		}
	}
}
