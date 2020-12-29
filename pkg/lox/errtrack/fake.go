package errtrack

import (
	"bytes"
)

type FakeTracker struct {
	Tracker *Tracker
	Buffer  bytes.Buffer
}

// NewFake creates a fake tracker.
// TODO Trigger test errors directly from the tracker
func NewFake() *FakeTracker {
	fake := FakeTracker{}
	fake.Tracker = &Tracker{
		hadError: false,
		output:   &fake.Buffer,
	}
	return &fake
}

func (f *FakeTracker) Errors() []byte {
	return f.Buffer.Bytes()
}
