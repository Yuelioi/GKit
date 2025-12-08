package kv_test

import (
	"os"
	"testing"
	"time"

	"github.com/Yuelioi/gkit/utils/kv"
)

func TestKVStore_SetGet(t *testing.T) {
	path := "testdata/db1.json"
	os.Remove(path)

	store, err := kv.NewKVStore[string](path)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	store.Set("foo", "bar")
	val, ok := store.Get("foo")

	if !ok || val != "bar" {
		t.Fatalf("expected bar, got %v / %v", val, ok)
	}
}

func TestKVStore_TTL(t *testing.T) {
	path := "testdata/db2.json"
	os.Remove(path)

	store, err := kv.NewKVStore[string](path, kv.WithSaveInterval(20*time.Millisecond))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	// 设置 50ms 过期
	store.SetWithTTL("temp", "123", 50*time.Millisecond)

	// 初始应存在
	if v, ok := store.Get("temp"); !ok || v != "123" {
		t.Fatalf("expected 123 immediately, got %v %v", v, ok)
	}

	// 等待过期
	time.Sleep(80 * time.Millisecond)

	// 惰性获取应消失
	if _, ok := store.Get("temp"); ok {
		t.Fatal("expected temp expired, but got ok=true")
	}

	// 再等后台清理(确保 cleanup 执行)
	time.Sleep(50 * time.Millisecond)

	keys := store.Keys()
	if len(keys) != 0 {
		t.Fatalf("expected no keys after ttl cleanup, got %v", keys)
	}
}

func TestKVStore_Delete(t *testing.T) {
	path := "testdata/db3.json"
	os.Remove(path)

	store, err := kv.NewKVStore[int](path)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	store.Set("k1", 100)
	store.Delete("k1")

	if _, ok := store.Get("k1"); ok {
		t.Fatal("expected key deleted")
	}
}

func TestKVStore_PersistAndReload(t *testing.T) {
	path := "testdata/db4.json"
	os.Remove(path)

	// 创建 store 并写入数据
	store, err := kv.NewKVStore[int](path)
	if err != nil {
		t.Fatal(err)
	}

	store.Set("a", 1)
	store.Set("b", 2)
	store.Close() // 强制写盘

	// 重新加载
	store2, err := kv.NewKVStore[int](path)
	if err != nil {
		t.Fatal(err)
	}
	defer store2.Close()

	if v, ok := store2.Get("a"); !ok || v != 1 {
		t.Fatalf("expected 1, got %v %v", v, ok)
	}
	if v, ok := store2.Get("b"); !ok || v != 2 {
		t.Fatalf("expected 2, got %v %v", v, ok)
	}
}

func TestKVStore_Keys(t *testing.T) {
	path := "testdata/db5.json"
	os.Remove(path)

	store, err := kv.NewKVStore[string](path)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	store.Set("k1", "a")
	store.SetWithTTL("k2", "b", 50*time.Millisecond) // 50ms 过期

	if len(store.Keys()) != 2 {
		t.Fatalf("expected 2 keys, got %v", store.Keys())
	}

	time.Sleep(80 * time.Millisecond)

	keys := store.Keys()
	if len(keys) != 1 || keys[0] != "k1" {
		t.Fatalf("expected only k1 after ttl, got %v", keys)
	}
}

func TestKVStore_CloseFlush(t *testing.T) {
	path := "testdata/db6.json"
	os.Remove(path)

	store, err := kv.NewKVStore[int](path)
	if err != nil {
		t.Fatal(err)
	}

	store.Set("x", 999)

	// Close 会强制 save()
	store.Close()

	// 重新打开
	store2, err := kv.NewKVStore[int](path)
	if err != nil {
		t.Fatal(err)
	}
	defer store2.Close()

	val, ok := store2.Get("x")
	if !ok || val != 999 {
		t.Fatalf("expected 999 after Close flush, got %v %v", val, ok)
	}
}
