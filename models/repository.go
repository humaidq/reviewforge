package models

import "errors"

// Repository represents a repository being audited
type Repository struct {
	RepositoryID int64 `xorm:"pk autoincr"`
	Name         string
	GitRemote    string
	OwnerUserID  string
	Auditors     []string
	Assignment   map[string]string
}

func GetRepository(id int64) (*Repository, error) {
	r := new(Repository)
	has, err := engine.ID(id).Get(r)
	if err != nil {
		return r, err
	} else if !has {
		return r, errors.New("Repository does not exist")
	}
	return r, nil
}

func AddRepository(r *Repository) (err error) {
	_, err = engine.Insert(r)
	return
}

func GetRepositories() (repos []Repository, err error) {
	err = engine.Find(&repos)
	return
}
