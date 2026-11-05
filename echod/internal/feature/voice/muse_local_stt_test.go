package voice

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/HuskerMinion/techo5/echod/internal/config"
	"github.com/HuskerMinion/techo5/echod/internal/lib/muse"
)

type localSTTLinkFake struct {
	voiceCalls int
	textCalls  int
	text       string
	lastWAV    []byte
}

func (*localSTTLinkFake) Ready() bool  { return true }
func (*localSTTLinkFake) Paired() bool { return true }
func (f *localSTTLinkFake) Ask(_ context.Context, wav []byte, _ func(muse.ReplyEvent)) (muse.Reply, error) {
	f.voiceCalls++
	f.lastWAV = append([]byte(nil), wav...)
	return muse.Reply{Text: "voice response"}, nil
}
func (f *localSTTLinkFake) AskText(_ context.Context, text string, _ func(muse.ReplyEvent)) (muse.Reply, error) {
	f.textCalls++
	f.text = text
	return muse.Reply{Text: "text response"}, nil
}

func TestMuseLocalSTTSendsTextNotAudio(t *testing.T) {
	link := &localSTTLinkFake{}
	m := &viaMuse{link: link, recognize: func(_ context.Context, b config.Brain, pcm []byte) (string, error) {
		if b.Language != "ja" || len(pcm) != 4 {
			t.Errorf("unexpected STT input: %+v, %d", b, len(pcm))
		}
		return "おはようございます", nil
	}}
	var heard string
	reply, err := m.askWithBrain(context.Background(), config.Brain{STT: "127.0.0.1:50024", Language: "ja"}, []byte{1, 0, 2, 0}, func(e muse.ReplyEvent) {
		if h, ok := e.(muse.Heard); ok {
			heard = h.Text
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	if reply.Text != "text response" || heard != "おはようございます" || link.text != "おはようございます" || link.textCalls != 1 || link.voiceCalls != 0 {
		t.Fatalf("text=%q heard=%q reply=%q textCalls=%d voiceCalls=%d", link.text, heard, reply.Text, link.textCalls, link.voiceCalls)
	}
}

func TestMuseConfiguredSTTFailClosed(t *testing.T) {
	for _, tc := range []struct {
		name string
		text string
		err  error
	}{
		{name: "offline", err: errors.New("offline")},
		{name: "empty", text: "   "},
	} {
		t.Run(tc.name, func(t *testing.T) {
			link := &localSTTLinkFake{}
			m := &viaMuse{link: link, recognize: func(context.Context, config.Brain, []byte) (string, error) { return tc.text, tc.err }}
			reply, err := m.askWithBrain(context.Background(), config.Brain{STT: "127.0.0.1:50024", Language: "ja"}, []byte{1, 0}, func(muse.ReplyEvent) {})
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(reply.Text, "もう一度") || link.textCalls != 0 || link.voiceCalls != 0 {
				t.Fatalf("fallback must stay on device: response=%q textCalls=%d voiceCalls=%d", reply.Text, link.textCalls, link.voiceCalls)
			}
		})
	}
}

func TestMuseWithoutSTTUsesVoiceNote(t *testing.T) {
	link := &localSTTLinkFake{}
	m := &viaMuse{link: link, recognize: func(context.Context, config.Brain, []byte) (string, error) {
		t.Fatal("should not invoke recognizer")
		return "", nil
	}}
	_, err := m.askWithBrain(context.Background(), config.Brain{}, []byte{1, 0, 2, 0}, func(muse.ReplyEvent) {})
	if err != nil {
		t.Fatal(err)
	}
	if link.voiceCalls != 1 || link.textCalls != 0 || len(link.lastWAV) < 12 || string(link.lastWAV[:4]) != "RIFF" {
		t.Fatalf("expected one WAV voice note: %+v", link)
	}
}
