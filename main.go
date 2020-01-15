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
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/go-ini/ini"
)

type Recipe struct {
	Title       string            `yaml:"title"`
	Preptime    int               `yaml:"preptime"`
	Cooktime    int               `yaml:"cooktime"`
	Image       string            `yaml:"image"`
	Origin      string            `yaml:"origin"`
	Tags        []string          `yaml:"tags"`
	Ingredients map[string]string `yaml:"ingredients"`
	Body        string            `yaml:"body"`
}

func main() {

	// central ini file which sets the workdir -> makes it possible to run exec from everywhere

	dir := setWorkDir()

	// load all the recipes

	recipes := loadRecipes(dir)

	fmt.Println(recipes)

	// Define all of the commands

	createCommand := flag.NewFlagSet("create", flag.ExitOnError)
	titleCreateFlag := createCommand.String("title", "DEFAULT TITLE", "Give a title to your recipe.")

	searchCommand := flag.NewFlagSet("search", flag.ExitOnError)
	verboseSearchFlag := searchCommand.Bool("v", false, "Set the output verbosity. Default is false.")
	textSearchFlag := searchCommand.String("text", "", "Search text. Default is empty.")
	tagSearchFlag := searchCommand.String("tag", "", "Search for entries with a specific tag. Default is empty.")
	ingredientSearchFlag := searchCommand.String("ingredient", "", "Search for entries with a specific ingredient. Default is empty.")

	// Input Switch

	if len(os.Args) == 1 {

		fmt.Println("Please provide arguments. Type 'enchiridion help' for more information.")

	} else {

		// define command switch

		switch os.Args[1] {
		case "create":
			createCommand.Parse(os.Args[2:])
		// case "render":
		// 	renderCommand.Parse(os.Args[2:])
		case "search":
			searchCommand.Parse(os.Args[2:])
		// case "tag":
		// 	tagCommand.Parse(os.Args[2:])
		// case "stat":
		// 	statistics(diary.wd, tStr)
		// case "help":
		// 	printHelp()
		default:
			fmt.Printf("%q is not valid command.\n", os.Args[1])
			fmt.Println("Use the help command 'enchiridion help' for help.")
			os.Exit(2)
		}
	}

	// Parse commands

	if createCommand.Parsed() {
		createEntry(dir, *titleCreateFlag)
	}

	if searchCommand.Parsed() {
		searchEntry(dir, *verboseSearchFlag, *textSearchFlag, *tagSearchFlag, *ingredientSearchFlag)
	}
}

func setWorkDir() string {
	usr, err := user.Current()
	check(err)

	sd := filepath.Join(usr.HomeDir, ".enchiridion")
	cfgFile := filepath.Join(sd, "config.ini")

	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		_ = os.MkdirAll(sd, 0755)
		os.Create(cfgFile)
		var cfg, _ = ini.LooseLoad(cfgFile)
		_, _ = cfg.Section("general").NewKey("home", "")
		err = cfg.SaveTo(cfgFile)
	}

	var cfg, _ = ini.LooseLoad(cfgFile)

	if len(cfg.Section("general").Key("home").String()) == 0 {
		fmt.Printf("Home directory not set, please add to config.ini. Config file is located here: %s \n", sd)
		os.Exit(2)
	}

	return filepath.Clean(cfg.Section("general").Key("home").String())
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
