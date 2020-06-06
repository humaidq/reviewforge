package analysers

type Analyser interface {
	GetInfo() AnalyserInfo
	HasTool() bool
	Run(path string) ([]Issue, error)
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
	LineNumber   uint64
	ColumnNumber uint64
	CheckName    string
	Category     string // Optional
	Description  string
	CVE          string // Optional, for exploits
	Serverity    int    // Optional, for exploits TODO make this an enum
}
