// gohSignals.go

/*
	Source file auto-generated on Fri, 02 Apr 2021 08:27:58 using Gotk3 Objects Handler v1.7.5 ©2018-21 hfmrow
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2019-21 H.F.M - Functions & Library Manager v1.1.4 github.com/hfmrow/go-func-lib-mgr
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

/********************************************************/
/* This section preserve user modifications on update. */
/* Signals & Property implementations:                */
/* initialise signals used by gtk objects ...        */
/****************************************************/
func signalsPropHandler() {
	mainObjects.ButtonCreatVendor.Connect("clicked", ButtonCreatVendorClicked)
	mainObjects.ButtonExit.Connect("clicked", windowDestroy)
	mainObjects.ButtonLoadProj.Connect("clicked", ButtonLoadProjClicked)
	mainObjects.ButtonRefreshLibraryData.Connect("clicked", ButtonRefreshLibraryDataClicked)
	mainObjects.ButtonSourceOk.Connect("clicked", ButtonSourceOkClicked)
	mainObjects.ButtonStackLibrarySelection.Connect("clicked", ButtonStackLibrarySelectionClicked)
	mainObjects.ButtonStackSearch.Connect("clicked", ButtonStackSearchClicked)
	mainObjects.ButtonStackVendoring.Connect("clicked", ButtonStackVendoringClicked)
	mainObjects.CheckBoxAddShortcuts.Connect("toggled", CheckBoxGenericToggled)
	mainObjects.CheckBoxIncludeExported.Connect("toggled", CheckBoxGenericToggled)
	mainObjects.CheckBoxIncludeFunctions.Connect("toggled", CheckBoxGenericToggled)
	mainObjects.CheckBoxIncludeStructures.Connect("toggled", CheckBoxGenericToggled)
	mainObjects.EntrySearchFor.Connect("changed", EntrySearchForChanged)
	mainObjects.EvenboxTop.Connect("button-release-event", ProjEvenboxButtonReleaseEvent)
	mainObjects.FileChooserSelectDir.Connect("file-set", FileChooserSelectDirFileSet)
	mainObjects.SourceToggleButtonMapWidth.Connect("toggled", SourceToggleButtonMapWidthToggled)
	mainObjects.SourceToggleButtonWrap.Connect("toggled", SourceToggleButtonWrapToggled)
	mainObjects.SpinButtonScoreThreshold.Connect("value-changed", CheckBoxGenericToggled)
	mainObjects.TreeViewExclude.Connect("button-press-event", popupLibExc.CheckRMBFromTreeView)
	mainObjects.TreeViewFound.Connect("row-activated", TreeViewFoundRowActivated)
	mainObjects.TreeViewFound.Connect("button-press-event", TreeViewFoundButtonPressEvent)
	mainObjects.TreeViewInclude.Connect("button-press-event", popupLibInc.CheckRMBFromTreeView)
	mainObjects.WindowSource.Connect("check-resize", WindowSourceCheckResize)
}
