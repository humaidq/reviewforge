package routes

import (
	"net/http"

	macaron "gopkg.in/macaron.v1"
)

func ContextInit() macaron.Handler {
	return func(ctx *macaron.Context) {
		ctx.Data["SiteTitle"] = "reviewforge"
		ctx.Data["User"] = "humaid"
	}
}

func DashboardHandler(ctx *macaron.Context) {
	ctx.HTML(http.StatusOK, "index")
}

func AddRepoHandler(ctx *macaron.Context) {
	ctx.HTML(http.StatusOK, "add_repo")
}
