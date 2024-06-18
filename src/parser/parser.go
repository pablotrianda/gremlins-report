package parser

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type GremlinReport struct {
	GoModule          string  `json:"go_module"`
	Files             []File  `json:"files"`
	TestEfficacy      int     `json:"test_efficacy"`
	MutationsCoverage float64 `json:"mutations_coverage"`
	MutantsTotal      int     `json:"mutants_total"`
	MutantsKilled     int     `json:"mutants_killed"`
	MutantsLived      int     `json:"mutants_lived"`
	MutantsNotViable  int     `json:"mutants_not_viable"`
	MutantsNotCovered int     `json:"mutants_not_covered"`
	ElapsedTime       float64 `json:"elapsed_time"`
}

type File struct {
	FileName  string     `json:"file_name"`
	Mutations []Mutation `json:"mutations"`
}

type Mutation struct {
	Type      string `json:"type"`
	Status    string `json:"status"`
	Line      int    `json:"line"`
	Column    int    `json:"column"`
	Code      string
	ClassName string
	Url       string
}

// Return a json with the gremlins output report
func LoadGremlinsOutput(outputGremlin string) GremlinReport {
	jsonOput, err := os.ReadFile(outputGremlin)
	if err != nil {
		log.Panic(err)
	}

	var repo GremlinReport

	err = json.Unmarshal(jsonOput, &repo)
	if err != nil {
		log.Println(err)
	}

	loadFileExtraSections(&repo)
	return repo
}

func loadFileExtraSections(report *GremlinReport) {
	for i, f := range report.Files {
		for j, mut := range f.Mutations {
			report.Files[i].Mutations[j].Code = extractLineCode(f.FileName, mut.Line)
			report.Files[i].Mutations[j].ClassName = strings.ToLower(strings.Replace(mut.Status, " ", "_", -1))
			report.Files[i].Mutations[j].Url = fmt.Sprintf("https://gremlins.dev/latest/usage/mutations/%s", strings.ToLower(mut.Type))
		}
	}
}

func extractLineCode(fileName string, line int) string {
	command := fmt.Sprintf("sed -n %d,%dp %s", line, line, fileName)
	lineCode, _ := exec.Command("bash", "-c", command).CombinedOutput()
	return string(lineCode)
}
