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
	"strings"

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

func initPopupLibraryTreeView(tvs *gitw.TreeViewStructure, list *[]libs) (popMenu *gimc.PopupMenuIconStruct) {
	popMenu = PopupMenuIconStructNew()
	popMenu.AddItem("_Remove entry", func() {
		tvs.RemoveRows(tvs.GetSelectedIters()...)
		*list = retrieveTreeView(tvs)

		initSourceDirectories()
	}, popMenu.OPT_ICON, crossIcon48)
	popMenu.MenuBuild()
	return
}

// initPopupTreeView:
func initPopupTreeView() {
	popupMenu = PopupMenuIconStructNew()
	popupMenu.AddItem("_Copy clipboard", func() { toClipboard() }, popupMenu.OPT_ICON, actionsEditCopy)
	popupMenu.AddItem("Open _directory", func() { openDir() }, popupMenu.OPT_ICON, folderOpen20)
	popupMenu.AddItem("Open _file", func() { openFile() }, popupMenu.OPT_ICON, mimetypeSourceIconGolang48)
	popupMenu.MenuBuild()
}

// openFile:
func openFile(index ...int) {

	descr, ok := getDesc(index...)
	if ok {
		open(filepath.Join(desc.RootLibs, descr.File))
	}
}

// toClipboard:
func toClipboard(index ...int) {

	var content string

	descr, ok := getDesc(index...)
	if ok {
		if descr.Type != "method" {
			if mainObjects.CheckBoxAddShortcuts.GetActive() {
				content = strings.Join([]string{descr.Shortcut, `"` + filepath.Dir(descr.File) + `"`}, " ")
			} else {
				content = `"` + filepath.Dir(descr.File) + `"`
			}
			clipboard.SetText(content)
			clipboard.Store()
		}
	}
}

// openDir:
func openDir(index ...int) {

	descr, ok := getDesc(index...)
	if ok {
		open(filepath.Dir(filepath.Join(desc.RootLibs, descr.File)))
	}
}

// getDesc:
func getDesc(index ...int) (descr shortDescription, ok bool) {

	if len(index) > 0 {
		descr, ok = declIdexes.GetDescr(index[0])
	} else {
		var iters []*gtk.TreeIter
		if iters = tvsTreeSearch.GetSelectedIters(); len(iters) > 0 {
			idx := tvsTreeSearch.GetColValue(iters[0], mapListStoreColumns["idx"])
			descr, ok = declIdexes.GetDescr(idx.(int))
		}
	}
	return
}

// open: show file or dir depending on "path".
func open(path string) {

	glib.IdleAdd(func() { // IdleAdd to permit gtk3 working right during goroutine
		// Using goroutine to permit the usage of another thread and freeing the current one.
		go glts.ExecCommand([]string{mainOptions.AppLauncher, path})
	})
}

// popupTextViewPopulateMenu: Append some items to contextmenu of the popup textview
func popupTextViewPopulateMenu(txtView *source.SourceView, popup *gtk.Widget) {
	// Convert gtk.Widget to gtk.Menu object
	pop := &gtk.Menu{gtk.MenuShell{gtk.Container{*popup}}}
	// create new menuitems
	popMenuTextView = PopupMenuIconStructNew()
	popMenuTextView.AddItem("", nil, popupMenu.OPT_SEPARATOR)
	popMenuTextView.AddItem("Open _directory", func() { openDir(indexCurrText) }, popupMenu.OPT_ICON, folderOpen20)
	popMenuTextView.AddItem("Open _file", func() { openFile(indexCurrText) }, popupMenu.OPT_ICON, mimetypeSourceIconGolang48)
	// Append them to the existing menu
	popMenuTextView.AppendToExistingMenu(pop)
}
