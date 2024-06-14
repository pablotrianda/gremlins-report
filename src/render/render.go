package render

import (
	"os"
	"os/exec"
	"text/template"

	"github.com/pablotrianda/gremlinks-report/src/parser"
)

func BuildHTMLReport(templateFile []byte, report parser.GremlinReport) {
	tmpl, err := template.New("report").Parse(string(templateFile))

	if err != nil {
		panic(err)
	}

	tempFile, err := os.CreateTemp("", "reporte_*.html")
	if err != nil {
		panic(err)
	}
	defer tempFile.Close()

	err = tmpl.Execute(tempFile, report)
	if err != nil {
		panic(err)
	}

	tempFilePath := tempFile.Name()

	cmd := exec.Command("xdg-open", tempFilePath)
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
}
