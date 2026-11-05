package voice

import (
	"context"
	"sync"
)

// utterance buffers a single Muse request. The official direct pipeline owns its
// own buffer now; keeping this separate avoids regressing its lifecycle.
type utterance struct {
	mu     sync.Mutex
	pcm    []byte
	cancel context.CancelFunc
	ctx    context.Context
}

func (u *utterance) Start(string) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.cancel != nil {
		u.cancel()
	}
	u.ctx, u.cancel = context.WithCancel(context.Background())
	u.pcm = u.pcm[:0]
	return nil
}

func (u *utterance) Audio(frame []byte) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.ctx == nil || len(u.pcm)+len(frame) > maxUtterance {
		return nil
	}
	u.pcm = append(u.pcm, frame...)
	return nil
}

func (u *utterance) take() (context.Context, []byte) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.ctx == nil {
		return nil, nil
	}
	pcm := append([]byte(nil), u.pcm...)
	u.pcm = u.pcm[:0]
	return u.ctx, pcm
}

func (u *utterance) Stop() error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.cancel != nil {
		u.cancel()
	}
	u.cancel, u.ctx = nil, nil
	u.pcm = u.pcm[:0]
	return nil
}

// samplesOf and bytesOf convert the microphone's little-endian PCM16 frames.
func samplesOf(pcm []byte) []int16 {
	samples := make([]int16, len(pcm)/2)
	for i := range samples {
		samples[i] = int16(uint16(pcm[2*i]) | uint16(pcm[2*i+1])<<8)
	}
	return samples
}

func bytesOf(samples []int16) []byte {
	out := make([]byte, 2*len(samples))
	for i, sample := range samples {
		out[2*i], out[2*i+1] = byte(sample), byte(sample>>8)
	}
	return out
}
