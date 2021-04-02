// main.go

/*
	Source file auto-generated on Sat, 27 Mar 2021 10:13:00 using Gotk3 Objects Handler v1.6.8 ©2018-20 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2019-21 H.F.M - Functions & Library Manager v1.1.4 github.com/hfmrow/go-func-lib-mgr
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

/*
	This software use also:

	- Go library that provides fuzzy string matching,
	  under the MIT License: https://github.com/sahilm/fuzzy/blob/master/LICENSE
*/

import (
	"fmt"
	"log"
	"path/filepath"
)

// main: And at the beginning ... this part is not modified on update.
// Build options informations:
// devMode: is used in some functions to control the behavior of the program
// When software is ready to be published, this flag must be set at "false"
// that means:
// - options file will be stored in $HOME/.config/[Creat]/[softwareName],
// - translate function if used, will no more auto-update "sts" map sentences,
// - all built-in assets will be used instead of the files themselves.
//   Be aware to update assets via "Goh" and translations with "Got" before all.
func main() {

	devMode = true
	absoluteRealPath, optFilename = getAbsRealPath()

	/* Logger init. */
	Logger = Log2FileStructNew(optFilename, devMode)
	defer Logger.CloseLogger()

	// Initialization of assets according to the chosen mode (devMode).
	// you can set this flag to your liking without reference to devMode.
	assetsDeclarationsUseEmbedded(!devMode)

	// Create temp directory .. or not
	doTempDir = false

	/* Init & read options file */
	mainOptions = new(MainOpt) // Assignate options' structure.
	mainOptions.Read()         // Read values from options' file if exists.

	/* Init gtk display */
	mainStartGtk(fmt.Sprintf("%s %s  %s %s %s",
		Name,
		Vers,
		"©"+YearCreat,
		Creat,
		LicenseAbrv),
		mainOptions.MainWinWidth,
		mainOptions.MainWinHeight, true)
}

/*************************/
/* YOUR CODE START HERE */
/***********************/
func mainApplication() {

	var err error

	/* Init AboutBox */
	mainOptions.About.InitFillInfos(
		mainObjects.MainWindow,
		"About "+Name,
		Name,
		Vers,
		Creat,
		YearCreat,
		LicenseAbrv,
		LicenseShort,
		Repository,
		Descr,
		"",
		tickIcon48)

	/* Translate init. */
	translate = MainTranslateNew(filepath.Join(absoluteRealPath, mainOptions.LanguageFilename), devMode)

	/* Init Statusbar	*/
	statusbar = StatusBarStructureNew(mainObjects.Statusbar, []string{sts["entries"], sts["found"], sts["fileSet"], sts["status"]})

	/* init Clipboard */
	clipboard, err = ClipboardNew()
	Logger.Log(err, "mainApplication/ClipboardNew")

	/* init TreeViews */
	treeViewInit()

	/* Init Popup */
	initPopupTreeView()

	popupLibInc = initPopupLibraryTreeView(tvsLibInc, &mainOptions.SourceLibs)
	popupLibExc = initPopupLibraryTreeView(tvsLibExc, &mainOptions.SubDirToSkip)

	err = fillTreeView(tvsLibInc, &mainOptions.SourceLibs, nil, false)
	Logger.Log(err, "mainApplication/fillTreeView")
	err = fillTreeView(tvsLibExc, &mainOptions.SubDirToSkip, nil, false)
	Logger.Log(err, "mainApplication/fillTreeView")

	/* Init Spinbutton */
	_, err = SpinScaleSetNew(mainObjects.SpinButtonScoreThreshold,
		mainOptions.MinScoreThreshold,
		mainOptions.MaxScoreThreshold,
		mainOptions.ScoreThreshold, 1, nil)
	Logger.Log(err, "mainApplication/SpinScaleSetNew")

	// initSourceDirectories()
	mainObjects.EntrySearchFor.GrabFocus() // Set focus to search entry

}

/*************************************\
/* Executed just before closing all. */
/************************************/
func onShutdown() bool {
	var err error
	// Update mainOptions with GtkObjects and save it
	if err = mainOptions.Write(); err == nil {
		// What you want to execute before closing the app.
		// Return:
		// true for exit applicaton
		// false does not exit application
	}
	if err != nil {
		log.Fatalf("Unexpected error on exit: %s", err.Error())
	}
	return true
}
