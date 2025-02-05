package main

import "testing"

var compareTests = []struct {
	have string
	want string
}{
	{"C:\\Windows\\Fonts\\OperatorMono-Bold.otf", "Operator Mono"},
	{"C:\\Windows\\Fonts\\OperatorMono-BoldItalic.otf", "Operator Mono"},
	{"C:\\Users\\spddl\\AppData\\Local\\Microsoft\\Windows\\Fonts\\FZCuYuan-M03.ttf", "FZCuYuan-M03"},
	{"C:\\Windows\\Fonts\\DavidCLM-Medium.otf", "David CLM"},
	{"C:\\Windows\\Fonts\\DejaVuSans-ExtraLight.ttf", "DejaVu Sans Light"},
	{"C:\\Windows\\Fonts\\FiraCode-Bold.ttf", "Fira Code"},
	{"C:\\Windows\\Fonts\\FrankRuehlCLM-Bold.ttf", "Frank Ruehl CLM"},
	{"C:\\Windows\\Fonts\\HelveticaNeueBoldItalic.otf", "Helvetica Neue"},
	{"C:\\Windows\\Fonts\\impact.ttf", "Impact"},
	{"C:\\Windows\\Fonts\\lucon.ttf", "Lucida Console"},
	{"C:\\Windows\\Fonts\\msyi.ttf", "Microsoft Yi Baiti"},
	{"C:\\Windows\\Fonts\\Rubik-BoldItalic.ttf", "Rubik"},
	{"C:\\Windows\\Fonts\\OperatorMonoLig-LightItalic.otf", "Operator Mono Lig"},
}

func TestAbs(t *testing.T) {
	for _, test := range compareTests {
		got := GetFontName(test.have)
		if got != test.want {
			t.Errorf("%s != %s; %s", got, test.want, test.have)
		}
	}
}
