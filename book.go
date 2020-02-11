package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func parseBook(dir string, book string, ab string, rb string) {

	var cookbook Book
	recipes := loadRecipesTitles(dir)

	// create filename

	filename := book + ".yaml"
	file := filepath.Join(dir, "books", book, filename)

	if _, err := os.Stat(file); err != nil {
		_ = os.MkdirAll(filepath.Join(dir, "books", book, "templates"), 0755)
		os.Create(file)
		fmt.Println("Created ", file, ".")

		for k := range templ {
			_, err = copy(filepath.Join(dir, "templates", k), filepath.Join(dir, "books", book, "templates", k))
		}
	}

	data, err := ioutil.ReadFile(file)
	err = yaml.Unmarshal([]byte(data), &cookbook)

	if len(cookbook.Title) == 0 {
		cookbook.Title = book
	}

	if len(cookbook.Author) == 0 {
		cookbook.Author = "Enchiridion"
	}

	if len(cookbook.Coverpic) == 0 {
		cookbook.Coverpic = ""
	}

	if contains(recipes, ab) && !contains(cookbook.Recipes, ab) {
		cookbook.Recipes = append(cookbook.Recipes, ab)
	}

	if contains(recipes, rb) && contains(cookbook.Recipes, rb) {
		cookbook.Recipes = unappend(cookbook.Recipes, rb)
	}

	// marshal cookbook in the yaml format

	content, err := yaml.Marshal(&cookbook)

	// write the file

	err = ioutil.WriteFile(file, []byte(content), 0644)
	check(err)

}

func loadRecipesTitles(dir string) []string {
	var r []string

	files, err := getFileList(filepath.Join(dir, "recipes"))

	for _, v := range files {
		recipe := Recipe{}

		content, err := ioutil.ReadFile(v)

		front, _ := getFront(content)

		err = yaml.Unmarshal([]byte(front), &recipe)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		r = append(r, recipe.Title)

		check(err)
	}

	check(err)

	return r
}

func unappend(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
