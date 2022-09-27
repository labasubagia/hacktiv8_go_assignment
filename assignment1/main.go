// Assignment 1
package main

import (
	"errors"
	"fmt"
)

type Person struct {
	No      int
	Name    string
	Address string
	Job     string
	Reason  string
}

func FindPersonByNo(people *[]Person, no int) (*Person, error) {
	for _, person := range *people {
		if person.No != no {
			continue
		}
		return &person, nil
	}
	return nil, errors.New("person not found")
}

func main() {
	people := []Person{
		{1, "Agus", "Badung", "Student", "Many Big company use it"},
		{2, "Yudha", "Badung", "Student", "Interest in backend"},
		{3, "Aditya", "Denpasar", "Student", "Try new things"},
		{4, "Restu", "Klungkung", "Student", "Broader my knowledge"},
		{5, "Dipa", "Denpasar", "Student/Freelancer", "Want to become backend engineer"},
		{6, "Alex", "Denpasar", "Student/Freelancer", "Too much free time"},
		{7, "Diqey", "Bandung", "Student", "Want to learn static typing language"},
		{8, "Junda", "Bantul", "Software Engineer", "Work prospect"},
		{9, "Wafda", "Semarang", "Engineering Manager", "Keep knowledge sharp"},
		{10, "Sasono", "Magelang", "Backend Developer", "Try new language"},
	}
	person, err := FindPersonByNo(&people, 1)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%#v\n", person)
}
