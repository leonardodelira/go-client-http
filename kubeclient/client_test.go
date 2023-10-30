package kubeclient

import (
	"testing"
)

func TestDefaultClient(t *testing.T) {
	c, err := NewClient()
	if err != nil {
		t.Errorf("should not fail on create new client: %s", err)
	}

	if c.httpClient == nil {
		t.Error("expected httpClient to be set")
	}
}

func TestWithURL(t *testing.T) {
	url := "http://localhost:8080"
	c, err := NewClient(WithURL(url))
	if err != nil {
		t.Errorf("should not fail on create new client: %s", err)
	}

	if c.url != url {
		t.Errorf("expected url: %s", url)
	}
}

func TestWithInvalidURL(t *testing.T) {
	url := "nourl"
	_, err := NewClient(WithURL(url))
	if err == nil {
		t.Errorf("should fail on create new client with invalid url: %s", err)
	}
}
