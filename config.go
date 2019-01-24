package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"queuekeeper/qs"
	"strconv"
	"strings"

	"github.com/kylelemons/go-gypsy/yaml"
)

const ENV_PREFIX = "QUEUEKEEPER_"

type logConfiguration struct {
	level        int
	engine       string
	parsedEngine url.URL
}

type configuartion struct {
	queuePath  string
	maxWorkers int
	httpPort   int
	logConf    logConfiguration
	from       string
}

func (this configuartion) String() string {
	return fmt.Sprintf(
		"queuePath: %s, maxWorkers: %d, httpPort: %d, log level: %d, log engine: %s, load from: %s",
		this.queuePath, this.maxWorkers, this.httpPort, this.logConf.level, this.logConf.engine, this.from)
}

func readEnv(name string) string {
	return os.Getenv(ENV_PREFIX + name)
}

func readGlobalConfig() configuartion {
	path := readEnv("CONFIG_PATH")
	if "" == path {
		path = "qk.config.yml"
	}

	conf := configuartion{queuePath: "./", maxWorkers: 5, httpPort: 8088, from: "default"}

	if "ENV" == path {
		return readGlobalConfigFromEnv(conf)
	}

	config, err := yaml.ReadFile(path)

	if nil == err {
		conf.from = "yaml file "+path
		if qp, err := config.Get("queue_config_path"); nil == err {
			conf.queuePath = qp
		}

		if mw, err := config.GetInt("max_workers"); nil == err {
			conf.maxWorkers = int(mw)
		}

		if hp, err := config.GetInt("http_port"); nil == err {
			conf.httpPort = int(hp)
		}
		if logLevel, err := config.Get("log.level"); nil == err {
			conf.logConf.level = parseLogLevel(logLevel)
		} else {
			conf.logConf.level = QK_DEFAULT_LOG_LEVEL
		}

		conf.logConf.engine = QK_DEFAULT_LOG_ENGINE
		if logUri, err := config.Get("log.engine"); nil == err {
			if u, err := url.Parse(logUri); nil == err {
				conf.logConf.engine = logUri
				conf.logConf.parsedEngine = *u
			}
		}
	}
	return conf
}

func readGlobalConfigFromEnv(conf configuartion) configuartion {
	if qp := readEnv("QUEUE_CONFIG_PATH"); "" != qp {
		conf.queuePath = qp
	}

	if mw := os.Getenv("GOMAXPROCS"); "" != mw {
		mw, err := strconv.ParseInt(mw, 10, 0)
		if nil == err {
			conf.maxWorkers = int(mw)
		}
	}

	if hp := readEnv("HTTP_PORT"); "" != hp {
		if hp, err := strconv.ParseInt(hp, 10, 0); nil == err {
			conf.httpPort = int(hp)
		}
	}

	if logLevel := readEnv("LOG_LEVEL"); "" != logLevel {
		conf.logConf.level = parseLogLevel(logLevel)
	}

	if logUri := readEnv("LOG_ENGINE"); "" != logUri {
		if u, err := url.Parse(logUri); err == nil {
			conf.logConf.engine = logUri
			conf.logConf.parsedEngine = *u
		}
	}
	conf.from = "environment"
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

	delayDelivery, err := conf.GetInt("flags.delay_delivery")
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
