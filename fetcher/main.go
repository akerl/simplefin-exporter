package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/akerl/simplefin-exporter/config"

	"github.com/akerl/metrics/metrics"
	"github.com/akerl/metrics/server"
	"github.com/akerl/timber/v2/log"
)

const (
	accountsEndpoint = "/accounts"
)

var logger = log.NewLogger("simplefin-exporter.fetcher")

// Fetcher defines the ticker fetching engine
type Fetcher struct {
	Interval  int
	AccessURL string
	Cache     *server.Cache
}

type accounts struct {
	Errors   []string `json:"errors"`
	Accounts []struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Balance string `json:"balance"`
	} `json:"accounts"`
}

// NewFetcher creates a new syslog engine from the given config
func NewFetcher(conf config.Config, cache *server.Cache) *Fetcher {
	return &Fetcher{
		Interval:  conf.Interval,
		AccessURL: conf.AccessURL,
		Cache:     cache,
	}
}

// RunAsync launches the fetcher engine in the background
func (f *Fetcher) RunAsync() {
	go f.Run()
}

// Run launches the fetcher engine in the foreground
func (f *Fetcher) Run() {
	for {
		logger.DebugMsg("running fetcher loop")
		ms, err := f.fetchAccounts()
		if err != nil {
			panic(err)
		}

		f.Cache.MetricSet = ms
		time.Sleep(time.Duration(f.Interval) * time.Second)
	}
}

func (f *Fetcher) fetchAccounts() (metrics.MetricSet, error) {
	url := f.AccessURL + accountsEndpoint

	resp, err := http.Get(url, "application/json", nil)
	if err != nil {
		return metrics.MetricSet{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return metrics.MetricSet{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return metrics.MetricSet{}, fmt.Errorf("api error %d: %s", resp.StatusCode, body)
	}

	var a accounts
	err = json.Unmarshal(body, &a)
	if err != nil {
		return metrics.MetricSet{}, err
	}
	if len(a.Errors) > 0 {
		return metrics.MetricSet{}, fmt.Errorf("api returned errors: %v+", a.Errors)
	}

	ms := metrics.MetricSet{
		metrics.Metric{
			Name:  "last_updated",
			Type:  "gauge",
			Value: fmt.Sprintf("%d", time.Now().Unix()),
		},
	}
	for _, item := range a.Accounts {
		ms = append(ms, metrics.Metric{
			Name: "simplefin_balance",
			Type: "gauge",
			Tags: map[string]string{
				"account": item.Name,
				"id":      item.ID,
			},
			Value: item.Balance,
		})
	}
	return ms, nil
}
