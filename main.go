// main.go

/*
	Source file auto-generated on Thu, 05 Nov 2020 07:28:28 using Gotk3ObjHandler v1.6.5 ©2018-20 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2019-20 H.F.M - Functions & Library Manager v1.0 github.com/hfmrow/funcLibraryMgr
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

func main() {

	/* Build options */
	// devMode: is used in some functions to control the behavior of the program
	// When software is ready to be published, this flag must be set at "false"
	// that mean:
	// - options file will be stored in $HOME/.config/[Creat]/[softwareName],
	// - translate function if used, will no more auto-update "sts" map sentences,
	// - all built-in assets will be used instead of the files themselves.
	//   Be aware to update assets via "Goh" and translations with "Got" before all.
	devMode = true
	absoluteRealPath, optFilename = getAbsRealPath()

	// Initialization of assets according to the chosen mode (devMode).
	// you can set this flag to your liking without reference to devMode.
	assetsDeclarationsUseEmbedded(!devMode)

	// Create temp directory .. or not
	doTempDir = false

	/* Init & read options file */
	mainOptions = new(MainOpt) // Assignate options' structure.
	mainOptions.Init()         // Init with default values.
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

	/* init TreeViews */
	treeViewInit()

	/* Init Popup */
	initPopupTreeView()

	popupLibInc = initPopupLibraryTreeView(tvsLibInc, &mainOptions.SourceLibs)
	popupLibExc = initPopupLibraryTreeView(tvsLibExc, &mainOptions.SubDirToSkip)

	fillTreeView(tvsLibInc, &mainOptions.SourceLibs, nil, false)
	fillTreeView(tvsLibExc, &mainOptions.SubDirToSkip, nil, false)

	/* Init Spinbutton */
	SpinbuttonSetValues(mainObjects.SpinButtonScoreThreshold,
		mainOptions.MinScoreThreshold,
		mainOptions.MaxScoreThreshold,
		mainOptions.ScoreThreshold)

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
