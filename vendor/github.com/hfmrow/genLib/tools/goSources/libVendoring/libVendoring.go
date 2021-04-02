// libVendoring.go

/*
	Copyright (c) 2019 H.F.M
	See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	SOFTWARE.

	- The purpose of this library is to build a "vendor" directory in a project
	  containing all the imports necessary for a successful compilation.

	- "*.go" files containing the directive: "// +build ignore" are included in
	  the libraries retrieving process.

	- Voluntarily, the objective is not to worry about the version information
	  of libraries or the implementation of "GO111MODULE". So, these informations
	  are ignored.

	- It allows to exclude certain specific imports. (Useful for large libraries
	  as gotk3 which is always installed when in use and does not need to be
	  added in most cases.)

	- It use a self versionning implementation that only warn on files
	  modification inside used libraries (only "*.go" files are checked, but
	  except the hidden file/dir, the whole directory content is copied during
	  vendoring operation).
*/

package libVendoring

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	glfs "github.com/hfmrow/genLib/files"
	glfsfs "github.com/hfmrow/genLib/files/files-operations"
	glfssf "github.com/hfmrow/genLib/files/scanFileDir"
	glss "github.com/hfmrow/genLib/slices"
)

var name = "MIT License, libVendoringV1.0 Copyright (c) 2019-20 H.F.M, https://github.com/hfmrow"

type LibsVendor struct {
	Author          string
	SourcePathName  string   // Directory to be analysed
	SourceFilenames []string // Files to be anaysed
	Changed         bool

	ImportPaths     []string
	NotExist        [][]string
	UsedFiles       []string // Analysed library files
	UsedFilesMD5    string
	ExludedImports  []string
	ExludedFiles    []string
	IncludeSymlinks bool

	goSrc string
	goPkg string

	fos *glfsfs.FilesOpStruct
}

// LibVendorNew: The purpose of this library is to build a "vendor"
// directory in a project that hold all the necessary imports for a
// successful compilation even if the officials libs have changed.
func LibVendorNew() (lv *LibsVendor, err error) {
	lv = new(LibsVendor)
	lv.init()
	lv.fos, err = glfsfs.FilesOpStructNew()
	return
}

func (lv *LibsVendor) scanDirInfo(path string, ext []string) (fInf []os.FileInfo, err error) {
	var lInf []os.FileInfo
	fInf, _, lInf, err = glfssf.ScanDirFileInfoMask(path, ext)
	if lv.IncludeSymlinks {
		fInf = append(fInf, lInf...)
	}
	return
}

// RunForDir:
func (lv *LibsVendor) RunForDir(path string, skipImports ...[]string) (err error) {
	var skipedImports []string
	if len(skipImports) > 0 {
		skipedImports = skipImports[0]
	}

	lv.SourcePathName = path

	if err = lv.buildSkipImportsList(skipedImports); err == nil {

		var fInf []os.FileInfo
		if fInf, err = lv.scanDirInfo(path, []string{"*.go"}); err == nil {
			for _, fi := range fInf {
				filename := filepath.Join(path, fi.Name())
				if err = lv.getImportsFromSrc(filename); err != nil {
					return
				}
			}
		}
		// Sort results
		sort.SliceStable(lv.ImportPaths, func(i, j int) bool { return lv.ImportPaths[i] < lv.ImportPaths[j] })
		lv.UsedFilesMD5, err = lv.MakeMd5()
	}

	if err != nil {
		err = errors.New(fmt.Sprintf("From: %s, Error: %v\n", "RunForDir", err))
	}
	return
}

// RunForFiles:
func (lv *LibsVendor) RunForFiles(filenames []string, skipImports ...[]string) (err error) {
	var skipedImports []string
	if len(skipImports) > 0 {
		skipedImports = skipImports[0]
	}

	lv.SourceFilenames = filenames

	if err = lv.buildSkipImportsList(skipedImports); err == nil {

		for _, file := range lv.SourceFilenames {
			if err = lv.getImportsFromSrc(file); err != nil {
				return
			}
		}
		// Sort results
		sort.SliceStable(lv.ImportPaths, func(i, j int) bool { return lv.ImportPaths[i] < lv.ImportPaths[j] })
		lv.UsedFilesMD5, err = lv.MakeMd5()
	}

	if err != nil {
		err = errors.New(fmt.Sprintf("From: %s, Error: %v\n", "RunForFiles", err))
	}
	return
}

// Init: init some variables
func (lv *LibsVendor) init() {
	lv.Author = name
	lv.goSrc = filepath.Join(os.Getenv("GOPATH"), "src") // Full path to src directory based on GOPATH environment variable
	lv.goPkg = filepath.Join(os.Getenv("GOROOT"), "src") // Full path to root directory based on GOROOT environment variable - target native libs
	lv.ExludedImports = []string{                        // Some usual excluded directories, this list will be added to user's list
		".git",
		"TEST", /*
			"C",*/
	}
	lv.ExludedFiles = []string{"*.debug"} // Some usual excluded directories, this list will be added to user's list
}

// getSrcFromImports: Retrieve all "*.go" source files contained in an import path.
func (lv *LibsVendor) getSrcFromImports(importPath string) (err error) {
	var fInf []os.FileInfo
	if fInf, err = lv.scanDirInfo(filepath.Join(lv.goSrc, importPath), []string{"*.go"}); err == nil {
		for _, fi := range fInf {
			tmpFilename := filepath.Join(importPath, fi.Name())
			if glss.IsExistSl(lv.UsedFiles, tmpFilename) {
				continue
			}
			lv.UsedFiles = append(lv.UsedFiles, tmpFilename)
			if err = lv.getImportsFromSrc(filepath.Join(lv.goSrc, tmpFilename)); err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return
			}
		}
	}

	if os.IsNotExist(err) {
		lv.NotExist = append(lv.NotExist, []string{"getSrcFromImports", filepath.Join(lv.goSrc, importPath)})
	} else if err != nil {
		err = errors.New(fmt.Sprintf("From: %s, Error: %v\n", "getSrcFromImports", err))
	}
	return
}

