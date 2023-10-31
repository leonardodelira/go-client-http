package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/leonardodelira/go-lib-http/deployment"
	"github.com/leonardodelira/go-lib-http/errors"
	"github.com/leonardodelira/go-lib-http/kubeclient"
	"golang.org/x/exp/slices"
)

const url = "http://localhost:3000"

func TestCreateDeployment(t *testing.T) {
	c, err := kubeclient.NewClient(
		kubeclient.WithURL(url),
	)
	if err != nil {
		t.Errorf("should not fail on create new client: %s", err)
	}

	deployment := deployment.Deployment{
		ID:       uuid.New(),
		Replicas: 1,
		Image:    "nginx",
		Ports: []deployment.Port{
			{
				Name:   "http",
				Number: 80,
			},
		},
	}

	createdDeploy, err := c.Deployment.Create(context.Background(), deployment)
	if err != nil {
		t.Errorf("should not fail to create deployment: %s", err)
		return
	}

	assertDeployment(t, &deployment, createdDeploy)
}

func TestCreateDeploymentWithShortTimeout(t *testing.T) {
	c, err := kubeclient.NewClient(
		kubeclient.WithURL(url),
	)
	if err != nil {
		t.Errorf("should not fail on create new client: %s", err)
	}

	deployment := deployment.Deployment{
		ID:       uuid.New(),
		Replicas: 1,
		Image:    "nginx",
		Ports: []deployment.Port{
			{
				Name:   "http",
				Number: 80,
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	_, err = c.Deployment.Create(ctx, deployment)
	if err == nil {
		t.Errorf("should fail due to a timeout")
		return
	}
}

func TestCreateNonValidDeployment(t *testing.T) {
	c, err := kubeclient.NewClient(
		kubeclient.WithURL(url),
	)
	if err != nil {
		t.Errorf("should not fail on create new client: %s", err)
	}

	cases := map[string]struct {
		deployment   deployment.Deployment
		failedFields []string
	}{
		"empty deployment": {
			deployment:   deployment.Deployment{},
			failedFields: []string{"id", "replicas", "image", "ports"},
		},
		"invalid image": {
			deployment: deployment.Deployment{
				ID:       uuid.New(),
				Image:    "", // invalid image
				Replicas: 1,
				Ports: []deployment.Port{
					{
						Name:   "http",
						Number: 80,
					},
				},
			},
			failedFields: []string{"image"},
		},
		"replicas empty": {
			deployment: deployment.Deployment{
				ID:       uuid.New(),
				Image:    "nginx", // invalid image
				Replicas: 0,
				Ports: []deployment.Port{
					{
						Name:   "http",
						Number: 80,
					},
				},
			},
			failedFields: []string{"replicas"},
		},
	}

	for title, v := range cases {
		t.Run(title, func(t *testing.T) {
			_, err := c.Deployment.Create(context.Background(), v.deployment)
			if err == nil {
				t.Errorf("should fail with invalid body")
				return
			}

			if _, ok := err.(errors.InvalidResource); !ok {
				t.Errorf("should fail with InvalidResource error: %s", err)
				return
			}

			invalidResource := err.(errors.InvalidResource)
			if !slices.Equal(invalidResource.FailedFields, v.failedFields) {
				t.Errorf("should fail with %v failed fields: %v", v.failedFields, invalidResource.FailedFields)
				return
			}
		})
	}
}

func TestCreateDuplicatedDeployment(t *testing.T) {
	c, err := kubeclient.NewClient(
		kubeclient.WithURL(url),
	)
	if err != nil {
		t.Errorf("should not fail on create new client: %s", err)
	}

	deployment := deployment.Deployment{
		ID:       uuid.New(),
		Replicas: 1,
		Image:    "nginx",
		Ports: []deployment.Port{
			{
				Name:   "http",
				Number: 80,
			},
		},
	}

	_, err = c.Deployment.Create(context.Background(), deployment)
	if err != nil {
		t.Errorf("not should fail on create first deployment")
		return
	}

	_, err = c.Deployment.Create(context.Background(), deployment)
	if err == nil {
		t.Errorf("should fail on create duplicated deployment")
		return
	}
}

func TestDeleteDeployment(t *testing.T) {
	c, err := kubeclient.NewClient(
		kubeclient.WithURL(url),
	)
	if err != nil {
		t.Errorf("should not fail on create new client: %s", err)
	}

	deployment := deployment.Deployment{
		ID:       uuid.New(),
		Replicas: 1,
		Image:    "nginx",
		Ports: []deployment.Port{
			{
				Name:   "http",
				Number: 80,
			},
		},
	}

	_, err = c.Deployment.Create(context.Background(), deployment)
	if err != nil {
		t.Errorf("should not fail to create deployment: %s", err)
		return
	}

	err = c.Deployment.Delete(context.Background(), deployment.ID)
	if err != nil {
		t.Errorf("should not fail on delete a deployment")
	}
}

func TestGetDeployment(t *testing.T) {
	c, err := kubeclient.NewClient(
		kubeclient.WithURL(url),
	)
	if err != nil {
		t.Errorf("should not fail on create new client: %s", err)
	}

	deployment := deployment.Deployment{
		ID:       uuid.New(),
		Replicas: 1,
		Image:    "nginx",
		Ports: []deployment.Port{
			{
				Name:   "http",
				Number: 80,
			},
		},
	}

	_, err = c.Deployment.Create(context.Background(), deployment)
	if err != nil {
		t.Errorf("should not fail to create deployment: %s", err)
		return
	}

	cases := map[string]struct {
		ID    uuid.UUID
		Error string
	}{
		"success get deployment": {
			ID:    deployment.ID,
			Error: "",
		},
		"error on get deployment": {
			ID:    uuid.New(),
			Error: fmt.Sprintf("StatusCode=%d, Message=%s", 404, "not found"),
		},
	}

	for title, v := range cases {
		t.Run(title, func(t *testing.T) {
			_, err = c.Deployment.Get(context.Background(), v.ID)

			if v.Error == "" {
				if err != nil {
					t.Errorf("should not fail on get a deployment")
					return
				}
			}

			if v.Error != "" {
				if err == nil {
					t.Errorf("want: %v", v.Error)
					return
				}
			}
		})
	}
}

func assertDeployment(t *testing.T, expectedDeploy, foundDeploy *deployment.Deployment) {
	if foundDeploy.ID != expectedDeploy.ID {
		t.Errorf("should have same ID: %s", foundDeploy.ID)
		return
	}

	if foundDeploy.Image != expectedDeploy.Image {
		t.Errorf("should have same Image: %s", foundDeploy.Image)
		return
	}

	if foundDeploy.Replicas != expectedDeploy.Replicas {
		t.Errorf("should have same Replicas: %d", foundDeploy.Replicas)
		return
	}

	if !foundDeploy.CreatedAt.IsZero() {
		t.Errorf("should have CreatedAt set")
		return
	}

	if len(foundDeploy.Ports) != len(expectedDeploy.Ports) {
		t.Error("should have same number of ports")
		return
	}

	for i, p := range foundDeploy.Ports {
		if p.Number != expectedDeploy.Ports[i].Number {
			t.Errorf("should have same port number: %d", p.Number)
			return
		}
	}
}
