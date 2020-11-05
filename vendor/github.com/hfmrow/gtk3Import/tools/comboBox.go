// comboBox.go

/*
	©2019 H.F.M. MIT license
*/

package gtk3Import

import (
	"errors"
	"fmt"
	"log"

	"github.com/gotk3/gotk3/glib"

	"github.com/gotk3/gotk3/gtk"
)

// type ComboBoXStruct struct {
// 	TreeModel *gtk.TreeModel
// 	ListStore *gtk.ListStore
// 	Err       error
// 	CBT       *gtk.ComboBoxText
// 	CB        *gtk.ComboBox

// 	cbxX interface{}
// }

// // ComboBoXStructNew: Fill object with []sting. The GtkListStore returned
// // may be used to add items (using "ComboBoXListStoreAdd") to a GtkComboBox
// // a better way will be to use a GtkComboBoxText.
// // handle both, GtkComboBox or GtkComboBoxText
// func ComboBoXStructNew(cbxX interface{}, in []string) (cbxs *ComboBoXStruct) {
// 	var err error
// 	var iMdl gtk.ITreeModel
// 	cbxs = new(ComboBoXStruct)
// 	cbxs.cbxX = cbxX

// 	if cbxs.ListStore, err = gtk.ListStoreNew(glib.TYPE_STRING); err == nil {
// 		for _, item := range in {
// 			cbxs.ComboBoXListStoreAdd(item)
// 		}
// 		switch cbx := cbxs.cbxX.(type) {
// 		case *gtk.ComboBox:
// 			if mdl, _ := cbx.GetModel(); mdl == nil {
// 				// First time then build rendercell ...
// 				var cRndTxt *gtk.CellRendererText
// 				if cRndTxt, cbxs.Err = gtk.CellRendererTextNew(); cbxs.Err == nil {
// 					cbx.PackStart(cRndTxt, true)
// 					cbx.AddAttribute(cRndTxt, "text", 0)
// 					cbx.SetIDColumn(0)
// 				} else {
// 					return
// 				}
// 			}
// 			cbx.SetModel(cbxs.ListStore.ToTreeModel())
// 			cbxs.CB = cbx
// 			iMdl, cbxs.Err = cbx.GetModel()
// 		case *gtk.ComboBoxText:
// 			cbx.SetModel(cbxs.ListStore.ToTreeModel())
// 			cbxs.CBT = cbx
// 			iMdl, cbxs.Err = cbx.GetModel()
// 		default:
// 			cbxs.Err = errors.New("Unable to fill ComboBox")
// 		}
// 	}

// 	if cbxs.Err == nil {
// 		cbxs.TreeModel = iMdl.ToTreeModel()
// 	}
// 	return
// }

// // ComboBoXFill: Append or Prepend an entry if it not already exist.
// // handle both, GtkComboBox or GtkComboBoxText
// func (cbxs *ComboBoXStruct) ComboBoXAdd(in string, prepend ...interface{}) {
// 	var prependEntry bool
// 	switch len(prepend) {
// 	case 1:
// 		prependEntry = prepend[0].(bool)
// 	}

// 	switch cbx := cbxs.cbxX.(type) {
// 	case *gtk.ComboBox:
// 		fmt.Println("ComboBoXAdd: a GtkComboBox cannot add entry, try using GtkComboBoxText instead.")
// 	case *gtk.ComboBoxText:
// 		if pos := cbxs.ComboBoXFind(in); pos == -1 {
// 			if prependEntry {
// 				cbx.PrependText(in)
// 			} else {
// 				cbx.AppendText(in)
// 			}
// 		}
// 	default:
// 		cbxs.Err = errors.New("Unable to fill ComboBox")
// 	}
// }

// // ComboBoXFind: find string value. Return -1 if nothing found
// // handle both, GtkComboBox or GtkComboBoxText
// func (cbxs *ComboBoXStruct) ComboBoXFind(toFind string) (pos int) {
// 	cbxs.Err = nil
// 	var err error
// 	var val *glib.Value
// 	var valStr string
// 	pos = -1 // Default value returned if nothing found

// 	cbxs.TreeModel.ForEach(func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter, userData ...interface{}) bool {
// 		if val, err = model.GetValue(iter, 0); err == nil {
// 			if valStr, err = val.GetString(); err == nil {
// 				if valStr == toFind {
// 					pos = path.GetIndices()[0]
// 					return true
// 				}
// 			}
// 		}
// 		if err != nil {
// 			return true
// 		} else {
// 			return false
// 		}
// 	})
// 	if err != nil {
// 		log.Printf("ComboBoXDump: %s", err)
// 		cbxs.Err = err
// 	}
// 	return
// }

