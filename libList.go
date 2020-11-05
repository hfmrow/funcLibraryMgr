// libList.go

// Source file auto-generated on Sun, 06 Oct 2019 23:05:32 using Gotk3ObjHandler v1.3.8 ©2018-19 H.F.M
/*
	Copyright ©2019 H.F.M - Functions Library Manager
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	glco "github.com/hfmrow/genLib/crypto"
	glfssf "github.com/hfmrow/genLib/files/scanFileDir"
	gltsgssw "github.com/hfmrow/genLib/tools/goSources/sourceWalker"
)

// declIdx index that point to descriptions, used in treestore to track them.
type declIdx struct {
	Idx     int
	DescPtr *shortDescription
}

// DeclIndexes: Contain all indexes pointing to descriptions.
type DeclIndexes struct {
	Indexes []declIdx
	Count   int
}

// DeclIndxesNew:
func DeclIndxesNew(inDesc shortDescriptions) (di *DeclIndexes) {
	di = new(DeclIndexes)
	di.init(inDesc)
	return
}

// init: build Indexes for all existing declarations.
func (di *DeclIndexes) init(inDesc shortDescriptions) {
	for idx, d := range inDesc {
		di.Indexes = append(di.Indexes, declIdx{
			Idx:     d.Idx,
			DescPtr: &inDesc[idx]})
		for idxM, dM := range d.Methods {
			di.Indexes = append(di.Indexes, declIdx{
				Idx:     dM.Idx,
				DescPtr: &d.Methods[idxM]})
		}
	}
	di.Count = len(di.Indexes)
	return
}

// GetDecl: Get declaration from index.
func (di *DeclIndexes) GetDescr(index int) (outDesc shortDescription, ok bool) {
	for _, d := range di.Indexes {
		if d.Idx == index {
			return *d.DescPtr, true
		}
	}
	return
}

// LibsInfos: Raw version of the analysed libraries.
type LibsInfos struct {
	Shortcut, ImportPath string
	Ast                  *gltsgssw.GoSourceFileStructure
}

// shortDescription: contain description of a declaration (function, method, structure).
type shortDescription struct {
	Name        string
	NameFromSrc string
	Shortcut    string
	File        string
	LineStart   int
	LineEnd     int
	Exported    bool
	Type        string
	Methods     shortDescriptions
	Comment     string
	Idx         int
}

// Description: contain all specified libs that have been analysed.
type Description struct {
	Libs        []string
	LibsMd5     string
	ExludedDirs []string
	changed     bool
	Desc        shortDescriptions
	RootLibs    string
}

// shortDescriptions: define a structure to hold multiples descriptions.
type shortDescriptions []shortDescription

// String: function to complies with fuzzy search into structure.
func (e shortDescriptions) String(i int) string {
	if len(e[i].NameFromSrc) > 0 {
		return e[i].NameFromSrc
	}
	return e[i].Name
}

// Len: function to complies with fuzzy search into structure.
func (e shortDescriptions) Len() int {
	return len(e)
}

// Read: Descriptions from file.
func (stru *Description) Read(filename string) (err error) {
	var textFileBytes []byte
	if textFileBytes, err = ioutil.ReadFile(filename); err == nil {
		if err = json.Unmarshal(textFileBytes, &stru); err == nil {
			if md5, _ := stru.getMd5Libs(); md5 != stru.LibsMd5 {
				log.Printf("Some files in the library path have been modified\nBuilding a new AST data file ...\n")
				stru.changed = true
			}
		}
	}
	if err != nil {
		log.Printf("Error while reading AST data file: %s\nBuilding a new one ...\n", err.Error())
	}
	return
}

// Write: Descriptions to file.
func (stru *Description) Write(filename string) (err error) {
	var jsonData []byte
	var out bytes.Buffer
	if jsonData, err = json.Marshal(&stru); err == nil {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			err = ioutil.WriteFile(filename, out.Bytes(), os.ModePerm)
		}
	}
	return err
}

// initSourceLibs: Initiate and store the contents of the libraries
// or check if the files have been modified since the last time.
func initSourceLibs(sources, subDirToSkip []string) (err error) {

	// Compute description filename
	descFilename := filepath.Join(filepath.Dir(optFilename), mainOptions.LastDescFilename)

	// Try to read description file
	if err = desc.Read(descFilename); err != nil || len(desc.Desc) == 0 {
		desc.changed = true
	}
	// Compare libraries requested with saved ones.
	if !reflect.DeepEqual(desc.Libs, sources) ||
		!reflect.DeepEqual(desc.ExludedDirs, subDirToSkip) {
		desc.changed = true
	}
	// If one of the two previous test fail, create new description file.
	if desc.changed {
		desc = Description{}
		if len(sources) > 0 {
			if sourcesFromAst, desc.LibsMd5, err = buildLibList(sources, append(subDirToSkip, mainOptions.DefaultExclude...)); err == nil { /*
					if err == io.EOF {
						return
					} else {*/
				var globalIdx int
				for _, sfa := range sourcesFromAst {
					for _, f := range sfa.Ast.Func {
						desc.Desc = append(desc.Desc, shortDescription{
							Name:        f.Ident.Name,
							NameFromSrc: "",
							Shortcut:    sfa.Shortcut,
							File:        f.File,
							LineStart:   f.Content.LineStart,
							LineEnd:     f.Content.LineEnd,
							Exported:    f.Exported,
							Type:        "func",
							Comment:     f.Content.Comment,
							Idx:         globalIdx})
						globalIdx++
					}
					for _, f := range sfa.Ast.Struct {
						var methods shortDescriptions
						for _, m := range f.Methods { // Get methods.
							methods = append(methods, shortDescription{
								Name:        m.Ident.Name,
								NameFromSrc: f.Ident.Name,
								File:        m.File,
								LineStart:   m.Content.LineStart,
								LineEnd:     m.Content.LineEnd,
								Exported:    m.Exported,
								Type:        "method",
								Comment:     m.Content.Comment,
								Idx:         globalIdx})
							globalIdx++
						}
						desc.Desc = append(desc.Desc, shortDescription{
							Name:        f.Ident.Name,
							NameFromSrc: "",
							Shortcut:    sfa.Shortcut,
							File:        f.File,
							LineStart:   f.Content.LineStart,
							LineEnd:     f.Content.LineEnd,
							Exported:    f.Exported,
							Type:        "struct",
							Methods:     methods,
							Comment:     f.Content.Comment,
							Idx:         globalIdx})
						globalIdx++
					}
				}
			} /*else {	// Case where there is an issuewith file, like bad go-formatting
				err = nil
			}*/
			// }
			if err == io.EOF {
				return
			}
			if err == nil {
				desc.Libs = sources
				desc.ExludedDirs = subDirToSkip
				err = desc.Write(descFilename)
			}
		} else {
			err = errors.New("There is no library to explore ...")
		}
	}

	declIdexes = DeclIndxesNew(desc.Desc)
	return
}

