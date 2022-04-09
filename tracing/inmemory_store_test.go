package tracing

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_can_save_and_read(t *testing.T) {
	store, err := NewInMemoryStore(&Config{
		KeepOriginalPayload: true,
	})
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

func Test_time_range(t *testing.T) {
	store, err := NewInMemoryStore(EmptyConfig())
	assert.Nil(t, err)
	from := time.Now().UTC().Add(-10 * 24 * time.Hour)
	to := time.Now().UTC().Add(-1 * 24 * time.Hour)
	for _, trc := range getRandomTraces(100, &from, &to) {
		err = store.Store(trc, "")
		assert.Nil(t, err)
	}

	fromBefore := from.Add(-1 * time.Hour)
	toAfter := to.Add(1 * time.Hour)
	trcBefore := getTrace(&fromBefore)
	trcAfter := getTrace(&toAfter)

	_ = store.Store(trcBefore, "")
	_ = store.Store(trcAfter, "")

	tracesBack, err := store.ListByTimeRange(42, &from, &to)
	assert.Nil(t, err)
	assert.Equal(t, 42, len(tracesBack))

}

func Test_time_range_reverse(t *testing.T) {
	store, err := NewInMemoryStore(EmptyConfig())
	assert.Nil(t, err)
	from := time.Now().UTC().Add(-10 * 24 * time.Hour)
	to := time.Now().UTC().Add(-1 * 24 * time.Hour)
	for _, trc := range getRandomTraces(100, &from, &to) {
		err = store.Store(trc, "")
		assert.Nil(t, err)
	}

	fromBefore := from.Add(-1 * time.Hour)
	toAfter := to.Add(1 * time.Hour)
	trcBefore := getTrace(&fromBefore)
	trcAfter := getTrace(&toAfter)

	afterAfter := toAfter.Add(1 * time.Hour)

	_ = store.Store(trcBefore, "")
	_ = store.Store(trcAfter, "")

	tracesBack, err := store.ListByTimeRange(42, nil, &afterAfter)
	assert.Nil(t, err)
	assert.Equal(t, 42, len(tracesBack))
	assert.Equal(t, trcAfter.TraceId, tracesBack[0].TraceId)

}

func getRandomTime(from, to *time.Time) time.Time {
	var u1, u2 int64
	if from == nil {
		u1 = 0
	} else {
		u1 = from.Unix()
	}
	if from == nil {
		u2 = time.Now().UTC().Unix()
	} else {
		u2 = to.Unix()
	}

	return time.Unix(rand.Int63n(u2-u1)+u1, 0)
}

func getRandomTraces(n int, from, to *time.Time) []*Trace {
	traces := make([]*Trace, n)
	for i := 0; i < n; i++ {
		traces[i] = getRandomTrace(from, to)
	}
	return traces
}

func getRandomTrace(from, to *time.Time) *Trace {
	t := getRandomTime(from, to)
	trc := getTrace(&t)
	return trc
}

func getTrace(t *time.Time) *Trace {
	trc := NewTrace(*t, "hello", "12345", "info")
	return trc
}
