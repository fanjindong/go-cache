package cache

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"testing"
	"time"
)

var c ICache

func TestMain(m *testing.M) {
	mockCache()
	os.Exit(m.Run())
}

func mockCache() {
	c = NewMemCache()
	cw = c.GetCleanupWorker()
	c.Set("int", 1)
	c.Set("int32", int32(1))
	c.Set("int64", int64(1))
	c.Set("string", "a")
	c.Set("float64", 1.1)
	c.Set("float32", float32(1.1))
	c.Set("ex", 1, WithEx(1*time.Second))
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
	mockCache()
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

func TestMemCache_Decr(t *testing.T) {
	type args struct {
		k string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{name: "int", args: args{k: "int"}, want: 0},
		{name: "int32", args: args{k: "int32"}, want: 0},
		{name: "int64", args: args{k: "int64"}, want: 0},
		{name: "string", args: args{k: "string"}, want: 0, wantErr: true},
		{name: "null", args: args{k: "null"}, want: -1},
	}
	mockCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Decr(tt.args.k)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Decr() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemCache_DecrBy(t *testing.T) {
	type args struct {
		k string
		v int64
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{name: "int", args: args{k: "int", v: 1}, want: 0},
		{name: "int32", args: args{k: "int32", v: 2}, want: -1},
		{name: "int64", args: args{k: "int64", v: 1}, want: 0},
		{name: "string", args: args{k: "string", v: 1}, want: 0, wantErr: true},
		{name: "null", args: args{k: "null", v: 2}, want: -2},
	}
	mockCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.DecrBy(tt.args.k, tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecrBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DecrBy() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemCache_DecrByFloat(t *testing.T) {
	type args struct {
		k string
		v float64
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{name: "float32", args: args{k: "float32", v: 0.5}, want: 0.6},
		{name: "float64", args: args{k: "float64", v: 1.1}, want: 0},
		{name: "string", args: args{k: "string", v: 3.0}, want: 0, wantErr: true},
		{name: "null", args: args{k: "null", v: 2.1}, want: -2.1},
	}
	mockCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.DecrByFloat(tt.args.k, tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecrByFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if fmt.Sprintf("%2.f", got) != fmt.Sprintf("%2.f", tt.want) {
				t.Errorf("DecrByFloat() got = %v, want %v", got, tt.want)
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
	mockCache()
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
	mockCache()
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
	mockCache()
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
	mockCache()
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
		{name: "null", args: args{k: "null"}, want: nil, want1: false},
	}
	mockCache()
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
	mockCache()
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
	}
	mockCache()
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

func TestMemCache_Incr(t *testing.T) {
	type args struct {
		k string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{name: "int", args: args{k: "int"}, want: 2, wantErr: false},
		{name: "int32", args: args{k: "int32"}, want: 2, wantErr: false},
		{name: "int64", args: args{k: "int64"}, want: 2, wantErr: false},
		{name: "float32", args: args{k: "float32"}, want: 0, wantErr: true},
		{name: "float64", args: args{k: "float64"}, want: 0, wantErr: true},
		{name: "string", args: args{k: "string"}, want: 0, wantErr: true},
		{name: "null", args: args{k: "null"}, want: 1, wantErr: false},
	}
	mockCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Incr(tt.args.k)
			if (err != nil) != tt.wantErr {
				t.Errorf("Incr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Incr() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemCache_IncrAndGet(t *testing.T) {
	type args struct {
		k string
	}
	tests := []struct {
		name  string
		args  args
		want  interface{}
		want1 bool
	}{
		{name: "int", args: args{k: "int"}, want: 2, want1: true},
		{name: "int32", args: args{k: "int32"}, want: int32(2), want1: true},
		{name: "int64", args: args{k: "int64"}, want: int64(2), want1: true},
		{name: "null", args: args{k: "null"}, want: 1, want1: true},
	}
	mockCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.Incr(tt.args.k)
			got, got1 := c.Get(tt.args.k)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IncrAndGet() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("IncrAndGet() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMemCache_IncrBy(t *testing.T) {
	type args struct {
		k string
		v int64
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{name: "int", args: args{k: "int", v: 1}, want: 2, wantErr: false},
		{name: "int32", args: args{k: "int32", v: 2}, want: 3, wantErr: false},
		{name: "int64", args: args{k: "int64", v: 1}, want: 2, wantErr: false},
		{name: "float32", args: args{k: "float32", v: 1}, want: 0, wantErr: true},
		{name: "string", args: args{k: "string", v: 1}, want: 0, wantErr: true},
		{name: "null", args: args{k: "null", v: 5}, want: 5, wantErr: false},
	}
	mockCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.IncrBy(tt.args.k, tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("IncrBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IncrBy() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemCache_IncrByFloat(t *testing.T) {
	type args struct {
		k string
		v float64
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{name: "int", args: args{k: "int", v: 1.1}, want: 0, wantErr: true},
		{name: "int32", args: args{k: "int32", v: 2}, want: 0, wantErr: true},
		{name: "int64", args: args{k: "int64", v: 1}, want: 0, wantErr: true},
		{name: "string", args: args{k: "string", v: 2.1}, want: 0, wantErr: true},
		{name: "float32", args: args{k: "float32", v: 1.1}, want: 2.2, wantErr: false},
		{name: "float64", args: args{k: "float64", v: 1.2}, want: 2.3, wantErr: false},
		{name: "null", args: args{k: "null", v: -1.2}, want: -1.2, wantErr: false},
	}
	mockCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.IncrByFloat(tt.args.k, tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("IncrByFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if fmt.Sprintf("%2.f", got) != fmt.Sprintf("%2.f", tt.want) {
				t.Errorf("IncrByFloat() got = %v, want %v", got, tt.want)
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
	mockCache()
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
	mockCache()
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
	mockCache()
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
	mockCache()
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
	mockCache()
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
			mockCache()
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Set(tt.args.k, tt.args.v, tt.args.opts...); got != tt.want {
				t.Errorf("Set() = %v, want %v", got, tt.want)
			}
		})
	}
	mockCache()
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
	mockCache()
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

//func BenchmarkMemCache_Set(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		c.Set(fmt.Sprintf("%d", i), "a", WithEx(1*time.Millisecond))
//	}
//}
