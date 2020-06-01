package dependency

import "git.sr.ht/~humaid/reviewforge/analysers"

type DependencyCheckTool struct {
	ToolName string
}

func (d *DependencyCheckTool) GetInfo() analysers.AnalyserInfo {
	return analysers.AnalyserInfo{
		Name:    "Dependency Check",
		Version: "v5.3.2",
		URL:     "https://github.com/jeremylong/DependencyCheck",
	}
}
