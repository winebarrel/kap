package kap_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/kap"
)

func TestOptionsAfterApply_OK(t *testing.T) {
	assert := assert.New(t)

	options := kap.Options{
		Secret: []string{"foo", "bar"},
	}

	assert.NoError(options.AfterApply())
}

func TestOptionsAfterApply_Err(t *testing.T) {
	assert := assert.New(t)

	options := kap.Options{
		Secret: []string{"foo", "", "bar"},
	}

	assert.ErrorContains(options.AfterApply(), "contains an empty secret value")
}
