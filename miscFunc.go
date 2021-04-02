// miscFunc.go

// Source file auto-generated on Sun, 06 Oct 2019 23:05:32 using Gotk3ObjHandler v1.3.8 ©2018-19 H.F.M
/*
	Copyright ©2019 H.F.M - Functions Library Manager
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"fmt"
	"strings"
)

// initSourceDirectories:
func initSourceDirectories() {
	var err error
	// if len(mainOptions.SourceLibs) > 0 {
	err = initSourceLibs(mainOptions.SourceLibs, mainOptions.SubDirToSkip)
	DlgErr(sts["errEOF"], err)
	updateStatusBar()
	EntrySearchForChanged(mainObjects.EntrySearchFor)
	// }
}

// updateStatusBar:
func updateStatusBar(status ...string) {

	if declIdexes != nil {
		if len(desc.Desc) == 0 {
			declIdexes.Count = 0
		}
		statusbar.Set(fmt.Sprintf("%d", declIdexes.Count), 0)
		statusbar.Set(fmt.Sprintf("%d", tvsTreeSearch.CountRows()), 1)
		statusbar.Set(fmt.Sprintf("%s", mainOptions.LastDescFilename), 2)

		if len(status) > 0 {
			statusbar.Set(strings.Join(status, statusbar.Separator), 3)
		} else {
			statusbar.Set("", 3)
		}
	}
}

// IsExistSlice: if exist then  ...
func isExistSlice(slice []libs, item libs) bool {
	for _, mainRow := range slice {
		if mainRow.Path == item.Path {
			return true
		}
	}
	return false
}
