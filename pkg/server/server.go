package server

import (
	"fmt"
	"github.com/hasangenc0/corona/pkg/configuration"
	"github.com/hasangenc0/corona/pkg/corona"
	"github.com/hasangenc0/corona/pkg/db"
	"github.com/hasangenc0/corona/pkg/environment"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
)

type Server struct {
	Env    environment.Environment
	Conf   *configuration.Config
	Db     *db.DB
	Corona *corona.Corona
}

func (s *Server) Bootstrap() *Server {
	s.Env = environment.Get()
	s.Conf = configuration.Read(s.Env.Env)
	s.Db = &db.DB{Host: s.Conf.Db.Host, Timeout: s.Conf.Db.Timeout}
	s.Corona = &corona.Corona{Api: s.Conf.Api.Corona}
	return s
}

func (s *Server) Start() {
	router := routing.New()
	if s.Env.Env != "production" {
		router.Get("/<path>", s.swagger)
	}
	api := router.Group("/api")
	api.Get("/country", s.allCountriesHandler).Post(s.countryPostHandler)
	api.Get("/country/<country>", s.countryHandler)

	port := s.Conf.Server.Port
	if s.Env.Port != "" {
		port = fmt.Sprintf(":%s", s.Env.Port)
	}
	log.Println(fmt.Sprintf("🏃Server started: http://localhost%s/", port))
	log.Fatal(fasthttp.ListenAndServe(port, router.HandleRequest))
}
