package main

import (
	"io/ioutil"
	"math"
	"os"
	"path/filepath"

	"./cfg"
	"./fontconfig"
	"./helper"
	"github.com/golang/freetype/truetype"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type MyMainWindow struct {
	*walk.MainWindow
	prevFilePath string
}

func main() {
	Config := cfg.Init()

	fontsDir := filepath.Join(os.Getenv("windir"), "Fonts")
	fontList, err := helper.ReadDir(fontsDir, ".ttf")
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
	var LineEditPath *walk.LineEdit
	var ComboBoxFont *walk.ComboBox
	var NumberEditPixelSize *walk.NumberEdit
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
						// AssignTo: &PushButtonFileExplorer,
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

			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text: "PixelSize:",
					},
					NumberEdit{
						Enabled:  false,
						AssignTo: &NumberEditPixelSize,
						Value:    Config.PixelSize,
						MinValue: 0,
						MaxValue: math.Inf(+1),
						Suffix:   " pt",
						Decimals: 2,
						OnValueChanged: func() {
							Config.PixelSize = NumberEditPixelSize.Value()
						},
					},
				},
			},

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
						Value:              "Font",
						Model:              fontList,
						OnCurrentIndexChanged: func() {
							Config.File = ComboBoxFont.Text()
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
							err = fontconfig.WriteSpddlConf(def)
							if err != nil {
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

							fontfile := filepath.Join(fontsDir, Config.File)
							data, err := ioutil.ReadFile(fontfile)
							if err != nil {
								walk.MsgBox(mw, "ErrorReadFile", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
							}
							f, err := truetype.Parse(data)
							if err != nil {
								walk.MsgBox(mw, "ErrorTrueTypeParse", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
							}
							Config.Font = f.Name(truetype.NameIDFontFullName)

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

							// Save Config FIXME: wrong save location
							err = cfg.SaveConfigFile(Config)
							if err != nil {
								walk.MsgBox(mw, "Error", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
								panic(err)
							}

							walk.MsgBox(mw, "FontConfig", "done.", walk.MsgBoxIconInformation)
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
	mw.Run() // https://github.com/lxn/walk/issues/103#issuecomment-278243090

	LineEditPath.SetText(Config.Path)
	ComboBoxFont.SetText(Config.File)
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
