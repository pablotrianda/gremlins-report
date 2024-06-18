package main

import (
	"embed"
	"log"
	"os"
	"os/exec"
	"strings"

	parser "github.com/pablotrianda/gremlinks-report/src/parser"
	"github.com/pablotrianda/gremlinks-report/src/render"
)

//go:embed template.html
var f embed.FS

func main() {
	jsonOutput, err := runGremlins(os.Args)
	if err != nil {
		log.Println(err)
	}

	gremlinsReport := parser.LoadGremlinsOutput(jsonOutput)

	templateFile, err := f.ReadFile("template.html")
	if err != nil {
		panic(err)
	}
	render.BuildHTMLReport(templateFile, gremlinsReport)
}

func runGremlins(args []string) (string, error) {
	gremlinsParams := strings.Join(args[1:], " ")
	tmpFile, err := os.CreateTemp("", "gremlins_output_*.json")
	if err != nil {
		return "", err
	}
	tmpFile.Close()
	command := "gremlins " + gremlinsParams + " --output=" + tmpFile.Name()
	exec.Command("bash", "-c", command).Run()

	return tmpFile.Name(), nil
}
