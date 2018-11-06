package main

import (
	"fmt"
	"os"

	"github.com/kylelemons/go-gypsy/yaml"
)

type configuartion struct {
	queuePath  string
	maxWorkers int
	httpPort   int
}

func (this configuartion) String() string {
	return fmt.Sprintf("queuePath: %s, maxWorkers: %d, httpPort: %d", this.queuePath, this.maxWorkers, this.httpPort)
}

func readGlobalConfig() configuartion {
	path := os.Getenv("QUEUEKEEPER_CONFIG_PATH")
	if "" == path {
		path = "qk.config.yml"
	}

	conf := configuartion{queuePath: "./", maxWorkers: 5, httpPort: 8088}

	config, err := yaml.ReadFile(path)

	if nil == err {
		qp, err := config.Get("queue_config_path")
		if nil == err {
			conf.queuePath = qp
		}
		mw, err := config.GetInt("max_workers")
		if nil == err {
			conf.maxWorkers = int(mw)
		}
		hp, err := config.GetInt("http_port")
		if nil == err {
			conf.httpPort = int(hp)
		}
	}
	return conf
}
