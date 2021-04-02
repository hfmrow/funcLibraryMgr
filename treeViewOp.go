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
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/gtk"

	gitw "github.com/hfmrow/gtk3Import/treeview"
)

var HeaderHeight int

func getHeaderHeight(tv *gtk.TreeView) (height int) {
	for gtk.EventsPending() {
		gtk.MainIteration() // Wait for pending events (until the widget is redrawn)
	}
	// Getting header height
	backupVisibleHeader := tv.GetHeadersVisible()
	tv.SetHeadersVisible(true)
	withHeader, _ := tv.GetPreferredHeight()
	tv.SetHeadersVisible(false)
	withoutHeader, _ := tv.GetPreferredHeight()
	tv.SetHeadersVisible(backupVisibleHeader)
	return withHeader - withoutHeader
}

// DisplayTooltipFunc: display tooltip with func/struct informations
func DisplayTooltipFunc(iter *gtk.TreeIter, path *gtk.TreePath, column *gtk.TreeViewColumn, tooltip *gtk.Tooltip) bool {

	index := tvsTreeSearch.GetColValue(iter, mapListStoreColumns["idx"]).(int)

	descr, ok := declIdexes.GetDescr(index)
	if !ok {
		log.Printf("TreeViewFoundQueryTooltip: Unable to get description\n")
		return false
	}
	comment := strings.TrimSpace(descr.Comment)
	if len(comment) > 2 {
		tooltip.SetText(comment)
		tvsTreeSearch.TreeView.SetTooltipRow(tooltip, path)
		return true // ok to display
	}
	return false
}

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

		/********************************************************************************/
		// TODO tooltip freeze For test, the signal callback function is deported here. /
		/******************************************************************************/
		tvsTreeSearch.TreeView.Connect("query-tooltip",
			// treeViewQueryTooltip: function to display tooltip according to rows currently hovered
			// Note: return TRUE if tooltip should be shown right now, FALSE otherwise
			func(tv *gtk.TreeView, x, y int, KeyboardMode bool, tooltip *gtk.Tooltip) bool {

				// if tvs.CountRows() > 0 {
				// we need to substract header height to "y" position to get the correct path.
				if path, column, _, _, isBlank := tv.IsBlankAtPos(x, y-HeaderHeight); !isBlank {
					if iter, err := tvsTreeSearch.Model.GetIter(path); err == nil {

						return DisplayTooltipFunc(iter, path, column, tooltip)
					} else {

						log.Printf("treeViewQueryTooltip:GetIter: %s\n", err.Error())
					}
				}
				// }
				return false
			})

		// Set function to display tooltip according to rows currently hovered
		// tvsTreeSearch.DisplayTooltipFunc = DisplayTooltipFunc

		/*******************************************/

		err = tvsTreeSearch.StoreSetup(new(gtk.TreeStore))
	}
	if err != nil {
		DlgErr("mainApplication:TreeViewStructureNew(listStoreColumns)", err)
		return
	}
	// Get header height for tooltip position synchronization
	HeaderHeight = getHeaderHeight(tvsTreeSearch.TreeView)

	/* Init ListStore Library include */
	if tvsLibInc, err = TreeViewStructureNew(mainObjects.TreeViewInclude, true, false); err == nil {
		tvsLibInc.AddColumns(mainOptions.libInclude, false, false, true, true, true, true)

		tvsLibInc.Columns[includeMap["chk"]].Editable = true

		tvsLibInc.CallbackOnSetColValue = func(iter *gtk.TreeIter, col int, value interface{}) {
			if col == includeMap["chk"] {
				// Refresh 'search function' tab
				ButtonRefreshLibraryDataClicked()
				EntrySearchForChanged(mainObjects.EntrySearchFor)
			}
		}

		if err = tvsLibInc.StoreSetup(new(gtk.ListStore)); err == nil {

			/* Init Drag and drop */
			dndLibInc = DragNDropNew(mainObjects.TreeViewInclude, nil,
				func() {
					var (
						err     error
						tmpLibs []libs
					)
					// Callaback function on Drag and drop operations
					for _, lib := range *dndLibInc.FilesList {
						tmpLibs = append(tmpLibs, libs{Path: lib, Active: true})
					}
					err = fillTreeView(tvsLibInc, &tmpLibs, &mainOptions.SourceLibs, true)
					Logger.Log(err, "TreeViewStructureNew/tvsLibInc/fillTreeView")
				})
		}
	}
	if err != nil {
		DlgErr("mainApplication:TreeViewStructureNew(libInclude)", err)
		return
	}

	/* Init ListStore Library include */
	if tvsLibExc, err = TreeViewStructureNew(mainObjects.TreeViewExclude, true, false); err == nil {
		tvsLibExc.AddColumns(mainOptions.libInclude, false, false, true, true, true, true)

		tvsLibExc.Columns[includeMap["chk"]].Editable = true

		tvsLibExc.CallbackOnSetColValue = func(iter *gtk.TreeIter, col int, value interface{}) {
			if col == includeMap["chk"] {
				// Refresh 'search function' tab
				ButtonRefreshLibraryDataClicked()
				EntrySearchForChanged(mainObjects.EntrySearchFor)
			}
		}

		if err = tvsLibExc.StoreSetup(new(gtk.ListStore)); err == nil {

			/* Init Drag and drop */
			dndLibExc = DragNDropNew(mainObjects.TreeViewExclude, nil,
				func() {
					var (
						err     error
						tmpLibs []libs
					)
					// Callaback function on Drag and drop operations
					for _, lib := range *dndLibInc.FilesList {
						tmpLibs = append(tmpLibs, libs{Path: lib, Active: true})
					}
					err = fillTreeView(tvsLibExc, &tmpLibs, &mainOptions.SubDirToSkip, true)
					Logger.Log(err, "TreeViewStructureNew/tvsLibExc/fillTreeView")
				})
		}
	}
	if err != nil {
		DlgErr("mainApplication:TreeViewStructureNew(libExclude)", err)
	}
}

// fillTreeView:
func fillTreeView(tvs *gitw.TreeViewStructure, list, existList *[]libs, makeRel bool) (err error) {

	libRootDir := filepath.Join(os.Getenv("GOPATH"), "src")

	for _, value := range *list {

		if makeRel {
			if fi, err := os.Stat(value.Path); err == nil {
				if fi.IsDir() {

					value.Path, _ = filepath.Rel(libRootDir, value.Path)

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

		name := filepath.Base(value.Path)
		path := filepath.Dir(value.Path)

		if _, err = tvs.AddRow(nil, value.Active, name, path); err == nil {

			mainOptions.SourceLibs = retrieveTreeView(tvsLibInc)
			mainOptions.SubDirToSkip = retrieveTreeView(tvsLibExc)
			initSourceDirectories()
		}
	}
	return
}

// retrieveTreeView:
func retrieveTreeView(tvs *gitw.TreeViewStructure) (out []libs) {
	if rows, err := tvs.StoreToIfaceSlice(); err == nil {
		for _, row := range rows {
			out = append(out, libs{
				Path: filepath.Join(
					row[includeMap["path"]].(string),
					row[includeMap["name"]].(string)),
				Active: row[includeMap["chk"]].(bool)})
		}
	} else {
		DlgErr("retrieveTreeView:", err)
	}
	return
}
