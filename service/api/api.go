/*
Package api exposes the main API engine. All HTTP APIs are handled here - so-called "business logic" should be here, or
in a dedicated package (if that logic is complex enough).

To use this package, you should create a new instance with New() passing a valid Config. The resulting Router will have
the Router.Handler() function that returns a handler that can be used in a http.Server (or in other middlewares).

Example:

	// Create the API router
	apirouter, err := api.New(api.Config{
		Logger:   logger,
		Database: appdb,
	})
	if err != nil {
		logger.WithError(err).Error("error creating the API server instance")
		return fmt.Errorf("error creating the API server instance: %w", err)
	}
	router := apirouter.Handler()

	// ... other stuff here, like middleware chaining, etc.

	// Create the API server
	apiserver := http.Server{
		Addr:              cfg.Web.APIHost,
		Handler:           router,
		ReadTimeout:       cfg.Web.ReadTimeout,
		ReadHeaderTimeout: cfg.Web.ReadTimeout,
		WriteTimeout:      cfg.Web.WriteTimeout,
	}

	// Start the service listening for requests in a separate goroutine
	apiserver.ListenAndServe()

See the `main.go` file inside the `cmd/webapi` for a full usage example.
*/
package api

import (
	"errors"
	"github.com/ciottolomaggico/wasatext/service/app"
	"github.com/ciottolomaggico/wasatext/service/middlewares"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

// Config is used to provide dependencies and configuration to the New function.
type Config struct {
	// Logger where log entries are sent
	Logger          logrus.FieldLogger
	StaticFilesUrl  string
	StaticFilesPath string
}

// Router is the package API interface representing an API handler builder
type Router interface {
	// Handler returns an HTTP handler for APIs provided in this package
	//Handler(routers []routers.ControllerRouter) http.Handler
	Handler() http.Handler

	// Close terminates any resource used in the package
	Close() error
}

type _router struct {
	router         *httprouter.Router
	authMiddleware middlewares.AuthMiddleware
	baseLogger     logrus.FieldLogger
	app            app.Application
}

// New returns a new Router instance
func New(authMiddleware middlewares.AuthMiddleware, cfg Config, app app.Application) (Router, error) {
	// Check if the configuration is correct
	if cfg.Logger == nil {
		return nil, errors.New("logger is required")
	}
	if _, err := os.Stat(cfg.StaticFilesPath); errors.Is(err, os.ErrNotExist) {
		return nil, errors.New("static files path does not exist")
	}
	// Create a new router where we will register HTTP endpoints. The server will pass requests to this router to be
	// handled.
	router := httprouter.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false
	router.ServeFiles(cfg.StaticFilesUrl, http.Dir(cfg.StaticFilesPath))

	return &_router{
		router,
		authMiddleware,
		cfg.Logger,
		app,
	}, nil
}
