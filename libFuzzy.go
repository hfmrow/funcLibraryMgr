// libFuzzy.go

// Source file auto-generated on Sun, 06 Oct 2019 23:05:32 using Gotk3ObjHandler v1.3.8 ©2018-19 H.F.M
/*
	Copyright ©2019 H.F.M - Functions Library Manager
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php

	This software use also:
	- Go library that provides fuzzy string matching,
	  under the MIT License: https://github.com/sahilm/fuzzy/blob/master/LICENSE
*/

package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sahilm/fuzzy"
)

type toDispTreeStore struct {
	Name     string
	SortName string
	NameDisp string
	Type     string
	Exported bool
	Path     string
	File     string
	Score    int
	Idx      int
	Methods  []toDispTreeStore
}

// findInLibs: Find similar expression
func findInLibs(toFind string, scoreThreshold int) (out []toDispTreeStore) {
	// Find using Fuzzy
	if results := fuzzy.FindFrom(toFind, desc.Desc); len(results) > 0 {
		var name, path, blueCol, redCol string
		// Select color for Exported or not definition
		var getColor = func(exported bool) (blueCol, redCol string) {
			if exported {
				blueCol = `"#2222AA"`
				redCol = `"#AA2222"`
			} else {
				// Using alpha chanel for unexported items
				blueCol = `"#2222AABB"`
				redCol = `"#AA2222BB"`
			}
			return
		}
		// Sorting results (case insensitive)
		sort.SliceStable(results, func(i, j int) bool {
			return strings.ToUpper(results[i].Str) < strings.ToUpper(results[j].Str)
		})
		for _, found := range results {
			if scoreThreshold <= found.Score {
				disp := false
				if desc.Desc[found.Index].Exported && mainObjects.CheckBoxIncludeExported.GetActive() ||
					!mainObjects.CheckBoxIncludeExported.GetActive() {
					disp = true
				}
				if disp && ((mainObjects.CheckBoxIncludeFunctions.GetActive() && desc.Desc[found.Index].Type == "func") ||
					(mainObjects.CheckBoxIncludeStructures.GetActive() && desc.Desc[found.Index].Type == "struct")) {

					blueCol, redCol = getColor(desc.Desc[found.Index].Exported)
					if mainOptions.MarkResult {
						for i := 0; i < len(found.Str); i++ {

							if contains(i, found.MatchedIndexes) {
								name += `<span foreground=` + redCol + `><b>` + string(found.Str[i]) + `</b></span>`
							} else {
								name += `<span foreground=` + blueCol + `><b>` + string(found.Str[i]) + `</b></span>`
							}
						}
					} else {
						name = `<span foreground=` + blueCol + `><b>` + found.Str + `</b></span>`
					}

					if mainObjects.CheckBoxAddShortcuts.GetActive() {
						path = fmt.Sprintf("%s \"%s\"", desc.Desc[found.Index].Shortcut, filepath.Dir(desc.Desc[found.Index].File))
					} else {
						path = fmt.Sprintf("\"%s\"", filepath.Dir(desc.Desc[found.Index].File))
					}
					// Try to get methods if exists
					var methods []toDispTreeStore
					if desc.Desc[found.Index].Type == "struct" && len(desc.Desc[found.Index].Methods) > 0 {
						for _, method := range desc.Desc[found.Index].Methods {
							// Check for exported or not depending on choosen option
							if method.Exported && mainObjects.CheckBoxIncludeExported.GetActive() ||
								!mainObjects.CheckBoxIncludeExported.GetActive() {

								blueCol, _ = getColor(method.Exported)
								methods = append(methods,
									toDispTreeStore{
										Name:     `<span foreground=` + blueCol + `><b>` + method.Name + `</b></span>`,
										SortName: method.Name,
										Type:     method.Type,
										Path:     method.File,
										Exported: method.Exported,
										Idx:      method.Idx})
							}
						}
						// Sorting methods (case insensitive)
						sort.SliceStable(methods, func(i, j int) bool {
							return strings.ToUpper(methods[i].SortName) < strings.ToUpper(methods[j].SortName)
						})
					}
					out = append(out,
						toDispTreeStore{
							Name:     name,
							Type:     desc.Desc[found.Index].Type,
							Exported: desc.Desc[found.Index].Exported,
							Path:     path,
							Score:    found.Score,
							Idx:      desc.Desc[found.Index].Idx,
							Methods:  methods})
				}
				name = ""
			}
		}
	}
	return
}

// contains: used with fuzzy search lib.
func contains(needle int, haystack []int) bool {
	for _, i := range haystack {
		if needle == i {
			return true
		}
	}
	return false
}
