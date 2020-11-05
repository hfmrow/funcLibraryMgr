// popup.go

/*
	Source file auto-generated on Fri, 08 Nov 2019 03:45:35 using Gotk3ObjHandler v1.4 ©2018-19 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2019 H.F.M - Functions & Library Manager github.com/...
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"path/filepath"

	"github.com/hfmrow/gotk3_gtksource/source"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	glts "github.com/hfmrow/genLib/tools"

	gimc "github.com/hfmrow/gtk3Import/misc"
	gitw "github.com/hfmrow/gtk3Import/treeview"
)

/*
* Init Popup
 */

func initPopupLibraryTreeView(tvs *gitw.TreeViewStructure, list *[]string) (popMenu *gimc.PopupMenuStruct) {
	popMenu = gimc.PopupMenuStructNew()
	popMenu.WithIcons = true
	popMenu.AddItem("_Remove entry", func() {
		tvs.RemoveRow(tvs.GetSelectedIters()...)
		*list = retrieveTreeView(tvs)

		initSourceDirectories()
	}, crossIcon48)
	popMenu.MenuBuild()
	return
}

// initPopupTreeView:
func initPopupTreeView() {
	popupMenu = gimc.PopupMenuStructNew()
	popupMenu.WithIcons = true
	popupMenu.AddItem("Open _directory", func() { openDir() }, folderOpen20)
	popupMenu.AddItem("Open _file", func() { openFile() }, mimetypeSourceIconGolang48)
	popupMenu.MenuBuild()
}

// openFile:
func openFile(index ...int) {
	var ok bool
	var descr shortDescription
	if len(index) > 0 {
		descr, ok = declIdexes.GetDescr(index[0])
	} else {
		var iters []*gtk.TreeIter
		if iters = tvsTreeSearch.GetSelectedIters(); len(iters) > 0 {
			idx := tvsTreeSearch.GetColValue(iters[0], 5)
			descr, ok = declIdexes.GetDescr(idx.(int))
		}
	}
	if ok {
		open(filepath.Join(desc.RootLibs, descr.File))
	}
}

// openDir:
func openDir(index ...int) {
	var ok bool
	var descr shortDescription
	if len(index) > 0 {
		descr, ok = declIdexes.GetDescr(index[0])
	} else {
		var iters []*gtk.TreeIter
		if iters = tvsTreeSearch.GetSelectedIters(); len(iters) > 0 {
			idx := tvsTreeSearch.GetColValue(iters[0], 5)
			descr, ok = declIdexes.GetDescr(idx.(int))
		}
	}
	if ok {
		open(filepath.Dir(filepath.Join(desc.RootLibs, descr.File)))
	}
}

// open: show file or dir depending on "path".
func open(path string) {

	glib.IdleAdd(func() { // IdleAdd to permit gtk3 working right during goroutine
		// Using goroutine to permit the usage of another thread and freeing the current one.
		go glts.ExecCommand(mainOptions.AppLauncher, path)
	})
}

// popupTextViewPopulateMenu: Append some items to contextmenu of the popup textview
func popupTextViewPopulateMenu(txtView *source.SourceView, popup *gtk.Widget) {
	// Convert gtk.Widget to gtk.Menu object
	pop := &gtk.Menu{gtk.MenuShell{gtk.Container{*popup}}}
	// create new menuitems
	popMenuTextView = gimc.PopupMenuStructNew()
	popMenuTextView.WithIcons = true
	popMenuTextView.AddSeparator()
	popMenuTextView.AddItem("Open _directory", func() { openDir(indexCurrText) }, folderOpen20)
	popMenuTextView.AddItem("Open _file", func() { openFile(indexCurrText) }, mimetypeSourceIconGolang48)
	// Append them to the existing menu
	popMenuTextView.AppendToExistingMenu(pop)
}
