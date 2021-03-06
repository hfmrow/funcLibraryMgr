// libSelection.go

// Source file auto-generated on Wed, 09 Oct 2019 16:17:45 using Gotk3ObjHandler v1.3.8 ©2018-19 H.F.M
/*
	Copyright ©2019 H.F.M - Functions Library Manager
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"os"
	"strings"

	"github.com/gotk3/gotk3/gtk"

	gitw "github.com/hfmrow/gtk3Import/treeview"
)

// buildVendorList:
func buildVendorList(root string) (err error) {

	if libsVendor, err = LibVendorNew(); err == nil {
		libsVendor.IncludeSymlinks = true
		if err = libsVendor.RunForDir(root); err == nil {

			if tvsTreeVendor != nil {
				tvsTreeVendor.ClearAll()
			}
			err = displayTreeViewVendor(libsVendor.ImportPaths)
		}
	}
	Logger.Log(err, "buildVendorList")
	return
}

// displayTreeViewVendor: fill the treeview vith imported libraries
func displayTreeViewVendor(libList []string) (err error) {
	if tvsTreeVendor, err = gitw.TreeViewStructureNew(mainObjects.TreeViewVendor, false, false); err == nil { // Create Structure With his columns
		var iter *gtk.TreeIter

		tvsTreeVendor.AddColumn("Selected", "active", true, false, false, false, false, true)
		tvsTreeVendor.AddColumn("Library", "markup", false, true, false, false, true, true)
		tvsTreeVendor.AddColumn("Path", "markup", false, true, false, false, true, true)

		tvsTreeVendor.StoreSetup(new(gtk.TreeStore)) // Setup structure with desired TreeModel
		tvsTreeVendor.StoreDetach()                  // Free TreeStore from TreeView before fill it. (useful for very large entries)
		defer tvsTreeVendor.StoreAttach()

		var splitted []string
		var iSplitted []interface{}

		for _, path := range libList {

			splitted = strings.Split(path, string(os.PathSeparator))

			iSplitted = tvsTreeVendor.ColValuesStringSliceToIfaceSlice(splitted...)

			if iter, err = tvsTreeVendor.AddTree(includeMap["chk"], includeMap["name"], true, nil, iSplitted...); err == nil {

				err = tvsTreeVendor.SetColValue(iter, includeMap["path"], path)
			}
			if err != nil {
				break
			}
		}
		tvsTreeVendor.StoreAttach()
	}
	return
}

// buildVendorDir: create vendor directory and copy all libs choosen to it
func buildVendorDir(libList []string) (err error) {
	libsVendor.ImportPaths = libList
	return libsVendor.CopyLibsToVendor()
}
