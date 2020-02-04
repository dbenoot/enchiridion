package main

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func render(dir string, b string, r string) {

	var body string
	var selectedRecipes []Recipe

	// Load all recipes

	recipes := loadRecipes(dir)

	// Load the book info

	var cookbook Book

	filename := b + ".yaml"
	file := filepath.Join(dir, "books", filename)

	if _, err := os.Stat(file); err != nil {
		fmt.Println("Book does not exist.")
	}

	data, err := ioutil.ReadFile(file)
	err = yaml.Unmarshal([]byte(data), &cookbook)
	check(err)

	// get the recipes in the book

	for _, v := range recipes {
		if contains(cookbook.Recipes, v.Title) {
			selectedRecipes = append(selectedRecipes, v)
		}
	}

	// render body

	t, err := template.ParseFiles(filepath.Join(dir, "templates", "page.html"))
	check(err)

	for _, v := range selectedRecipes {

		u := bytes.NewBufferString("")

		t.Execute(u, map[string]string{"recipetitle": v.Title, "recipeingredients": string(v.Preptime), "recipebody": v.Body})

		body = body + u.String()
	}

	// render book

	tb, err := template.ParseFiles(filepath.Join(dir, "templates", "book.html"))
	check(err)

	u := bytes.NewBufferString("")

	tb.Execute(u, map[string]interface{}{"author": cookbook.Author, "coverpic": cookbook.Coverpic, "title": cookbook.Title, "pages": template.HTML(html.UnescapeString(body))})

	// write u to file

	f, err := os.Create(filepath.Join(dir, "rendered", b+".html"))
	check(err)
	defer f.Close()

	_, err = f.Write(u.Bytes())
	check(err)

}
