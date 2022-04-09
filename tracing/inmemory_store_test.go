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
	trc, err := store.GetById(trace.TraceId)
	assert.Nil(t, err)
	assert.Equal(t, "hello", trc.Message)
}

func Test_non_existent(t *testing.T) {
	store, err := NewInMemoryStore(EmptyConfig())
	assert.Nil(t, err)

	store.Store(NewTrace(time.Now(), "hello", "12345", "info"), "")

	trc, err := store.GetById("NonExistent")
	assert.Nil(t, err)
	assert.Nil(t, trc)
}
