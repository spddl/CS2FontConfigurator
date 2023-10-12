package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

func ReplGlobalConf(c *Config) error {
	fontconfig := &Fontconfig{}
	for _, fs := range config.Fonts {
		var font = xMatch{
			Target: "font",
			Test: []xTest{
				{
					Name:   "family",
					String: fs.ValveFont,
				},
			},
			Edit: []xEdit{
				{
					Name:   "family",
					Mode:   "assign",
					String: fs.ReplaceFont,
				},
			},
		}
		var pattern = xMatch{
			Target: "pattern",
			Test: []xTest{
				{
					Name:   "family",
					String: fs.ValveFont,
				},
			},
			Edit: []xEdit{
				{
					Name:    "family",
					Mode:    "prepend",
					Binding: "strong",
					String:  fs.ReplaceFont,
				},
			},
		}

		if fs.Pixelsize != 0 {
			var pixelsizeEdit = &xEdit{
				Name: "pixelsize",
				Mode: "assign",
				Times: &xTimes{
					Name:   "pixelsize",
					Double: floatToStr(fs.Pixelsize / 10),
				},
			}

			font.Edit = append(font.Edit,
				*pixelsizeEdit,
			)
			pattern.Edit = append(pattern.Edit,
				*pixelsizeEdit,
			)
			// <edit name="pixelsize" mode="assign">
			// 	<times>
			// 		<name>pixelsize</name>
			// 		<double>0.9</double>
			// 	</times>
			// </edit>
		}

		if fs.Weight != -1 {
			var weight string
			switch fs.Weight {
			case 0:
				weight = "Light"
			case 1:
				weight = "Regular"
			case 2:
				weight = "Medium"
			case 3:
				weight = "Demibold"
			case 4:
				weight = "Bold"
			case 5:
				weight = "Black"
			}

			var weightEdit = &xEdit{
				Name:    "weight",
				Mode:    "assign",
				Binding: "strong",
				String:  weight,
			}

			font.Edit = append(font.Edit,
				*weightEdit,
			)
			pattern.Edit = append(pattern.Edit,
				*weightEdit,
			)
			// 	<edit name="weight" mode="assign" binding="strong">
			// 		<string>Regular</string> //  Light, medium, demibold, bold or black
			// 	</edit>
		}

		if fs.Dpi != -1 {
			var dpiEdit = &xEdit{
				Name:   "dpi",
				Mode:   "assign",
				Double: intToStr(fs.Dpi),
			}

			font.Edit = append(font.Edit,
				*dpiEdit,
			)
			pattern.Edit = append(pattern.Edit,
				*dpiEdit,
			)
			// <!-- Set DPI.  dpi should be set in ~/.Xresources to 96 -->
			// <!-- Setting to 72 here makes the px to pt conversions work better (Chrome) -->
			// <!-- Some may need to set this to 96 though -->
			// <match target="pattern">
			// 	<edit name="dpi" mode="assign">
			// 		<double>96</double>
			// 	</edit>
			// </match>
		}

		if fs.EmbeddedBitmap != config.EmbeddedBitmap {
			var EmbeddedBitmapEdit = &xEdit{
				Name: "embeddedbitmap",
				Mode: "assign",
				Bool: boolToStr(fs.EmbeddedBitmap),
			}

			font.Edit = append(font.Edit,
				*EmbeddedBitmapEdit,
			)
			pattern.Edit = append(pattern.Edit,
				*EmbeddedBitmapEdit,
			)
			// <!-- Globally use embedded bitmaps in fonts like Calibri? -->
			// <match target="font" >
			// 	<edit name="embeddedbitmap" mode="assign">
			// 		<bool>false</bool>
			// 	</edit>
			// </match>
		}

		if fs.PreferOutline != config.PreferOutline {
			var PreferOutlineEdit = &xEdit{
				Name: "prefer_outline",
				Mode: "assign",
				Bool: boolToStr(fs.PreferOutline),
			}

			font.Edit = append(font.Edit,
				*PreferOutlineEdit,
			)
			pattern.Edit = append(pattern.Edit,
				*PreferOutlineEdit,
			)
			// <!-- Substitute truetype fonts in place of bitmap ones? -->
			// <match target="pattern" >
			// 	<edit name="prefer_outline" mode="assign">
			// 		<bool>true</bool>
			// 	</edit>
			// </match>
		}

		if fs.DoSubstitutions != config.DoSubstitutions {
			var DoSubstitutionsEdit = &xEdit{
				Name: "do_substitutions",
				Mode: "assign",
				Bool: boolToStr(fs.DoSubstitutions),
			}

			font.Edit = append(font.Edit,
				*DoSubstitutionsEdit,
			)
			pattern.Edit = append(pattern.Edit,
				*DoSubstitutionsEdit,
			)
			// <!-- Do font substitutions for the set style? -->
			// <!-- NOTE: Custom substitutions in 42-repl-global.conf will still be done -->
			// <!-- NOTE: Corrective substitutions will still be done -->
			// <match target="pattern" >
			// 	<edit name="do_substitutions" mode="assign">
			// 		<bool>true</bool>
			// 	</edit>
			// </match>
		}

		if fs.BitmapMonospace != config.BitmapMonospace {
			var BitmapMonospaceEdit = &xEdit{
				Name: "bitmap_monospace",
				Mode: "assign",
				Bool: boolToStr(fs.BitmapMonospace),
			}

			font.Edit = append(font.Edit,
				*BitmapMonospaceEdit,
			)
			pattern.Edit = append(pattern.Edit,
				*BitmapMonospaceEdit,
			)
			// <!-- Make (some) monospace/coding TTF fonts render as bitmaps? -->
			// <!-- courier new, andale mono, monaco, etc. -->
			// <match target="pattern" >
			// 	<edit name="bitmap_monospace" mode="assign">
			// 		<bool>false</bool>
			// 	</edit>
			// </match>
		}

		if fs.ForceAutohint != config.ForceAutohint {
			var ForceAutohintEdit = &xEdit{
				Name: "force_autohint",
				Mode: "assign",
				Bool: boolToStr(fs.ForceAutohint),
			}

			font.Edit = append(font.Edit,
				*ForceAutohintEdit,
			)
			pattern.Edit = append(pattern.Edit,
				*ForceAutohintEdit,
			)
			// <!-- Force autohint always -->
			// <!-- Useful for debugging and for free software purists -->
			// <match target="font">
			// 	<edit name="force_autohint" mode="assign">
			// 		<bool>false</bool>
			// 	</edit>
			// </match>
		}

		if fs.QtUseSubpixelPositioning != config.QtUseSubpixelPositioning {
			var QtUseSubpixelPositioningEdit = &xEdit{
				Name: "qt_use_subpixel_positioning",
				Mode: "assign",
				Bool: boolToStr(fs.QtUseSubpixelPositioning),
			}

			font.Edit = append(font.Edit,
				*QtUseSubpixelPositioningEdit,
			)
			pattern.Edit = append(pattern.Edit,
				*QtUseSubpixelPositioningEdit,
			)
			// <!-- Use Qt subpixel positioning on autohinted fonts? -->
			// <!-- This only applies to Qt and autohinted fonts. Qt determines subpixel positioning based on hintslight vs. hintfull, -->
			// <!--   however infinality patches force slight hinting inside freetype, so this essentially just fakes out Qt. -->
			// <!-- Should only be set to true if you are not doing any stem alignment or fitting in environment variables -->
			// <match target="pattern" >
			// 	<edit name="qt_use_subpixel_positioning" mode="assign">
			// 		<bool>false</bool>
			// 	</edit>
			// </match>
		}

		fontconfig.Match = append(fontconfig.Match,
			font,
			pattern,
		)
	}

	bytes, err := xml.MarshalIndent(fontconfig, "", "\t")
	if err != nil {
		return err
	}

	confdDir := filepath.Join(c.Path, "game", "core", "panorama", "fonts", "conf.d")
	if exist, _ := FileExists(confdDir); !exist {
		os.MkdirAll(confdDir, os.ModePerm)
	}
	err = os.WriteFile(filepath.Join(confdDir, "42-repl-global.conf"), append(header, bytes...), 0644)
	if err == nil {
		fmt.Println(filepath.Join(confdDir, "42-repl-global.conf"), "saved")
		return nil
	} else {
		return err
	}
}
