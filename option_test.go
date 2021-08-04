package cache

import (
	"reflect"
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

func TestWithShards(t *testing.T) {
	type args struct {
		shards int
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "1", args: args{shards: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewMemCache(WithShards(tt.args.shards))
			if len(c.(*MemCache).shards) != tt.args.shards {
				t.Errorf("WithShards() = %v, want %v", len(c.(*MemCache).shards), tt.args.shards)
			}
		})
	}
}

func TestWithExpiredCallback(t *testing.T) {
	var c ICache
	type args struct {
		do func()
		ec ExpiredCallback
	}
	tests := []struct {
		name string
		args args
		want func() bool
	}{
		{name: "1",
			args: args{
				do: func() {
					c.Set("1", 1, WithEx(100*time.Millisecond))
					time.Sleep(1100 * time.Millisecond)
				},
				ec: func(k string, v interface{}) error {
					c.Set(k, 2)
					return nil
				}},
			want: func() bool {
				v, ok := c.Get("1")
				if !ok {
					return false
				}
				return 2 == v.(int)
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c = NewMemCache(WithExpiredCallback(tt.args.ec))
			tt.args.do()
			if !tt.want() {
				t.Errorf("WithExpiredCallback() = %v, want %v", tt.want(), true)
			}
		})
	}
}

type hash1 struct {
}

func (h hash1) Sum64(s string) uint64 {
	return 0
}

func TestWithHash(t *testing.T) {
	var c ICache
	type args struct {
		hash IHash
	}
	tests := []struct {
		name string
		args args
		do   func()
		got  func() int
		want int
	}{
		{name: "1", args: args{hash: hash1{}},
			do: func() {
				c.Set("a", 1)
				c.Set("b", 1)
				c.Set("c", 1)
				c.Set("d", 1)
			},
			got: func() int {
				return len(c.(*MemCache).shards[0].hashmap)
			},
			want: 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c = NewMemCache(WithHash(tt.args.hash))
			tt.do()
			if got := tt.got(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithClearInterval(t *testing.T) {
	var c ICache
	type args struct {
		d time.Duration
	}
	tests := []struct {
		name string
		args args
		got  func() interface{}
		want interface{}
	}{
		{name: "1", args: args{d: 100 * time.Millisecond},
			got: func() interface{} {
				c.Set("ex", 1, WithEx(10*time.Millisecond))
				time.Sleep(50 * time.Millisecond)
				return c.(*MemCache).shards[0].hashmap["ex"].v
			}, want: 1},
		{name: "nil", args: args{d: 100 * time.Millisecond},
			got: func() interface{} {
				c.Set("ex", 1, WithEx(10*time.Millisecond))
				time.Sleep(110 * time.Millisecond)
				return c.(*MemCache).shards[0].hashmap["ex"].v
			}, want: nil},
		{name: "no expire", args: args{d: 0},
			got: func() interface{} {
				c.Set("ex", 1, WithEx(10*time.Millisecond))
				time.Sleep(110 * time.Millisecond)
				return c.(*MemCache).shards[0].hashmap["ex"].v
			}, want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c = mockCache(WithHash(hash1{}), WithClearInterval(tt.args.d))
			if got := tt.got(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithClearInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}
