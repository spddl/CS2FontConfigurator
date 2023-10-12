package main

import "encoding/xml"

var header = []byte("<?xml version='1.0'?>\n<!DOCTYPE fontconfig SYSTEM 'fonts.dtd'>\n")

type Fontconfig struct {
	XMLName     xml.Name     `xml:"fontconfig"`
	Text        string       `xml:",chardata"`
	Dir         []xDir       `xml:"dir"`
	Fontpattern []string     `xml:"fontpattern"`
	Cachedir    []string     `xml:"cachedir"`
	Match       []xMatch     `xml:"match"`
	Selectfont  *xSelectfont `xml:"selectfont"`
	Include     string       `xml:"include,omitempty"`
}

type xSelectfont struct {
	Text       string      `xml:",chardata"`
	Rejectfont xRejectfont `xml:"rejectfont"`
}

type xRejectfont struct {
	Text    string   `xml:",chardata"`
	Pattern xPattern `xml:"pattern"`
}

type xPattern struct {
	Text   string   `xml:",chardata"`
	Patelt *xPatelt `xml:"patelt"`
}

type xPatelt struct {
	Text   string `xml:",chardata"`
	Name   string `xml:"name,attr"`
	String string `xml:"string"`
}

type xDir struct {
	Text   string `xml:",chardata"`
	Prefix string `xml:"prefix,attr,omitempty"`
}

type xMatch struct {
	Text   string  `xml:",chardata"`
	Target string  `xml:"target,attr,omitempty"`
	Test   []xTest `xml:"test"`
	Edit   []xEdit `xml:"edit"`
}

type xTest struct {
	Text    string `xml:",chardata"`
	Qual    string `xml:"qual,attr,omitempty"`
	Name    string `xml:"name,attr,omitempty"`
	Target  string `xml:"target,attr,omitempty"`
	Compare string `xml:"compare,attr,omitempty"`
	String  string `xml:"string"`
}

type xEdit struct {
	Text    string  `xml:",chardata"`
	Name    string  `xml:"name,attr"`
	Mode    string  `xml:"mode,attr"`
	Binding string  `xml:"binding,attr,omitempty"`
	Const   string  `xml:"const,omitempty"`
	String  string  `xml:"string,omitempty"`
	If      *xIf    `xml:"if"`
	Times   *xTimes `xml:"times"`
	Minus   *xMinus `xml:"minus"`
	Bool    string  `xml:"bool,omitempty"`
	Double  string  `xml:"double,omitempty"`
}

type xTimes struct {
	Text   string `xml:",chardata"`
	Name   string `xml:"name,omitempty"`
	Double string `xml:"double,omitempty"`
}

type xMinus struct {
	Text    string   `xml:",chardata"`
	Name    string   `xml:"name"`
	Charset xChatset `xml:"charset"`
}

type xChatset struct {
	Text string   `xml:",chardata"`
	Int  []string `xml:"int"`
}

type xIf struct {
	Text     string     `xml:",chardata"`
	Contains *xContains `xml:"contains,omitempty"`
	Int      string     `xml:"int,omitempty"`
	Name     string     `xml:"name,omitempty"`
	Or       *xOr       `xml:"or,omitempty"`
	Times    []xTimes   `xml:"times,omitempty"`
}

type xContains struct {
	Text   string `xml:",chardata"`
	Name   string `xml:"name"`
	String string `xml:"string"`
}

type xOr struct {
	Text     string    `xml:",chardata"`
	Contains xContains `xml:"contains"`
	LessEq   xLessEq   `xml:"less_eq"`
}

type xLessEq struct {
	Text string `xml:",chardata"`
	Name string `xml:"name"`
	Int  string `xml:"int"`
}
