package router

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/thedivinez/gothex/pages"
	"github.com/thedivinez/gothex/types"
)

type Configs struct {
	Port         string `json:"PORT"`
	AppTitle     string `json:"APP_TITLE"`
	SessionAge   string `json:"SESSION_AGE"`
	SignInPage   string `json:"SIGNIN_PAGE"`
	GothexSecret string `json:"GOTHEX_SECRET"`
	AfterSignin  string `json:"AFTER_SIGNIN_PAGE"`
}

var customErrorPageContent []pages.ErrorPageContent

func inCustomErrorPages(code int) *pages.ErrorPageContent {
	for _, page := range customErrorPageContent {
		if page.Code == code {
			return &page
		}
	}
	return nil
}

func getConfigs() (*Configs, error) {
	data, err := godotenv.Read()
	if err != nil {
		return nil, err
	}
	resultBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var configs Configs
	if err := json.Unmarshal(resultBytes, &configs); err != nil {
		return nil, err
	}
	return &configs, nil
}

func GetAuthSession(c echo.Context) (*sessions.Session, bool) {
	if session, err := session.Get("__auth__", c); err != nil {
		return nil, false
	} else {
		if auth, ok := session.Values["authenticated"].(bool); ok && auth {
			return session, true
		} else {
			return session, false
		}
	}
}

func SignIn(c echo.Context, redirect string, user map[string]any) error {
	session, _ := GetAuthSession(c)
	if jsonString, err := json.Marshal(user); err != nil {
		return err
	} else {
		session.Values = map[any]any{"authenticated": true, "user": jsonString}
	}
	if err := session.Save(c.Request(), c.Response()); err != nil {
		return err
	}
	c.Response().Header().Set("HX-Refresh", "true")
	c.Response().WriteHeader(http.StatusPermanentRedirect)
	return nil
}

func SignOut(c echo.Context, redirect string) error {
	session, _ := GetAuthSession(c)
	session.Values = map[any]any{"authenticated": false}
	if err := session.Save(c.Request(), c.Response()); err != nil {
		return err
	}
	c.Response().Header().Set("HX-Refresh", "true")
	c.Response().WriteHeader(http.StatusPermanentRedirect)
	return nil
}

type TemplBody func(ctx context.Context, options ...any) templ.Component

func Render(c echo.Context, body TemplBody, data ...any) error {
	ctx := context.WithValue(c.Request().Context(), types.ContextKey{Key: "X-Title"}, c.Get("X-Title"))
	return body(ctx, data...).Render(ctx, c.Response().Writer)
}

func RenderWithTitle(c echo.Context, title string, body TemplBody, data ...any) error {
	c.Set("X-Title", title)
	return Render(c, body, data...)
}

func ShowComponent(c echo.Context, cmp templ.Component) error {
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

func IsHxRequest(c echo.Context) bool {
	return c.Request().Header.Get("HX-Request") == "true"
}
