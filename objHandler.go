// objHandler.go

// Source file auto-generated on Sun, 06 Oct 2019 20:50:44 using Gotk3ObjHandler v1.3.8 ©2018-19 H.F.M
/*
	Copyright ©2019 H.F.M - Functions Library Manager
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"log"
	"os"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func TreeSelectionVendorChanged() {
	// if _, iter, ok := mainObjects.TreeSelectionVendor.GetSelected(); ok {
	// 	if path, err := tvsTree.TreeStore.GetPath(iter); err == nil {
	// 		fmt.Println(path.String())
	// 	}
	// }
}

// // TreeViewFoundButtonPressEvent: Popup menu integration
// func TreeViewFoundButtonPressEvent(tw *gtk.TreeView, event *gdk.Event) bool {
// 	// popup whether at least an element is selected
// 	if tvsTreeSearch.Selection.CountSelectedRows() > 0 {
// 		popupMenu.CheckRMBFromTreeView(tw, event)
// 	}
// 	return false // Propagate event
// }

// TreeViewFoundButtonPressEvent: Popup menu integration
func TreeViewFoundButtonPressEvent(tw *gtk.TreeView, event *gdk.Event) bool {
	// popup whether at least an element is selected
	if tvsTreeSearch.Selection.CountSelectedRows() > 0 {
		popupMenu.CheckRMBFromTreeView(tw, event)
	}
	return false // Propagate event
}

// CallbackTooltipFunc: display tooltip with func/struct informations
func CallbackTooltipFunc(iter *gtk.TreeIter, path *gtk.TreePath, column *gtk.TreeViewColumn, tooltip *gtk.Tooltip) bool {
	index := tvsTreeSearch.GetColValue(iter, 5).(int)
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

// TreeViewFoundRowActivated:
func TreeViewFoundRowActivated(tv *gtk.TreeView, path *gtk.TreePath, col *gtk.TreeViewColumn) {
	var err error
	var iter *gtk.TreeIter
	if iter, err = tvsTreeSearch.Model.GetIter(path); err == nil {

		tvsTreeSearch.GetRowNbIter(iter)

		popupSourceView(tvsTreeSearch.GetColValue(iter, 5).(int))
	}
}

// ButtonLoadProjClicked:
func ButtonLoadProjClicked() {
	FileChooserSelectDirFileSet(mainObjects.FileChooserSelectDir)
}

// EntrySearchForChanged:
func EntrySearchForChanged(e *gtk.Entry) {
	var err error
	srch := GetEntryText(e)
	if len(srch) >= mainOptions.SearchCharMinLen {
		err = displayTreeStore(findInLibs(srch, mainObjects.SpinButtonScoreThreshold.GetValueAsInt()))
		DlgErr("EntrySearchForChanged", err)
	}
}

// CheckBoxGenericToggled:
func CheckBoxGenericToggled() {
	EntrySearchForChanged(mainObjects.EntrySearchFor)
}

// ButtonCreatVendorClicked:
func ButtonCreatVendorClicked() {
	var err error
	var checked []string
	if tvsTreeVendor != nil {
		if tvsTreeVendor.CountRows() > 0 {
			if checked, _, err = tvsTreeVendor.GetTree(0, 2); err == nil {
				err = buildVendorDir(checked)
			}
		}
	}
	DlgErr("ButtonCreatVendorClicked", err)
}

// ButtonRefreshLibraryDataClicked:
func ButtonRefreshLibraryDataClicked() {

	mainOptions.SourceLibs = retrieveTreeView(tvsLibInc)
	mainOptions.SubDirToSkip = retrieveTreeView(tvsLibExc)

	initSourceDirectories()
}

// FileChooserSelectDirFileSet:
func FileChooserSelectDirFileSet(fcb *gtk.FileChooserButton) {
	var fi os.FileInfo
	var err error

	if fi, err = os.Stat(fcb.GetFilename()); !os.IsNotExist(err) && fi.IsDir() {
		mainOptions.LastProjFilename = fcb.GetFilename()
		buildVendorList(mainOptions.LastProjFilename)
	}

	if err != nil {
		log.Fatal(err)
	}
}

// Display AboutBox
func ProjEvenboxButtonReleaseEvent() {
	mainOptions.About.Show()
}

/*
 * SourceView
 */

// ButtonSourceOkClicked:
func ButtonSourceOkClicked() {

	mainOptions.SourceWinWidth, mainOptions.SourceWinHeight = mainObjects.WindowSource.GetSize()
	mainOptions.SourceWinPosX, mainOptions.SourceWinPosY = mainObjects.WindowSource.GetPosition()
	mainOptions.PanedWidth = mainOptions.SourceWinWidth - mainObjects.PanedSource.GetPosition()
	mainObjects.WindowSource.Hide()
}

func WindowSourceCheckResize(w *gtk.Window) {

	if mainObjects.SourceToggleButtonMapWidth.GetActive() {
		mainOptions.SourceWinWidth, mainOptions.SourceWinHeight = mainObjects.WindowSource.GetSize()
		mainObjects.PanedSource.SetPosition(mainOptions.SourceWinWidth - mainOptions.PanedWidth)
	}
}

func SourceToggleButtonMapWidthToggled() {
	mainOptions.PanedWidth = mainOptions.SourceWinWidth - mainObjects.PanedSource.GetPosition()
}

func SourceToggleButtonWrapToggled(t *gtk.ToggleButton) {

	if svs != nil {
		if t.GetActive() {
			svs.View.SetWrapMode(gtk.WRAP_WORD_CHAR)
		} else {
			svs.View.SetWrapMode(gtk.WRAP_NONE)
		}
	}
}

/*
 * stack
 */

func ButtonStackLibrarySelectionClicked() {
	mainObjects.Stack.SetVisibleChildName("PageLibSel")
}

func ButtonStackSearchClicked() {
	mainObjects.Stack.SetVisibleChildName("PageSearch")
}

func ButtonStackVendoringClicked() {
	mainObjects.Stack.SetVisibleChildName("PageVendoring")
}