// // ComboBoXListStoreAdd: Append or Prepend an entry does not check if exist.
// func (cbxs *ComboBoXStruct) ComboBoXListStoreAdd(in string, prepend ...interface{}) {
// 	var err error
// 	var iter *gtk.TreeIter
// 	var prependEntry bool

// 	switch len(prepend) {
// 	case 1:
// 		prependEntry = prepend[0].(bool)
// 	}
// 	if prependEntry {
// 		iter = cbxs.ListStore.Prepend()
// 	} else {
// 		iter = cbxs.ListStore.Append()
// 	}
// 	err = cbxs.ListStore.SetValue(iter, 0, in)
// 	if err != nil {
// 		log.Printf("ComboBoXListStoreAdd: %s", err)
// 		cbxs.Err = err
// 	}
// }

// comboBoXtoTreeModel: retrieve TreeModel
// handle both, GtkComboBox or GtkComboBoxText
func comboBoXtoTreeModel(comboBoX interface{}) *gtk.TreeModel {
	var iTM gtk.ITreeModel
	var err error

	switch cbx := comboBoX.(type) {
	case *gtk.ComboBox:
		iTM, err = cbx.GetModel()
		// if err==errors.New("cgo returned unexpected nil pointer"){

		// 	}
	case *gtk.ComboBoxText:
		iTM, err = cbx.GetModel()
	default:
		err = errors.New("Unable to get model")
	}
	if err != nil {
		log.Fatalf("comboBoXModel: %s", err)
	}
	return iTM.ToTreeModel()
}

// ComboBoXFill: Append or Prepend an entry if it not already exist.
// handle both, GtkComboBox or GtkComboBoxText
func ComboBoXAdd(comboBoX interface{}, in string, prepend ...interface{}) {
	var prependEntry bool
	switch len(prepend) {
	case 1:
		prependEntry = prepend[0].(bool)
	}

	switch cbx := comboBoX.(type) {
	case *gtk.ComboBox:
		fmt.Println("ComboBoXAdd: a GtkComboBox cannot add entry, try using GtkComboBoxText instead.")
	case *gtk.ComboBoxText:
		if pos := ComboBoXFind(comboBoX, in); pos == -1 {
			if prependEntry {
				cbx.PrependText(in)
			} else {
				cbx.AppendText(in)
			}
		}
	default:
		log.Fatalf("Unable to fill ComboBox")
	}
}

// ComboBoXListStoreAdd: Append or Prepend an entry does not check if exist.
func ComboBoXListStoreAdd(lStore *gtk.ListStore, in string, prepend ...interface{}) {
	var err error
	var iter *gtk.TreeIter
	var prependEntry bool

	switch len(prepend) {
	case 1:
		prependEntry = prepend[0].(bool)
	}
	if prependEntry {
		iter = lStore.Prepend()
	} else {
		iter = lStore.Append()
	}
	err = lStore.SetValue(iter, 0, in)
	if err != nil {
		log.Printf("ComboBoXListStoreAdd: %s", err)
	}
}

// ComboBoXFill: Fill object with []sting. The GtkListStore returned
// may be used to add items (using "ComboBoXListStoreAdd") to a GtkComboBox
// a better way will be to use a GtkComboBoxText.
// handle both, GtkComboBox or GtkComboBoxText
func ComboBoXFill(comboBoX interface{}, in []string) (lStore *gtk.ListStore) {
	var err error

	if lStore, err = gtk.ListStoreNew(glib.TYPE_STRING); err == nil {
		for _, item := range in {
			ComboBoXListStoreAdd(lStore, item)
		}
		switch cbx := comboBoX.(type) {
		case *gtk.ComboBox:
			if mdl, _ := cbx.GetModel(); mdl == nil {
				// First time then build rendercell ...
				var cRndTxt *gtk.CellRendererText
				cRndTxt, err = gtk.CellRendererTextNew()
				cbx.PackStart(cRndTxt, true)
				cbx.AddAttribute(cRndTxt, "text", 0)
				cbx.SetIDColumn(0)
			}
			cbx.SetModel(lStore.ToTreeModel())
		case *gtk.ComboBoxText:
			cbx.SetModel(lStore.ToTreeModel())
		default:
			err = errors.New("Unable to fill ComboBox")
		}
	}

	if err != nil {
		log.Fatalf("ComboBoXFill: %s", err)
	}
	return
}

