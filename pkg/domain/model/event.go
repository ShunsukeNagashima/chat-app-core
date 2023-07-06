package model

import "encoding/json"

type Event interface{}

type RawEvent struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
