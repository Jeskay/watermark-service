package endpoints

import (
	"image"
	"watermark-service/internal"
)

type GetRequest struct {
	Filters []internal.Filter `json:"filters,omitempty"`
}

type GetResponse struct {
	Documents []internal.Document `json:"documents"`
	Err       string              `json:"err,omitempty"`
}

type AddRequest struct {
	Logo  image.Image       `json:"logo"`
	Image image.Image       `json:"image"`
	Text  string            `json:"text"`
	Fill  bool              `json:"fill"`
	Pos   internal.Position `json:"pos"`
}

type AddResponse struct {
	TicketID string `json:"ticketID"`
	Err      string `json:"err,omitempty"`
}

type RemoveRequest struct {
	TicketID string `json:"ticketID"`
}

type RemoveResponse struct {
	Code int    `json:"code"`
	Err  string `json:"err,omitempty"`
}

type ServiceStatusRequest struct{}

type ServiceStatusResponse struct {
	Code int    `json:"code"`
	Err  string `json:"err,omitempty"`
}
