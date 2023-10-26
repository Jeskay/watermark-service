package internal

import (
	uuid "github.com/google/uuid"
)

type Document struct {
	ID       uuid.UUID `json:"id"`
	AuthorId uuid.UUID `json:"author_id"`
	Title    string    `json:"title"`
	ImageUrl string    `json:"image_url"`
}

type Filter struct {
	Key   string `json:"key"`
	Value string `json:"value,omitempty"`
}

type Status string

const (
	Pending    Status = "Pending"
	Started    Status = "Started"
	InProgress Status = "InProgress"
	Finished   Status = "Finished"
	Failed     Status = "Failed"
)
