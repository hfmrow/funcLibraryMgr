// gohOptions.go

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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gotk3/gotk3/gtk"

	gltsgslv "github.com/hfmrow/genLib/tools/goSources/libVendoring"
	gltsle "github.com/hfmrow/genLib/tools/log2file"

	gidg "github.com/hfmrow/gtk3Import/dialog"
	gimc "github.com/hfmrow/gtk3Import/misc"
	gits "github.com/hfmrow/gtk3Import/tools"
	gitw "github.com/hfmrow/gtk3Import/treeview"

	gitvsv "github.com/hfmrow/gtk3Import/textView/sourceView"
)

// App infos. Only this part can be modified during an update.
var (
	Name         = "Functions & Library Manager"
	Vers         = "v1.1.4"
	Descr        = "This software allows you to search for a function\nin specified libraries. It also allows you to create a\n'vendor' directory for the relevant project containing\n all the imported/selected libraries used to compile\nthe targeted project successfully."
	Creat        = "H.F.M"
	YearCreat    = "2019-21"
	LicenseShort = "This program comes with absolutely no warranty.\nSee the The MIT License (MIT) for details:\nhttps://opensource.org/licenses/mit-license.php"
	LicenseAbrv  = "License (MIT)"
	Repository   = "github.com/hfmrow/go-func-lib-mgr"
)

// Vars declarations
var (
	absoluteRealPath,
	optFilename,
	tempDir string
	mainOptions *MainOpt
	devMode,
	doTempDir bool
)

type libs struct {
	Path   string
	Active bool
}

var (
	sourcesFromAst []LibsInfos
	desc           Description

	// Current description index
	indexCurrText int
	declIdexes    *DeclIndexes // Global indexes

	/*
	 * Library mapping
	 */
	// Treeview
	tvsTreeSearch,
	tvsTreeVendor,
	tvsLibInc,
	tvsLibExc *gitw.TreeViewStructure
	TreeViewStructureNew = gitw.TreeViewStructureNew

	// Misc
	GetEntryText    = gits.GetEntryText
	SpinScaleSetNew = gits.SpinScaleSetNew

	// Statusbar
	StatusBarStructureNew = gimc.StatusBarStructureNew
	statusbar             *gimc.StatusBar

	// Vendoring
	libsVendor   *gltsgslv.LibsVendor
	LibVendorNew = gltsgslv.LibVendorNew

	// error logging
	Logger            *gltsle.Log2FileStruct
	Log2FileStructNew = gltsle.Log2FileStructNew

	// Dialog
	DlgMsg       = gidg.DialogMessage
	dialogText   *gidg.DialogBoxStructure
	DialogBoxNew = gidg.DialogBoxNew
	chk          *gtk.CheckButton

	// SourceView
	svs                 *gitvsv.SourceViewStruct
	SourceViewStructNew = gitvsv.SourceViewStructNew
	markFound           *gtk.TextTag

	DlgErr = func(dsc string, err error) (yes bool) {
		yes = gidg.DialogError(mainObjects.MainWindow, sts["issue"], dsc, err, devMode, true)
		Logger.Log(err, dsc)
		return
	}

	// Popup
	PopupMenuIconStructNew = gimc.PopupMenuIconStructNew
	popupMenu,
	popupLibInc,
	popupLibExc,
	popMenuTextView *gimc.PopupMenuIconStruct

	// D&D
	dndLibInc, dndLibExc *gimc.DragNDropStruct
	DragNDropNew         = gimc.DragNDropNew

	// Clipboard
	clipboard    *gimc.Clipboard
	ClipboardNew = gimc.ClipboardNew

	// TreeView conlumns mapping
	includeMap = map[string]int{
		`chk`:  0,
		`name`: 1,
		`path`: 2,
	}

	mapListStoreColumns = map[string]int{
		"found": 0,
		"type":  1,
		"exprt": 2,
		"libry": 3,
		"score": 4,
		"idx":   5}
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
	SubDirToSkip []libs
	DefaultExclude []string
	ListSep        string

	SearchCharMinLen int

	ScoreThreshold,
	MinScoreThreshold,
	MaxScoreThreshold float64

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
		{"", "active"},
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

	opt.ScoreThreshold = mainObjects.SpinButtonScoreThreshold.GetValue()

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

	opt.About.DlgBoxStruct = nil // remove dialog object before saving

	if jsonData, err = json.Marshal(&opt); err == nil {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			err = ioutil.WriteFile(optFilename, out.Bytes(), os.ModePerm)
		}
	}
	return err
}
