package models

import "errors"

// Report represents a report
type Report struct {
	ReportID      int64 `xorm:"pk autoincr"`
	RepositoryID  int64
	ToolGenerated string
}

type Issue struct {
	IssueID      int64 `xorm:"pk autoincr"`
	ReportID     int64
	FilePath     string
	LineNumber   uint
	ColumnNumber uint
	CheckName    string
	Description  string
	CVE          string
	Serverity    int
}

func GetReport(id int64) (*Report, error) {
	r := new(Report)
	has, err := engine.ID(id).Get(r)
	if err != nil {
		return r, err
	} else if !has {
		return r, errors.New("Report does not exist")
	}
	return r, nil
}

func AddReport(r *Report) (err error) {
	_, err = engine.Insert(r)
	return
}

func GetReports() (report []Report, err error) {
	err = engine.Find(&report)
	return
}
