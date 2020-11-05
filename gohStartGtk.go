// gohStartGtk.go

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
	"log"
	"os"

	"github.com/gotk3/gotk3/gtk"
)

/*******************************/
/* Gtk3 Window Initialisation */
/*****************************/
func mainStartGtk(winTitle string, width, height int, center bool) {
	mainObjects = new(MainControlsObj)
	gtk.Init(nil)
	if err := newBuilder(mainGlade); err == nil {
		/* Init tempDir and plan to delete it when leaving. */
		if doTempDir {
			tempDir = tempMake(Name)
			defer os.RemoveAll(tempDir)
		}
		/* Parse Gtk objects */
		gladeObjParser()
		/* Update gtk conctrols with stored values into mainOptions */
		mainOptions.UpdateObjects()
		/* Fill control with images */
		assignImages()
		/* Start main application ... */
		mainApplication()
		/* Objects Signals initialisations */
		signalsPropHandler()
		/* Set Window Properties */
		if center {
			mainObjects.MainWindow.SetPosition(gtk.WIN_POS_CENTER)
		}
		mainObjects.MainWindow.SetTitle(winTitle)
		mainObjects.MainWindow.SetDefaultSize(width, height)
		mainObjects.MainWindow.Connect("delete-event", windowDestroy)
		mainObjects.MainWindow.ShowAll()
		/* Start Gui loop */
		gtk.Main()
	} else {
		log.Fatal("Builder initialisation error.", err.Error())
	}
}

// windowDestroy: on closing/destroying the gui window.
func windowDestroy() {
	if onShutdown() {
		gtk.MainQuit()
	}
}
