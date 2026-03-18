package ungo

/*
#include <stdlib.h>
#include <stdint.h>
#include <string.h>

typedef struct {
    int64_t key;
    int64_t value;
    uint8_t occupied;
} Entry;

typedef struct {
    Entry* entries;
    size_t capacity;
    size_t size;
} CMap;

static uint64_t hash(int64_t k) {
    uint64_t x = (uint64_t)k;
    x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9ULL;
    x = (x ^ (x >> 27)) * 0x94d049bb133111ebULL;
    x = x ^ (x >> 31);
    return x;
}

CMap* create_map(size_t cap) {
    CMap* m = malloc(sizeof(CMap));
    m->capacity = cap;
    m->size = 0;
    m->entries = calloc(cap, sizeof(Entry));
    return m;
}

void free_map(CMap* m) {
    free(m->entries);
    free(m);
}

void set(CMap* m, int64_t key, int64_t value) {
    uint64_t h = hash(key) % m->capacity;
    while (m->entries[h].occupied) {
        if (m->entries[h].key == key) {
            m->entries[h].value = value;
            return;
        }
        h = (h + 1) % m->capacity;
    }
    m->entries[h].key = key;
    m->entries[h].value = value;
    m->entries[h].occupied = 1;
    m->size++;
}

int get(CMap* m, int64_t key, int64_t* out_val) {
    uint64_t h = hash(key) % m->capacity;
    uint64_t start = h;
    while (m->entries[h].occupied) {
        if (m->entries[h].key == key) {
            *out_val = m->entries[h].value;
            return 1;
        }
        h = (h + 1) % m->capacity;
        if (h == start) break;
    }
    return 0;
}

void delete_key(CMap* m, int64_t key) {
    uint64_t h = hash(key) % m->capacity;
    while (m->entries[h].occupied) {
        if (m->entries[h].key == key) {
            m->entries[h].occupied = 0;
            m->size--;

            // Re-hash the cluster to ensure linear probing integrity
            uint64_t i = h;
            while (1) {
                h = (h + 1) % m->capacity;
                if (!m->entries[h].occupied) break;

                int64_t temp_k = m->entries[h].key;
                int64_t temp_v = m->entries[h].value;
                m->entries[h].occupied = 0;
                m->size--;
                set(m, temp_k, temp_v);
            }
            return;
        }
        h = (h + 1) % m->capacity;
    }
}
*/
import "C"
import (
	"unsafe"
)

type FastMap struct {
	ptr *C.CMap
}

func NewFastMap(capacity int) *FastMap {
	// Ensure capacity is non-zero to avoid modulo zero in C
	if capacity <= 0 {
		capacity = 16
	}
	return &FastMap{
		ptr: C.create_map(C.size_t(capacity)),
	}
}

func (fm *FastMap) Set(key int64, value int64) {
	C.set(fm.ptr, C.int64_t(key), C.int64_t(value))
}

func (fm *FastMap) Get(key int64) (int64, bool) {
	var val C.int64_t
	if C.get(fm.ptr, C.int64_t(key), &val) == 1 {
		return int64(val), true
	}
	return 0, false
}

func (fm *FastMap) Delete(key int64) {
	C.delete_key(fm.ptr, C.int64_t(key))
}

func (fm *FastMap) Size() int {
	return int(fm.ptr.size)
}

func (fm *FastMap) Contains(key int64) bool {
	var val C.int64_t
	return C.get(fm.ptr, C.int64_t(key), &val) == 1
}

func (fm *FastMap) ForEach(f func(key int64, value int64)) {
	entries := unsafe.Slice(fm.ptr.entries, fm.ptr.capacity)
	for i := 0; i < int(fm.ptr.capacity); i++ {
		if entries[i].occupied == 1 {
			f(int64(entries[i].key), int64(entries[i].value))
		}
	}
}

func (fm *FastMap) Keys() []int64 {
	keys := make([]int64, 0, fm.ptr.size)
	entries := unsafe.Slice(fm.ptr.entries, fm.ptr.capacity)
	for i := 0; i < int(fm.ptr.capacity); i++ {
		if entries[i].occupied == 1 {
			keys = append(keys, int64(entries[i].key))
		}
	}
	return keys
}

func (fm *FastMap) Values() []int64 {
	values := make([]int64, 0, fm.ptr.size)
	entries := unsafe.Slice(fm.ptr.entries, fm.ptr.capacity)
	for i := 0; i < int(fm.ptr.capacity); i++ {
		if entries[i].occupied == 1 {
			values = append(values, int64(entries[i].value))
		}
	}
	return values
}

func (fm *FastMap) Clear() {
	cap := fm.ptr.capacity
	C.free_map(fm.ptr)
	fm.ptr = C.create_map(cap)
}

func (fm *FastMap) Close() {
	if fm.ptr != nil {
		C.free_map(fm.ptr)
		fm.ptr = nil
	}
}
