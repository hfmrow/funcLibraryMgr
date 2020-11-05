// entry.go

/*
	©2019 H.F.M. MIT license
*/

package gtk3Import

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"

	glsg "github.com/hfmrow/genLib/strings"
)

// GetTextView: Retrieve text from TextView as []string
func GetTextView(tv *gtk.TextView, removeEmpty ...bool) (out []string) {
	var re bool
	var tmpTxt string
	var err error
	var buff *gtk.TextBuffer

	if len(removeEmpty) > 0 {
		re = removeEmpty[0]
	}
	if buff, err = tv.GetBuffer(); err == nil {
		if tmpTxt, err = buff.GetText(buff.GetStartIter(), buff.GetEndIter(), false); err == nil {
			out = strings.Split(tmpTxt, glsg.GetTextEOL([]byte(tmpTxt)))
			if re {
				for idx, line := range out {
					if len(line) == 0 {
						out = append(out[:idx], out[idx+1:]...)
					}
				}
			}
		}
	}
	if err != nil {
		fmt.Printf("GetTextView: %s", err.Error())
	}
	return
}

// SetTextView: Set []string to TextView
func SetTextView(tv *gtk.TextView, in []string, removeEmpty ...bool) {
	var re bool
	var err error
	var buff *gtk.TextBuffer

	if len(removeEmpty) > 0 {
		re = removeEmpty[0]
	}
	if buff, err = tv.GetBuffer(); err == nil {
		if re {
			for idx, line := range in {
				if len(line) == 0 {
					in = append(in[:idx], in[idx+1:]...)
				}
			}
		}
		buff.SetText(strings.Join(in, glsg.GetOsLineEnd()))
	}
	if err != nil {
		fmt.Printf("SetTextView: %s", err.Error())
	}
	return
}

// GetSepEntry: Sanitize and get separated entries
func GetSepEntry(e *gtk.Entry, separator string) (out []string) {
	tmpOut := strings.Split(strings.TrimSpace(GetEntryText(e)), separator)
	if len(tmpOut) > 0 {
		for _, item := range tmpOut {
			if len(item) > 0 {
				out = append(out, strings.TrimSpace(item))
			}
		}
	}
	return
}

// SetSepEntry: set separated entries
func SetSepEntry(e *gtk.Entry, separator string, in []string) {
	e.SetText(strings.Join(in, separator+" "))
}

// GetExtEntry: Sanitize and get extension entries
func GetExtEntry(e *gtk.Entry, separator string) (out []string) {
	tmpOut := strings.Split(strings.TrimSpace(GetEntryText(e)), separator)
	if len(tmpOut) > 0 {
		for _, ext := range tmpOut {
			if len(ext) > 0 {
				tmp := strings.Split(strings.TrimSpace(ext), ".")
				var tmp1 []string
				for _, s := range tmp {
					if len(s) == 0 {
						s = "*"
					}
					tmp1 = append(tmp1, s)
				}
				out = append(out, strings.Join(tmp1, "."))
			}
		}
	}
	return
}

// SetExtEntry: set extension entries
func SetExtEntry(e *gtk.Entry, separator string, in []string) {
	if len(in) > 0 {
		e.SetText(strings.Join(in, separator+" "))
	} else {
		e.SetText("")
	}
}

// GetEntryText: retrieve value of an entry control.
func GetEntryText(entry *gtk.Entry) (outString string) {
	var err error
	if outString, err = entry.GetText(); err != nil {
		fmt.Printf("GetEntryText: %s", err.Error())
	}
	return
}

// GetEntryTextAsInt: retrieve value of an entry control as integer
func GetEntryTextAsInt(entry *gtk.Entry) (outint int) {
	var err error
	var outString string
	if outString, err = entry.GetText(); err == nil {
		if outint, err = strconv.Atoi(outString); err == nil {
			return
		}
	}
	if err != nil {
		fmt.Printf("GetEntryTextAsInt: %s", err.Error())
	}
	return
}