// IsDirOrSymlinkDir: File is a directory or a symlinked directory ?
func IsDirOrSymlinkDir(slRoot string, slStat os.FileInfo) (slIsDir bool) {
	var err error
	var fName string
	if slStat.IsDir() {
		return true
	} else if slStat.Mode()&os.ModeSymlink != 0 {
		if fName, err = os.Readlink(filepath.Join(slRoot, slStat.Name())); err == nil {
			if slStat, err = os.Stat(fName); err == nil {
				if slStat.IsDir() {
					return true
				}
			}
		}
	}
	if err != nil {
		log.Printf("Unable to scan: %s\n%s\n", fName, err.Error())
	}
	return
}

// The purpose of this function is to generate a list of
// libraries contained in a specific directory and
// creating shortcut name to access them.
func buildLibList(sources, subDirToSkip []string) (sourcesFromAst []LibsInfos, md5 string, err error) {
	var root string
	var existing []string
	var data []byte

	if len(sources) > 0 {
		for idx := 0; idx < len(sources); idx++ {
			desc.RootLibs = filepath.Join(os.Getenv("GOPATH"), "src")
			root = filepath.Join(desc.RootLibs, strings.TrimSpace(sources[idx]))
			if _, err = os.Stat(root); err == nil {
				rootPath := splitPath(root)
				if err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
					if err == nil {
						if IsDirOrSymlinkDir(filepath.Dir(path), info) { // Is dir.
							for _, toSkip := range subDirToSkip { // Skip unwanted directories.
								if info.Name() == strings.TrimSpace(filepath.Base(toSkip)) {
									return filepath.SkipDir
								}
							}
							var infosFiles []os.FileInfo
							// Scan files inside directory
							if infosFiles, err = glfssf.ScanDirFileInfo(path); err != nil {
								return err
							}
							// Get ast infos for each files.
							gsfs := new(gltsgssw.GoSourceFileStructure)
							for idx, osFile := range infosFiles {
								if !osFile.IsDir() {
									if filename := filepath.Join(path, osFile.Name()); filepath.Ext(filename) == ".go" {
										if idx == 0 {
											if err = gsfs.GoSourceFileStructureSetup(filename); err != nil {
												return err
											}
										} else {
											if err = gsfs.AppendFile(filename); err != nil {
												return err
											}
										}
										// Create md5 (same as "getMd5Libs" function) included here
										// since when collecting informations we walk files too.
										if data, err = ioutil.ReadFile(filename); err == nil {
											md5 += glco.Md5String(string(data))
										}
									}
								}
							}
							// Build shortcut and libs path (import style).
							shortcut := computeShort(rootPath, splitPath(path), existing)
							libInf := LibsInfos{
								Ast:        gsfs,
								Shortcut:   shortcut,
								ImportPath: filepath.Join(removePathBefore(splitPath(path), "src", true)...),
							}
							sourcesFromAst = append(sourcesFromAst, libInf)
							existing = append(existing, libInf.Shortcut)

							return nil
						}
					}
					return err
				}); err != nil { // issue with formatted source
					return
				}
			}
		}
	} else {
		DlgErr(sts["missing"], errors.New(sts["noLibsToScan"]))
	}
	return sourcesFromAst, glco.Md5String(md5), err // generate global md5.
}

