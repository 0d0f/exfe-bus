package model

import (
	"fmt"
	"time"
)

type Relationship struct {
	URI      string `json:"uri"`
	Relation string `json:"relation"`
}

func (r Relationship) String() string {
	return fmt.Sprintf("{{%s:%s}}", r.Relation, r.URI)
}

type Meta struct {
	URI          string         `json:"uri"`
	CreatedAt    time.Time      `json:"created_at"`
	By           Identity       `json:"by"`
	Relationship []Relationship `json:"relationship"`
}
