package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
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

		err = copydir(filepath.Join(dir, "books", "templates"), filepath.Join(dir, "books", book, "templates"))
		check(err)
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

func copydir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = copydir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = copyfile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func copyfile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}
