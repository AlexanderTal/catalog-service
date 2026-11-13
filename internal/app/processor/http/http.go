package rprocessor

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/AlexanderTal/catalog-service/internal/app/config/section"
	rhandler "github.com/AlexanderTal/catalog-service/internal/app/handler/http"
	"github.com/AlexanderTal/catalog-service/internal/app/processor"
	"github.com/AlexanderTal/catalog-service/internal/app/util"
	httph "github.com/AlexanderTal/catalog-service/internal/pkg/http/httph"
	mzerolog "github.com/AlexanderTal/catalog-service/internal/pkg/http/mzerolog"
)

type httpProc struct {
	addr   string
	server *http.Server
}

func NewHTTP(
	hHealth rhandler.Health,
	hCategory rhandler.Category,
	hProduct rhandler.Product,
	cfg section.ProcessorWebServer,
) processor.Processor {
	r := mux.NewRouter()

	r.Use(
		httph.NewErrorMiddleware(),
		mzerolog.NewMiddleware(
			mzerolog.WithSkipper(util.IsFilteredHttpRoute),
		),
	)

	r.HandleFunc("/health", hHealth.LastCheck).Methods(http.MethodGet)
	log.Info().Strs("methods", []string{http.MethodGet}).Str("path", "/health").Msg("Route registered")

	r.HandleFunc("/v1/category/create", hCategory.Create).Methods(http.MethodPost)
	log.Info().Strs("methods", []string{http.MethodPost}).Str("path", "/v1/category/create").Msg("Route registered")

	r.HandleFunc("/v1/category/{guid}", hCategory.Update).Methods(http.MethodPatch)
	log.Info().Strs("methods", []string{http.MethodPatch}).Str("path", "/v1/category/{guid}").Msg("Route registered")

	r.HandleFunc("/v1/category/{guid}", hCategory.Delete).Methods(http.MethodDelete)
	log.Info().Strs("methods", []string{http.MethodDelete}).Str("path", "/v1/category/{guid}").Msg("Route registered")

	r.HandleFunc("/v1/category/list", hCategory.List).Methods(http.MethodPost)
	log.Info().Strs("methods", []string{http.MethodPost}).Str("path", "/v1/category/list").Msg("Route registered")

	r.HandleFunc("/v1/product/create", hProduct.Create).Methods(http.MethodPost)
	log.Info().Strs("methods", []string{http.MethodPost}).Str("path", "/v1/product/create").Msg("Route registered")

	r.HandleFunc("/v1/product/{guid}", hProduct.Update).Methods(http.MethodPatch)
	log.Info().Strs("methods", []string{http.MethodPatch}).Str("path", "/v1/product/{guid}").Msg("Route registered")

	r.HandleFunc("/v1/product/{guid}", hProduct.Delete).Methods(http.MethodDelete)
	log.Info().Strs("methods", []string{http.MethodDelete}).Str("path", "/v1/product/{guid}").Msg("Route registered")

	r.HandleFunc("/v1/product/list", hProduct.List).Methods(http.MethodPost)
	log.Info().Strs("methods", []string{http.MethodPost}).Str("path", "/v1/product/list").Msg("Route registered")

	p := httpProc{addr: fmt.Sprintf(":%d", cfg.ListenPort)}
	p.server = &http.Server{
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &p
}

func (p *httpProc) StartAsync(ctx context.Context, wg *sync.WaitGroup) {
	var lc net.ListenConfig
	l, err := lc.Listen(ctx, "tcp", p.addr)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start listening TCP addr")
	}

	log.Info().Str("listen_addr", p.addr).Msg("Listening of TCP addr for HTTP server has been started")

	go p.serve(l)

	processor.WatchForShutdown(ctx, wg, processor.CloserFunc(l.Close))
	processor.WatchForShutdown(ctx, wg, processor.NewCloserContextFunc(
		p.server.Shutdown,
		context.Background(), 5*time.Second,
	))
}

func (p *httpProc) serve(l net.Listener) {
	_ = p.server.Serve(l)
}
