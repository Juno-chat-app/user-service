package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReader_LoadConfiguration(t *testing.T) {
	conf, err := LoadConfiguration("../user-service_test.yml")
	require.Nil(t, err)
	require.Equal(t, conf.GRPCConfig.Host, "localhost")
}
