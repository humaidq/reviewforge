package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-macaron/binding"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"github.com/urfave/cli/v2"
	macaron "gopkg.in/macaron.v1"

	"git.sr.ht/~humaid/reviewforge/analysers"
	"git.sr.ht/~humaid/reviewforge/analysers/checkstyle"
	"git.sr.ht/~humaid/reviewforge/models"
	"git.sr.ht/~humaid/reviewforge/models/forms"
	"git.sr.ht/~humaid/reviewforge/routes"
)

// CmdStart represents a command-line command
// which starts the code review system
var CmdStart = &cli.Command{
	Name:    "start",
	Aliases: []string{"run"},
	Usage:   "Start the code review system's server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "port",
			Value: "8080",
			Usage: "the web server port",
		},
		&cli.BoolFlag{
			Name:  "dev",
			Value: false,
			Usage: "enables development mode (for templates)",
		},
	},
	Action: start,
}

func getMacaron(dev bool) *macaron.Macaron {
	m := macaron.Classic()
	m.Use(macaron.Renderer())
	m.Use(session.Sessioner())
	m.Use(csrf.Csrfer())
	m.Use(routes.ContextInit())

	m.Get("/", routes.DashboardHandler)
	m.Get("/new", routes.AddRepoHandler)
	m.Post("/new", binding.BindIgnErr(forms.AddRepositoryForm{}), routes.AddRepoPostHandler)
	m.Get("/:id", routes.RepoHandler)
	m.Get("/:id/assign", routes.AssignRepoHandler)
	m.Get("/:id/*", routes.RepoHandler)
	return m
}

var analyserList = []interface{}{
	checkstyle.CheckstyleTool{},
}

func start(clx *cli.Context) (err error) {
	// TODO load config
	// TODO setup database
	if _, err := os.Stat("./repos"); os.IsNotExist(err) {
		err := os.Mkdir("./repos", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	for _, a := range analyserList {
		ana := a.(analysers.Analyser)
		if !ana.HasTool() {
			log.Printf("%s has no tool installed.\n", ana.GetInfo().Name)
		}
	}

	if (analyserList[0].(analysers.Analyser)).HasTool() {

	}

	e := models.SetupEngine()
	defer e.Close()

	log.Printf("Starting TLS web server on :%s\n", clx.String("port"))
	m := getMacaron(clx.Bool("dev"))
	server := &http.Server{Addr: fmt.Sprintf(":%s", clx.String("port")), Handler: m}
	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			panic(err)
		}
		time.Sleep(500 * time.Millisecond)
	}()
	// defer close whatever here

	// Capture system interrupt
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	return nil
}
