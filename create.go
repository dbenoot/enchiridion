package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

func createEntry(wd string, title string) {

	var r Recipe

	r.Title = title

	// create filename

	t := time.Now()

	filename := t.Format("20060102") + "_" + title + ".md"
	file := filepath.Join(wd, filename)

	// Create the markdown file as YYYYMMDD_title.md

	if _, err := os.Stat(file); err == nil {

		fmt.Println("A recipe entry for this date and title already exists.")

	} else {

		os.Create(file)
		fmt.Println("Created ", file, ".")

		// Add title and tags to file
		// Entry title has md frontmatter delimiter ---

		d, err := yaml.Marshal(&r)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		frontmatter := strings.Replace(string(d), "body: \"\"", "", -1)
		content := "---\n" + frontmatter + "\n---\n"

		// write the file

		err = ioutil.WriteFile(file, []byte(content), 0644)
		check(err)

	}
}
