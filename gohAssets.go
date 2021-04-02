// gohAssets.go

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
	"embed"
	"log"
)

//go:embed assets/glade
//go:embed assets/images
var embeddedFiles embed.FS

// This functionality does not require explicit encoding of the files, at each
// compilation, the files are inserted into the resulting binary. Thus, updating
// assets is only required when new files are added to be embedded in order to
// create and declare the variables to which the files are linked.
// assetsDeclarationsUseEmbedded: Use native Go 'embed' package to include files
// content at runtime.
func assetsDeclarationsUseEmbedded(embedded ...bool) {
	mainGlade = readEmbedFile("assets/glade/main.glade")
	actionsEditCopy = readEmbedFile("assets/images/Actions-edit-copy.png")
	crossIcon48 = readEmbedFile("assets/images/Cross-icon-48.png")
	folderOpen20 = readEmbedFile("assets/images/Folder-open-20.png")
	funcfinderlibmgr372x32 = readEmbedFile("assets/images/FuncFinderLibMgr-372x32.png")
	funcfinderlibmgr558x48 = readEmbedFile("assets/images/FuncFinderLibMgr-558x48.png")
	libraryManager = readEmbedFile("assets/images/library-manager.png")
	logout48 = readEmbedFile("assets/images/logout-48.png")
	mimetypeSourceIconGolang48 = readEmbedFile("assets/images/Mimetype-source-icon-golang-48.png")
	pack48 = readEmbedFile("assets/images/pack-48.png")
	refresh = readEmbedFile("assets/images/refresh.png")
	searchsource48 = readEmbedFile("assets/images/searchSource-48.png")
	tickIcon48 = readEmbedFile("assets/images/Tick-icon-48.png")
}

// readEmbedFile: read 'embed' file system and return []byte data.
func readEmbedFile(filename string) (out []byte) {
	var err error
	out, err = embeddedFiles.ReadFile(filename)
	if err != nil {
		log.Printf("Unable to read embedded file: %s, %v\n", filename, err)
	}
	return
}
