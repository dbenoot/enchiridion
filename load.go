package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	header    = "---\n"   // front matter file header
	separator = "\n---\n" // front matter metadata/content separator
)

func loadRecipes(dir string) []Recipe {
	var R []Recipe

	files, err := getFileList(dir)

	for _, v := range files {
		recipe := Recipe{}

		content, err := ioutil.ReadFile(v)

		front, body := getFront(content)

		err = yaml.Unmarshal([]byte(front), &recipe)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		recipe.Body = body

		R = append(R, recipe)

		check(err)
	}

	check(err)

	return R
}

func getFileList(wd string) ([]string, error) {

	fileList := []string{}

	err := filepath.Walk(wd, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})

	fileList = filterFile(fileList, wd)

	return fileList, err
}

func filterFile(f []string, wd string) []string {

	var fo []string

	// use regexp [0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{4}_['\\w,\\s-\\.]*\\.md
	//
	// meaning
	// [0-9]{4} : 4 digits
	// - : the character -
	// [0-9]{2} : 2 digits
	// - : the character -
	// [0-9]{2} : 2 digits
	// _ : the character _
	// ['\\w,\\s-\\.]* : text consisting of  ', all word symbols, all whitespace symbols, dashes, dots
	// \\.md : .md
	//
	// Remark the double escape for w, s and . -> otherwise the string parser complains (and '' didn't work...)

	var r = regexp.MustCompile("[0-9]{4}[0-9]{2}[0-9]{2}_['\\w,\\s-\\.]*\\.md")

	for _, file := range f {
		if r.MatchString(file) { //&& strings.Contains(file, ".md") {
			fo = append(fo, file)
		} else {
			fi, _ := os.Stat(file)
			if fi.Mode().IsRegular() == true {
				fmt.Printf("File was not included in the filterlist %s. Please check filterFile function. \n", file)
			}
		}
	}

	return fo
}

func getFront(data []byte) (string, string) {

	var metadata string

	//always working as string
	txt := string(data)

	content := txt //by default

	if strings.HasPrefix(txt, header) { // there is a header, therefore there MUST be a front matter

		//we remove the prefix
		txt = strings.TrimPrefix(txt, header)
		// nice trick: we split using the separator
		// and we hope the get: metadata (valid yaml) and the content
		// all the rest is check

		splitted := strings.SplitN(txt, separator, 2)

		if len(splitted) != 2 {
			return "", "found a heading --- without separator ---"
		}

		metadata = splitted[0]
		content = splitted[1]

	}

	return metadata, content
}
