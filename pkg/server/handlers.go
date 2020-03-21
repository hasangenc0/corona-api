package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hasangenc0/corona/pkg/corona"
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
	result, err := corona.Country(ctx, country)

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

	result, err := corona.AllCountries(ctx)

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

	countries, err := corona.AllCountries(ctx)
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
