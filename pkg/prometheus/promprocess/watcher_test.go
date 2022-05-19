package promprocess

import (
	"testing"
	"time"
)

func TestPromWatcher(t *testing.T) {
	prom := NewPrometheusProcess()
	watcher := NewPrometheusWatcher(prom)

	err := watcher.WatchPrometheus()

	if err != nil {
		t.Errorf("watch error: %v", err)
	}
	time.Sleep(5 * time.Second)
}

func TestPromHealthCheck(t *testing.T) {
	prom := NewPrometheusProcess()
	watcher := NewPrometheusWatcher(prom)
	s, err := watcher.checkPrometheusHealthy()
	if err != nil {
		t.Errorf("Healthy check error")
	}

	t.Logf("server %v", s)

}
