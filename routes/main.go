package routes

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"git.sr.ht/~humaid/reviewforge/models"
	"git.sr.ht/~humaid/reviewforge/models/forms"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/go-git/go-git/v5"
	git_http "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-macaron/binding"
	"github.com/go-macaron/session"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	gm_html "github.com/yuin/goldmark/renderer/html"
	macaron "gopkg.in/macaron.v1"
)

func ContextInit() macaron.Handler {
	return func(ctx *macaron.Context) {
		ctx.Data["SiteTitle"] = "reviewforge"
		ctx.Data["User"] = "humaid"
	}
}

func DashboardHandler(ctx *macaron.Context) {
	ctx.Data["Title"] = "Dashboard"
	repos, err := models.GetRepositories()
	if err != nil {
		panic(err)
	}
	ctx.Data["Repos"] = repos
	ctx.HTML(http.StatusOK, "index")
}

func AddRepoHandler(ctx *macaron.Context) {
	ctx.Data["Title"] = "Add a new repository"
	ctx.HTML(http.StatusOK, "add_repo")
}

type DirEntry struct {
	Mode, Name, Complexity, Issues, Size string
	IsDir                                bool
}

func RepoHandler(ctx *macaron.Context, f *session.Flash) {
	repo, err := models.GetRepository(ctx.ParamsInt64("id"))
	if err != nil {
		ctx.Redirect("/")
		return
	}
	if !strings.HasSuffix(ctx.Req.RequestURI, "/") {
		ctx.Redirect(ctx.Req.RequestURI + "/")
		return
	}
	var entries []DirEntry

	ctx.Data["Repo"] = repo
	var files []os.FileInfo
	if len(ctx.Params("*")) > 0 {
		// directory
		file := "./repos/" + repo.Name + "/" + path.Clean(ctx.Params("*"))
		st, err := os.Stat(file)
		if err != nil {
			ctx.Redirect(fmt.Sprintf("/%d", repo.ID))
			return
		}
		if st.IsDir() {
			ctx.Data["Path"] = "/" + path.Clean(ctx.Params("*"))
			files, err = ioutil.ReadDir(file)
			entries = append(entries, DirEntry{
				Mode:       st.Mode().String(),
				Name:       "..",
				Complexity: "",
				Issues:     "",
				Size:       fmt.Sprintf("%d", st.Size()),
				IsDir:      true,
			})
		} else {
			// We are viewing a file...
			contents, err := ioutil.ReadFile(file)
			if err != nil {
				panic(err)
			}
			lexer := lexers.Match(file)
			if lexer == nil {
				lexer = lexers.Fallback
			}
			lexer = chroma.Coalesce(lexer)
			style := styles.Get("pygments")
			if style == nil {
				style = styles.Fallback
			}
			formatter := html.New(html.WithClasses(true))
			iterator, err := lexer.Tokenise(nil, string(contents))
			if err != nil {
				panic(err)
			}

			var w strings.Builder
			err = formatter.Format(&w, style, iterator)
			ctx.Data["File"] = template.HTML(w.String())
			var css strings.Builder
			err = formatter.WriteCSS(&css, style)
			ctx.Data["CSS"] = template.CSS(css.String())

			ctx.HTML(http.StatusOK, "file")
			return
		}
	} else {
		ctx.Data["Path"] = "/"
		files, err = ioutil.ReadDir("./repos/" + repo.Name)
	}
	if err != nil {
		f.Error("Failed to load repository")
		log.Println("Cannot load repo: ", err)
		ctx.Redirect("/")
		return
	}

	for _, f := range files {
		if f.Name() == ".git" {
			continue
		}
		entries = append(entries, DirEntry{
			Mode:       f.Mode().String(),
			Name:       f.Name(),
			Complexity: "0",
			Issues:     "0",
			Size:       fmt.Sprintf("%d", f.Size()),
			IsDir:      f.IsDir(),
		})
	}
	ctx.Data["Files"] = entries

	var readme string
	if len(ctx.Params("*")) > 0 {
		readme = "./repos/" + repo.Name + "/" + path.Clean(ctx.Params("*")) + "/README.md"
	} else {
		readme = "./repos/" + repo.Name + "/README.md"
	}
	readStat, err := os.Stat(readme)
	if err == nil && !readStat.IsDir() {
		md := goldmark.New(
			goldmark.WithExtensions(extension.GFM),
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
			goldmark.WithRendererOptions(
				gm_html.WithXHTML(),
				gm_html.WithUnsafe(),
			),
		)
		readme, err := ioutil.ReadFile(readme)
		if err == nil {
			var buf bytes.Buffer
			if err := md.Convert([]byte(readme), &buf); err == nil {
				ctx.Data["Readme"] = template.HTML(string(bluemonday.UGCPolicy().SanitizeBytes(buf.Bytes())))
			}
		}
	}

	ctx.HTML(http.StatusOK, "repo")
}

func AddRepoPostHandler(ctx *macaron.Context, form forms.AddRepositoryForm,
	errs binding.Errors, f *session.Flash) {
	if len(errs) > 0 {
		f.Error("Make sure to fill the required fields.")
		ctx.Redirect("/add")
		return
	}
	// TODO check if exists first
	repo := models.Repository{
		Name:      form.Name,
		GitRemote: form.GitRepoURL,
	}

	cloneOpts := git.CloneOptions{
		URL:      form.GitRepoURL,
		Depth:    1, // We only fetch latest commit
		Progress: os.Stdout,
	}

	if len(form.Username) > 0 {
		cloneOpts.Auth = &git_http.BasicAuth{
			Username: form.Username,
			Password: form.Password,
		}
	}

	_, err := git.PlainClone("./repos/"+form.Name, false, &cloneOpts)
	if err != nil {
		f.Error("Failed to clone. Git returned: " + err.Error())
		ctx.Redirect("/")
		return
	}

	models.AddRepository(&repo)
	ctx.Redirect(fmt.Sprintf("/%d", repo.ID))
}
