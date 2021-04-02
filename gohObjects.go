// gohObjects.go

/*
	Source file auto-generated on Fri, 02 Apr 2021 08:27:58 using Gotk3 Objects Handler v1.7.5 ©2018-21 hfmrow
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2019-21 H.F.M - Functions & Library Manager v1.1.4 github.com/hfmrow/go-func-lib-mgr
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/hfmrow/gotk3_gtksource/source"
)

/* Control over all used objects from glade. */
var mainObjects *MainControlsObj

/********************************************************/
/* This section preserve user modifications on update. */
/* Main structure Declaration: You may add your own   */
/* declarations (gotk3 objects only) here.           */
/****************************************************/
type MainControlsObj struct {
	ButtonCreatVendor           *gtk.Button
	ButtonExit                  *gtk.Button
	ButtonLoadProj              *gtk.Button
	ButtonRefreshLibraryData    *gtk.Button
	ButtonSourceOk              *gtk.Button
	ButtonStackLibrarySelection *gtk.Button
	ButtonStackSearch           *gtk.Button
	ButtonStackVendoring        *gtk.Button
	CheckBoxAddShortcuts        *gtk.CheckButton
	CheckBoxIncludeExported     *gtk.CheckButton
	CheckBoxIncludeFunctions    *gtk.CheckButton
	CheckBoxIncludeStructures   *gtk.CheckButton
	ComboboxSourceLanguage      *gtk.ComboBoxText
	ComboboxSourceStyle         *gtk.ComboBoxText
	EntrySearchFor              *gtk.Entry
	EvenboxTop                  *gtk.EventBox
	FileChooserSelectDir        *gtk.FileChooserButton
	ImageTop                    *gtk.Image
	mainUiBuilder               *gtk.Builder
	MainWindow                  *gtk.Window
	PanedSource                 *gtk.Paned
	ScrolledWindowExclude       *gtk.ScrolledWindow
	ScrolledWindowInclude       *gtk.ScrolledWindow
	Source                      *source.SourceView
	SourceMap                   *source.SourceMap
	SourceToggleButtonMapWidth  *gtk.ToggleButton
	SourceToggleButtonWrap      *gtk.ToggleButton
	SpinButtonScoreThreshold    *gtk.SpinButton
	Stack                       *gtk.Stack
	Statusbar                   *gtk.Statusbar
	TreeViewExclude             *gtk.TreeView
	TreeViewFound               *gtk.TreeView
	TreeViewInclude             *gtk.TreeView
	TreeViewVendor              *gtk.TreeView
	WindowSource                *gtk.Window
}

/******************************************************************/
/* This section preserve user modification on update.            */
/* GtkObjects initialisation: You may add your own declarations */
/* as you  wish, best way is to add them by grouping  same     */
/* objects names (below first declaration).                   */
/*************************************************************/
func gladeObjParser() {
	mainObjects.ButtonCreatVendor = loadObject("ButtonCreatVendor").(*gtk.Button)
	mainObjects.ButtonExit = loadObject("ButtonExit").(*gtk.Button)
	mainObjects.ButtonLoadProj = loadObject("ButtonLoadProj").(*gtk.Button)
	mainObjects.ButtonRefreshLibraryData = loadObject("ButtonRefreshLibraryData").(*gtk.Button)
	mainObjects.ButtonSourceOk = loadObject("ButtonSourceOk").(*gtk.Button)
	mainObjects.ButtonStackLibrarySelection = loadObject("ButtonStackLibrarySelection").(*gtk.Button)
	mainObjects.ButtonStackSearch = loadObject("ButtonStackSearch").(*gtk.Button)
	mainObjects.ButtonStackVendoring = loadObject("ButtonStackVendoring").(*gtk.Button)
	mainObjects.CheckBoxAddShortcuts = loadObject("CheckBoxAddShortcuts").(*gtk.CheckButton)
	mainObjects.CheckBoxIncludeExported = loadObject("CheckBoxIncludeExported").(*gtk.CheckButton)
	mainObjects.CheckBoxIncludeFunctions = loadObject("CheckBoxIncludeFunctions").(*gtk.CheckButton)
	mainObjects.CheckBoxIncludeStructures = loadObject("CheckBoxIncludeStructures").(*gtk.CheckButton)
	mainObjects.ComboboxSourceLanguage = loadObject("ComboboxSourceLanguage").(*gtk.ComboBoxText)
	mainObjects.ComboboxSourceStyle = loadObject("ComboboxSourceStyle").(*gtk.ComboBoxText)
	mainObjects.EntrySearchFor = loadObject("EntrySearchFor").(*gtk.Entry)
	mainObjects.EvenboxTop = loadObject("EvenboxTop").(*gtk.EventBox)
	mainObjects.FileChooserSelectDir = loadObject("FileChooserSelectDir").(*gtk.FileChooserButton)
	mainObjects.ImageTop = loadObject("ImageTop").(*gtk.Image)
	mainObjects.MainWindow = loadObject("MainWindow").(*gtk.Window)
	mainObjects.PanedSource = loadObject("PanedSource").(*gtk.Paned)
	mainObjects.ScrolledWindowExclude = loadObject("ScrolledWindowExclude").(*gtk.ScrolledWindow)
	mainObjects.ScrolledWindowInclude = loadObject("ScrolledWindowInclude").(*gtk.ScrolledWindow)
	mainObjects.Source = loadObject("Source").(*source.SourceView)
	mainObjects.SourceMap = loadObject("SourceMap").(*source.SourceMap)
	mainObjects.SourceToggleButtonMapWidth = loadObject("SourceToggleButtonMapWidth").(*gtk.ToggleButton)
	mainObjects.SourceToggleButtonWrap = loadObject("SourceToggleButtonWrap").(*gtk.ToggleButton)
	mainObjects.SpinButtonScoreThreshold = loadObject("SpinButtonScoreThreshold").(*gtk.SpinButton)
	mainObjects.Stack = loadObject("Stack").(*gtk.Stack)
	mainObjects.Statusbar = loadObject("Statusbar").(*gtk.Statusbar)
	mainObjects.TreeViewExclude = loadObject("TreeViewExclude").(*gtk.TreeView)
	mainObjects.TreeViewFound = loadObject("TreeViewFound").(*gtk.TreeView)
	mainObjects.TreeViewInclude = loadObject("TreeViewInclude").(*gtk.TreeView)
	mainObjects.TreeViewVendor = loadObject("TreeViewVendor").(*gtk.TreeView)
	mainObjects.WindowSource = loadObject("WindowSource").(*gtk.Window)
}
