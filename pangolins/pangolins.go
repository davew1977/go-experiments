package main

import (
	"os/user"
	"log"
	"path/filepath"
	"os"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"bufio"
	"regexp"
	"strings"
)

var confDir = filepath.Join(homeDir(), ".config", "davew1977", "pangolins")
var questionDb *QuestionNode
var dbFile = filepath.Join(confDir, "db.json")

func main() {
	fmt.Printf("data stored at %s\n", dbFile)
	fmt.Println("Think of an animal, any animal, and I will guess it")
	loadQuestions();
	done := false
	current := questionDb
	reader := bufio.NewReader(os.Stdin)
	for !done  {
		if(!current.isLeaf()) {
			fmt.Println(current.Question)
		} else {
			fmt.Printf("Is it %s?\n", current.Question)
		}
		text, _ := reader.ReadString('\n')
		if isYes(text) {
			if(!current.isLeaf()) {
				current = current.YesBranch
			} else {
				fmt.Println("I knew it!")
				break;
			}
		} else {
			if(!current.isLeaf()) {
				current = current.NoBranch
			} else {
				fmt.Println("What is it then?")
				newAnimal, _ := reader.ReadString('\n')
				newAnimal = strings.TrimSpace(newAnimal)
				fmt.Printf("Enter a question to distinguish between %s and %s\n", newAnimal, current.Question)
				newQuestionTxt, _ := reader.ReadString('\n')
				fmt.Printf("What is the answer for %s?\n", newAnimal)
				newAnswer, _ := reader.ReadString('\n')
				var yb, nb *QuestionNode
				if(isYes(newAnswer)) {
					yb = createLeaf(newAnimal)
					nb = createLeaf(current.Question)
				} else {
					yb = createLeaf(current.Question)
					nb = createLeaf(newAnimal)
				}
				current.Question = strings.TrimSpace(newQuestionTxt)
				current.YesBranch = yb
				current.NoBranch = nb
				save()
				break;
			}
		}
	}

}
func isYes(s string) (bool){
	res, _ :=regexp.MatchString("^y.*|^Y.*", s)
	return res
}

func createLeaf(s string) (*QuestionNode) {
	return &QuestionNode{
		Question: s,
	}
}

func loadQuestions() {
	//try load question db
	os.MkdirAll(confDir, 0777)
	b, err := ioutil.ReadFile(dbFile) // just pass the file name
	if err != nil {
		questionDb = &QuestionNode{
			Question: "Does it live in the sea?",
			YesBranch: &QuestionNode{
				Question: "a whale",
			},
			NoBranch: &QuestionNode{
				Question: "Is it scaley?",
				YesBranch:  &QuestionNode{
					Question: "a pangolin",
				},
				NoBranch: &QuestionNode{
					Question: "a dog",
				},
			},
		}
		save()
	} else {
		questionDb = &QuestionNode{}
		json.Unmarshal(b, &questionDb)
	}
	//create default question struct

}

func save() {
	json_s, _ := json.MarshalIndent(questionDb, "", "  ")
	err := ioutil.WriteFile(dbFile, json_s, 0777);
	if err != nil {
		log.Fatal(err)
	}
}

type QuestionNode struct {
	Question  string
	YesBranch *QuestionNode
	NoBranch  *QuestionNode
}

func (q *QuestionNode) isLeaf() (bool) {
	return q.YesBranch == nil
}

func homeDir() (string) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

