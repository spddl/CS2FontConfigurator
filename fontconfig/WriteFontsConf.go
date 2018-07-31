package fontconfig

// https://www.freedesktop.org/software/fontconfig/fontconfig-user.html#AEN134

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"../cfg"
	"../helper"
)

func WriteFontsConf(fontsDir string, c cfg.Config) error {
	fontconfig := &Sfontconfig{}
	fontconfig.Sdir = []Sdir{
		{Value: "WINDOWSFONTDIR"},
		{Value: "~/.fonts"},
		{Value: "/usr/share/fonts"},
		{Value: "/usr/local/share/fonts"},
		{Attrprefix: "xdg", Value: "fonts"},
	}
	if fontsDir != "" {
		fontconfig.Sdir = append(fontconfig.Sdir, Sdir{Value: fontsDir})
	}

	fontconfig.Sfontpattern = []Sfontpattern{
		{Value: "Arial"},
		{Value: ".vfont"},
		{Value: "notosans"},
		{Value: "notoserif"},
		{Value: "notomono-regular"},
	}
	if c.Font != "" {
		file, ext := helper.FileAndExt(c.File)
		fontconfig.Sfontpattern = append(fontconfig.Sfontpattern, Sfontpattern{Value: c.Font})
		fontconfig.Sfontpattern = append(fontconfig.Sfontpattern, Sfontpattern{Value: ext})
		fontconfig.Sfontpattern = append(fontconfig.Sfontpattern, Sfontpattern{Value: file})
	}
	fontconfig.Scachedir = []Scachedir{
		{Value: "WINDOWSTEMPDIR_FONTCONFIG_CACHE"},
		{Value: "~/.fontconfig"},
	}
	fontconfig.Smatch = []Smatch{
		{
			Stest: &Stest{
				Attrname: "family",
				Sstring: &Sstring{
					Value: "Stratum2 Bold Monodigit",
				},
			},
			Sedit: &Sedit{
				Attrname:    "family",
				Attrmode:    "append",
				Attrbinding: "strong",
				Sstring: &Sstring{
					Value: "Stratum2 Bold",
				},
			},
		},
		{
			Stest: &Stest{
				Attrname: "family",
				Sstring: &Sstring{
					Value: "Stratum2 Regular Monodigit",
				},
			},
			Sedit: &Sedit{
				Attrname:    "family",
				Attrmode:    "append",
				Attrbinding: "strong",
				Sstring: &Sstring{
					Value: "Stratum2",
				},
			},
		},
		{
			Attrtarget: "font",
			Sedit: &Sedit{
				Attrname: "embeddedbitmap",
				Attrmode: "assign",
				Sbool: &Sbool{
					Value: "false",
				},
			},
		},
		{
			Attrtarget: "pattern",
			Sedit: &Sedit{
				Attrname: "prefer_outline",
				Attrmode: "assign",
				Sbool: &Sbool{
					Value: "true",
				},
			},
		},
		{
			Attrtarget: "pattern",
			Sedit: &Sedit{
				Attrname: "do_substitutions",
				Attrmode: "assign",
				Sbool: &Sbool{
					Value: "true",
				},
			},
		},
		{
			Attrtarget: "pattern",
			Sedit: &Sedit{
				Attrname: "bitmap_monospace",
				Attrmode: "assign",
				Sbool: &Sbool{
					Value: "false",
				},
			},
		},
		{
			Attrtarget: "font",
			Sedit: &Sedit{
				Attrname: "force_autohint",
				Attrmode: "assign",
				Sbool: &Sbool{
					Value: "false",
				},
			},
		},
		{
			Attrtarget: "pattern",
			Sedit: &Sedit{
				Attrname: "force_autohint",
				Attrmode: "assign",
				Sbool: &Sbool{
					Value: "false",
				},
			},
		},
		{
			Attrtarget: "pattern",
			Sedit: &Sedit{
				Attrname: "dpi",
				Attrmode: "assign",
				Sdouble: &Sdouble{
					Value: "96",
				},
			},
		},
		{
			Attrtarget: "pattern",
			Sedit: &Sedit{
				Attrname: "qt_use_subpixel_positioning",
				Attrmode: "assign",
				Sbool: &Sbool{
					Value: "false",
				},
			},
		},
	}

	fontconfig.Sselectfont = &Sselectfont{
		Srejectfont: &Srejectfont{
			Spattern: &Spattern{
				Spatelt: &Spatelt{
					Attrname: "fontformat",
					Sstring: &Sstring{
						Value: "Type 1",
					},
				},
			},
		},
	}

	fontconfig.Sinclude = &Sinclude{
		Value: "conf.d",
	}

	bytes, err := xml.Marshal(fontconfig)
	if err != nil {
		return err
	}

	data := append(header, bytes...)
	err = ioutil.WriteFile(filepath.Join(filepath.Dir(c.Path), "csgo", "panorama", "fonts", "fonts.conf"), data, 0644)
	if err == nil {
		fmt.Println("fonts.conf saved")
		return nil
	} else {
		return err
	}
}
