package main

// https://www.freedesktop.org/software/fontconfig/fontconfig-user.html#AEN134

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

func WriteFontsConf(c *Config) error {
	fontconfig := &Fontconfig{}

	fontconfig.Dir = []xDir{
		{Prefix: "cwd", Text: "../../csgo/panorama/fonts"},
		{Text: "WINDOWSFONTDIR"},
		{Text: "~/.fonts"},
		{Text: "/usr/share/fonts"},
		{Text: "/usr/local/share/fonts"},
		{Prefix: "xdg", Text: "fonts"},
	}

	// fontconfig.Fontpattern = []string{
	// 	"Arial",
	// 	".uifont",
	// 	"notosans",
	// 	"notoserif",
	// 	"notomono-regular",
	// }

	fontconfig.Cachedir = []string{
		"WINDOWSTEMPDIR_FONTCONFIG_CACHE",
		"~/.fontconfig",
	}

	targetMatch := xMatch{
		Target: "pattern",
		Edit: []xEdit{
			{
				Name:    "family",
				Mode:    "assign",
				Binding: "strong",
				String:  config.Font,
			},
		},
	}

	if config.Pixelsize != 1.0 {
		targetMatch.Edit = append(targetMatch.Edit, xEdit{
			Name: "pixelsize",
			Mode: "assign",
			Times: &xTimes{
				Name:   "pixelsize",
				Double: fmt.Sprintf("%g", config.Pixelsize),
			},
		})

	}

	fontconfig.Match = []xMatch{
		{
			Test: []xTest{
				{
					Name:   "family",
					String: "Stratum2 Bold Monodigit",
				},
			},
			Edit: []xEdit{
				{
					Name:    "family",
					Mode:    "append",
					Binding: "strong",
					String:  "Stratum2",
				},
				{
					Name:    "style",
					Mode:    "assign",
					Binding: "strong",
					String:  "Bold",
				},
			},
		},

		{
			Test: []xTest{
				{
					Name:   "family",
					String: "Stratum2 Regular Monodigit",
				},
			},
			Edit: []xEdit{
				{
					Name:    "family",
					Mode:    "append",
					Binding: "strong",
					String:  "Stratum2",
				},
				{
					Name:    "weight",
					Mode:    "assign",
					Binding: "strong",
					// String:  "Regular",
					Int: 90,
				},
			},
		},

		{
			Test: []xTest{
				{
					Name:   "lang",
					String: "vi-vn",
				},
				{
					Name:    "family",
					Compare: "contains",
					String:  "Stratum2",
				},
				{
					Qual:    "all",
					Name:    "family",
					Compare: "not_contains",
					String:  "TF",
				},
				{
					Qual:    "all",
					Name:    "family",
					Compare: "not_contains",
					String:  "Mono",
				},
				{
					Qual:    "all",
					Name:    "family",
					Compare: "not_contains",
					String:  "ForceStratum2",
				},
			},
			Edit: []xEdit{
				{
					Name: "weight",
					Mode: "assign",
					If: &xIf{
						Contains: &xContains{
							Name:   "family",
							String: "Stratum2 Black",
						},
						Int:  "210",
						Name: "weight",
					},
				},
				{
					Name: "slant",
					Mode: "assign",
					If: &xIf{
						Contains: &xContains{
							Name:   "family",
							String: "Italic",
						},
						Int:  "100",
						Name: "slant",
					},
				},
				{
					Name: "pixelsize",
					Mode: "assign",
					If: &xIf{
						Or: &xOr{
							Contains: xContains{
								Name:   "family",
								String: "Condensed",
							},
							LessEq: xLessEq{
								Name: "width",
								Int:  "75",
							},
						},
						Times: []xTimes{
							{
								Name:   "pixelsize",
								Double: "0.7",
							},
							{
								Name:   "pixelsize",
								Double: "0.9",
							},
						},
					},
				},
				{
					Name:    "family",
					Mode:    "assign",
					Binding: "same",
					String:  "notosans",
				},
			},
		},

		{
			Test: []xTest{
				{
					Name:   "lang",
					String: "vi-vn",
				},

				{
					Name:   "family",
					String: "ForceStratum2",
				},
			},
			Edit: []xEdit{
				{
					Name:    "family",
					Mode:    "assign",
					Binding: "same",
					String:  "Stratum2",
				},
			},
		},

		{
			Target: "font",
			Test: []xTest{
				{
					Name:    "family",
					Target:  "pattern",
					Compare: "contains",
					String:  "Stratum2",
				},
				{
					Name:    "family",
					Target:  "font",
					Compare: "contains",
					String:  "Arial",
				},
			},
			Edit: []xEdit{
				{
					Name: "pixelsize",
					Mode: "assign",
					Times: &xTimes{
						Name:   "pixelsize",
						Double: "0.9",
					},
				},
			},
		},

		{
			Target: "font",
			Test: []xTest{
				{
					Name:    "family",
					Target:  "pattern",
					Compare: "contains",
					String:  "Stratum2",
				},
				{
					Name:    "family",
					Target:  "font",
					Compare: "contains",
					String:  "Noto",
				},
			},
			Edit: []xEdit{
				{
					Name: "pixelsize",
					Mode: "assign",
					Times: &xTimes{
						Name:   "pixelsize",
						Double: "0.9",
					},
				},
			},
		},

		{
			Target: "scan",
			Test: []xTest{
				{
					Name:   "family",
					String: "Stratum2",
				},
			},
			Edit: []xEdit{
				{
					Name: "charset",
					Mode: "assign",
					Minus: &xMinus{
						Name: "charset",
						Charset: xChatset{
							Int: []string{
								"0x0394",
								"0x03A9",
								"0x03BC",
								"0x03C0",
								"0x2202",
								"0x2206",
								"0x220F",
								"0x2211",
								"0x221A",
								"0x221E",
								"0x222B",
								"0x2248",
								"0x2260",
								"0x2264",
								"0x2265",
								"0x25CA",
							},
						},
					},
				},
			},
		},

		{
			Target: "font",
			Edit: []xEdit{
				{
					Name: "embeddedbitmap",
					Mode: "assign",
					Bool: "false",
				},
			},
		},

		{
			Target: "pattern",
			Edit: []xEdit{
				{
					Name: "prefer_outline",
					Mode: "assign",
					Bool: "true",
				},
			},
		},

		{
			Target: "pattern",
			Edit: []xEdit{
				{
					Name: "do_substitutions",
					Mode: "assign",
					Bool: "true",
				},
			},
		},

		{
			Target: "pattern",
			Edit: []xEdit{
				{
					Name: "bitmap_monospace",
					Mode: "assign",
					Bool: "false",
				},
			},
		},

		{
			Target: "font",
			Edit: []xEdit{
				{
					Name: "force_autohint",
					Mode: "assign",
					Bool: "false",
				},
			},
		},

		{
			Target: "pattern",
			Edit: []xEdit{
				{
					Name:   "dpi",
					Mode:   "assign",
					Double: "96",
				},
			},
		},

		{
			Target: "pattern",
			Edit: []xEdit{
				{
					Name: "qt_use_subpixel_positioning",
					Mode: "assign",
					Bool: "false",
				},
			},
		},

		targetMatch,
	}

	fontconfig.Selectfont = &xSelectfont{
		Rejectfont: xRejectfont{
			Pattern: xPattern{
				Patelt: &xPatelt{
					Name:   "fontformat",
					String: "Type 1",
				},
			},
		},
	}

	fontconfig.Include = `../../../core/panorama/fonts/conf.d`

	bytes, err := xml.MarshalIndent(fontconfig, "", "\t")
	if err != nil {
		return err
	}

	csgofontsDir := filepath.Join(c.Path, "game", "csgo", "panorama", "fonts")
	if exist, _ := FileExists(csgofontsDir); !exist {
		os.MkdirAll(csgofontsDir, os.ModePerm)
	}
	err = os.WriteFile(filepath.Join(csgofontsDir, "fonts.conf"), append(header, bytes...), 0o644)
	if err == nil {
		return nil
	} else {
		return err
	}
}
