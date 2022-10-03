package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/pkg/errors"
	"github.com/restlesswhy/video-merger/config"
	v1 "github.com/restlesswhy/video-merger/internal/delivery/http/v1"
	"github.com/restlesswhy/video-merger/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"

	"net/http"
	_ "net/http/pprof"
)

type server struct {
	log   logger.Logger
	cfg   *config.Config
	pool  *pgxpool.Pool
	fiber *fiber.App
}

func New(log logger.Logger, cfg *config.Config, pool *pgxpool.Pool) *server {
	return &server{log: log, cfg: cfg, pool: pool}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	tmp := filepath.Join(".", "tmp")
	if _, err := os.Stat(tmp); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(tmp, os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "create tmp dir error")
		}
	}

	f := fiber.New(fiber.Config{
		BodyLimit: 50 << 20,
	})
	s.fiber = f

	controller := v1.New(nil)
	controller.SetupRoutes(s.fiber)

	go func() {
		if err := s.runHttp(); err != nil {
			s.log.Errorf("(runHttp) err: %v", err)
			cancel()
		}
	}()
	s.log.Infof("%s is listening on PORT: %v", s.getMicroserviceName(), s.cfg.Http.Port)

	go func() {
		s.log.Error(http.ListenAndServe(":6060", nil))
	}()

	<-ctx.Done()

	err := s.fiber.Shutdown()
	if err != nil {
		s.log.Warnf("(Shutdown) err: %v", err)
		return err
	}

	s.log.Info("Service gracefully closed.")
	return nil
}

func (s server) getMicroserviceName() string {
	return fmt.Sprintf("(%s)", strings.ToUpper(s.cfg.ServiceName))
}
