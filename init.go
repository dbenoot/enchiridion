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
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/go-ini/ini"
)

func initEnchiridion(dir string) {

	var templ = map[string]string{
		filepath.Join(dir, "books", "templates", "book.html"):     "<html>\n<head>\n<title>{{.title}}</title>\n<meta name=\"author\" content=\"{{.author}}}\"/>\n<link rel=\"stylesheet\" type=\"text/css\" href=\"book.css\"/>\n</head>\n<body>\n\n<section class=\"frontcover\">\n<img src=\"{{.coverpic}}\"/>\n<h1>{{.title}}</h1>\n</section>\n\n<section class=\"imprint\">\n<p>{{.author}}</p>\n</section>\n\n{{.pages}}\n\n</body>\n</html>",
		filepath.Join(dir, "books", "templates", "page.html"):     "<section class=\"chapter\" id=\"html-h-1\">\n\n\t{{.recipeimage}}\n\n\t<h1>{{.recipetitle}}</h1>\n\t<h2>Preparation time: {{.preptime}} - Cooking time: {{.cooktime}}</h2>\n\t<h3>Origin: {{.origin}}</h3>\n\t<h3>Tags: {{.tags}}</h3>\n\t<p class=\"sidenote\">{{.recipeingredients}}</p>\n\t<p>{{.recipebody}}</p>\n</section>",
		filepath.Join(dir, "books", "templates", "book.css"):      "",
		filepath.Join(dir, "recipes", "templates", "recipe.html"): "<html><head><title>{{.title}}</title><meta name=\"author\" content=\"{{.origin}}}\"/><link rel=\"stylesheet\" type=\"text/css\" href=\"recipe.css\"/></head><body><section class=\"chapter\">{{.recipeimage}}<h1>{{.recipetitle}}</h1><h2>Preparation time:{{.preptime}}- Cooking time:{{.cooktime}}</h2><h3>Author:{{.origin}}</h3><h3>Tags:{{.tags}}</h3><div class=\"sidebar sidenote\"><h2>Ingredients</h2><p>{{.recipeingredients}}</p></div><p>{{.recipebody}}</p></section></body></html>",
		filepath.Join(dir, "recipes", "templates", "recipe.css"):  "",
	}

	if checkSiteNotExist(dir) == false {
		fmt.Println("Directory not empty.")
	} else {

		usr, err := user.Current()
		check(err)

		sd := filepath.Join(usr.HomeDir, ".enchiridion")
		cfgFile := filepath.Join(sd, "config.ini")

		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			_ = os.MkdirAll(sd, 0755)
			os.Create(cfgFile)
			var cfg, _ = ini.LooseLoad(cfgFile)
			_, _ = cfg.Section("general").NewKey("home", dir)
			err = cfg.SaveTo(cfgFile)
		}

		createFolders(dir)
		createTemplates(dir, templ)

		fmt.Println("Congratulations! You have initiated your new recipe book.")

	}
}

func checkSiteNotExist(dir string) bool {
	ok, err := isDirEmpty(dir)
	check(err)

	return ok
}

func createFolders(dir string) {
	for i := 0; i < len(dirs); i++ {
		os.MkdirAll(filepath.Join(dir, dirs[i]), 0755)
	}
}

func createTemplates(dir string, templ map[string]string) {
	for key, value := range templ {
		initTemplate(key, value, dir)
	}
}

func isDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdir(1)

	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func initTemplate(filename string, filecontent string, dir string) {
	err := ioutil.WriteFile(filename, []byte(filecontent), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

var dirs = []string{
	filepath.Join("recipes", "images"),
	filepath.Join("recipes", "templates"),
	"rendered",
	filepath.Join("books", "templates"),
}
