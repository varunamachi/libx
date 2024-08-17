package httpx

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/errx"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type noopWriter struct{}
type userKeyType string

const UserKey = userKeyType("requestUser")

func (nw *noopWriter) Write(b []byte) (int, error) {
	return 0, nil
}

type Endpoint struct {
	Method      string
	Path        string
	Category    string
	Desc        string
	Version     string
	Role        auth.Role
	Permissions []string
	Route       *echo.Route
	Handler     echo.HandlerFunc
}

func (ep *Endpoint) NeedsAuth() bool {
	return len(ep.Permissions) != 0 || (ep.Role != auth.None && ep.Role != "")
}

type Server struct {
	apiCatg         map[string][]*Endpoint
	pageCatg        map[string][]*Endpoint
	apiEps          []*Endpoint
	pageEps         []*Endpoint
	echo            *echo.Echo
	printer         io.Writer
	userRetriever   auth.UserRetriever
	rootMiddlewares []echo.MiddlewareFunc
	printAllAccess  bool
}

func NewServer(printer io.Writer, userGetter auth.UserRetriever) *Server {
	if printer == nil {
		printer = &noopWriter{}
	}
	return &Server{
		apiCatg:       make(map[string][]*Endpoint),
		pageCatg:      make(map[string][]*Endpoint),
		apiEps:        make([]*Endpoint, 0, 100),
		pageEps:       make([]*Endpoint, 0, 100),
		echo:          echo.New(),
		printer:       printer,
		userRetriever: userGetter,
	}
}

func (s *Server) WithAPIs(ep ...*Endpoint) *Server {
	s.apiEps = append(s.apiEps, ep...)
	return s
}

func (s *Server) WithPages(ep ...*Endpoint) *Server {
	s.pageEps = append(s.pageEps, ep...)
	return s
}

func (s *Server) WithRootMiddlewares(mws ...echo.MiddlewareFunc) *Server {
	s.rootMiddlewares = append(s.rootMiddlewares, mws...)
	return s
}

func (s *Server) PrintAllAccess(enable bool) *Server {
	s.printAllAccess = enable
	return s
}

func (s *Server) Start(port uint32) error {
	s.configure()
	s.Print()
	if err := s.echo.Start(fmt.Sprintf(":%d", port)); err != nil {
		return errx.Wrap(err)
	}
	log.Info().Uint32("port", port).Msg("server started")
	return nil
}

func (s *Server) configure() {
	s.echo.HTTPErrorHandler = errorHandlerFunc
	s.echo.HideBanner = true
	s.echo.HidePort = true
	// s.echo.Use(getAccessMiddleware())
	s.echo.Use(accessMiddleware(s.printAllAccess))
	s.echo.Use(s.rootMiddlewares...)

	groups := map[string]*echo.Group{}

	for _, ep := range s.apiEps {
		ep := ep
		grp := groups[ep.Version]
		if grp == nil {
			grp = s.echo.Group("api/" + ep.Version + "/")
			groups[ep.Version] = grp
		}

		path := data.Qop(
			strings.HasPrefix(ep.Path, "/"), ep.Path[1:], ep.Path)
		if ep.NeedsAuth() {
			ep.Route = grp.Add(
				ep.Method,
				path,
				ep.Handler,
				getAuthzMiddleware(ep, s))

		} else {
			ep.Route = grp.Add(
				ep.Method, path, ep.Handler)
		}

		if _, found := s.apiCatg[ep.Category]; !found {
			s.apiCatg[ep.Category] = make([]*Endpoint, 0, 20)
		}
		s.apiCatg[ep.Category] = append(s.apiCatg[ep.Category], ep)
	}

	// For simplicity we just duplicate the loop for pages
	for _, ep := range s.pageEps {
		if ep.NeedsAuth() {
			ep.Route = s.echo.Add(
				ep.Method,
				ep.Path,
				ep.Handler,
				getAuthzMiddleware(ep, s))

		} else {
			ep.Route = s.echo.Add(
				ep.Method, ep.Path, ep.Handler)
		}

		if _, found := s.apiCatg[ep.Category]; !found {
			s.pageCatg[ep.Category] = make([]*Endpoint, 0, 20)
		}
		s.pageCatg[ep.Category] = append(s.pageCatg[ep.Category], ep)
	}
}

func (s *Server) Print() {
	caser := cases.Upper(language.English)
	for cat, eps := range s.apiCatg {
		fmt.Fprintln(s.printer, caser.String(cat))
		for idx, ep := range eps {

			sym := "\u251c"
			if idx == len(eps)-1 {
				sym = "\u2514"
			}

			fmt.Fprintf(s.printer,
				" %s\u2500 %-3s %-10s %-60s %s\n",
				sym, ep.Version, ep.Route.Method, ep.Route.Path, ep.Desc)
		}
		fmt.Fprintln(s.printer)
	}

	if len(s.pageEps) != 0 {

		fmt.Fprintln(s.printer, "\nPAGES:")
		for idx, ep := range s.pageEps {
			cat := ep.Category
			if len(cat) > 14 {
				cat = ep.Category[:14]
			}

			sym := "\u251c"
			if idx == len(s.pageEps)-2 {
				sym = "\u2514"
			}

			fmt.Fprintf(s.printer,
				" %s\u2500 %-3s %-10s %-60s %-15s %s\n",
				sym, ep.Version, ep.Route.Method, ep.Route.Path, cat, ep.Desc)
		}
	}
	fmt.Fprintln(s.printer, "")
}

func (s *Server) Close() error {
	return s.echo.Close()
}

func SendJSON(etx echo.Context, data interface{}) error {

	// Following is required for flutter client
	etx.Response().Header().Set(
		echo.HeaderContentType, "application/json; charset=UTF-8")
	if err := etx.JSON(http.StatusOK, data); err != nil {
		log.Error().Err(err).Msg("failed to write JSON response")
		return errx.Wrap(err)
	}
	return nil
}
