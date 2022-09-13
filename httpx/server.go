package httpx

import (
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/varunamachi/libx/auth"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type noopWriter struct{}

func (nw *noopWriter) Write(b []byte) (int, error) {
	return 0, nil
}

type Endpoint struct {
	Method     string
	Path       string
	Category   string
	Desc       string
	Version    string
	Role       auth.Role
	Permission string
	Route      *echo.Route
	Handler    echo.HandlerFunc
}

func (ep *Endpoint) NeedsAuth() bool {
	return ep.Permission != "" || (ep.Role != auth.None && ep.Role != "")
}

type Server struct {
	apiCatg       map[string][]*Endpoint
	pageCatg      map[string][]*Endpoint
	apiEps        []*Endpoint
	pageEps       []*Endpoint
	echo          *echo.Echo
	printer       io.Writer
	userRetriever auth.UserRetrieverFunc
}

func NewServer(printer io.Writer, userGetter auth.UserRetrieverFunc) *Server {
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

func (s *Server) Start(port uint32) error {
	s.configure()
	s.Print()
	return s.echo.Start(fmt.Sprintf(":%d", port))
}

func (s *Server) configure() {

	s.echo.HideBanner = true
	s.echo.Use(getAccessMiddleware())
	groups := map[string]*echo.Group{}

	for _, ep := range s.apiEps {
		ep := ep
		grp := groups[ep.Version]
		if grp == nil {
			grp = s.echo.Group("api/" + ep.Version + "/")
			groups[ep.Version] = grp
		}

		if ep.NeedsAuth() {
			ep.Route = grp.Add(
				ep.Method,
				ep.Path,
				ep.Handler,
				getAuthzMiddleware(ep, s))

		} else {
			ep.Route = grp.Add(
				ep.Method, ep.Path, ep.Handler)
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
	// for _, ep := range s.apiEps {
	// 	cat := ep.Category
	// 	if len(cat) > 14 {
	// 		cat = ep.Category[:14]
	// 	}
	// 	fmt.Fprintf(s.printer,
	// 		"%-3s %-5s %-40s %-15s %s\n",
	// 		ep.Version, ep.Route.Method, ep.Route.Path, cat, ep.Desc)
	// }
	// fmt.Print("\n\n")

	caser := cases.Upper(language.English)
	for cat, eps := range s.apiCatg {
		fmt.Fprintln(s.printer, caser.String(cat))
		for _, ep := range eps {
			fmt.Fprintf(s.printer,
				"    |--  %-3s %-5s %-40s %s\n",
				ep.Version, ep.Route.Method, ep.Route.Path, ep.Desc)
		}
		fmt.Fprintln(s.printer)
	}
	fmt.Fprintln(s.printer, "\nPAGES:")
	for _, ep := range s.pageEps {
		cat := ep.Category
		if len(cat) > 14 {
			cat = ep.Category[:14]
		}

		fmt.Fprintf(s.printer,
			"%-3s %-5s %-40s %-15s %s\n",
			ep.Version, ep.Route.Method, ep.Route.Path, cat, ep.Desc)
	}
	fmt.Fprintln(s.printer, "")
}

func (s *Server) Close() error {
	return s.echo.Close()
}

func SendJSON(etx echo.Context, data interface{}) error {
	if err := etx.JSON(http.StatusOK, data); err != nil {
		log.Error().Err(err).Msg("failed to write JSON response")
		return err
	}
	return nil
}
