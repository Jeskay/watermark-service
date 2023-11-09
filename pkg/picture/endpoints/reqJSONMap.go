package endpoints

import (
	"image"
	"watermark-service/internal"
)

type CreateRequest struct {
	Image image.Image       `json:"image"`
	Logo  image.Image       `json:"logo"`
	Text  string            `json:"text"`
	Fill  bool              `json:"fill"`
	Pos   internal.Position `json:"position"`
}

type CreateResponse struct {
	Image image.Image `json:"image"`
	Err   string      `json:"err,omitempty"`
}

type ServiceStatusRequest struct{}

type ServiceStatusResponse struct {
	Code int64  `json:"status"`
	Err  string `json:"err,omitempty"`
}