// getImportsfromSrc: Retrieve all imported libs from a "*.go" source file.
func (lv *LibsVendor) getImportsFromSrc(filename string) (err error) {
	var astFile *ast.File
	if astFile, err = parser.ParseFile(token.NewFileSet(), filename, nil, parser.ImportsOnly); err == nil {
		for _, val := range astFile.Imports {
			var unQuoted string
			if unQuoted, err = strconv.Unquote(val.Path.Value); err == nil {
				if glss.IsExistSl(lv.ImportPaths, unQuoted) ||
					glss.IsExistSl(lv.ExludedImports, unQuoted) ||
					lv.isNativeLib(unQuoted) {
					continue
				}
				lv.ImportPaths = append(lv.ImportPaths, unQuoted)
				if err = lv.getSrcFromImports(unQuoted); err != nil {
					if os.IsNotExist(err) {
						lv.ImportPaths = lv.ImportPaths[:len(lv.ImportPaths)-1]
						continue
					}
					return
				}
			}
		}
	}

	if os.IsNotExist(err) {
		err = nil
	}
	if err != nil {
		err = errors.New(fmt.Sprintf("From: %s, Error: %v\n", "getImportsFromSrc", err))
	}
	return
}

// isNativeLib: Check if given library is a native lib (Golang out of the box)
func (lv *LibsVendor) isNativeLib(lib string) bool {
	if _, err := os.Stat(filepath.Join(lv.goPkg, lib)); !os.IsNotExist(err) {
		return true
	}
	return false
}

// buildSkipImportsList: Build a list of directories included in
// folders that we does not want to add to the vendor directory.
func (lv *LibsVendor) buildSkipImportsList(skipImports []string) (err error) {
	var tmpPath string
	for _, path := range skipImports {
		var dInf []os.FileInfo
		tmpPath = filepath.Join(lv.goSrc, path)
		if _, dInf, _, err = glfssf.ScanDirFileInfoMask(tmpPath, []string{"*"}); err == nil {
			for _, fi := range dInf {
				importPath := filepath.Join(path, fi.Name())
				if glss.IsExistSl(lv.ExludedImports, fi.Name()) {
					continue
				}
				lv.ExludedImports = append(lv.ExludedImports, importPath)
				if err = lv.buildSkipImportsList([]string{importPath}); err != nil {
					return
				}
			}
		}
	}

	if os.IsNotExist(err) {
		lv.NotExist = append(lv.NotExist, []string{"buildSkipImportsList", tmpPath})
		err = nil
	}
	if err != nil {
		err = errors.New(fmt.Sprintf("From: %s, Error: %v\n", "buildSkipImportsList", err))
	}
	return
}

// LibsToVendor: copy libraries to vendor directory of "SourcePathName"
func (lv *LibsVendor) CopyLibsToVendor() (err error) {
	if _, err = os.Stat(lv.SourcePathName); err == nil {
		vendorDir := filepath.Join(lv.SourcePathName, "vendor")
		if _, err = os.Stat(vendorDir); !os.IsNotExist(err) {
			os.RemoveAll(vendorDir) // Remove "vendor" directory if already exist
		}
		for _, impt := range lv.ImportPaths {
			var fInf []os.FileInfo
			if fInf, err = lv.scanDirInfo(filepath.Join(lv.goSrc, impt),
				[]string{"*.*", "*"}); err == nil {
				for _, fi := range fInf {
					if !glfs.ExtFileMatch(fi.Name(), lv.ExludedFiles) { // skip *.debug files
						src := filepath.Join(lv.goSrc, impt, fi.Name())
						dst := filepath.Join(vendorDir, impt, fi.Name())
						if err = lv.fos.CopyFile(src, dst); err != nil {
							break
						}
					}
				}
			}
		}
	}
	if err != nil {
		err = errors.New(fmt.Sprintf("CopyLibsToVendor: %s\n", err.Error()))
	}
	return
}

// MD5: make global md5 hash for all used files.
func (lv *LibsVendor) MakeMd5() (out string, err error) {
	var file *os.File
	hasher := md5.New()
	for _, filename := range lv.UsedFiles {
		if file, err = os.Open(filepath.Join(lv.goSrc, filename)); err == nil {
			defer file.Close()
			if _, err = io.Copy(hasher, file); err == nil {
				out += hex.EncodeToString(hasher.Sum(nil))
			}
		}
	}
	if err == nil {
		if _, err = hasher.Write([]byte(out)); err == nil {
			out = hex.EncodeToString(hasher.Sum(nil))
		}
	}
	return
}

// Read vendor informations structure from file
func (lv *LibsVendor) Read(filename string) (err error) {
	var textFileBytes []byte
	if textFileBytes, err = ioutil.ReadFile(filename); err == nil {
		err = json.Unmarshal(textFileBytes, &lv)
		var md5 string
		if md5, err = lv.MakeMd5(); err != nil || lv.UsedFilesMD5 != md5 {
			lv.Changed = true
			fmt.Println("The previous set of imported libraries has been changed from\nthe ones contained in the current \"vendor\" directory.")
		}
	}
	return
}

// Write vendor informations structure to file
func (lv *LibsVendor) Write(filename string) (err error) {
	var jsonData []byte
	var out bytes.Buffer
	if jsonData, err = json.Marshal(&lv); err == nil {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			err = lv.fos.WriteFile(filename, out.Bytes(), lv.fos.Perms.File)
		}
	}
	return
}
