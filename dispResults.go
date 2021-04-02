// dispResults.go

/*
	Source file auto-generated on Wed, 23 Oct 2019 18:40:49 using Gotk3ObjHandler v1.3.9 ©2018-19 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2019 H.F.M - Functions & Library Manager
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"fmt"
	"path/filepath"

	"github.com/gotk3/gotk3/gtk"
	// gltsbh "github.com/hfmrow/genLib/tools/bench"
)

/*
	Main TreeStore (results)
*/

// displayTreeStore: Fill TreeViewFound with found results
func displayTreeStore(in []toDispTreeStore) (err error) {

	var (
		iter *gtk.TreeIter

		onExitFunc = func() {
			// Attach listStore
			tvsTreeSearch.StoreAttach()
			// Update statusbar
			if tvsTreeSearch.CountRows() == 0 {
				updateStatusBar(fmt.Sprintf(sts["noResult"]+"\"%s\"", GetEntryText(mainObjects.EntrySearchFor)))
			}
			updateStatusBar()
		}
	)

	// Detach & clear listStore
	tvsTreeSearch.StoreDetach()
	tvsTreeSearch.TreeStore.Clear()
	defer onExitFunc()

	if len(in) > 0 {

		// Add parents
		for _, row := range in {
			if iter, err = tvsTreeSearch.AddRow(nil, row.Name, row.Type, row.Exported, row.Path, row.Score, row.Idx); err != nil {
				DlgErr("displayTreeStore:AddParents", err)
				return
			}
			// Add childs if exists (structures' methods)
			for _, subRow := range row.Methods {
				if _, err = tvsTreeSearch.AddRow(iter, subRow.Name, subRow.Type, subRow.Exported, subRow.Path, subRow.Score, subRow.Idx); err != nil {
					DlgErr("displayTreeStore:AddChilds", err)
					return
				}
			}
		}
	}

	return
}

/*
	Popup window
*/

// popupTreeview: Display content as TextView
func popupSourceView(index int) {

	var err error

	if svs == nil {

		svs, err = SourceViewStructNew(mainObjects.Source, mainObjects.SourceMap, mainObjects.WindowSource)
		DlgErr("popupSourceView/SourceViewStructNew", err)
		svs.View.SetEditable(false)

		// Handling "populate-popup" signal to add some personal entries
		svs.View.Connect("populate-popup", popupTextViewPopulateMenu)

		// TODO Think to use search in preview window
		// Make a tag to indicate found element (when HighlightFound not checked)
		tag := make(map[string]interface{})
		tag["background"] = "#ABF6FF"
		markFound = svs.Buffer.CreateTag("markFound", tag)

		// Language & style, add a personal version for Golang (directory content)
		svs.UserStylePath = filepath.Join(absoluteRealPath, mainOptions.HighlightUserDefined)
		svs.UserLanguagePath = filepath.Join(absoluteRealPath, mainOptions.HighlightUserDefined)
		svs.DefaultLanguageId = mainOptions.DefaulLanguage
		svs.DefaultStyleShemeId = mainOptions.DefaultStyle

		svs.ComboboxHandling(
			mainObjects.ComboboxSourceLanguage,
			mainObjects.ComboboxSourceStyle,
			&mainOptions.DefaulLanguage,
			&mainOptions.DefaultStyle)

		mainObjects.WindowSource.Resize(mainOptions.SourceWinWidth, mainOptions.SourceWinHeight)
		mainObjects.WindowSource.Move(mainOptions.SourceWinPosX, mainOptions.SourceWinPosY)
		mainObjects.PanedSource.SetPosition(mainOptions.MainWinWidth - mainOptions.PanedWidth)
	}

	indexCurrText = index // needed to always get the last selected choice
	dispTextView(index)
}

// dispTreeStore:
func dispTextView(index int) {
	var err error

	descr, ok := declIdexes.GetDescr(index)
	if ok {
		if err = svs.LoadSource(filepath.Join(desc.RootLibs, descr.File)); err == nil {

			mainObjects.WindowSource.SetTitle(descr.File)

			svs.TxtBgCol = mainOptions.TxtBgCol
			svs.TxtFgCol = mainOptions.TxtFgCol

			svs.ColorBgRangeSet = mainOptions.DefRangeCol
			svs.SelBgCol = mainOptions.SelBgCol

			svs.ColorBgRange(descr.LineStart+1, descr.LineEnd+1)
			svs.SelectRange(descr.LineStart+1, descr.LineEnd+1)

			svs.RunAfterEvents(func() {
				svs.BringToFront()
				svs.ScrollToLine(descr.LineStart)
			})

		}
		DlgErr("ReadFile", err)
		// tvn.TextView.GrabFocus()
		return
	}
	DlgErr("Description not found", fmt.Errorf("#%d", index))
}
