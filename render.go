package main

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/russross/blackfriday.v2"
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

	if _, err := os.Stat(filepath.Join(dir, "rendered", b, "images")); os.IsNotExist(err) {
		_ = os.MkdirAll(filepath.Join(dir, "rendered", b, "images"), 0755)
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

		// instance a bufferstring

		u := bytes.NewBufferString("")

		// preprocess the ingredientslist to a table

		var ing = "<table>"
		for k, v := range v.Ingredients {
			ing = ing + "<tr><td>" + k + "</td><td>" + v + "</td></tr>"
		}
		ing = ing + "</table>"

		// get image path

		var img string
		if len(v.Image) > 0 {
			img = filepath.Join("images", v.Image)
			img = "<img src=\"" + img + "\">"
			_, err = copy(filepath.Join(dir, "recipes", "images", v.Image), filepath.Join(dir, "rendered", b, "images", v.Image))
			check(err)
		}
		// create the page

		t.Execute(u, map[string]interface{}{
			"recipetitle":       v.Title,
			"recipeingredients": ing,
			"recipeimage":       img,
			"recipebody":        string(blackfriday.Run([]byte(v.Body))),
			"preptime":          v.Preptime,
			"cooktime":          v.Cooktime,
			"origin":            v.Origin,
			"tags":              v.Tags,
		})

		// add the rendered page to the HTML body of the book

		body = body + u.String()
	}

	// render book

	tb, err := template.ParseFiles(filepath.Join(dir, "templates", "book.html"))
	check(err)

	u := bytes.NewBufferString("")

	tb.Execute(u, map[string]interface{}{"author": cookbook.Author, "coverpic": cookbook.Coverpic, "title": cookbook.Title, "pages": template.HTML(html.UnescapeString(body))})

	// write u to file

	bookdir := filepath.Join(dir, "rendered", b)

	if _, err := os.Stat(bookdir); os.IsNotExist(err) {
		err = os.MkdirAll(bookdir, 0755)
		check(err)
	}

	f, err := os.Create(filepath.Join(bookdir, b+".html"))
	check(err)
	defer f.Close()

	_, err = f.Write(u.Bytes())
	check(err)

	if _, err := os.Stat(filepath.Join(bookdir, "book.css")); os.IsNotExist(err) {
		_, err = copy(filepath.Join(dir, "templates", "book.css"), filepath.Join(bookdir, "book.css"))
		check(err)
	}

}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