// ComboBoXDump: get all entries
// handle both, GtkComboBox or GtkComboBoxText
func ComboBoXDump(comboBoX interface{}) (out []string) {
	var err error
	var val *glib.Value
	var valStr string

	tModel := comboBoXtoTreeModel(comboBoX)
	tModel.ForEach(func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter, userData ...interface{}) bool {
		if val, err = model.GetValue(iter, 0); err == nil {
			if valStr, err = val.GetString(); err == nil {
				out = append(out, valStr)
				return false
			}
		}
		return true
	})
	if err != nil {
		log.Fatalf("ComboBoXDump: %s", err)
	}
	return
}

// ComboBoXFind: find string value. Return -1 if nothing found
// handle both, GtkComboBox or GtkComboBoxText
func ComboBoXFind(comboBoX interface{}, toFind string) (pos int) {
	var err error
	var val *glib.Value
	var valStr string
	pos = -1 // Default value returned if nothing found

	tModel := comboBoXtoTreeModel(comboBoX)
	tModel.ForEach(func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter, userData ...interface{}) bool {
		if val, err = model.GetValue(iter, 0); err == nil {
			if valStr, err = val.GetString(); err == nil {
				if valStr == toFind {
					pos = path.GetIndices()[0]
					return true
				}
			}
		}
		if err != nil {
			return true
		} else {
			return false
		}
	})
	if err != nil {
		log.Fatalf("ComboBoXDump: %s", err)
	}
	return
}

// ComboBoxTextGetAllEntries: Retrieve all ComboBoxText entries to a string slice.
func ComboBoxTextGetAllEntries(cbxEntry *gtk.ComboBoxText) (outSlice []string) {
	iTreeModel, err := cbxEntry.GetModel()
	model := iTreeModel.ToTreeModel()
	iter, ok := model.GetIterFirst()
	for ok {
		if glibValue, err := model.GetValue(iter, 0); err == nil {
			if entry, err := glibValue.GetString(); err == nil {
				outSlice = append(outSlice, entry)
				ok = model.IterNext(iter)
			}
		}
		if err != nil {
			fmt.Errorf("ComboBoxTextGetAllEntries: %s", err.Error())
		}
	}
	return outSlice
}

// Fill / Clean comboBoxText
func ComboBoxTextFill(cbxEntry *gtk.ComboBoxText, entries []string, options ...bool) {
	var prepend, removeAll bool
	switch len(options) {
	case 1:
		prepend = options[0]
	case 2:
		prepend = options[0]
		removeAll = options[1]
	}
	if !removeAll {
		for _, word := range entries {
			ComboBoxTextAddSetEntry(cbxEntry, word, prepend)
		}
		return
	}
	cbxEntry.RemoveAll()
}

// ComboBoxTextAddSetEntry: Add newEntry if not exist to ComboBoxText, Option: prepend:bool.
// Get index and set cbxText at it if already exist.
func ComboBoxTextAddSetEntry(cbxEntry *gtk.ComboBoxText, newEntry string, prepend ...bool) (existAtPos int) {
	var prependEntry bool
	var count int
	var iter *gtk.TreeIter
	var ok bool
	existAtPos = -1
	if len(prepend) > 0 {
		prependEntry = prepend[0]
	}
	iTreeModel, err := cbxEntry.GetModel()
	model := iTreeModel.ToTreeModel()
	iter, ok = model.GetIterFirst()
	for ok {
		if glibValue, err := model.GetValue(iter, 0); err == nil {
			if entry, err := glibValue.GetString(); err == nil {
				if entry == newEntry {
					existAtPos = count
					break
				}
				count++
				ok = model.IterNext(iter)
			}
		}
		if err != nil {
			fmt.Errorf("ComboBoxTextAddSetEntry: %s", err.Error())
		}
	}
	if existAtPos == -1 {
		switch {
		case prependEntry:
			cbxEntry.PrependText(newEntry)
		default:
			cbxEntry.AppendText(newEntry)
		}
	} else {
		cbxEntry.SetActiveIter(iter)
	}
	return existAtPos
}

func ComboBoxTextClearAll(cbxEntry *gtk.ComboBoxText) {
	cbxEntry.PrependText("")
	cbxEntry.SetActive(0)
	cbxEntry.RemoveAll()

}
