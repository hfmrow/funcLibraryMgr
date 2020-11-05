// treeViewOp.go

/*
	Source file auto-generated on Wed, 04 Nov 2020 12:12:12 using Gotk3ObjHandler v1.6.5 ©2018-20 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2019-20 H.F.M - Functions & Library Manager v0.8 github.com/...
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gotk3/gotk3/gtk"

	gitw "github.com/hfmrow/gtk3Import/treeview"
)

func treeViewInit() {

	var err error
	/* Init Treestore Functions */
	if tvsTreeSearch, err = TreeViewStructureNew(mainObjects.TreeViewFound, false, false); err == nil {
		tvsTreeSearch.AddColumns(mainOptions.listStoreColumns, false, true, true, true, true, true)

		tvsTreeSearch.Columns[3].Editable = true

		// Define selection changed function .
		tvsTreeSearch.SelectionChangedFunc = func() {
			if tvsTreeSearch.CountRows() > 0 {
				updateStatusBar()
			}
		}
		tvsTreeSearch.CallbackTooltipFunc = CallbackTooltipFunc // Set function to display tooltip according to rows currently hovered
		err = tvsTreeSearch.StoreSetup(new(gtk.TreeStore))
	}

	DlgErr("mainApplication:TreeViewStructureNew(listStoreColumns)", err)

	/* Init ListStore Library include */
	if tvsLibInc, err = TreeViewStructureNew(mainObjects.TreeViewInclude, false, false); err == nil {
		tvsLibInc.AddColumns(mainOptions.libInclude, false, true, true, true, true, true)

		err = tvsLibInc.StoreSetup(new(gtk.ListStore))

		/* Init Drag and drop */
		dndLibInc = DragNDropNew(mainObjects.TreeViewInclude, nil,
			func() {
				// Callaback function on Drag and drop operations
				fillTreeView(tvsLibInc, dndLibInc.FilesList, &mainOptions.SourceLibs, true)
			})
	}

	DlgErr("mainApplication:TreeViewStructureNew(libInclude)", err)

	/* Init ListStore Library include */
	if tvsLibExc, err = TreeViewStructureNew(mainObjects.TreeViewExclude, false, false); err == nil {
		tvsLibExc.AddColumns(mainOptions.libInclude, false, true, true, true, true, true)

		err = tvsLibExc.StoreSetup(new(gtk.ListStore))

		/* Init Drag and drop */
		dndLibExc = DragNDropNew(mainObjects.TreeViewExclude, nil,
			func() {
				// Callaback function on Drag and drop operations
				fillTreeView(tvsLibExc, dndLibExc.FilesList, &mainOptions.SubDirToSkip, true)
			})
	}

	DlgErr("mainApplication:TreeViewStructureNew(libExclude)", err)

}

// fillTreeView:
func fillTreeView(tvs *gitw.TreeViewStructure, list, existList *[]string, makeRel bool) {

	var modified bool

	libRootDir := filepath.Join(os.Getenv("GOPATH"), "src")

	for _, value := range *list {

		if makeRel {
			if fi, err := os.Stat(value); err == nil {

				if fi.IsDir() {
					value, _ = filepath.Rel(libRootDir, value)

					if isExistSlice(*existList, value) {
						continue
					}
				} else {
					DlgErr("File type mistake", fmt.Errorf("Only Directory are allowed"))
					continue
				}
			} else {
				DlgErr("fillTreeView:", err)
				break
			}
		}

		name := filepath.Base(value)
		path := filepath.Dir(value)

		tvs.AddRow(nil, name, path)
		modified = true
	}

	if modified {
		mainOptions.SourceLibs = retrieveTreeView(tvsLibInc)
		mainOptions.SubDirToSkip = retrieveTreeView(tvsLibExc)
		initSourceDirectories()
	}
}

// retrieveTreeView:
func retrieveTreeView(tvs *gitw.TreeViewStructure) (out []string) {
	if values, err := tvs.StoreToStringSlice(); err == nil {
		for i := 1; i < len(values); i++ {
			value := values[i]
			path := filepath.Join(value[1], value[0])
			out = append(out, path)
		}
	} else {
		DlgErr("retrieveTreeView:", err)
	}
	return
}
