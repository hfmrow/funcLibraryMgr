// gohOptions.go

/*
	Source file auto-generated on Thu, 05 Nov 2020 07:28:28 using Gotk3ObjHandler v1.6.5 ©2018-20 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2019-20 H.F.M - Functions & Library Manager v1.0 github.com/hfmrow/funcLibraryMgr
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gotk3/gotk3/gtk"

	gltsgslv "github.com/hfmrow/genLib/tools/goSources/libVendoring"

	gidg "github.com/hfmrow/gtk3Import/dialog"
	gimc "github.com/hfmrow/gtk3Import/misc"
	gits "github.com/hfmrow/gtk3Import/tools"
	gitw "github.com/hfmrow/gtk3Import/treeview"

	gitvsv "github.com/hfmrow/gtk3Import/textView/sourceView"
)

// App infos. Only this part can be modified during an update.
var Name = "Functions & Library Manager"
var Vers = "v1.0"
var Descr = "This software allows you to search for a function\nin specified libraries. It also allows you to create a\n'vendor' directory for the relevant project containing\n all the imported/selected libraries used to compile\nthe targeted project successfully."
var Creat = "H.F.M"
var YearCreat = "2019-20"
var LicenseShort = "This program comes with absolutely no warranty.\nSee the The MIT License (MIT) for details:\nhttps://opensource.org/licenses/mit-license.php"
var LicenseAbrv = "License (MIT)"
var Repository = "github.com/hfmrow/funcLibraryMgr"

// Vars declarations
var (
	absoluteRealPath,
	optFilename,
	tempDir string
	mainOptions *MainOpt
	devMode,
	doTempDir bool
)

var (
	sourcesFromAst []LibsInfos
	desc           Description

	// Current description index
	indexCurrText, indexCurrTree int          // TODO Check if always needed
	declIdexes                   *DeclIndexes // Global indexes

	/*
	 * Library mapping
	 */

	tvsTreeSearch,
	tvsTreeVendor,
	tvsLibInc,
	tvsLibExc *gitw.TreeViewStructure
	TreeViewStructureNew = gitw.TreeViewStructureNew

	statusbar  *gimc.StatusBar
	libsVendor *gltsgslv.LibsVendor

	dialogText *gidg.DialogBoxStructure

	chk *gtk.CheckButton

	svs                 *gitvsv.SourceViewStruct
	SourceViewStructNew = gitvsv.SourceViewStructNew
	markFound           *gtk.TextTag

	DialogBoxNew          = gidg.DialogBoxNew
	StatusBarStructureNew = gimc.StatusBarStructureNew
	SpinbuttonSetValues   = gits.SpinbuttonSetValues

	DlgErr = func(dsc string, err error) bool {
		return gidg.DialogError(mainObjects.MainWindow, sts["issue"], dsc, err, devMode, true)
	}

	DlgMsg       = gidg.DialogMessage
	GetEntryText = gits.GetEntryText

	popupMenu,
	popupLibInc,
	popupLibExc,
	popMenuTextView *gimc.PopupMenuStruct

	dndLibInc, dndLibExc *gimc.DragNDropStruct
	DragNDropNew         = gimc.DragNDropNew
)

type MainOpt struct {
	/* Public, will be saved and restored */
	About *gidg.AboutInfos

	MainWinWidth,
	MainWinHeight,
	MainWinPosX,
	MainWinPosY int

	SourceWinWidth,
	SourceWinHeight,
	SourceWinPosX,
	SourceWinPosY,
	PanedWidth int

	LanguageFilename, // In case GOTranslate is used.
	LastDescFilename,
	LastProjFilename,

	AppLauncher,
	HighlightUserDefined string

	SourceLibs,
	SubDirToSkip,
	DefaultExclude []string
	ListSep string

	SearchCharMinLen,
	ScoreThreshold,
	MinScoreThreshold,
	MaxScoreThreshold int

	FixedMapWidth,
	Wrap,
	MarkResult,
	Functions,
	Structures,
	Exported,
	AddShortcuts bool

	CurrentPage,
	DefaulLanguage,
	DefaultStyle,
	TxtBgCol,
	TxtFgCol,
	SelBgCol,
	DefRangeCol string

	/* Private, will NOT be saved */
	listStoreColumns,
	dispResultsColumns,
	libInclude [][]string
}

