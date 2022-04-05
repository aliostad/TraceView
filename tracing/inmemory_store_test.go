package tracing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_can_save_and_read(t *testing.T) {
	store, err := NewInMemoryStore(EmptyConfig())
	assert.Nil(t, err)
	trace := NewTrace(time.Now(), "hello", "12345", "info")
	err = store.Store(trace, "")
	assert.Nil(t, err)

}
