package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"unsafe"

	"github.com/tailscale/walk"
	. "github.com/tailscale/walk/declarative"
	"github.com/tailscale/win"
	"golang.org/x/sys/windows"
)

type MyMainWindow struct {
	*walk.Dialog
	AppPath string
}

var (
	config = new(Config)
	mw     *MyMainWindow
	DPI    int
	model  *CustomFontModel
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	config.Init()
	model = NewCustomFontModel()
}

func main() {
	if config.Path == "" {
		CSGOPath, err := CSPath()
		if err != nil {
			walk.MsgBox(nil, "CS2 Path", err.Error(), walk.MsgBoxIconError)
			panic(err)
		}
		config.Path = CSGOPath
	}

	mw = new(MyMainWindow)
	mw.AppPath, _ = filepath.Abs("./")

	var ComboBoxWeight *walk.ComboBox
	var LabelFontPreview *walk.Label
	var NumberEditPixelSize *walk.NumberEdit
	var CheckBoxCustomize *walk.CheckBox
	var CompositeCustomize *walk.Composite
	var GroupBoxSimpleSelection *walk.GroupBox
	var tv *walk.TableView
	var appIcon, _ = walk.NewIconFromResourceId(7)

	if err := (Dialog{
		AssignTo: &mw.Dialog,
		Title:    "CS2FontConfigurator",
		Size:     Size{Width: 1, Height: 1},
		Icon:     appIcon,
		Layout:   Grid{Columns: 1},
		Children: []Widget{
			GroupBox{
				AssignTo: &GroupBoxSimpleSelection,
				Layout: Grid{
					Columns: 3,
				},
				Title: "Simple font selection",
				Children: []Widget{

					Label{
						Text: "Font:",
					},
					Label{
						AssignTo: &LabelFontPreview,
						Text:     "press Select",
					},
					PushButton{
						Text: "Select",
						OnClicked: func() {
							name := [32]uint16{}
							if config.Font != "" {
								copy(name[:], windows.StringToUTF16(config.Font))
							}

							// https://learn.microsoft.com/en-us/windows/win32/api/commdlg/ns-commdlg-choosefonta
							cfont := CHOOSEFONT{
								HwndOwner: mw.Handle(),
								LpLogFont: &win.LOGFONT{ // https://learn.microsoft.com/en-us/windows/win32/api/wingdi/ns-wingdi-logfontw
									LfFaceName: name,
									LfCharSet:  win.DEFAULT_CHARSET,
								},
								Flags:     CF_SCREENFONTS | CF_INITTOLOGFONTSTRUCT | CF_NOSCRIPTSEL | CF_FORCEFONTEXIST,
								NFontType: SCREEN_FONTTYPE,
							}
							cfont.LStructSize = uint32(unsafe.Sizeof(cfont))
							ChooseFontW(&cfont)

							// Replace Font
							config.Font = windows.UTF16ToString(cfont.LpLogFont.LfFaceName[:])

							// Pixelsize
							config.Pixelsize = float64(cfont.IPointSize / 10)
							if cfont.IPointSize != 0 {
								NumberEditPixelSize.SetValue(float64(cfont.IPointSize / 10))
							}

							// Weight
							switch cfont.LpLogFont.LfWeight {
							case 0:
								// FW_DONTCARE
							case 100, 200, 300:
								config.Weight = WeightLight
							case 400:
								config.Weight = WeightRegular
							case 500:
								config.Weight = WeightMedium
							case 600:
								config.Weight = WeightDemibold
							case 700:
								config.Weight = WeightBold
							case 800, 900:
								config.Weight = WeightBlack
							}
							ComboBoxWeight.SetCurrentIndex(config.Weight)

							LabelFontPreview.SetText(config.Font)

							config.Default()
							model.ResetRows()
						},
					},
					Label{
						Text: "PixelSize",
					},
					NumberEdit{
						AssignTo:           &NumberEditPixelSize,
						Value:              config.Pixelsize,
						ColumnSpan:         2,
						Suffix:             "px",
						MinValue:           0,
						SpinButtonsVisible: true,
						OnValueChanged: func() {
							config.Pixelsize = NumberEditPixelSize.Value()
						},
					},

					Label{
						Text: "Weight",
					},
					ComboBox{
						AssignTo:      &ComboBoxWeight,
						ColumnSpan:    2,
						Model:         KnownWeight(),
						BindingMember: "Id",
						DisplayMember: "Name",
						OnCurrentIndexChanged: func() {
							config.Weight = ComboBoxWeight.CurrentIndex()
						},
					},
				},
			},

			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					HSpacer{},
					PushButton{
						Text: "Options",
						OnClicked: func() {
							OpenOptions(mw)
						},
					},
				},
			},

			Composite{ // https://github.com/tailscale/walk/tree/55ec69fff1171242ca677d7f083a962dcc51c875/examples/tableview
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					CheckBox{
						AssignTo: &CheckBoxCustomize,
						Text:     "customize each font individually",
						OnCheckStateChanged: func() {
							config.Customize = CheckBoxCustomize.Checked()
							mw.Synchronize(func() {
								CompositeCustomize.SetVisible(config.Customize)
								if config.Customize {
									GroupBoxSimpleSelection.SetEnabled(false)
									config.Default()
									model.ResetRows()
								} else {
									GroupBoxSimpleSelection.SetEnabled(true)
									mw.Dialog.SetSize(walk.Size{Width: 0, Height: 0})
								}
							})
						},
					},
				},
			},

			Composite{
				AssignTo: &CompositeCustomize,
				Visible:  false,
				Layout: VBox{
					MarginsZero: true,
				},
				MinSize: Size{
					Width:  200,
					Height: 200,
				},
				Children: []Widget{

					TableView{
						AssignTo:            &tv,
						AlternatingRowBG:    true,
						ColumnsOrderable:    true,
						MultiSelection:      true,
						LastColumnStretched: true,
						Columns: []TableViewColumn{
							{
								Title: "ValveFont",
								Width: 150,
							},
							{
								DataMember: "ReplaceFont",
								Title:      "Replace",
								Width:      150,
							},
							{
								Title: "Weight",
								Width: 70,
								FormatFunc: func(value interface{}) string {
									switch value.(int) {
									case 0:
										return "Light"
									case 1:
										return "Regular"
									case 2:
										return "Medium"
									case 3:
										return "Demibold"
									case 4:
										return "Bold"
									case 5:
										return "Black"
									}
									return ""
								},
							},
							{
								DataMember: "Pixelsize",
								Title:      "PxSize",
								Width:      50,
								FormatFunc: func(value interface{}) string {
									return fmt.Sprintf("%d", int(value.(float64)))
								},
							},
							{
								DataMember: "Dpi",
								Title:      "DPI",
								Width:      50,
							},
							{
								DataMember: "PreferOutline",
								Title:      "Outline",
							},
							{
								Title: "DoSubstitutions",
							},
							{
								Title: "BitmapMonospace",
							},
							{
								Title: "ForceAutohint",
							},
							{
								Title: "QtUseSubpixelPositioning",
							},
							{
								Title: "EmbeddedBitmap",
							},
						},
						Model: model,
						OnItemActivated: func() {
							i := tv.CurrentIndex()
							newItem := &model.Items[i]
							orgItem := *newItem
							result := OpenCustomFontPtr(mw, newItem)

							if result == 0 || result == 2 { // cancel
								model.Items[i] = orgItem
								return
							} else {
								model.PublishRowsChanged(0, len(model.Items)-1)
								model.PublishRowsReset()
							}
						},
						ContextMenuItems: []MenuItem{
							Action{
								Text: "&Add",
								OnTriggered: func() {
									trackLatest := tv.ItemVisible(len(model.Items)-1) && len(tv.SelectedIndexes()) <= 1

									var newVal = &FontStruct{
										ValveFont:   "Valve Font",
										ReplaceFont: "Replace Font",
									}
									if ret := OpenCustomFontPtr(mw, newVal); ret == 1 {
										model.Items = append(model.Items, *newVal)
									}

									model.PublishRowsReset()

									if trackLatest {
										tv.EnsureItemVisible(len(model.Items) - 1)
									}
								},
							},
							Action{
								Text: "&Remove",
								OnTriggered: func() {
									indexes := tv.SelectedIndexes()
									for i := len(indexes) - 1; i >= 0; i-- {
										model.Items = append(model.Items[:indexes[i]], model.Items[indexes[i]+1:]...)

									}

									if len(indexes) != 0 {
										model.PublishRowsChanged(indexes[0], len(model.Items)-1)
										model.PublishRowsReset()
									}
								},
							},
						},
					},
				},
			},

			Composite{
				Layout: HBox{
					MarginsZero: true,
				},
				Children: []Widget{
					PushButton{
						Text: "Preview",
						OnClicked: func() {
							// create demo.html
							createPreview(mw.AppPath)
							// https://stackoverflow.com/a/12076082
							exec.Command("rundll32.exe", []string{"url.dll,FileProtocolHandler", filepath.Join(mw.AppPath, "demo.htm")}...).Start()

						},
					},

					PushButton{
						Text: "Default",
						OnClicked: func() {
							def := new(Config)
							def.newConfig()
							def.Path = config.Path
							config = def

							if err := WriteFontsConf(config); err != nil {
								walk.MsgBox(mw, "Error", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
								panic(err)
							}

							if err := ReplGlobalConf(config); err != nil {
								walk.MsgBox(mw, "Error", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
								panic(err)
							}

							walk.MsgBox(mw, "FontConfig", "done.", walk.MsgBoxIconInformation)
						},
					},
					HSpacer{},
					PushButton{
						Text: "Apply",
						OnClicked: func() {
							if !config.Customize {
								config.Default()
							}

							if err := WriteFontsConf(config); err != nil {
								walk.MsgBox(mw, "Error", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
								panic(err)
							}

							if err := ReplGlobalConf(config); err != nil {
								walk.MsgBox(mw, "Error", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
								panic(err)
							}

							if err := config.SaveConfigFile(mw.AppPath); err != nil {
								walk.MsgBox(mw, "Error", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
								panic(err)
							}

							walk.MsgBox(mw, "FontConfig", "done.", walk.MsgBoxIconInformation)
						},
					},
					PushButton{
						Text: "Exit",
						OnClicked: func() {
							os.Exit(0)
						},
					},
				},
			},
		},
	}.Create(nil)); err != nil {
		panic(err)
	}

	DPI = mw.DPI()

	if config.Font != "" {
		LabelFontPreview.SetText(config.Font)
	}

	if config.Pixelsize != 0 {
		NumberEditPixelSize.SetValue(float64(config.Pixelsize))
	}

	if config.Weight != 0 {
		ComboBoxWeight.SetCurrentIndex(config.Weight)
	}
	if config.Customize && len(config.Fonts) != 0 {
		// model = NewCustomFontModel()
		CheckBoxCustomize.SetChecked(true)
	}

	mw.Run() // https://github.com/lxn/walk/issues/103#issuecomment-278243090

	// cleanup
	os.Remove(filepath.Join(mw.AppPath, "demo.htm"))
}

type CustomFontModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	Items      []FontStruct
}

func NewCustomFontModel() *CustomFontModel {
	m := new(CustomFontModel)
	m.ResetRows()
	return m
}

// Called by the TableView from SetModel and every time the model publishes a
// RowsReset event.
func (m *CustomFontModel) RowCount() int {
	return len(m.Items)
}

// Called by the TableView when it needs the text to display for a given cell.
func (m *CustomFontModel) Value(row, col int) interface{} {
	item := m.Items[row]

	switch col {
	case 0:
		return item.ValveFont
	case 1:
		return item.ReplaceFont
	case 2:
		return item.Weight
	case 3:
		return item.Pixelsize
	case 4:
		return item.Dpi
	case 5:
		return item.PreferOutline
	case 6:
		return item.DoSubstitutions
	case 7:
		return item.BitmapMonospace
	case 8:
		return item.ForceAutohint
	case 9:
		return item.QtUseSubpixelPositioning
	case 10:
		return item.EmbeddedBitmap
	}

	panic("unexpected col")
}

// Called by the TableView to sort the model.
func (m *CustomFontModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order

	sort.SliceStable(m.Items, func(i, j int) bool {
		a, b := m.Items[i], m.Items[j]
		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}

			return !ls
		}

		switch m.sortColumn {
		case 0:
			return c(a.ValveFont < b.ValveFont)
		case 1:
			return c(a.ReplaceFont < b.ReplaceFont)
		case 2:
			return c(a.Weight < b.Weight)
		case 3:
			return c(a.Pixelsize < b.Pixelsize)
		case 4:
			return c(a.Dpi < b.Dpi)
		case 5:
			return a.PreferOutline != b.PreferOutline
		case 6:
			return a.DoSubstitutions != b.DoSubstitutions
		case 7:
			return a.BitmapMonospace != b.BitmapMonospace
		case 8:
			return a.ForceAutohint != b.ForceAutohint
		case 9:
			return a.QtUseSubpixelPositioning != b.QtUseSubpixelPositioning
		case 10:
			return a.EmbeddedBitmap != b.EmbeddedBitmap

		}
		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)
}

func (m *CustomFontModel) ResetRows() {
	m.Items = config.Fonts

	// Notify TableView and other interested parties about the reset.
	m.PublishRowsReset()

	m.Sort(m.sortColumn, m.sortOrder)
}
