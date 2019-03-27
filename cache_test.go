package gocache

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	var cache *Cache
	cache = NewCache(2, 5*time.Second)
	values := []string{"test1", "test2", "test3"}
	key := "key1"
	for _, v := range values {
		cache.Set(v, v, 3*time.Second)
		val, ok := cache.Get(v)
		if !ok {
			t.Fatalf("expect key:%v ,value:%v", key, v)
		} else if ok && val != v {
			t.Fatalf("expect key:%v , expect value:%v, get value:%v", key, v, val)
		}
		t.Logf("value:%v ", val)
	}
}

func TestGet(t *testing.T) {
	var cache *Cache
	cache = NewCache(2, 5*time.Second)
	key := "key1"
	val := "test1"
	cache.Set(key, val, 3*time.Second)
	fmt.Println(111)
	time.Sleep(5 * time.Second)
	fmt.Println(cache.Get(key))
	fmt.Println(222)
}

func BenchmarkGet(b *testing.B) {
	b.ResetTimer()
	var cache *Cache
	cache = NewCache(b.N, 5*time.Second)
	for i := 0; i < b.N; i++ {
		cache.Set(strconv.Itoa(i), i, 3*time.Second)
		cache.Get(strconv.Itoa(i))
	}
}

func BenchmarkGet2(b *testing.B) {
	b.ResetTimer()
	var cache *Cache
	cache = NewCache(b.N, 5*time.Second)
	for i := 0; i < b.N; i++ {
		cache.Set(strconv.Itoa(i), i, 3*time.Second)
		cache.Get(strconv.Itoa(i))
	}
}

func TestDelete(t *testing.T) {
	var cache *Cache
	cache = NewCache(2, 5*time.Second)
	cache.Set("myKey", 1234, 2*time.Second)
	if val, ok := cache.Get("myKey"); !ok {
		t.Fatal("TestRemove returned no match")
	} else if val != 1234 {
		t.Fatalf("TestRemove failed.  Expected %d, got %v", 1234, val)
	} else {
		fmt.Println(val)
	}

	cache.Delete("myKey")
	if _, ok := cache.Get("myKey"); ok {
		t.Fatal("TestRemove returned a removed item")
	}
}

func TestRun(t *testing.T) {
	var cache *Cache
	cache = NewCache(2, 5*time.Second)
	cache.Set("myKey", 1234, 3*time.Second)
	time.Sleep(5 * time.Second)
	fmt.Println(cache.Get("myKey"))
}

