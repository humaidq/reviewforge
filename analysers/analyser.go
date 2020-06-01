package analysers

type Analyser interface {
	GetInfo() AnalyserInfo
	HasTool() bool
	Run(path string) []Issue
}

type ProgrammingLanguage int

const (
	Any = iota + 1
	Go
	C
	Java
	JavaScript
	PHP
	Python
	Rust
	CSharp
	CPP
	Swift
	ObjC
	Ruby
	Perl
	D
	FSharp
	COBOL
	Delphi
	Ada
	Erlang
	Haskell
	Kotlin
	Shell
	TypeScript
)

type AnalyserInfo struct {
	Name, Version      string
	URL                string
	LanguagesSupported []ProgrammingLanguage
}

type Issue struct {
	FilePath     string
	LineNumber   uint
	ColumnNumber uint
	CheckName    string
	Description  string
	CVE          string
	Serverity    int
}
