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
		createTemplates(dir)
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

func createTemplates(dir string) {
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
	err := ioutil.WriteFile(filepath.Join(dir, "templates", filename), []byte(filecontent), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

var dirs = []string{
	"templates",
	filepath.Join("recipes", "images"),
	"rendered",
	"books",
}
