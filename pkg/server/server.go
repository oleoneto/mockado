package server

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cleopatrio/mockado/pkg/envloader"
	"github.com/cleopatrio/mockado/pkg/helpers"
	"github.com/cleopatrio/mockado/pkg/httpclient"
	"github.com/cleopatrio/mockado/pkg/logger"
	"github.com/cleopatrio/mockado/pkg/redactor"
	"github.com/gofiber/fiber/v3"
	flogger "github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

var once sync.Once

var instance *Server = &Server{}

type Server struct {
	app    *fiber.App
	redact *redactor.Redactor
	loader *envloader.Loader
}

func New() *Server {
	once.Do(func() {
		app := fiber.New(fiber.Config{
			Immutable: true,
		})

		app.Use(flogger.New(
			flogger.Config{
				ForceColors: true,
				Format:      flogger.JSONFormat,
				TimeZone:    "America/Sao_Paulo",
				// TimeFormat: flogger.,
			},
		))

		app.Use(recover.New(
			recover.Config{EnableStackTrace: false},
		))

		instance.app = app
		instance.redact = redactor.New()
		instance.loader = envloader.New()

		log.
			WithField("package", "server").
			Tracef("Called New()")
	})

	return instance
}

func (s *Server) Listen(port int) error {
	return s.app.Listen(
		fmt.Sprintf("127.0.0.1:%d", port),
		fiber.ListenConfig{
			EnablePrintRoutes:     false,
			DisableStartupMessage: true,
		},
	)
}

func (s *Server) Add(routes ...httpclient.Route) *Server {
	for _, r := range routes {
		s.app.Add(r.Methods, r.Path, s.route(r))
	}

	return s
}

func (s *Server) route(r httpclient.Route) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		c.Response().Header.Add("X-Request-Time", time.Now().String())

		if r.Accepts != "" {
			c.Response().Header.Add("Accepts", r.Accepts)
		}

		for key, value := range r.Response.Headers {
			c.Response().Header.Add(key, s.loader.LoadVar(value))
		}

		latency := helpers.UnwrapOrDefault(r.Response.Latency, 0)
		c.Response().Header.Add("X-Latency", latency.String())

		var b []byte
		var variables map[string]string

		if r.Response.EnvAsBody {
			variables = make(map[string]string, 0)

			for _, pair := range os.Environ() {
				a := strings.SplitN(pair, "=", 2)
				key, value := a[0], a[1]

				if v, matched := s.redact.RedactKey(key); matched {
					variables[key] = v
					continue
				}

				variables[key] = s.redact.Redact(value)
			}
		}

		time.Sleep(latency)

		c.Response().Header.Add("X-Response-Time", time.Now().String())

		c.Status(r.Response.StatusCode)
		if r.Response.EnvAsBody {
			return c.JSON(variables)
		}

		return c.Send(b)
	}
}

var log = logger.GetLogger("server")
