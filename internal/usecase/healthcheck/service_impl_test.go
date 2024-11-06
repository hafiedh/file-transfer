package healthcheck

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var healthCheckService = NewService()

func TestHealthCheck(t *testing.T) {
	_, err := healthCheckService.HealthCheck(context.Background())
	if err != nil {
		assert.Error(t, err)
	}

	assert.NoError(t, err)
	assert.NotNil(t, healthCheckService)
}
