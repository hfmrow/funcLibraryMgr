// gohImages.go

/*
	Source file auto-generated on Fri, 02 Apr 2021 08:27:58 using Gotk3 Objects Handler v1.7.5 ©2018-21 hfmrow
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2019-21 H.F.M - Functions & Library Manager v1.1.4 github.com/hfmrow/go-func-lib-mgr
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

/**********************************************************/
/* This section preserve user modifications on update.   */
/* Images declarations, used to initialize objects with */
/* The SetPict() func, accept both kind of variables:  */
/* filename or []byte content in case of using        */
/* embedded binary data. The variables names are the */
/* same. "assetsDeclarationsUseEmbedded(bool)" func */
/* could be used to toggle between filenames and   */
/* embedded binary type. See SetPict()            */
/* declaration to learn more on how to use it.   */
/************************************************/
func assignImages() {
	size := 18
	SetPict(mainObjects.ButtonCreatVendor, libraryManager, size)
	SetPict(mainObjects.ButtonExit, logout48, size)
	SetPict(mainObjects.ButtonLoadProj, folderOpen20, size)
	SetPict(mainObjects.ButtonRefreshLibraryData, refresh, size)
	SetPict(mainObjects.ButtonSourceOk, tickIcon48, size)
	SetPict(mainObjects.ButtonStackLibrarySelection, mimetypeSourceIconGolang48, size)
	SetPict(mainObjects.ButtonStackSearch, searchsource48, size)
	SetPict(mainObjects.ButtonStackVendoring, pack48, size-2)
	SetPict(mainObjects.ImageTop, funcfinderlibmgr558x48)
	SetPict(mainObjects.MainWindow, libraryManager, size)
	SetPict(mainObjects.SourceToggleButtonMapWidth, "")
	SetPict(mainObjects.SourceToggleButtonWrap, "")
	SetPict(mainObjects.SpinButtonScoreThreshold, "")
	SetPict(mainObjects.WindowSource, searchsource48, size)
}

/**********************************************************/
/* This section is rewritten on assets update.           */
/* Assets var declarations, this step permit to make a  */
/* bridge between the differents types used, string or */
/* []byte, and to simply switch from one to another.  */
/*****************************************************/
var mainGlade interface{}                  // assets/glade/main.glade
var actionsEditCopy interface{}            // assets/images/Actions-edit-copy.png
var crossIcon48 interface{}                // assets/images/Cross-icon-48.png
var folderOpen20 interface{}               // assets/images/Folder-open-20.png
var funcfinderlibmgr372x32 interface{}     // assets/images/FuncFinderLibMgr-372x32.png
var funcfinderlibmgr558x48 interface{}     // assets/images/FuncFinderLibMgr-558x48.png
var libraryManager interface{}             // assets/images/library-manager.png
var logout48 interface{}                   // assets/images/logout-48.png
var mimetypeSourceIconGolang48 interface{} // assets/images/Mimetype-source-icon-golang-48.png
var pack48 interface{}                     // assets/images/pack-48.png
var refresh interface{}                    // assets/images/refresh.png
var searchsource48 interface{}             // assets/images/searchSource-48.png
var tickIcon48 interface{}                 // assets/images/Tick-icon-48.png
