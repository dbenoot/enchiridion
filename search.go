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
	"strings"
)

func searchEntry(dir string, v bool, text string, tag string, ingredients string) {

	// load all the recipes

	recipes := loadRecipes(dir)

	for _, v := range recipes {

		var il []string

		for j := range v.Ingredients {
			il = append(il, j)
		}

		// search for ingredients
		fmt.Println(strings.Split(ingredients, " "))
		fmt.Println(il)
		if subslice(strings.Split(ingredients, " "), il) == true {
			fmt.Println(v.Title)
		}

	}
}

func subslice(s1 []string, s2 []string) bool {
	if len(s1) > len(s2) {
		return false
	}
	for _, e := range s1 {
		if !contains(s2, e) {
			return false
		}
	}
	return true
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
