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
	"strings"
)

func searchEntry(dir string, v bool, text string, tag string, ingredients string) {

	var result []string

	// load all the recipes

	recipes := loadRecipes(dir)

	for _, v := range recipes {

		// INGREDIENTS
		// get a list of all the ingredients of this particular recipe

		var il []string
		for j := range v.Ingredients {
			il = append(il, j)
		}

		var tl []string
		for _, k := range v.Tags {
			tl = append(tl, k)
		}

		// search for all ingredients in the recipe. Only full matches!

		if subslice(strings.Split(ingredients, " "), il) == true && ubslice(strings.Split(tag, " "), tl) == true {
			result = appendIfMissing(result, v.Title)
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

func appendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}
