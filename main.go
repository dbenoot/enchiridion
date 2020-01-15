// Legalese

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
	Title       string
	Preptime    int
	Cooktime    int
	Image       string
	Origin      string
	Tags        []string
	Ingredients []string
	Body        string
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
		// case "search":
		// 	searchCommand.Parse(os.Args[2:])
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
