package forms

// AddRepositoryForm represents a form to add a new repository.
type AddRepositoryForm struct {
	Name       string `form:"name" binding:"Required"`
	GitRepoURL string `form:"repo" binding:"Required"`
	Username   string `form:"user"`
	Password   string `form:"pass"`
}
