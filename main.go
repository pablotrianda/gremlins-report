package main

import (
	"log"
	"os"
	"os/exec"
	"strings"

	parser "github.com/pablotrianda/gremlinks-report/src/parser"
	"github.com/pablotrianda/gremlinks-report/src/render"
)

func main() {
	output, err := runGremlins(os.Args)
	if err != nil {
		log.Println(err)
	}

	gremlinsReport := parser.GremlinsOutput(output)

	// jsonData, err := json.Marshal(gremlinsReport)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(string(jsonData))

	templateFile, err := os.ReadFile("./template.html")
	if err != nil {
		panic(err)
	}
	render.BuildHTMLReport(templateFile, gremlinsReport)
}

func runGremlins(args []string) (string, error) {
	gremlinsParams := strings.Join(args[1:], " ")
	outputGremlins, err := exec.Command("bash", "-c", "gremlins "+gremlinsParams).CombinedOutput()
	if err != nil {
		return "", nil
	}

	return string(outputGremlins), nil
}
