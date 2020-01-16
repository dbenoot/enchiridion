//   Copyright 2020 The enchiridion Authors
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

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
	file := filepath.Join(wd, "recipes", filename)

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
