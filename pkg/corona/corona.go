package corona

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Result struct {
	Country            string `json:"country"`
	Cases              int32  `json:"cases"`
	TodayCases         int32  `json:"todayCases"`
	Deaths             int32  `json:"deaths"`
	TodayDeaths        int32  `json:"todayDeaths"`
	Recovered          int32  `json:"recovered"`
	Active             int32  `json:"active"`
	Critical           int32  `json:"critical"`
	CasesPerOneMillion int32  `json:"casesPerOneMillion"`
}

func Search(ctx context.Context, country string) (Result, error) {
	url := fmt.Sprintf("https://corona.lmao.ninja/countries/%s", country)
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
