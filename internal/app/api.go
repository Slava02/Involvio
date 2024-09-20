package app

import (
	"crypto/tls"
	"fmt"
	"github.com/Slava02/Involvio/config"
	"github.com/Slava02/Involvio/internal/controller"
	"github.com/Slava02/Involvio/internal/route/api"
	"github.com/Slava02/Involvio/internal/usecase"
	"github.com/Slava02/Involvio/internal/usecase/repository"
	"github.com/Slava02/Involvio/pkg/database"
	"github.com/Slava02/Involvio/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"log"
	"log/slog"
	"net/http"
	"sync"
)

type Router struct {
	*controller.Impl
}

func New() *Router {
	return &Router{}
}

func (r *Router) SetupGlobalMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logrus.Infof("Received request on path: %q with method %q", req.URL.String(), req.Method)
		handler.ServeHTTP(w, req)
	})
}

func (r *Router) ConfigureFlags(api *api.InvolvioAPI) {
	return
}

func (*Router) ConfigureTLS(_ *tls.Config) {}

func (*Router) ConfigureServer(_ *http.Server, _, _ string) {}

func (*Router) SetupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

func (r *Router) CustomConfigure(api *api.InvolvioAPI) {
	// Load environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	// Initialize the logger
	logger.SetupLogger(cfg)

	r.Impl = newImpl(cfg)
}

func newImpl(cfg *config.Config) *controller.Impl {
	// TODO: как тут закрыть коннект тогда?
	pg, err := database.New(cfg, database.MaxPoolSize(cfg.DB.PoolMax), database.Isolation(pgx.ReadCommitted))
	if err != nil {
		log.Fatal("postgres connection failed", slog.String("error", err.Error()))
	}

	err = applyMigrations(cfg.DB)
	if err != nil {
		log.Fatal("apply migrations failed", slog.String("error", err.Error()))
	}

	o := sync.Once{}
	repo := repository.NewRepository(&o, pg)
	uc := usecase.NewUseCase(repo)

	return controller.NewImpl(uc)
}
