package main

import "encoding/json"

type message struct {
	msg   string `json: "message"`
	delay string `json: "delay"`
}

type postJson struct {
	msgs   []message `json: "messages"`
	cnt    int       `json: "count"`
	action string    `json: "method"`
}

func resolveJson(jsonString []byte) (postJson, error) {
	var result postJson
	err := json.Unmarshal(jsonString, &result)
	if err == nil {
		return result, nil
	}
	return result, err
}
