package model

import (
	"time"
)

type Relationship struct {
	URI      string `json:"uri"`
	Relation string `json:"relation"`
}

type Meta struct {
	ID           string         `json:"id"`
	Type         string         `json:"type"`
	CreatedAt    time.Time      `json:"created_at"`
	By           Identity       `json:"by"`
	Relationship []Relationship `json:"relationship"`
}
