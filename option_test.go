package cache

import (
	"reflect"
	"testing"
	"time"
)

func TestWithEx(t *testing.T) {
	type args struct {
		key string
		v interface{}
		opt SetIOption
	}
	tests := []struct {
		name  string
		args  args
		sleep time.Duration
		want  bool
	}{
		{name: "int", args: args{key: "intWithEx", v: 1,  opt: WithEx(10 * time.Millisecond)}, sleep: 0, want: true},
		{name: "int", args: args{key: "intWithEx", v: 1,  opt: WithEx(10 * time.Millisecond)}, sleep: 10 * time.Millisecond, want: false},
		{name: "int", args: args{key: "intWithEx",  v: 1, opt: WithEx(100 * time.Millisecond)}, sleep: 50 * time.Millisecond, want: true},
	}
	mockCache()
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
	mockCache()
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

func TestWithNx(t *testing.T) {
	type args struct {
		key string
		v   interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "nullWithNx", args: args{key: "nullWithNx"}, want: true},
		{name: "int", args: args{key: "int", v: 1}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Set(tt.args.key, "v", WithNx()); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithNx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithXx(t *testing.T) {
	type args struct {
		key string
		v   interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "int", args: args{key: "int", v: 1}, want: true},
		{name: "nullWithXx", args: args{key: "nullWithXx"}, want: false},
	}
	mockCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Set(tt.args.key, tt.args.v, WithXx()); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithXx() = %v, want %v", got, tt.want)
			}
		})
	}
}
