// spinButton.go

/// +build ignore

/*
	©2019 H.F.M. MIT license
*/

package gtk3Import

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

// SpinbuttonSetValues: Configure spin button
// min, max, value, stepIncrement, pageIncrement, pageSize
// nil value send as object mean just return a configured *gtk.Adjustment
func SpinbuttonSetValues(sb interface{}, min, max, value int, step ...int) (adjustment *gtk.Adjustment, err error) {
	incStep, pageIncrement, pageSize := 1, 0, 0
	switch len(step) {
	case 1:
		incStep = step[0]
	case 2:
		incStep = step[0]
		pageIncrement = step[1]
	case 3:
		incStep = step[0]
		pageIncrement = step[1]
		pageSize = step[2]
	}
	if adjustment, err = gtk.AdjustmentNew(float64(value), float64(min), float64(max),
		float64(incStep), float64(pageIncrement), float64(pageSize)); err == nil {
		if sb != nil {
			sb.(*gtk.SpinButton).Configure(adjustment, 1, 0)
		}
	} else {
		err = fmt.Errorf("SpinbuttonSetValues: %vi", err)
	}
	return
}
