package fontconfig

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"../cfg"
)

func WriteSpddlConf(c cfg.Config) error {
	// file, _ := helper.FileAndExt(c.Font)

	fontconfig := &Sfontconfig{}
	fontconfig.Smatch = []Smatch{
		{
			Attrtarget: "font",
			Sedit: &Sedit{
				Attrname: "hintstyle",
				Attrmode: "assign",
				Sconst: &Sconst{
					Value: "hintnone",
				},
			},
		},

		{
			Attrtarget: "font",
			Stest: &Stest{
				Attrname: "family",
				Sstring: &Sstring{
					Value: "Stratum2",
				},
			},
			Sedit: &Sedit{
				Attrname: "family",
				Attrmode: "assign",
				Sstring: &Sstring{
					Value: c.Font,
				},
			},
		},

		{
			Attrtarget: "pattern",
			Stest: &Stest{
				Attrname: "family",
				Sstring: &Sstring{
					Value: "Stratum2",
				},
			},
			Sedit: &Sedit{
				Attrname:    "family",
				Attrmode:    "prepend",
				Attrbinding: "strong",
				Sstring: &Sstring{
					Value: c.Font,
				},
			},
		},
		{
			Attrtarget: "font",
			Stest: &Stest{
				Attrname: "family",
				Sstring: &Sstring{
					Value: "Stratum2 Bold",
				},
			},
			Sedit: &Sedit{
				Attrname: "family",
				Attrmode: "assign",
				Sstring: &Sstring{
					Value: c.Font,
				},
			},
		},
		{
			Attrtarget: "pattern",
			Stest: &Stest{
				Attrname: "family",
				Sstring: &Sstring{
					Value: "Stratum2 Bold",
				},
			},
			Sedit: &Sedit{
				Attrname:    "family",
				Attrmode:    "prepend",
				Attrbinding: "strong",
				Sstring: &Sstring{
					Value: c.Font,
				},
			},
		},
		{
			Attrtarget: "font",
			Stest: &Stest{
				Attrname: "family",
				Sstring: &Sstring{
					Value: "Arial",
				},
			},
			Sedit: &Sedit{
				Attrname: "family",
				Attrmode: "assign",
				Sstring: &Sstring{
					Value: c.Font,
				},
			},
		},
		{
			Attrtarget: "font",
			Stest: &Stest{
				Attrname: "family",
				Sstring: &Sstring{
					Value: "Arial",
				},
			},
			Sedit: &Sedit{
				Attrname:    "family",
				Attrmode:    "prepend",
				Attrbinding: "strong",
				Sstring: &Sstring{
					Value: c.Font,
				},
			},
		},
	}

	bytes, err := xml.Marshal(fontconfig)
	if err != nil {
		return err
	}

	data := append(header, bytes...)
	err = ioutil.WriteFile(filepath.Join(filepath.Dir(c.Path), "csgo", "panorama", "fonts", "conf.d", "99-spddl.conf"), data, 0644)
	if err == nil {
		fmt.Println("99-spddl.conf saved")
		return nil
	} else {
		return err
	}
}
