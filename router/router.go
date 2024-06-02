package router

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/thedivinez/gothex/pages"
)

type GothexRouter struct {
	*echo.Echo
	configs     *Configs
	CookieStore *sessions.CookieStore
	Protected   func(next echo.HandlerFunc) echo.HandlerFunc
}

func createRouter(configs *Configs) *GothexRouter {
	gothex := echo.New()
	gothex.Static("/", "public")
	sessionAge, _ := strconv.Atoi(configs.SessionAge)
	cookieStore := sessions.NewCookieStore([]byte(configs.GothexSecret))
	cookieStore.Options.HttpOnly = true
	cookieStore.Options.MaxAge = sessionAge
	cookieStore.Options.Path = configs.AfterSignin
	gothex.Use(session.Middleware(cookieStore))
	gothex.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("X-Title", configs.AppTitle)
			return next(c)
		}
	})

	gothex.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			if !strings.HasPrefix(c.Request().RequestURI, "/api") {
				errPage := pages.ErrorPageContent{Code: code}
				errPage.ErrorType = he.Message.(string)
				errPage.Title = fmt.Sprintf("Error (%d)", code)
				if page := inCustomErrorPages(code); page != nil {
					errPage = *page
				}
				c.Set("X-Title", errPage.Title)
				c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
				ShowComponent(c, pages.Error(c.Request().Context(), errPage))
				return
			}
		}
		c.JSON(code, err)
	}

	return &GothexRouter{
		gothex,
		configs,
		cookieStore,
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				session, authenticated := GetAuthSession(c)
				if !authenticated {
					return c.Redirect(http.StatusPermanentRedirect, configs.SignInPage)
				}
				sessionAge, _ := strconv.Atoi(configs.SessionAge)
				session.Options.MaxAge = sessionAge
				if err := session.Save(c.Request(), c.Response()); err != nil {
					return err
				}
				return next(c)
			}
		},
	}
}

func NewGothexRouter() *GothexRouter {
	configs, err := getConfigs()
	if err != nil {
		log.Fatal("failed to get configs", err)
	}
	return createRouter(configs)
}

func NewRouterWithConfigs(configs *Configs) *GothexRouter {
	return createRouter(configs)
}

func (r *GothexRouter) WithCustomErrorPageContent(pages ...pages.ErrorPageContent) *GothexRouter {
	customErrorPageContent = pages
	return r
}

func (r *GothexRouter) Run() error {
	return r.Start(fmt.Sprintf(":%s", r.configs.Port))
}
