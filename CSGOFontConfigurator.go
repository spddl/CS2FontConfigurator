package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"

	"./cfg"
	"./fontconfig"
	"./helper"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type Reader interface {
	Read(p []byte) (n int, err error)
}

type MyMainWindow struct {
	*walk.MainWindow
	prevFilePath string
	AppPath      string
}

func main() {
	Config := cfg.Init()
	fontsDir := filepath.Join(strings.Replace(os.Getenv("windir"), "WINDOWS", "Windows", 1), "Fonts") // TODO: IDK
	fontList, err := helper.ReadDir(fontsDir)
	if err != nil {
		panic(err)
	}

	if Config.Path == "" {
		CSGOPath, err := helper.CSGOPath()
		if err != nil {
			panic(err)
		}
		Config.Path = CSGOPath
	}

	// FIXME: walk.ValidationErrorEffect, _ = walk.NewBorderGlowEffect(walk.RGB(255, 0, 0))

	mw := new(MyMainWindow)
	mw.AppPath, _ = filepath.Abs("./")
	var LineEditPath *walk.LineEdit
	var ComboBoxFont *walk.ComboBox
	var CheckBoxEmbeddedBitmap,
		CheckBoxPreferOutline,
		CheckBoxDoSubstitutions,
		CheckBoxBitmapMonospace,
		CheckBoxForceAutohint,
		CheckBoxQtUseSubpixelPositioning *walk.CheckBox

	var NumberEditdpi *walk.NumberEdit
	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "CSGOFontConfigurator",
		MinSize:  Size{400, 180},
		Layout:   VBox{},
		Children: []Widget{
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text: "Path:",
					},
					LineEdit{
						AssignTo: &LineEditPath,
						Text:     Config.Path,
						OnKeyDown: func(key walk.Key) {
							Config.Path = LineEditPath.Text()
						},
					},
					PushButton{
						Text: "📁",
						OnClicked: func() {
							path, err := mw.openFileExplorer()
							if err != nil {
								panic(err)
							}
							Config.Path = path
							LineEditPath.SetText(path)
						},
					},
				},
			},

			// Composite{ TODO:
			// 	Layout: HBox{},
			// 	Children: []Widget{
			// 		Label{
			// 			Text: "PixelSize:",
			// 		},
			// 		NumberEdit{
			// 			Enabled:  false,
			// 			AssignTo: &NumberEditPixelSize,
			// 			Value:    Config.PixelSize,
			// 			MinValue: 0,
			// 			MaxValue: math.Inf(+1),
			// 			Suffix:   " pt",
			// 			Decimals: 2,
			// 			OnValueChanged: func() {
			// 				Config.PixelSize = NumberEditPixelSize.Value()
			// 			},
			// 		},
			// 	},
			// },

			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text: "Font:",
					},
					ComboBox{
						AlwaysConsumeSpace: true,
						StretchFactor:      50,
						AssignTo:           &ComboBoxFont,
						Model:              fontList,
						OnCurrentIndexChanged: func() {
							fmt.Println("OnCurrentIndexChanged", ComboBoxFont.CurrentIndex())
							Config.File = ComboBoxFont.Text()
						},
					},
				},
			},

			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text:        "Embedded Bitmap:",
						ToolTipText: "Bitmap fonts are sometimes used as fallbacks for missing fonts, which may cause text to be rendered pixelated or too large.", // TODO:
					},
					CheckBox{
						AssignTo:    &CheckBoxEmbeddedBitmap,
						Checked:     Config.EmbeddedBitmap,
						ToolTipText: "Default: false",
						OnCheckStateChanged: func() {
							Config.EmbeddedBitmap = CheckBoxEmbeddedBitmap.Checked()
						},
					},

					Label{
						Text:        "Prefer Outline:",
						ToolTipText: "", // TODO:
					},
					CheckBox{
						AssignTo:    &CheckBoxPreferOutline,
						Checked:     Config.PreferOutline,
						ToolTipText: "Default: true",
						OnCheckStateChanged: func() {
							Config.PreferOutline = CheckBoxPreferOutline.Checked()
						},
					},

					Label{
						Text:        "Do Substitutions:",
						ToolTipText: "", // TODO:
					},
					CheckBox{
						AssignTo:    &CheckBoxDoSubstitutions,
						Checked:     Config.DoSubstitutions,
						ToolTipText: "Default: true",
						OnCheckStateChanged: func() {
							Config.DoSubstitutions = CheckBoxDoSubstitutions.Checked()
						},
					},

					Label{
						Text:        "Bitmap Monospace:",
						ToolTipText: "", // TODO:
					},
					CheckBox{
						AssignTo:    &CheckBoxBitmapMonospace,
						Checked:     Config.BitmapMonospace,
						ToolTipText: "Default: false",
						OnCheckStateChanged: func() {
							Config.BitmapMonospace = CheckBoxBitmapMonospace.Checked()
						},
					},

					Label{
						Text:        "Force Autohint:",
						ToolTipText: "", // TODO:
					},
					CheckBox{
						AssignTo:    &CheckBoxForceAutohint,
						Checked:     Config.ForceAutohint,
						ToolTipText: "Default: false",
						OnCheckStateChanged: func() {
							Config.ForceAutohint = CheckBoxForceAutohint.Checked()
						},
					},

					Label{
						Text:        "Qt Use Subpixel Positioning:",
						ToolTipText: "", // TODO:
					},
					CheckBox{
						AssignTo:    &CheckBoxQtUseSubpixelPositioning,
						Checked:     Config.QtUseSubpixelPositioning,
						ToolTipText: "Default: false",
						OnCheckStateChanged: func() {
							Config.QtUseSubpixelPositioning = CheckBoxQtUseSubpixelPositioning.Checked()
						},
					},

					Label{
						Text:        "dpi:",
						ToolTipText: "", // TODO:
					},
					NumberEdit{
						AssignTo:    &NumberEditdpi,
						Value:       Config.Dpi,
						MinValue:    0,
						MaxValue:    math.Inf(+1),
						Decimals:    0,
						ToolTipText: "Default: 96",
						OnValueChanged: func() {
							Config.Dpi = NumberEditdpi.Value()
						},
					},
				},
			},

			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						Text: "Default",
						OnClicked: func() {
							def := cfg.NewConfig()
							def.Path = Config.Path

							err = fontconfig.WriteFontsConf(fontsDir, def)
							if err != nil {
								walk.MsgBox(mw, "Error", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
								panic(err)
							}

							spddlconfpath := filepath.Join(filepath.Dir(def.Path), "csgo", "panorama", "fonts", "conf.d", "99-spddl.conf")
							exists, err := helper.FileExists(spddlconfpath)
							if err != nil {
								walk.MsgBox(mw, "Error", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
								panic(err)
							}
							if exists {
								err := os.Remove(spddlconfpath)
								if err != nil {
									walk.MsgBox(mw, "Error", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
									panic(err)
								}
							}

							walk.MsgBox(mw, "FontConfig", "done.", walk.MsgBoxIconInformation)
						},
					},
					HSpacer{},
					PushButton{
						Text: "Apply",
						OnClicked: func() {
							fontfile := filepath.Join(fontsDir, Config.File)
							data, err := ioutil.ReadFile(fontfile)
							if err != nil {
								walk.MsgBox(mw, "ErrorReadFile", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
							}

							fontData, err := helper.NewFontData(Config.File, data)
							if err != nil {
								walk.MsgBox(mw, "Error Parse", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
							}
							// fmt.Println("Name", fontData.Name)         // Name AdobeHebrew-Italic
							// fmt.Println("Family", fontData.Family)     // Family Adobe Hebrew
							// fmt.Println("FileName", fontData.FileName) // FileName AdobeHebrew-Italic.otf
							Config.Font = fontData.Name

							//
							// https://www.freedesktop.org/software/fontconfig/fontconfig-user.html#AEN134
							//
							err = fontconfig.WriteFontsConf(fontsDir, Config)
							if err != nil {
								walk.MsgBox(mw, "Error", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
								panic(err)
							}
							err = fontconfig.WriteSpddlConf(Config)
							if err != nil {
								walk.MsgBox(mw, "Error", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
								panic(err)
							}

							err = cfg.SaveConfigFile(mw.AppPath, Config)
							if err != nil {
								walk.MsgBox(mw, "Error", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
								panic(err)
							}

							walk.MsgBox(mw, "FontConfig", "done.", walk.MsgBoxIconInformation)
							os.Exit(0)
						},
					},
					PushButton{
						Text: "Exit",
						OnClicked: func() {
							os.Exit(1)
						},
					},
				},
			},
		},
	}.Create()); err != nil {
		panic(err)
	}

	if Config.File != "" {
		i := helper.IndexOf(Config.File, fontList)
		fmt.Println("Config.File", Config.File, i)
		err = ComboBoxFont.SetCurrentIndex(i)
		if err != nil {
			panic(err)
		}
	}

	mw.Run() // https://github.com/lxn/walk/issues/103#issuecomment-278243090
}

func (mw *MyMainWindow) openFileExplorer() (FilePath string, err error) {
	dlg := new(walk.FileDialog)

	dlg.Title = "Find csgo.exe"
	dlg.FilePath = mw.prevFilePath

	dlg.Filter = "csgo.exe"

	if ok, err := dlg.ShowOpen(mw); err != nil {
		return "", err
	} else if !ok {
		return "", nil
	}

	mw.prevFilePath = dlg.FilePath

	return dlg.FilePath, nil
}
