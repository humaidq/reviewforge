package models

type Repository struct {
	ID          int64 `xorm:"pk autoincr"`
	Name        string
	GitRemote   string
	OwnerUserID string
	Auditors    []string
}
