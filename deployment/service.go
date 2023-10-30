package deployment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/leonardodelira/go-lib-http/errors"
)

type Service struct {
	client *http.Client
	url    string
}

func NewService(client *http.Client, url string) Service {
	return Service{
		client: client,
		url:    url,
	}
}

func (s *Service) Create(ctx context.Context, deployment Deployment) (*Deployment, error) {
	j, err := json.Marshal(deployment)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/deployments", s.url)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return nil, errors.FromBadRequest(resp)
	}

	if resp.StatusCode == http.StatusConflict {
		return nil, errors.FromHTTPResponse(resp)
	}

	createdDeploy := Deployment{}
	if err := json.NewDecoder(resp.Body).Decode(&createdDeploy); err != nil {
		return nil, err
	}

	return &createdDeploy, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return nil
}
