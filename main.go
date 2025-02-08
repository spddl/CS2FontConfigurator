package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/tailscale/walk"

	//lint:ignore ST1001 standard behavior tailscale/walk
	. "github.com/tailscale/walk/declarative"
)

type MyMainWindow struct {
	*walk.Dialog
	cb      *walk.ComboBox
	AppPath string
}

var (
	config    = new(Config)
	mw        *MyMainWindow
	fontslist []*Fontslist // Windows
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	config.Init()
}

type Fontslist struct {
	Id       int
	Dir      string
	Name     string
	Filename string
	Fontname string
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

	fontslist = []*Fontslist{}
	for _, dir := range []string{
		".",
		"C:\\Windows\\Fonts",
		filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Microsoft", "Windows", "Fonts"),
	} {

		entries, err := os.ReadDir(dir)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, e := range entries {
			filename := e.Name()

			if filepath.Ext(filename) != ".ttf" && filepath.Ext(filename) != ".otf" {
				continue
			}

			var fontname string = GetFontName(filepath.Join(dir, filename))
			if fontname == "" {
				continue
			}

			fontslist = append(fontslist, &Fontslist{
				Name:     fmt.Sprintf("%s (%s)", fontname, filename),
				Dir:      dir,
				Fontname: fontname,
				Filename: filename,
			})
		}
	}

	sort.Slice(fontslist, func(i, j int) bool {
		if fontslist[i].Fontname != fontslist[j].Fontname {
			return fontslist[i].Fontname < fontslist[j].Fontname
		}
		return fontslist[i].Filename < fontslist[j].Filename
	})

	// reduces fonts that are installed for the user and for the computer
	fontslist = slices.CompactFunc(fontslist, func(b *Fontslist, a *Fontslist) bool {
		return a.Name == b.Name
	})

	// creates the index for the combo box
	for i := range fontslist {
		fontslist[i].Id = i
	}

	var GroupBoxSimpleSelection *walk.GroupBox
	var appIcon, _ = walk.NewIconFromResourceId(7)
	var testfont *walk.Label
	var testfontinput *walk.LineEdit
	var defaultButton *walk.PushButton
	var PixelsizeEdit *walk.NumberEdit

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
					Columns: 2,
				},
				Title: "Simple font selection",
				Children: []Widget{

					Label{
						Text: "Font:",
					},

					ComboBox{
						AssignTo:      &mw.cb,
						BindingMember: "Id",
						DisplayMember: "Name",
						Model:         fontslist,
						OnCurrentIndexChanged: func() {
							font := fontslist[mw.cb.CurrentIndex()]

							config.Font = font.Fontname
							config.FontFile = font.Filename

							var style walk.FontStyle
							if strings.Contains(strings.ToLower(font.Name), "bold") {
								style |= walk.FontBold
							}
							if strings.Contains(strings.ToLower(font.Name), "italic") {
								style |= walk.FontItalic
							}
							walkfont, err := walk.NewFont(font.Fontname, 20, style)
							if err != nil {
								log.Println(err)
							}

							testfont.SetFont(walkfont)
						},
					},

					Label{
						Text: "Size (1 is Default):",
					},
					NumberEdit{
						AssignTo:           &PixelsizeEdit,
						MinValue:           0,
						Value:              config.Pixelsize,
						Decimals:           2,
						Increment:          0.1,
						SpinButtonsVisible: true,
						OnValueChanged: func() {
							config.Pixelsize = PixelsizeEdit.Value()
						},
					},

					Label{
						AssignTo:   &testfont,
						Text:       config.TestCase,
						ColumnSpan: 2,
					},

					LineEdit{
						AssignTo:   &testfontinput,
						ColumnSpan: 2,
						Text:       config.TestCase,
						OnKeyUp: func(key walk.Key) {
							testcase := testfontinput.Text()
							testfont.SetText(testcase)
							config.TestCase = testcase
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
						AssignTo: &defaultButton,
						Text:     "Default",
						Enabled:  false,
						OnClicked: func() {
							backupFile := filepath.Join(config.Path, "game", "csgo", "panorama", "fonts", "fonts_backup.conf")
							fontFile := filepath.Join(config.Path, "game", "csgo", "panorama", "fonts", "fonts.conf")

							if err := os.Remove(fontFile); err != nil {
								walk.MsgBox(mw, "Error", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
								panic(err)
							}

							if err := os.Rename(backupFile, fontFile); err != nil {
								walk.MsgBox(mw, "Error", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
								panic(err)
							}

							config.ClearUpFontsFolder()

							// Delete Cache
							os.RemoveAll(filepath.Join(os.Getenv("TEMP"), "fontconfig"))

							walk.MsgBox(mw, "FontConfig", "done.", walk.MsgBoxIconInformation)
						},
					},
					HSpacer{},
					PushButton{
						Text: "Apply",
						OnClicked: func() {
							// creates a backup of the font.conf if it does not already exist
							backupFile := filepath.Join(config.Path, "game", "csgo", "panorama", "fonts", "fonts_backup.conf")

							if exist, _ := FileExists(backupFile); !exist {
								fontFile := filepath.Join(config.Path, "game", "csgo", "panorama", "fonts", "fonts.conf")
								copyFile(fontFile, backupFile)
								defaultButton.SetEnabled(true)
							}

							config.ClearUpFontsFolder()

							font := fontslist[mw.cb.CurrentIndex()]
							copyFile(filepath.Join(font.Dir, font.Filename), filepath.Join(config.Path, "game", "csgo", "panorama", "fonts", font.Filename))

							if err := WriteFontsConf(config); err != nil {
								walk.MsgBox(mw, "Error", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
								panic(err)
							}

							if err := config.SaveConfigFile(mw.AppPath); err != nil {
								walk.MsgBox(mw, "Error", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
								panic(err)
							}

							// Delete Cache
							os.RemoveAll(filepath.Join(os.Getenv("TEMP"), "fontconfig"))

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

	// selects the font from the config
	if len(config.Font) == 0 {
		mw.cb.SetCurrentIndex(0)
	} else {
		for i, v := range fontslist {
			if v.Filename == config.FontFile {
				mw.cb.SetCurrentIndex(i)
				break
			}
		}
	}

	// Activates the default button if there is also a backup
	if exist, _ := FileExists(filepath.Join(config.Path, "game", "csgo", "panorama", "fonts", "fonts_backup.conf")); exist {
		defaultButton.SetEnabled(true)
	}

	mw.Run()
}
