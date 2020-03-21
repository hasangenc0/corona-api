package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hasangenc0/corona/pkg/helpers"
	routing "github.com/qiangxue/fasthttp-routing"
	"log"
	"time"
)

func (s *Server) countryHandler(c *routing.Context) error {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	timeout, err := time.ParseDuration(s.Conf.Server.Timeout)
	if err == nil {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	country := c.Param("country")
	result, err := s.Corona.Country(ctx, country)

	if err != nil {
		return err
	}

	data, err := json.Marshal(result)
	err = c.WriteData(data)
	return err
}

func (s *Server) allCountriesHandler(c *routing.Context) error {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	timeout, err := time.ParseDuration(s.Conf.Server.Timeout)
	if err == nil {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}

	defer cancel()

	result, err := s.Corona.AllCountries(ctx)

	if err != nil {
		return err
	}

	data, err := json.Marshal(result)
	err = c.WriteData(data)
	return err
}

func (s *Server) countryPostHandler(c *routing.Context) error {
	client := s.Db.Connect()
	collection := client.Database("corona").Collection("country")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	countries, err := s.Corona.AllCountries(ctx)
	if err != nil {
		return err
	}

	for _, country := range countries {
		country.Date = time.Now()
		res, err := collection.InsertOne(ctx, country)

		if err == nil {
			log.Println(fmt.Sprintf("Inserted: %v", res.InsertedID))
		}
	}

	return err
}

func (s *Server) swagger(c *routing.Context) error {
	file := c.Param("path")

	if file == "" {
		file = "/cmd/swagger/index.html"
	} else {
		file = fmt.Sprintf("/cmd/swagger/%s", file)
	}

	path := helpers.GetPath(file)
	c.SendFile(path)
	return nil
}
