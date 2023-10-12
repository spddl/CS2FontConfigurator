package main

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

func createPreview(appPath string) {
	const tpl = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>demo</title>
</head>
<body>
</body>
<script>
	function add(varclass) {
		const newDiv = document.createElement("pre")
		newDiv.style.fontFamily = varclass[1]
		newDiv.appendChild(document.createTextNode(varclass[0] + " => " + varclass[1] + ":\n1234567890 abcdefghijklmnopqrstuvwxyz ABCDEFGHIJKLMNOPQRSTUVWXYZ"))
		document.body.appendChild(newDiv)
	}
	[{{ range $_, $value := .Fonts }}
		["{{ $value.ValveFont }}", "{{ $value.ReplaceFont }}"],{{ end }}
	].forEach(ele => {
		add(ele)
	})
</script>
</html>`

	t := template.Must(template.New("preview").Parse(tpl))
	if !config.Customize {
		config.Default()
	}

	data := struct {
		Fonts []FontStruct
	}{
		Fonts: config.Fonts,
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := os.WriteFile(filepath.Join(appPath, "demo.htm"), buf.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}
}