// Main options initialisation
func (opt *MainOpt) Init() {
	opt.About = new(gidg.AboutInfos)

	opt.MainWinWidth = 800
	opt.MainWinHeight = 600

	opt.PanedWidth = 120

	opt.LanguageFilename = "assets/lang/eng.lang"

	opt.DefaultExclude = []string{"TEST", ".git"}
	opt.LastDescFilename = "descList.lst"
	opt.listStoreColumns = [][]string{
		{"Found", "markup"},
		{"Type", "text"},
		{"Exported", "text"},
		{"Library", "text"},
		{"Score", "integer"}, // -> Will be hidden since it is an (int) column
		// used to track original place of the information.
		{"idx", "integer"}} // -> Will be hidden since it is an (int) column
	opt.dispResultsColumns = [][]string{
		{"", "active"},
		{"", "markup"},
		{"idx", "integer"}} // -> Will be hidden since it is an (int) column

	opt.libInclude = [][]string{
		{"Name", "text"},
		{"Path", "text"}}

	opt.ListSep = ";"

	opt.SourceWinWidth = 800
	opt.SourceWinHeight = 480
	opt.SourceWinPosX = -1
	opt.SourceWinPosY = -1

	opt.SearchCharMinLen = 2
	opt.ScoreThreshold = -200
	opt.MinScoreThreshold = -1000
	opt.MaxScoreThreshold = 1000

	opt.DefaulLanguage = "go-hfmrow"
	opt.DefaultStyle = "hfmrow"
	opt.HighlightUserDefined = "assets/langAndstyle"

	opt.TxtBgCol = "#F2F2F2"
	opt.TxtFgCol = "#1A1A1A"

	opt.DefRangeCol = "#FFFFFF"
	opt.SelBgCol = "#B5F1F1"

	opt.AppLauncher = "xdg-open"

	opt.MarkResult = true
	opt.Functions = true
	opt.Structures = true
	opt.Exported = true
	opt.AddShortcuts = false
}

// Variables -> Objects.
func (opt *MainOpt) UpdateObjects() {
	mainObjects.MainWindow.Resize(opt.MainWinWidth, opt.MainWinHeight)
	mainObjects.MainWindow.Move(opt.MainWinPosX, opt.MainWinPosY)

	mainObjects.SpinButtonScoreThreshold.SetValue(float64(opt.ScoreThreshold))
	mainObjects.FileChooserSelectDir.SetFilename(opt.LastProjFilename)

	mainObjects.CheckBoxIncludeFunctions.SetActive(opt.Functions)
	mainObjects.CheckBoxIncludeStructures.SetActive(opt.Structures)
	mainObjects.CheckBoxIncludeExported.SetActive(opt.Exported)
	mainObjects.CheckBoxAddShortcuts.SetActive(opt.AddShortcuts)

	mainObjects.SourceToggleButtonWrap.SetActive(opt.Wrap)
	mainObjects.SourceToggleButtonMapWidth.SetActive(opt.FixedMapWidth)
	mainObjects.WindowSource.Resize(opt.SourceWinWidth, opt.SourceWinHeight)
	mainObjects.WindowSource.Move(opt.SourceWinPosX, opt.SourceWinPosY)
	mainObjects.PanedSource.SetPosition(opt.SourceWinWidth - opt.PanedWidth)

	mainObjects.Stack.SetVisibleChildName(opt.CurrentPage)
}

// Objects -> Variables.
func (opt *MainOpt) UpdateOptions() {
	opt.MainWinWidth, opt.MainWinHeight = mainObjects.MainWindow.GetSize()
	opt.MainWinPosX, opt.MainWinPosY = mainObjects.MainWindow.GetPosition()

	opt.ScoreThreshold = mainObjects.SpinButtonScoreThreshold.GetValueAsInt()

	opt.Functions = mainObjects.CheckBoxIncludeFunctions.GetActive()
	opt.Structures = mainObjects.CheckBoxIncludeStructures.GetActive()
	opt.Exported = mainObjects.CheckBoxIncludeExported.GetActive()
	opt.AddShortcuts = mainObjects.CheckBoxAddShortcuts.GetActive()

	opt.Wrap = mainObjects.SourceToggleButtonWrap.GetActive()
	opt.FixedMapWidth = mainObjects.SourceToggleButtonMapWidth.GetActive()
	opt.SourceWinWidth, opt.SourceWinHeight = mainObjects.WindowSource.GetSize()
	opt.SourceWinPosX, opt.SourceWinPosY = mainObjects.WindowSource.GetPosition()
	opt.PanedWidth = opt.SourceWinWidth - mainObjects.PanedSource.GetPosition()

	opt.SourceLibs = retrieveTreeView(tvsLibInc)
	opt.SubDirToSkip = retrieveTreeView(tvsLibExc)

	opt.CurrentPage = mainObjects.Stack.GetVisibleChildName()
}

// Read Options from file
func (opt *MainOpt) Read() (err error) {
	var textFileBytes []byte
	opt.Init()
	if textFileBytes, err = ioutil.ReadFile(optFilename); err == nil {
		err = json.Unmarshal(textFileBytes, &opt)
	}
	if err != nil {
		fmt.Printf("Error while reading options file: %s\n", err.Error())
	}
	return
}

// Write Options to file
func (opt *MainOpt) Write() (err error) {
	var jsonData []byte
	var out bytes.Buffer
	opt.UpdateOptions()
	if jsonData, err = json.Marshal(&opt); err == nil {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			err = ioutil.WriteFile(optFilename, out.Bytes(), os.ModePerm)
		}
	}
	return err
}
