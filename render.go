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
	var allRecipes []string
	lookupFilename := make(map[string]string)

	// Load all recipes

	recipes := loadRecipes(dir)

	// create a map to lookup the filename for a given recipe

	files, err := getFileList(filepath.Join(dir, "recipes"))
	check(err)

	for _, v := range files {
		recipe := Recipe{}

		content, err := ioutil.ReadFile(v)
		check(err)

		front, _ := getFront(content)

		err = yaml.Unmarshal([]byte(front), &recipe)
		check(err)

		lookupFilename[recipe.Title] = v

	}

	// Load the book info

	if len(b) > 0 {

		var cookbook Book

		filename := b + ".yaml"
		file := filepath.Join(dir, "books", b, filename)

		if _, err := os.Stat(file); err != nil {
			fmt.Println("Book does not exist.")
		}

		if _, err := os.Stat(filepath.Join(dir, "rendered", "books", b, "images")); os.IsNotExist(err) {
			_ = os.MkdirAll(filepath.Join(dir, "rendered", "books", b, "images"), 0755)
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

		// // Check if we need this file for the recipes

		// files, err := getFileList(filepath.Join(dir, "recipes"))
		// check(err)

		// for _, v := range files {
		// 	recipe := Recipe{}

		// 	content, err := ioutil.ReadFile(v)
		// 	check(err)

		// 	front, _ := getFront(content)

		// 	err = yaml.Unmarshal([]byte(front), &recipe)
		// 	check(err)

		// 	lookupFilename[recipe.Title] = v

		// }
		// render body

		t, err := template.ParseFiles(filepath.Join(dir, "books", b, "templates", "page.html"))
		check(err)

		for _, v := range selectedRecipes {

			// instance a bufferstring

			u := bytes.NewBufferString("")

			// preprocess the ingredientslist to a table

			m := yaml.MapSlice{}

			content, err := ioutil.ReadFile(lookupFilename[v.Title])
			front, _ := getFront(content)

			err = yaml.Unmarshal([]byte(front), &m)
			check(err)

			yamlVals, _ := m[6].Value.(yaml.MapSlice)

			var ing = "<table>"

			for _, kv := range yamlVals {
				yamlKey, _ := kv.Key.(string)
				// fmt.Println(yamlKey, v.Ingredients[yamlKey])

				ing = ing + "<tr><td>" + yamlKey + "</td><td>" + v.Ingredients[yamlKey] + "</td></tr>"
			}

			ing = ing + "</table>"

			// get image path

			var img string
			if len(v.Image) > 0 {
				img = filepath.Join("images", v.Image)
				img = "<img src=\"" + img + "\">"
				_, err = copy(filepath.Join(dir, "recipes", "images", v.Image), filepath.Join(dir, "rendered", "books", b, "images", v.Image))
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

		var cover string
		if len(cookbook.Coverpic) > 0 {
			if _, err := os.Stat(filepath.Join(dir, "recipes", "images", cookbook.Coverpic)); !os.IsNotExist(err) {
				_, err = copy(filepath.Join(dir, "recipes", "images", cookbook.Coverpic), filepath.Join(dir, "rendered", "books", b, "images", cookbook.Coverpic))
				check(err)

				cover = "<img src=\"" + filepath.Join("images", cookbook.Coverpic) + "\"/>"
			}
		}

		tb, err := template.ParseFiles(filepath.Join(dir, "books", b, "templates", "book.html"))
		check(err)

		u := bytes.NewBufferString("")

		tb.Execute(u, map[string]interface{}{"author": cookbook.Author, "coverpic": template.HTML(cover), "title": cookbook.Title, "pages": template.HTML(html.UnescapeString(body))})

		// write u to file

		bookdir := filepath.Join(dir, "rendered", "books", b)

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
			_, err = copy(filepath.Join(dir, "books", b, "templates", "book.css"), filepath.Join(bookdir, "book.css"))
			check(err)
		}

	}

	for _, v := range recipes {
		allRecipes = append(allRecipes, v.Title)
	}

	if len(r) > 0 && contains(allRecipes, r) {

		var recipe Recipe

		if _, err := os.Stat(filepath.Join(dir, "rendered", "recipes", r, "images")); os.IsNotExist(err) {
			_ = os.MkdirAll(filepath.Join(dir, "rendered", "recipes", r, "images"), 0755)
		}

		// select the correct recipe

		for _, v := range recipes {
			if v.Title == r {
				recipe = v
			}
		}

		// render recipe

		t, err := template.ParseFiles(filepath.Join(dir, "recipes", "templates", "recipe.html"))
		check(err)

		// instance a bufferstring

		u := bytes.NewBufferString("")

		// preprocess the ingredientslist to a table

		m := yaml.MapSlice{}

		content, err := ioutil.ReadFile(lookupFilename[recipe.Title])
		front, _ := getFront(content)

		err = yaml.Unmarshal([]byte(front), &m)
		check(err)

		yamlVals, _ := m[6].Value.(yaml.MapSlice)

		var ing = "<table>"

		for _, kv := range yamlVals {
			yamlKey, _ := kv.Key.(string)
			// fmt.Println(yamlKey, v.Ingredients[yamlKey])

			ing = ing + "<tr><td>" + yamlKey + "</td><td>" + recipe.Ingredients[yamlKey] + "</td></tr>"
		}

		ing = ing + "</table>"

		// get image path

		var img string
		if len(recipe.Image) > 0 {
			img = filepath.Join("images", recipe.Image)
			img = "<img src=\"" + img + "\">"
			_, err = copy(filepath.Join(dir, "recipes", "images", recipe.Image), filepath.Join(dir, "rendered", "recipes", r, "images", recipe.Image))
			check(err)
		}
		// create the page

		t.Execute(u, map[string]interface{}{
			"recipetitle":       recipe.Title,
			"recipeingredients": ing,
			"recipeimage":       img,
			"recipebody":        string(blackfriday.Run([]byte(recipe.Body))),
			"preptime":          recipe.Preptime,
			"cooktime":          recipe.Cooktime,
			"origin":            recipe.Origin,
			"tags":              recipe.Tags,
		})

		// write u to file

		recipedir := filepath.Join(dir, "rendered", "recipes", r)

		if _, err := os.Stat(recipedir); os.IsNotExist(err) {
			err = os.MkdirAll(recipedir, 0755)
			check(err)
		}

		f, err := os.Create(filepath.Join(recipedir, r+".html"))
		check(err)
		defer f.Close()

		_, err = f.Write(u.Bytes())
		check(err)

		if _, err := os.Stat(filepath.Join(recipedir, "book.css")); os.IsNotExist(err) {
			_, err = copy(filepath.Join(dir, "recipes", "templates", "book.css"), filepath.Join(recipedir, "book.css"))
			check(err)
		}

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
