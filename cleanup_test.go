package cache

import (
	"testing"
	"time"
)

var cw ICleanupWorker

func TestRingBufferWheel_Register(t *testing.T) {
	now := time.Now()
	type args struct {
		key      string
		expireAt time.Time
	}
	tests := []struct {
		name        string
		args        args
		wantIndex   int
		wantCounter int
	}{
		{name: "1", args: args{key: "1", expireAt: time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 1, 0, now.Location())}, wantIndex: 2, wantCounter: 0},
		{name: "2", args: args{key: "2", expireAt: time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 59, 0, now.Location())}, wantIndex: 0, wantCounter: 0},
	}
	cw := NewRingBufferWheel(c)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cw.Register(tt.args.key, tt.args.expireAt)
			if cw.buffers[tt.wantIndex].next.key != tt.args.key {
				t.Errorf("Register() got key = %v, want %v", cw.buffers[tt.wantIndex].next.key, tt.args.key)
			}
			if cw.buffers[tt.wantIndex].next.counter != tt.wantCounter {
				t.Errorf("Register() got counter = %v, want %v", cw.buffers[tt.wantIndex].next.counter, tt.wantCounter)
			}
		})
	}
}

func TestRingBufferWheel_Run(t *testing.T) {
	type args struct {
		k string
		d time.Duration
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "1", args: args{k: "run1", d: 0 * time.Second}, want: false},
		{name: "2", args: args{k: "run2", d: 1 * time.Second}, want: true},
	}
	c := NewMemCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()
			c.Set(tt.args.k, "v", WithEx(tt.args.d))
			time.Sleep(now.Add(1100 * time.Millisecond).Sub(time.Now()))
			c.rw.RLock()
			_, got := c.m[tt.args.k]
			c.rw.RUnlock()
			if got != tt.want {
				t.Errorf("Run() got = %v, want %v", got, tt.want)
			}
		})
	}
}
