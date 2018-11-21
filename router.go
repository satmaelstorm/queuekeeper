package main

import (
	"github.com/julienschmidt/httprouter"
)

func route() *httprouter.Router {
	router := httprouter.New()
	router.GET("/q/:queue", getFromQueueHandler)
	router.PUT("/q/:queue", putToQueueHandler)
	router.POST("/q/:queue", postQueueHandler)
	router.GET("/admin/reload/queues", adminReloadQueueConfigHandler)
	return router
}
