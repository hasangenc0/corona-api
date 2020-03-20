package main

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/hasangenc0/corona/pkg/corona"
	"github.com/hasangenc0/corona/pkg/userip"
)

type RenderModel struct {
	Result  corona.Result
	Timeout time.Duration
	Elapsed time.Duration
	Start   time.Time
}

func main() {
	http.HandleFunc("/api", handleSearch)
	log.Println("🏃Server started: http://localhost:8080/api")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleSearch(w http.ResponseWriter, req *http.Request) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	timeout, err := time.ParseDuration(req.FormValue("timeout"))
	if err == nil {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	query := req.FormValue("country")
	if query == "" {
		http.Error(w, "no query", http.StatusBadRequest)
		return
	}

	userIP, err := userip.FromRequest(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx = userip.NewContext(ctx, userIP)
	start := time.Now()
	result, err := corona.Search(ctx, query)
	elapsed := time.Since(start)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	model := RenderModel{
		result,
		timeout,
		elapsed,
		start,
	}

	//_ = renderJson(w, result)

	renderTemplate(w, model)
}

func renderJson(w http.ResponseWriter, result corona.Result) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(result)
}

func renderTemplate(w http.ResponseWriter, model RenderModel) {
	if err := resultsTemplate.Execute(w, model); err != nil {
		log.Print(err)
		return
	}
}

var resultsTemplate = template.Must(template.New("results").Parse(`
<html>
<head/>
<body>
	<ul>
	<li>Country: {{ .Result.Country }}</li>
	<li>Cases: {{ .Result.Cases }},</li>
	<li>TodayCases: {{ .Result.TodayCases }},</li>
	<li>Deaths: {{ .Result.Deaths }},</li>
	<li>TodayDeaths: {{ .Result.TodayDeaths }},</li>
	<li>Recovered: {{ .Result.Recovered }},</li>
	<li>Active: {{ .Result.Active }},</li>
	<li>Critical: {{ .Result.Critical }},</li>
	<li>CasesPerOneMillion: {{ .Result.CasesPerOneMillion }}</li>
	</ul>
</body>
</html>
`))