// getMd5Libs: Used to control integrity of already saved informations.
func (stru *Description) getMd5Libs() (md5 string, err error) {
	var root string
	var infosFiles []os.FileInfo
	var data []byte

	if len(stru.Libs) > 0 {
		for idx := 0; idx < len(stru.Libs); idx++ {
			root = filepath.Join(os.Getenv("GOPATH"), "src", stru.Libs[idx])
			if _, err = os.Stat(root); err == nil {
				err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
					if err == nil {
						if info.IsDir() { // Is dir
							for _, toSkip := range stru.ExludedDirs { // Skip unwanted directories.
								if info.Name() == toSkip {
									return filepath.SkipDir
								}
							}
							// Scan for "*.go" files inside directory.
							if infosFiles, err = glfssf.ScanDirFileInfo(path); err != nil {
								return err
							}
							for _, osFile := range infosFiles {
								if !osFile.IsDir() && osFile.Mode()&os.ModeSymlink == 0 { // not dir & not symlink ?
									if filename := filepath.Join(path, osFile.Name()); filepath.Ext(filename) == ".go" { // Is *.go file
										if data, err = ioutil.ReadFile(filename); err == nil {
											md5 += glco.Md5String(string(data)) // Concatenate md5 of each files.
										}
									}
								}
							}
							return nil
						}
					}
					return err
				})
			}
		}
	} else {
		err = errors.New("Missing directories to be analysed ...")
	}
	return glco.Md5String(md5), err // generate global md5.
}

// computeShort: Create shortcut name for the library
func computeShort(rootPath, newPath, existing []string) (outShortCut string) {
	unableBuildShortCut := "Error"
	upperChar := regexp.MustCompile(`([[:upper:]])`)
	newPath = removePathBefore(newPath, rootPath[len(rootPath)-1])
	for _, name := range newPath {
		outShortCut += name[:1]
		// Search for uppercase character in name, if it found, it
		// included in short name rather than the last char of the name.
		if found := upperChar.FindAllString(name, 1); len(found) > 0 {
			outShortCut += found[0]
		} else {
			outShortCut += name[len(name)-1:]
		}
		outShortCut = strings.ToLower(outShortCut)
	}
	// Search for duplicate shortcut and compute a new one if found.
	subChar := 0
	name := newPath[len(newPath)-1]
	for _, existingName := range existing {
		if outShortCut == existingName && !(outShortCut == unableBuildShortCut) {
			subChar++
			outShortCut = outShortCut[:len(outShortCut)-1]
			if a, b := len(name)-(1+subChar), len(name)-subChar; a >= 0 && b < len(name) {
				outShortCut += name[a:b]
			} else {
				outShortCut = unableBuildShortCut
			}
		}
	}
	return
}

// splitPath: make a slice from a string path.
func splitPath(path string) (outSlice []string) {
	// remove leading and ending PathSeparator.
	path = strings.Trim(path, string(os.PathSeparator))
	return strings.Split(path, string(os.PathSeparator))
}

// removePathBefore: remove directories before or after the chosen one.
func removePathBefore(path []string, at string, after ...bool) []string {
	var afterMark bool
	if len(after) > 0 {
		afterMark = after[0]
	}
	for idx := len(path) - 1; idx >= 0; idx-- {
		if path[idx] == at {
			if afterMark {
				path = path[idx+1:]
			} else {
				path = path[idx:]
			}
			break
		}
	}
	return path
}
