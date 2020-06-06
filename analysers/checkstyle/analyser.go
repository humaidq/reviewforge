package dependency

import (
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"git.sr.ht/~humaid/reviewforge/analysers"
)

type CheckstyleTool struct {
	JarPath    string
	ConfigPath string
}

func (c *CheckstyleTool) HasTool() bool {
	return true
}

var (
	// 0: Full, 1: [WARN] part, 2: Path, 3: Line number, 4: Column, 5: Message.
	errorExp = regexp.MustCompile(`\[([A-Z]+)\] ([a-zA-Z0-9 _.\/]+):(\d+)(?>:(\d+))?: (.+)`)
)

func parseCheckstyleOutput(output string, projPath string) (issues []analysers.Issue) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		matches := errorExp.FindAllString(line, -1)
		var issue analysers.Issue
		sp := strings.Split(matches[2], projPath)
		if len(sp) < 2 {
			log.Println("checkstyle parse: cannot get proj path from: ", matches[2])
			continue
		}
		issue.FilePath = sp[1]
		if len(matches[3]) > 0 {
			ui, err := strconv.ParseUint(matches[3], 10, 0)
			if err != nil {
				issue.LineNumber = ui
			}
		}
		if len(matches[4]) > 0 {
			ui, err := strconv.ParseUint(matches[4], 10, 0)
			if err != nil {
				issue.ColumnNumber = ui
			}
		}

	}

	return
}

func (c *CheckstyleTool) Run(path string) ([]analysers.Issue, error) {
	o, err := exec.Command("java", "-jar", c.JarPath, "-c", c.ConfigPath, path).Output()
	if err != nil {
		return []analysers.Issue{}, err
	}
	res := parseCheckstyleOutput(string(o), path)
	return res, nil
}

func (c *CheckstyleTool) GetInfo() analysers.AnalyserInfo {
	return analysers.AnalyserInfo{
		Name:    "checkstyle",
		Version: "8.33",
		URL:     "https://checkstyle.sourceforge.io/",
		LanguagesSupported: []analysers.ProgrammingLanguage{
			analysers.Java,
		},
	}
}
