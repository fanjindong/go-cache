package cache

import (
	"testing"
	"time"
)

func TestWithEx(t *testing.T) {
	type args struct {
		key string
		v   interface{}
		opt SetIOption
	}
	tests := []struct {
		name  string
		args  args
		sleep time.Duration
		want  bool
	}{
		{name: "int", args: args{key: "intWithEx", v: 1, opt: WithEx(10 * time.Millisecond)}, sleep: 0, want: true},
		{name: "int", args: args{key: "intWithEx", v: 1, opt: WithEx(10 * time.Millisecond)}, sleep: 10 * time.Millisecond, want: false},
		{name: "int", args: args{key: "intWithEx", v: 1, opt: WithEx(100 * time.Millisecond)}, sleep: 50 * time.Millisecond, want: true},
	}
	c := mockCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.Set(tt.args.key, "v", tt.args.opt)
			time.Sleep(tt.sleep)
			_, got := c.Get(tt.args.key)
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithExAt(t *testing.T) {
	type args struct {
		key string
		v   interface{}
		opt SetIOption
	}
	tests := []struct {
		name  string
		args  args
		sleep time.Duration
		want  bool
	}{
		{name: "int", args: args{key: "int", v: 1, opt: WithExAt(time.Now().Add(10 * time.Millisecond))}, sleep: 0, want: true},
		{name: "int", args: args{key: "int", v: 1, opt: WithExAt(time.Now().Add(10 * time.Millisecond))}, sleep: 10 * time.Millisecond, want: false},
		{name: "int", args: args{key: "int", v: 1, opt: WithExAt(time.Now().Add(100 * time.Millisecond))}, sleep: 50 * time.Millisecond, want: true},
	}
	c := mockCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.Set(tt.args.key, tt.args.v, tt.args.opt)
			time.Sleep(tt.sleep)
			_, got := c.Get(tt.args.key)
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type ringBufferWheelHello struct {
	*RingBufferWheel
}

func (r *ringBufferWheelHello) Register(key string, expireAt time.Time) {
	r.RingBufferWheel.Register(key, expireAt)
}

func TestWithCleanup(t *testing.T) {

	type args struct {
		cw ICleanupWorker
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "int", args: args{cw: &ringBufferWheelHello{NewRingBufferWheel()}}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewMemCache(WithCleanup(tt.args.cw))
			c.Set("a", 1, WithEx(100*time.Millisecond))
			time.Sleep(1 * time.Second)
		})
	}
}
