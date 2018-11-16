package main

import (
	"encoding/json"
)

type message struct {
	msg   string `json:"message"`
	delay string `json:"delay"`
}

type postJson struct {
	Msgs   []message `json:"messages"`
	Cnt    int       `json:"count"`
	Action string    `json:"method"`
}

func resolveJson(jsonString []byte) (postJson, error) {
	result := postJson{}
	err := json.Unmarshal(jsonString, &result)
	if err == nil {
		return result, nil
	}
	return result, err
}
