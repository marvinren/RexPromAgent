package promprocess

import (
	"errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type PrometheusWatcher struct {
	interval int
	prom     PrometheusProcess
}

func NewPrometheusWatcher(prom *PrometheusProcess) *PrometheusWatcher {
	heartBeatInterval := 5
	return &PrometheusWatcher{
		heartBeatInterval,
		*prom,
	}
}

func (w *PrometheusWatcher) checkPrometheusHealthy() (string, error) {
	client := &http.Client{}
	// Build the request

	request, err := http.NewRequest("GET", w.prom.baseUrl+"/-/healthy", nil)
	if err != nil {
		logrus.Warnf("build prometheus healthy heartbeat request error: %v", err)
		return "down", err
	}
	// Send the request
	response, err2 := client.Do(request)
	if err2 != nil {
		logrus.Warnf("send prometheus healthy heartbeat error: %v", err)
		return "down", err
	}
	// Check http response status
	if response.StatusCode != http.StatusOK {
		return "down", errors.New("prometheus heartBeat status error.")
	}
	s, _ := ioutil.ReadAll(response.Body)
	logrus.Infof("get %v for prometheus healthy, status code: %v, body: %v", w.prom.baseUrl+"/-/healthy",
		response.StatusCode, string(s))
	return "up", nil
}

func (w *PrometheusWatcher) WatchPrometheus() error {
	ticker := time.NewTicker(time.Duration(w.interval) * time.Minute)
	i := 0

	go func() {
		for {
			i++
			status, err := w.checkPrometheusHealthy()
			w.prom.state = status
			if err != nil {
				logrus.Warnf("prometheus process down/error, %s", err)
				ticker.Stop()
				err := w.restartPrometheus()
				if err == nil {
					w.prom.state = "up"
				}
				ticker = time.NewTicker(time.Duration(w.interval) * time.Minute)
			}
			logrus.Infof("prometheus heart beat at %v", <-ticker.C)

		}
	}()

	return nil
}

func (w *PrometheusWatcher) restartPrometheus() error {
	return nil
}
