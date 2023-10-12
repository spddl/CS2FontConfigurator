package main

import (
	"unsafe"

	"github.com/tailscale/walk"
	"github.com/tailscale/win"
	"golang.org/x/sys/windows"

	. "github.com/tailscale/walk/declarative"
)

type weight struct {
	Id   int
	Name string
}

func KnownWeight() []*weight {
	return []*weight{
		{0, "Light"},
		{1, "Regular"},
		{2, "Medium"},
		{3, "Demibold"},
		{4, "Bold"},
		{5, "Black"},
	}
}

const (
	WeightLight = iota
	WeightRegular
	WeightMedium
	WeightDemibold
	WeightBold
	WeightBlack
)

func OpenCustomFontPtr(owner walk.Form, customFont *FontStruct) int {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton
	var LineEditValveFont,
		LineEditReplaceFont *walk.LineEdit
	var NumberEditPixelSize,
		NumberEditDpi *walk.NumberEdit
	var ComboBoxWeight *walk.ComboBox
	var CheckBoxEmbeddedBitmap,
		CheckBoxPreferOutline,
		CheckBoxDoSubstitutions,
		CheckBoxBitmapMonospace,
		CheckBoxForceAutohint,
		CheckBoxQtUseSubpixelPositioning *walk.CheckBox

	if err := (Dialog{
		AssignTo:      &dlg,
		Title:         customFont.ValveFont,
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		MinSize:       Size{Width: 1, Height: 1},
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout: VBox{},
				Children: []Widget{

					Composite{
						Layout: Grid{Columns: 3, MarginsZero: true},
						Children: []Widget{

							Label{
								Text: "Valve Font:",
							},
							LineEdit{
								AssignTo:   &LineEditValveFont,
								ColumnSpan: 2,
								Text:       customFont.ValveFont,
								OnEditingFinished: func() {
									customFont.ValveFont = LineEditValveFont.Text()
								},
							},

							VSeparator{
								ColumnSpan: 3,
							},

							Label{
								Text: "Replace Font",
							},

							LineEdit{
								AssignTo: &LineEditReplaceFont,
								ReadOnly: true,
								Text:     "press Select",
							},

							PushButton{
								Text: "Select",
								OnClicked: func() {
									name := [32]uint16{}
									if customFont.ReplaceFont != "" {
										copy(name[:], windows.StringToUTF16(customFont.ReplaceFont))
									}

									// https://learn.microsoft.com/en-us/windows/win32/api/commdlg/ns-commdlg-choosefonta
									cfont := CHOOSEFONT{
										HwndOwner: dlg.Handle(),
										LpLogFont: &win.LOGFONT{LfFaceName: name}, // https://learn.microsoft.com/en-us/windows/win32/api/wingdi/ns-wingdi-logfontw
										Flags:     CF_SCREENFONTS | CF_INITTOLOGFONTSTRUCT | CF_NOSCRIPTSEL | CF_FORCEFONTEXIST,
										NFontType: SCREEN_FONTTYPE,
									}
									cfont.LStructSize = uint32(unsafe.Sizeof(cfont))
									ChooseFontW(&cfont)

									// Replace Font
									customFont.ReplaceFont = windows.UTF16ToString(cfont.LpLogFont.LfFaceName[:])

									// Pixelsize
									customFont.Pixelsize = float64(cfont.IPointSize / 10)
									if cfont.IPointSize != 0 {
										NumberEditPixelSize.SetValue(float64(cfont.IPointSize / 10))
									}

									// Weight
									switch cfont.LpLogFont.LfWeight {
									case 0:
										// FW_DONTCARE
									case 100, 200, 300:
										customFont.Weight = WeightLight
									case 400:
										customFont.Weight = WeightRegular
									case 500:
										customFont.Weight = WeightMedium
									case 600:
										customFont.Weight = WeightDemibold
									case 700:
										customFont.Weight = WeightBold
									case 800, 900:
										customFont.Weight = WeightBlack
									}
									ComboBoxWeight.SetCurrentIndex(customFont.Weight)

									LineEditReplaceFont.SetText(customFont.ReplaceFont)
								},
							},
						},
					},

					Composite{
						Layout: Grid{Columns: 2, MarginsZero: true},
						Children: []Widget{
							Label{
								Text: "PixelSize",
							},
							NumberEdit{
								AssignTo:           &NumberEditPixelSize,
								Value:              customFont.Pixelsize,
								Suffix:             "px",
								MinValue:           0,
								SpinButtonsVisible: true,
								OnValueChanged: func() {
									customFont.Pixelsize = NumberEditPixelSize.Value()
								},
							},

							Label{
								Text: "Weight",
							},
							ComboBox{
								AssignTo:      &ComboBoxWeight,
								Model:         KnownWeight(),
								BindingMember: "Id",
								DisplayMember: "Name",
								OnCurrentIndexChanged: func() {
									customFont.Weight = ComboBoxWeight.CurrentIndex()
								},
							},

							Label{
								Text: "DPI",
							},
							NumberEdit{
								AssignTo:           &NumberEditDpi,
								Value:              float64(customFont.Dpi),
								MinValue:           0,
								SpinButtonsVisible: true,
								OnValueChanged: func() {
									customFont.Dpi = int(NumberEditDpi.Value())
								},
							},

							CheckBox{
								AssignTo:    &CheckBoxEmbeddedBitmap,
								Text:        "Embedded Bitmap",
								ToolTipText: "Bitmap fonts are sometimes used as fallbacks for missing fonts, which may cause text to be rendered pixelated or too large.", // TODO:
								OnCheckedChanged: func() {
									customFont.EmbeddedBitmap = CheckBoxEmbeddedBitmap.Checked()
								},
							},
							CheckBox{
								AssignTo: &CheckBoxPreferOutline,
								Text:     "Prefer Outline",
								OnCheckedChanged: func() {
									customFont.PreferOutline = CheckBoxPreferOutline.Checked()
								},
							},
							CheckBox{
								AssignTo: &CheckBoxDoSubstitutions,
								Text:     "Do Substitutions",
								OnCheckedChanged: func() {
									customFont.DoSubstitutions = CheckBoxDoSubstitutions.Checked()
								},
							},
							CheckBox{
								AssignTo: &CheckBoxBitmapMonospace,
								Text:     "Bitmap Monospace",
								OnCheckedChanged: func() {
									customFont.BitmapMonospace = CheckBoxBitmapMonospace.Checked()
								},
							},
							CheckBox{
								AssignTo: &CheckBoxForceAutohint,
								Text:     "Force Autohint",
								OnCheckedChanged: func() {
									customFont.ForceAutohint = CheckBoxForceAutohint.Checked()
								},
							},
							CheckBox{
								AssignTo: &CheckBoxQtUseSubpixelPositioning,
								Text:     "Qt Use Subpixel Positioning",
								OnCheckedChanged: func() {
									customFont.QtUseSubpixelPositioning = CheckBoxQtUseSubpixelPositioning.Checked()
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

	if customFont.ReplaceFont != "" {
		LineEditReplaceFont.SetText(customFont.ReplaceFont)
	}

	if customFont.Pixelsize != 0 {
		NumberEditPixelSize.SetValue(float64(customFont.Pixelsize))
	}

	ComboBoxWeight.SetCurrentIndex(customFont.Weight)

	if customFont.Dpi == 0 {
		NumberEditDpi.SetValue(float64(DPI))
	}

	if customFont.EmbeddedBitmap {
		CheckBoxEmbeddedBitmap.SetChecked(true)
	}
	if customFont.PreferOutline {
		CheckBoxPreferOutline.SetChecked(true)
	}
	if customFont.DoSubstitutions {
		CheckBoxDoSubstitutions.SetChecked(true)
	}
	if customFont.BitmapMonospace {
		CheckBoxBitmapMonospace.SetChecked(true)
	}
	if customFont.ForceAutohint {
		CheckBoxForceAutohint.SetChecked(true)
	}
	if customFont.QtUseSubpixelPositioning {
		CheckBoxQtUseSubpixelPositioning.SetChecked(true)
	}

	return dlg.Run()
}
