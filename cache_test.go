package cache

import (
	"os"
	"reflect"
	"runtime"
	"sort"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func mockCache() ICache {
	c := NewMemCache()
	c.Set("int", 1)
	c.Set("int32", int32(1))
	c.Set("int64", int64(1))
	c.Set("string", "a")
	c.Set("float64", 1.1)
	c.Set("float32", float32(1.1))
	c.Set("ex", 1, WithEx(1*time.Second))
	return c
}

func TestItem_Expired(t *testing.T) {
	type fields struct {
		v      interface{}
		expire time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "int", fields: fields{v: 1, expire: time.Now().Add(0 * time.Second)}, want: true},
		{name: "int32", fields: fields{v: 1, expire: time.Now().Add(1 * time.Second)}, want: false},
		{name: "int64", fields: fields{v: 1, expire: time.Now().Add(-1 * time.Second)}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Item{
				v:      tt.fields.v,
				expire: tt.fields.expire,
			}
			if got := i.Expired(); got != tt.want {
				t.Errorf("Expired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemCache_Del(t *testing.T) {
	type args struct {
		ks []string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "int", args: args{ks: []string{"int"}}, want: 1},
		{name: "int32,int64", args: args{ks: []string{"int32", "int64"}}, want: 2},
		{name: "string,null", args: args{ks: []string{"string", "null"}}, want: 1},
		{name: "null", args: args{ks: []string{"null"}}, want: 0},
	}
	c := mockCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Del(tt.args.ks...); got != tt.want {
				t.Errorf("Del() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemCache_Exists(t *testing.T) {
	type args struct {
		ks []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "int", args: args{ks: []string{"int"}}, want: true},
		{name: "int32,int64", args: args{ks: []string{"int32", "int64"}}, want: true},
		{name: "int64,null", args: args{ks: []string{"int64", "null"}}, want: false},
		{name: "null", args: args{ks: []string{"null"}}, want: false},
	}
	c := mockCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Exists(tt.args.ks...); got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemCache_Expire(t *testing.T) {
	type args struct {
		k string
		d time.Duration
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "int", args: args{k: "int", d: 1 * time.Second}, want: true},
		{name: "int32", args: args{k: "int32", d: 1 * time.Second}, want: true},
		{name: "null", args: args{k: "null", d: 1 * time.Second}, want: false},
	}
	c := mockCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Expire(tt.args.k, tt.args.d); got != tt.want {
				t.Errorf("Expire() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemCache_ExpireAt(t *testing.T) {
	type args struct {
		k string
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "int", args: args{k: "int", t: time.Now().Add(1 * time.Second)}, want: true},
		{name: "int32", args: args{k: "int32", t: time.Now().Add(1 * time.Second)}, want: true},
		{name: "null", args: args{k: "null", t: time.Now().Add(1 * time.Second)}, want: false},
	}
	c := mockCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.ExpireAt(tt.args.k, tt.args.t); got != tt.want {
				t.Errorf("ExpireAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemCache_Get(t *testing.T) {
	type args struct {
		k string
	}
	tests := []struct {
		name  string
		args  args
		want  interface{}
		want1 bool
	}{
		{name: "int", args: args{k: "int"}, want: 1, want1: true},
		{name: "int32", args: args{k: "int32"}, want: int32(1), want1: true},
		{name: "int64", args: args{k: "int64"}, want: int64(1), want1: true},
		{name: "string", args: args{k: "string"}, want: "a", want1: true},
		{name: "ex", args: args{k: "ex"}, want: 1, want1: true},
		{name: "null", args: args{k: "null"}, want: nil, want1: false},
	}
	c := mockCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := c.Get(tt.args.k)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMemCache_GetDel(t *testing.T) {
	type args struct {
		k string
	}
	tests := []struct {
		name  string
		args  args
		want  interface{}
		want1 bool
	}{
		{name: "int", args: args{k: "int"}, want: 1, want1: true},
		{name: "int", args: args{k: "int"}, want: nil, want1: false},
	}
	c := mockCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := c.GetDel(tt.args.k)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDel() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetDel() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMemCache_GetSet(t *testing.T) {
	type args struct {
		k    string
		v    interface{}
		opts []SetIOption
	}
	tests := []struct {
		name  string
		args  args
		want  interface{}
		want1 bool
	}{
		{name: "int", args: args{k: "int", v: 0}, want: 1, want1: true},
		{name: "int", args: args{k: "int", v: 1}, want: 0, want1: true},
		{name: "null", args: args{k: "null", v: 1}, want: nil, want1: false},
		{name: "null", args: args{k: "null", v: 0}, want: 1, want1: true},
	}
	c := mockCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := c.GetSet(tt.args.k, tt.args.v, tt.args.opts...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSet() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetSet() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMemCache_Keys(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{name: "int.*", args: args{pattern: "int.*"}, want: []string{"int", "int32", "int64"}, wantErr: false},
		{name: "string.*", args: args{pattern: "string.*"}, want: []string{"string"}, wantErr: false},
		{name: "float.*", args: args{pattern: "float.*"}, want: []string{"float32", "float64"}, wantErr: false},
		{name: "null", args: args{pattern: "^a.*"}, want: nil, wantErr: false},
		{name: "int32", args: args{pattern: "int32"}, want: []string{"int32"}, wantErr: false},
	}
	c := mockCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Keys(tt.args.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("Keys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			sort.Strings(got)
			sort.Strings(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keys() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemCache_Persist(t *testing.T) {
	type args struct {
		k string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "int", args: args{k: "int"}, want: true},
		{name: "ex", args: args{k: "ex"}, want: true},
		{name: "null", args: args{k: "null"}, want: false},
	}
	c := mockCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Persist(tt.args.k); got != tt.want {
				t.Errorf("Persist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemCache_PersistAndTtl(t *testing.T) {
	type args struct {
		k string
	}
	tests := []struct {
		name  string
		args  args
		want  time.Duration
		want1 bool
	}{
		{name: "int", args: args{k: "int"}, want: 0, want1: false},
		{name: "ex", args: args{k: "ex"}, want: 0, want1: false},
	}
	c := mockCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.Persist(tt.args.k)
			got, got1 := c.Ttl(tt.args.k)
			if got != tt.want {
				t.Errorf("Ttl() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Ttl() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMemCache_RandomKey(t *testing.T) {
	tests := []struct {
		name  string
		forN  int
		want  string
		want1 bool
	}{
		{name: "int", want: "", want1: true},
	}
	c := mockCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := c.RandomKey()
			t.Log(got, got1)
			//if got != tt.want {
			//	t.Errorf("RandomKey() got = %v, want %v", got, tt.want)
			//}
			if got1 != tt.want1 {
				t.Errorf("RandomKey() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMemCache_Rename(t *testing.T) {
	type args struct {
		k  string
		nk string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "int", args: args{k: "int", nk: "intNew"}, want: true},
		{name: "int32", args: args{k: "int32", nk: "int64"}, want: true},
		{name: "int64", args: args{k: "int64", nk: "string"}, want: true},
		{name: "null", args: args{k: "null", nk: "int32"}, want: false},
	}
	c := mockCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Rename(tt.args.k, tt.args.nk); got != tt.want {
				t.Errorf("Rename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemCache_RenameAndGet(t *testing.T) {
	type args struct {
		k  string
		nk string
	}
	tests := []struct {
		name  string
		args  args
		want  interface{}
		want1 bool
	}{
		{name: "int2string", args: args{k: "int", nk: "string"}, want: 1, want1: true},
		{name: "null2string", args: args{k: "null", nk: "string"}, want: "a", want1: true},
		{name: "int2null", args: args{k: "int", nk: "null"}, want: 1, want1: true},
		{name: "null2null", args: args{k: "null", nk: "null"}, want: nil, want1: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := mockCache()
			c.Rename(tt.args.k, tt.args.nk)
			got, got1 := c.Get(tt.args.nk)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RenameAndGet() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("RenameAndGet() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMemCache_Set(t *testing.T) {
	type args struct {
		k    string
		v    interface{}
		opts []SetIOption
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "int", args: args{k: "int", v: 1}, want: true},
		{name: "int32", args: args{k: "int32", v: int32(2)}, want: true},
		{name: "int64", args: args{k: "int64", v: int64(3)}, want: true},
	}
	c := NewMemCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Set(tt.args.k, tt.args.v, tt.args.opts...); got != tt.want {
				t.Errorf("Set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemCache_Ttl(t *testing.T) {
	type args struct {
		k string
	}
	tests := []struct {
		name  string
		args  args
		want  time.Duration
		want1 bool
	}{
		{name: "int", args: args{k: "int"}, want: 0, want1: false},
		{name: "ex", args: args{k: "ex"}, want: 1 * time.Second, want1: true},
		{name: "null", args: args{k: "null"}, want: 0, want1: false},
	}
	c := mockCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := c.Ttl(tt.args.k)
			if tt.want != got && tt.want-got > 10*time.Millisecond {
				t.Errorf("Ttl() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Ttl() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMemCache_DelExpired(t *testing.T) {
	type args struct {
		k     string
		sleep time.Duration
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "int", args: args{k: "int"}, want: 0},
		{name: "ex", args: args{k: "ex"}, want: 0},
		{name: "ex1", args: args{k: "ex", sleep: time.Second}, want: 1},
		{name: "null", args: args{k: "null"}, want: 0},
	}
	c := mockCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time.Sleep(tt.args.sleep)
			if got := c.DelExpired(tt.args.k); got != tt.want {
				t.Errorf("DelExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemCache_AfterExpiration(t *testing.T) {
	c := mockCache()
	type args struct {
		middlewares []Middleware
	}
	tests := []struct {
		name     string
		args     args
		wantFunc func() bool
		want     bool
	}{
		{args: args{
			middlewares: []Middleware{func(key string, value interface{}) { c.Set(key, value) }}},
			wantFunc: func() bool {
				if c.Exists("ex") != true {
					return false
				}
				time.Sleep(1 * time.Second)
				if c.Exists("ex") != false {
					return false
				}
				return c.Exists("ex")
			},
			want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.AfterExpiration(tt.args.middlewares...)
			if got := tt.wantFunc(); got != tt.want {
				t.Errorf("AfterExpiration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemCache_Finalize(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "1"},
	}
	mc := NewMemCache()
	mc.Set("a", 1)
	mc.Set("b", 1, WithEx(1*time.Nanosecond))
	closed := mc.IsClosed()
	mc = nil
	runtime.GC()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(<-closed)
		})
	}
}
