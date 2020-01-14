// Legalese

package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/go-ini/ini"
)

type Recipe struct {
	title       string
	preptime    int
	cooktime    int
	origin      string
	tags        []string
	ingredients []string
	body        string
}

func main() {

	// central ini file which sets the workdir -> makes it possible to run exec from everywhere

	dir := setWorkDir()

	// load all the recipes

	recipes := loadRecipes(dir)

	fmt.Println(dir)
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
