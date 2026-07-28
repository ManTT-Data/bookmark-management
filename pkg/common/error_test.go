package common

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleError(t *testing.T) {
	HandleError(nil)
	assert.Panics(t, func() { HandleError(errors.New("test error")) }, "test error")
}
