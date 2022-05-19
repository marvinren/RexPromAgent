package promprocess

import (
	"RexPromAgent/pkg/config"
	"RexPromAgent/pkg/log"
	"testing"
)

func TestRunPrometheus(t *testing.T) {
	prom := NewPrometheusProcess()
	err := prom.restart()
	if err != nil {
		t.Errorf("start promethues error: %v", err)
	}
}

func TestMain(m *testing.M) {
	config.Initialize()
	log.Initialize()
	m.Run()
}
