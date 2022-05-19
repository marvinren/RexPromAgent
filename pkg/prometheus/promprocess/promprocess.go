package promprocess

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os/exec"
	"syscall"
	"time"
)

type PrometheusProcess struct {
	baseUrl    string
	basePath   string
	configPath string
	state      string
}

func NewPrometheusProcess() *PrometheusProcess {
	prometheusBaseUrl := viper.GetString("prometheus.prometheusUrl")
	prometheusBasePath := viper.GetString("prometheus.prometheusBasePath")
	prometheusConfigPath := viper.GetString("prometheus.prometheusConfigPath")

	return &PrometheusProcess{
		prometheusBaseUrl,
		prometheusBasePath,
		prometheusConfigPath,
		"down",
	}
}

func (p *PrometheusProcess) restart() error {
	err := p.kill()
	if err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	err = p.start()
	if err != nil {
		return err
	}
	return nil
}

func (p *PrometheusProcess) kill() error {
	killPromCommand := "killall"
	killPromCommandArgs := []string{"prometheus"}
	cmd := exec.Command(killPromCommand, killPromCommandArgs...)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.Dir = p.basePath
	cmd.SysProcAttr.Setsid = true
	cmd.SysProcAttr.Foreground = false
	cmd.Stdout = nil
	cmd.Stdin = nil

	err := cmd.Start()

	if err != nil {
		return errors.New(fmt.Sprintf("kill promtheus err: %s", err))
	}

	return nil
}

func (p *PrometheusProcess) start() error {

	prometheusCommand := fmt.Sprintf("./prometheus")
	prometheusCommandArgs := []string{
		fmt.Sprintf("--config.file=%v", p.configPath),
		"--web.enable-lifecycle",
		"--web.enable-admin-api",
	}

	logrus.Infof("exec command: %v", prometheusCommand)
	cmd := exec.Command(prometheusCommand, prometheusCommandArgs...)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.Dir = p.basePath
	cmd.SysProcAttr.Setsid = true
	cmd.SysProcAttr.Foreground = false
	cmd.Stdout = nil
	cmd.Stdin = nil

	err := cmd.Start()

	if err != nil {
		return errors.New(fmt.Sprintf("execution err: %s", err))
	}

	return nil

}
