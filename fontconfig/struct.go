package fontconfig

import "encoding/xml"

var header = []byte("<?xml version='1.0'?><!DOCTYPE fontconfig SYSTEM 'fonts.dtd'>")

// TODO:
// https://stackoverflow.com/a/27246583
// https://stackoverflow.com/questions/37504533/marshalling-xml-in-golang-field-is-empty-append-doesnt-work

type Sfontconfig struct {
	XMLName      xml.Name `xml:"fontconfig"`
	Sdir         []Sdir
	Sfontpattern []Sfontpattern
	Scachedir    []Scachedir
	Smatch       []Smatch
	Sselectfont  *Sselectfont `xml:"selectfont"`
	Sinclude     *Sinclude    `xml:"include"`
}

type Sfontpattern struct {
	XMLName xml.Name `xml:"fontpattern,omitempty"`
	Value   string   `xml:",chardata"`
}

type Sdir struct {
	XMLName    xml.Name `xml:"dir,omitempty"`
	Attrprefix string   `xml:"prefix,attr,omitempty"`
	Value      string   `xml:",chardata"`
}

type Scachedir struct {
	XMLName xml.Name `xml:"cachedir,omitempty"`
	Value   string   `xml:",chardata"`
}

type Smatch struct {
	XMLName    xml.Name `xml:"match,omitempty"`
	Attrtarget string   `xml:"target,attr,omitempty"`
	Stest      *Stest
	Sedit      *Sedit
}

type Stest struct {
	XMLName  xml.Name `xml:"test,omitempty"`
	Attrname string   `xml:"name,attr,omitempty"`
	Sstring  *Sstring `xml:"string,omitempty"`
}

type Sedit struct {
	XMLName     xml.Name `xml:"edit,omitempty"`
	Attrbinding string   `xml:"binding,attr,omitempty"`
	Attrmode    string   `xml:"mode,attr"`
	Attrname    string   `xml:"name,attr"`
	Sconst      *Sconst  `xml:"const,omitempty"`
	Sbool       *Sbool   `xml:"bool,omitempty"`
	Sdouble     *Sdouble `xml:"double,omitempty"`
	Sstring     *Sstring `xml:"string,omitempty"`
	Stimes      *Stimes  `xml:"times,omitempty"`
}

type Sbool struct {
	XMLName xml.Name `xml:"bool,omitempty"`
	Value   bool     `xml:",chardata"`
}

type Sdouble struct {
	XMLName xml.Name `xml:"double,omitempty"`
	Value   float64  `xml:",chardata"`
}

type Sstring struct {
	XMLName xml.Name `xml:"string,omitempty"`
	Value   string   `xml:",chardata"`
}

type Sselectfont struct {
	XMLName     xml.Name     `xml:"selectfont,omitempty"`
	Srejectfont *Srejectfont `xml:"rejectfont,omitempty"`
}

type Srejectfont struct {
	XMLName  xml.Name  `xml:"rejectfont,omitempty"`
	Spattern *Spattern `xml:"pattern,omitempty"`
}

type Spattern struct {
	XMLName xml.Name `xml:"pattern,omitempty"`
	Spatelt *Spatelt `xml:"patelt,omitempty"`
}

type Spatelt struct {
	XMLName  xml.Name `xml:"patelt,omitempty"`
	Attrname string   `xml:"name,attr"`
	Sstring  *Sstring `xml:"string,omitempty"`
}

type Sinclude struct {
	Value string `xml:",chardata"`
}

type Sconst struct {
	XMLName xml.Name `xml:"const,omitempty"`
	Value   string   `xml:",chardata"`
}

type Sname struct {
	XMLName xml.Name `xml:"name,omitempty"`
	Value   string   `xml:",chardata"`
}

type Stimes struct {
	XMLName xml.Name `xml:"times,omitempty"`
	Sdouble *Sdouble `xml:"double,omitempty"`
	Sname   *Sname   `xml:"name,omitempty"`
}
