package frontend

import (
	"encoding/json"
	"gothex/pages"
	"gothex/router"
	"gothex/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Pages struct{}

func NewPages() *Pages {
	return &Pages{}
}

func (*Pages) HandleSignUp(c echo.Context) error {
	if session, authenticated := router.GetAuthSession(c); authenticated {
		return c.Redirect(http.StatusPermanentRedirect, session.Options.Path)
	}
	return router.RenderWithTitle(c, "Sign Up", pages.SignUp)
}

func (*Pages) HandleSignIn(c echo.Context) error {
	if session, authenticated := router.GetAuthSession(c); authenticated {
		return c.Redirect(http.StatusPermanentRedirect, session.Options.Path)
	}
	return router.RenderWithTitle(c, "Sign In", pages.SignIn)
}

func (*Pages) HandleHome(c echo.Context) error {
	var user types.User
	session, _ := router.GetAuthSession(c)
	json.Unmarshal(session.Values["user"].([]byte), &user)
	return router.Render(c, pages.Home, user)
}
