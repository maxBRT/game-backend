package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()
	assert.NotNil(t, config)
	assert.Equal(t, "http://localhost:5108", config.PlayerServiceURL)
	assert.Equal(t, "http://localhost:8000", config.MatchMakingServiceURL)
	assert.Equal(t, 100, config.MaxIdleConns)
	assert.Equal(t, 100, config.MaxIdleConnsPerHost)
	assert.Equal(t, 0, config.MaxConnsPerHost)
}
