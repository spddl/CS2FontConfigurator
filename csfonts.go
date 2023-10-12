package main

func (c *Config) Default() {
	var csfont = []string{
		"Courier New",
		"Stratum2 Black Condensed",
		"Stratum2 Black Italic",
		"Stratum2 Black TF",
		"Stratum2 Black",
		"Stratum2 Bold Condensed",
		"Stratum2 Bold Monodigit",
		"Stratum2 Bold TF",
		"Stratum2 Condensed",
		"Stratum2 Light Condensed",
		"Stratum2 Light Italic",
		"Stratum2 Light TF",
		"Stratum2 Light",
		"Stratum2 Medium Condensed",
		"Stratum2 Medium Italic",
		"Stratum2 Medium TF",
		"Stratum2 Medium",
		"Stratum2 Mono Light",
		"Stratum2 Mono",
		"Stratum2 Regular Monodigit",
		"Stratum2 TF",
		"Stratum2 Thin Condensed",
		"Stratum2 Thin TF",
		"Stratum2",
		"Times New Roman",
		"noto",
		"notomono-regular",
		"notosans",
		"notoserif",
	}

	if len(c.Fonts) == 0 {
		for _, defFont := range csfont {
			c.Fonts = append(c.Fonts, FontStruct{
				ValveFont:                defFont,
				ReplaceFont:              c.Font,
				Pixelsize:                c.Pixelsize,
				Weight:                   c.Weight,
				EmbeddedBitmap:           c.EmbeddedBitmap,
				PreferOutline:            c.PreferOutline,
				DoSubstitutions:          c.DoSubstitutions,
				BitmapMonospace:          c.BitmapMonospace,
				ForceAutohint:            c.ForceAutohint,
				Dpi:                      c.Dpi,
				QtUseSubpixelPositioning: c.QtUseSubpixelPositioning,
			})
		}
	} else {
		for i := 0; i < len(c.Fonts); i++ {
			c.Fonts[i] = FontStruct{
				ValveFont:                c.Fonts[i].ValveFont,
				ReplaceFont:              c.Font,
				Pixelsize:                c.Pixelsize,
				Weight:                   c.Weight,
				EmbeddedBitmap:           c.EmbeddedBitmap,
				PreferOutline:            c.PreferOutline,
				DoSubstitutions:          c.DoSubstitutions,
				BitmapMonospace:          c.BitmapMonospace,
				ForceAutohint:            c.ForceAutohint,
				Dpi:                      c.Dpi,
				QtUseSubpixelPositioning: c.QtUseSubpixelPositioning,
			}
		}
	}
}
