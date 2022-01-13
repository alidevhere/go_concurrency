package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"text/template"
	"time"
)

const tmpl = `{{range .}}My name is {{.Name}}. I am {{.Age}} years old. You can contact me at my phone number {{.Phone}}
{{end}}`

type Person struct {
	Name  string
	Age   uint8
	Phone string
}

func main() {

	// Without Go routines
	start := time.Now()
	Names := make(chan string, 20000)
	Numbers := make(chan string, 20000)
	Ages := make(chan uint8, 20000)
	readNames("./names.txt", Names)
	readNumbers("./phone.txt", Numbers)
	readAges("./age.txt", Ages)
	renderTemplate(Names, Ages, Numbers)
	end := time.Since(start)

	println("Took ", end, "  Nano seconds to complete")

	//	END

	//n, ok := <-Names
	//println("VALUES++ ", n, ok)
	/*
		for i := 0; i < 10; i++ {
			name := <-Names
			no := <-Numbers
			age := <-Ages
			fmt.Println("Name=", name, "Number=", no, "Age=", age)
		}
	*/
}

func readNames(path string, name chan<- string) {
	defer close(name)
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var i int = 0
	for scanner.Scan() {
		i++
		//fmt.Print(i, scanner.Text(), "\n")
		name <- scanner.Text()

	}

}

func readNumbers(path string, number chan<- string) {
	file, err := os.Open(path)
	if err != nil {
		print(err.Error())
	}
	defer close(number)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		number <- scanner.Text()
		//	println(scanner.Text())
	}

}

func readAges(path string, age chan<- uint8) {
	file, err := os.Open(path)
	if err != nil {
		print(err.Error())
	}
	defer file.Close()
	//defer close(age)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		i, err := strconv.ParseUint(scanner.Text(), 10, 8)
		if err != nil {
			println(err.Error())
		}
		age <- uint8(i)

	}

}

func renderTemplate(names <-chan string, ages <-chan uint8, number <-chan string) {
	var persons []Person
	l := len(names)
	for i := 0; i < l; i++ {
		//println("Iteration=", i)
		persons = append(persons, Person{Name: <-names, Age: <-ages, Phone: <-number})
	}
	//print("Persons= ", len(persons))

	t := template.New("Persons")
	tmp, err := t.Parse(tmpl)

	if err != nil {
		println(err.Error())
	}
	if err := tmp.Execute(os.Stdout, persons); err != nil {
		println(err.Error())
	}
}
