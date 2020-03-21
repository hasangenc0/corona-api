package corona

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Result struct {
	Country            string    `bson:"country" json:"country"`
	Cases              int32     `bson:"cases" json:"cases"`
	TodayCases         int32     `bson:"todayCases" json:"todayCases"`
	Deaths             int32     `bson:"deaths" json:"deaths"`
	TodayDeaths        int32     `bson:"todayDeaths" json:"todayDeaths"`
	Recovered          int32     `bson:"recovered" json:"recovered"`
	Active             int32     `bson:"active" json:"active"`
	Critical           int32     `bson:"critical" json:"critical"`
	CasesPerOneMillion int32     `bson:"casesPerOneMillion" json:"casesPerOneMillion"`
	Date               time.Time `bson:"date" json:"date"`
}

type Corona struct {
	Api string
}

func (corona *Corona) AllCountries(ctx context.Context) ([]Result, error) {
	url := fmt.Sprintf("%s/countries", corona.Api)
	req, err := http.NewRequest("GET", url, nil)
	var result []Result

	if err != nil {
		return result, err
	}

	err = httpDo(ctx, req, func(resp *http.Response, err error) error {
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		return nil
	})
	return result, err
}

func (corona *Corona) Country(ctx context.Context, country string) (Result, error) {
	url := fmt.Sprintf("%s/countries/%s", corona.Api, country)
	req, err := http.NewRequest("GET", url, nil)
	var result Result

	if err != nil {
		return result, err
	}

	err = httpDo(ctx, req, func(resp *http.Response, err error) error {
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		return nil
	})
	return result, err
}

func httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	// Run the HTTP request in a goroutine and pass the response to f.
	c := make(chan error, 1)
	req = req.WithContext(ctx)
	go func() { c <- f(http.DefaultClient.Do(req)) }()
	select {
	case <-ctx.Done():
		<-c // Wait for f to return.
		return ctx.Err()
	case err := <-c:
		return err
	}
}
