package parser

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type GremlinReport struct {
	Header    string     `json:"header"`
	Mutations []Mutation `json:"mutations"`
	Resume    string     `json:"resume"`
}

type Mutation struct {
	Name         string   `json:"name"`
	MutationType string   `json:"mutation_type"`
	File         Location `json:"file"`
}

type Location struct {
	Name     string `json:"file"`
	Line     string `json:"line"`
	Position string `json:"position"`
	Code     string `json:"code"`
}

// Return a json with the gremlins output report
func GremlinsOutput(outputGremlin string) GremlinReport {
	lines := strings.Split(outputGremlin, "\n")
	header := getHeader(lines)
	mutations := getMutations(lines)
	resume := getResume(lines)

	return GremlinReport{
		Header:    header,
		Mutations: mutations,
		Resume:    resume,
	}
}

func getHeader(lines []string) string {
	header := ""
	for _, l := range lines {
		if strings.Contains(l, "Gathering") {
			header += l
		}
		if !strings.Contains(header, "done in") && strings.Contains(l, "done in") {
			header += " " + l
		}

		// To avoid loop over all string when already have the complete header
		if strings.Contains(header, "Gathering") && strings.Contains(header, "done in") {
			break
		}
	}

	return header
}

func getMutations(lines []string) []Mutation {
	mutations := []Mutation{}
	for _, e := range lines {
		ok, mut := inMutationList(e)
		if ok {
			mutationsType := extractMutationType(e)
			location := extractLocation(e)
			mutation := Mutation{Name: mut, MutationType: mutationsType, File: location}
			mutations = append(mutations, mutation)
		}
	}
	return mutations
}

func inMutationList(str string) (bool, string) {
	mutations := []string{"RUNNABLE",
		"NOT COVERED",
		"KILLED",
		"LIVED",
		"TIMED OUT",
		"NOT VIABLE",
	}

	for _, s := range mutations {
		re := regexp.MustCompile(`\b` + s + `\b`)
		if re.MatchString(str) {
			return true, s
		}
	}

	return false, ""
}

func extractMutationType(str string) string {
	pattern := `\b[A-Z_]+\b`
	regex := regexp.MustCompile(pattern)
	found := regex.FindAllString(str, -1)

	if len(found) == 2 {
		return found[1]
	}

	return found[2]
}

func extractLocation(str string) Location {
	fileName := extractFileName(str)
	fileNameData := strings.Split(fileName, ":")
	lineCode := extractLineCode(fileNameData[0], fileNameData[1])

	return Location{
		Name:     fileName,
		Line:     fileNameData[1],
		Position: fileNameData[2],
		Code:     lineCode,
	}
}

func extractLineCode(fileName, line string) string {
	lineCode, _ := exec.Command("bash", "-c", "sed -n "+line+","+line+"p "+fileName).CombinedOutput()
	fmt.Println(string(lineCode))

	return string(lineCode)
}

func extractFileName(str string) string {
	return strings.Split(str, "at")[1][1:]
}

func getResume(lines []string) string {
	//resume := ""
	// resume += lines[len(lines)-6] + " \n"
	// resume += lines[len(lines)-5] + " \n"
	// resume += lines[len(lines)-4] + " \n"
	// resume += lines[len(lines)-3] + " \n"
	// resume += lines[len(lines)-2] + " \n"

	return ""
}
