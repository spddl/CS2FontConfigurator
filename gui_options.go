package main

import (
	"github.com/tailscale/walk"

	. "github.com/tailscale/walk/declarative"
)

func OpenOptions(owner walk.Form) int {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton
	var NumberEditDpi *walk.NumberEdit
	var CheckBoxEmbeddedBitmap,
		CheckBoxPreferOutline,
		CheckBoxDoSubstitutions,
		CheckBoxBitmapMonospace,
		CheckBoxForceAutohint,
		CheckBoxQtUseSubpixelPositioning *walk.CheckBox

	if err := (Dialog{
		AssignTo:      &dlg,
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		MinSize:       Size{Width: 1, Height: 1},
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout: VBox{},
				Children: []Widget{

					Composite{
						Layout: HBox{},
						Children: []Widget{
							Label{
								Text: "Global Font Properties",
							},
						},
					},

					Composite{
						Layout: Grid{Columns: 1},
						Children: []Widget{
							CheckBox{
								AssignTo:    &CheckBoxEmbeddedBitmap,
								Checked:     config.EmbeddedBitmap,
								Text:        "Embedded Bitmap",
								ToolTipText: "Bitmap fonts are sometimes used as fallbacks for missing fonts, which may cause text to be rendered pixelated or too large.",
								OnCheckStateChanged: func() {
									config.EmbeddedBitmap = CheckBoxEmbeddedBitmap.Checked()
								},
							},

							CheckBox{
								AssignTo: &CheckBoxPreferOutline,
								Checked:  config.PreferOutline,
								Text:     "Prefer Outline",
								OnCheckStateChanged: func() {
									config.PreferOutline = CheckBoxPreferOutline.Checked()
								},
							},

							CheckBox{
								AssignTo: &CheckBoxDoSubstitutions,
								Checked:  config.DoSubstitutions,
								Text:     "Do Substitutions",
								OnCheckStateChanged: func() {
									config.DoSubstitutions = CheckBoxDoSubstitutions.Checked()
								},
							},

							CheckBox{
								AssignTo: &CheckBoxBitmapMonospace,
								Checked:  config.BitmapMonospace,
								Text:     "Bitmap Monospace",
								OnCheckStateChanged: func() {
									config.BitmapMonospace = CheckBoxBitmapMonospace.Checked()
								},
							},

							CheckBox{
								AssignTo: &CheckBoxForceAutohint,
								Checked:  config.ForceAutohint,
								Text:     "Force Autohint",
								OnCheckStateChanged: func() {
									config.ForceAutohint = CheckBoxForceAutohint.Checked()
								},
							},

							CheckBox{
								AssignTo: &CheckBoxQtUseSubpixelPositioning,
								Checked:  config.QtUseSubpixelPositioning,
								Text:     "Qt Use Subpixel Positioning",
								OnCheckStateChanged: func() {
									config.QtUseSubpixelPositioning = CheckBoxQtUseSubpixelPositioning.Checked()
								},
							},
						},
					},

					Composite{
						Layout: HBox{
							SpacingZero: true,
						},
						Children: []Widget{
							Label{
								Text: "DPI:",
							},
							NumberEdit{
								AssignTo:           &NumberEditDpi,
								Value:              float64(config.Dpi),
								MinValue:           0,
								SpinButtonsVisible: true,
								OnValueChanged: func() {
									config.Dpi = int(NumberEditDpi.Value())
								},
							},
						},
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo:  &acceptPB,
						Text:      "OK",
						OnClicked: func() { dlg.Accept() },
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      "Cancel",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
		},
	}.Create(owner)); err != nil {
		panic(err)
	}
	return dlg.Run()
}
