package deployment

import (
	"time"

	"github.com/google/uuid"
)

type Deployment struct {
	ID        uuid.UUID         `json:"id"`
	Replicas  int               `json:"replicas"`
	Image     string            `json:"image"`
	Labels    map[string]string `json:"labels"`
	Ports     []Port            `json:"ports"`
	CreatedAt time.Time
}

type Port struct {
	Name   string `json:"name"`
	Number int    `json:"port"`
}
