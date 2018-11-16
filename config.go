package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"queuekeeper/qs"
	"strconv"
	"strings"

	"github.com/kylelemons/go-gypsy/yaml"
)

const ENV_PREFIX = "QUEUEKEEPER_"

type configuartion struct {
	queuePath  string
	maxWorkers int
	httpPort   int
}

func (this configuartion) String() string {
	return fmt.Sprintf("queuePath: %s, maxWorkers: %d, httpPort: %d", this.queuePath, this.maxWorkers, this.httpPort)
}

func readEnv(name string) string {
	return os.Getenv(ENV_PREFIX + name)
}

func readGlobalConfig() configuartion {
	path := readEnv("CONFIG_PATH")
	if "" == path {
		path = "qk.config.yml"
	}

	conf := configuartion{queuePath: "./", maxWorkers: 5, httpPort: 8088}

	if "ENV" == path {
		return readGlobalConfigFromEnv(conf)
	}

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

func readGlobalConfigFromEnv(conf configuartion) configuartion {
	qp := readEnv("QUEUE_CONFIG_PATH")
	if "" != qp {
		conf.queuePath = qp
	}
	mw := os.Getenv("GOMAXPROCS")
	if "" != mw {
		mw, err := strconv.ParseInt(mw, 10, 0)
		if nil == err {
			conf.maxWorkers = int(mw)
		}
	}
	hp := readEnv("HTTP_PORT")
	if "" != hp {
		hp, err := strconv.ParseInt(mw, 10, 0)
		if nil == err {
			conf.httpPort = int(hp)
		}
	}
	return conf
}

func readQueuesConfigs(qm *qs.QueueManager, conf configuartion) *qs.QueueManager {
	files, err := ioutil.ReadDir(conf.queuePath)
	if err != nil {
		log.Fatal(err)
	}
	var fn string
	for _, file := range files {
		fn = file.Name()
		if strings.HasSuffix(fn, ".yml") || strings.HasSuffix(fn, ".yaml") {
			c, err := yaml.ReadFile(conf.queuePath + fn)
			if err != nil {
				log.Println(err)
				continue
			}
			processOneQueueConfig(fn, c, qm)

		}
	}
	return qm
}

func processOneQueueConfig(fn string, conf *yaml.File, qm *qs.QueueManager) *qs.QueueManager {
	var name string
	if strings.HasSuffix(fn, ".yml") {
		name = strings.Replace(fn, ".yml", "", -1)
	} else if strings.HasSuffix(fn, ".yaml") {
		name = strings.Replace(fn, ".yaml", "", -1)
	} else {
		name = fn
	}
	n, err := conf.Get("name")
	if nil == err {
		name = n
	}

	_, err = qm.GetQueue(name)

	if nil == err {
		//queue already exists
		return qm
	}

	flags := qs.NewQueueFlags()

	delayDelivery, err := conf.GetInt("delay_delivery")
	if nil == err {
		flags.SetDelayedDelivery(int(delayDelivery))
	}

	processFlag(conf, "deduplication", flags.SetDeduplicated)

	qm.CreateQueue(name, flags)
	return qm
}

func processFlag(conf *yaml.File, name string, setter qs.QueueFlagsSetter) bool {
	fl, err := conf.GetBool("flags." + name)
	if nil == err {
		setter(fl)
	}
	return fl
}
