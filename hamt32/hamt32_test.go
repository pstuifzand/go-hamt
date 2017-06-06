package hamt32_test

import (
	"log"
	"testing"
	"time"

	"github.com/lleo/go-hamt/hamt32"
)

func TestBuild(t *testing.T) {
	var name = "TestBuild"
	if Functional {
		name += ":functional:" + hamt32.TableOptionName[TableOption]
	} else {
		name += ":transient:" + hamt32.TableOptionName[TableOption]
	}
	log.Print(name)

	var h = hamt32.New(Functional, TableOption)

	for _, kv := range KVS[:30] {
		var k = kv.Key
		var v = kv.Val

		var inserted bool
		h, inserted = h.Put(k, v)
		if !inserted {
			log.Printf("%s: failed to insert k=%s, v=%d", name, k, v)
			t.Fatalf("%s: failed to insert k=%s, v=%d", name, k, v)
		}

		//log.Print(h.LongString(""))
	}
}

func TestHamt32Put(t *testing.T) {
	var name = "TestHamt32Put"
	if Functional {
		name += ":functional:" + hamt32.TableOptionName[TableOption]
	} else {
		name += ":transient:" + hamt32.TableOptionName[TableOption]
	}
	log.Print(name)

	StartTime[name] = time.Now()
	Hamt32 = hamt32.New(Functional, TableOption)
	for _, kv := range KVS {
		var k = kv.Key
		var v = kv.Val

		var inserted bool
		Hamt32, inserted = Hamt32.Put(k, v)
		if !inserted {
			log.Printf("%s: failed to Hamt32.Put(%s, %v)", name, k, v)
			t.Fatalf("%s: failed to Hamt32.Put(%s, %v)", name, k, v)
		}

		var val, found = Hamt32.Get(k)
		if !found {
			log.Printf("%s: failed to Hamt32.Get(%s)", name, k)
			//log.Print(Hamt32.LongString(""))
			t.Fatalf("%s: failed to Hamt32.Get(%s)", name, k)
		}
		if val != v {
			log.Printf("%s: returned val,%d != expected v,%d for k=%s", name, val, v, k)
			t.Fatalf("%s: returned val,%d != expected v,%d for k=%s", name, val, v, k)
		}
	}
	RunTime[name] = time.Since(StartTime[name])
}

func TestHamt32Get(t *testing.T) {
	var name = "TestHamt32Get"
	if Functional {
		name += ":functional:" + hamt32.TableOptionName[TableOption]
	} else {
		name += ":transient:" + hamt32.TableOptionName[TableOption]
	}
	log.Print(name)

	if Hamt32 == nil {
		var err error
		Hamt32, err = buildHamt32(name, KVS, Functional, TableOption)
		if err != nil {
			log.Printf("%s: failed buildHamt32(%q, KVS#%d, %t, %s) => %s", name,
				name, len(KVS), Functional,
				hamt32.TableOptionName[TableOption], err)
			t.Fatalf("%s: failed buildHamt32(%q, KVS#%d, %t, %s) => %s", name,
				name, len(KVS), Functional,
				hamt32.TableOptionName[TableOption], err)
		}
	}

	StartTime[name] = time.Now()
	for _, kv := range KVS {
		var k = kv.Key
		var v = kv.Val

		var val, found = Hamt32.Get(k)
		if !found {
			log.Printf("%s: Failed to Hamt32.Get(%s)", name, k)
			log.Print(Hamt32.LongString(""))
			t.Fatalf("%s: Failed to Hamt32.Get(%s)", name, k)
		}
		if val != v {
			log.Printf("%s: retrieved val,%d != expected v,%d for k=%s", name, val, v, k)
			t.Fatalf("%s: retrieved val,%d != expected v,%d for k=%s", name, val, v, k)
		}
	}
	RunTime[name] = time.Since(StartTime[name])
}

func TestHamt32Del(t *testing.T) {
	var name = "TestHamt32Del"
	if Functional {
		name += ":functional:" + hamt32.TableOptionName[TableOption]
	} else {
		name += ":transient:" + hamt32.TableOptionName[TableOption]
	}
	log.Print(name)

	if Hamt32 == nil {
		var err error
		Hamt32, err = buildHamt32(name, KVS, Functional, TableOption)
		if err != nil {
			log.Printf("%s: failed buildHamt32(%q, KVS#%d, %t, %s) => %s", name,
				name, len(KVS), Functional,
				hamt32.TableOptionName[TableOption], err)
			t.Fatalf("%s: failed buildHamt32(%q, KVS#%d, %t, %s) => %s", name,
				name, len(KVS), Functional,
				hamt32.TableOptionName[TableOption], err)
		}
	}

	StartTime[name] = time.Now()
	for _, kv := range KVS {
		var k = kv.Key
		var v = kv.Val

		var val interface{}
		var deleted bool
		Hamt32, val, deleted = Hamt32.Del(k)
		if !deleted {
			log.Printf("%s: Failed to Hamt32.Del(%s)", name, k)
			log.Print(Hamt32.LongString(""))
			t.Fatalf("%s: Failed to Hamt32.Del(%s)", name, k)
		}
		if val != v {
			log.Printf("%s: retrieved val,%d != expected v,%d for k=%s", name, val, v, k)
			t.Fatalf("%s: retrieved val,%d != expected v,%d for k=%s", name, val, v, k)
		}
	}
	RunTime[name] = time.Since(StartTime[name])
}

func BenchmarkHamt32Get(b *testing.B) {
	var name = "BenchmarkHamt32Get"
	if Functional {
		name += ":functional:" + hamt32.TableOptionName[TableOption]
	} else {
		name += ":transient:" + hamt32.TableOptionName[TableOption]
	}
	log.Print(name)

	var h hamt32.Hamt
	var err error
	h, err = buildHamt32(name, KVS, Functional, TableOption)
	if err != nil {
		log.Printf("%s: failed buildHamt32(%q, KVS#%d, %t, %s) => %s", name,
			name, len(KVS), false, hamt32.TableOptionName[TableOption], err)
		b.Fatalf("%s: failed buildHamt32(%q, KVS#%d, %t, %s) => %s", name,
			name, len(KVS), false, hamt32.TableOptionName[TableOption], err)
	}

	log.Printf("%s: Functional=%t; b.N=%d", name, Functional, b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = i % len(KVS)
		var k = KVS[j].Key
		var v = KVS[j].Val

		var val, found = h.Get(k)
		if !found {
			log.Printf("%s: Failed to h.Get(%s)", name, k)
			//log.Print(h.LongString(""))
			b.Fatalf("%s: Failed to h.Get(%s)", name, k)
		}
		if val != v {
			log.Printf("%s: retrieved val,%d != expected v,%d for k=%s", name, val, v, k)
			b.Fatalf("%s: retrieved val,%d != expected v,%d for k=%s", name, val, v, k)
		}
	}
}

var BenchHamt32_T2F hamt32.Hamt

func BenchmarkHamt32_T2F_Get(b *testing.B) {
	var name = "BenchmarkHamt32_T2F_Get"
	name += ":functional:" + hamt32.TableOptionName[TableOption]
	log.Print(name)

	if BenchHamt32_T2F == nil {
		var err error
		BenchHamt32_T2F, err = buildHamt32(name, KVS, false, TableOption)
		if err != nil {
			log.Printf("%s: failed buildHamt32(%q, KVS#%d, %t, %s) => %s", name,
				name, len(KVS), false, hamt32.TableOptionName[TableOption], err)
			b.Fatalf("%s: failed buildHamt32(%q, KVS#%d, %t, %s) => %s", name,
				name, len(KVS), false, hamt32.TableOptionName[TableOption], err)
		}
		BenchHamt32_T2F = BenchHamt32_T2F.ToFunctional()
	}

	log.Printf("%s: Transient-to-Functional; b.N=%d;", name, b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = i % len(KVS)
		var k = KVS[j].Key
		var v = KVS[j].Val

		var val, found = BenchHamt32_T2F.Get(k)
		if !found {
			log.Printf("%s: Failed to BenchHamt32_T2F.Get(%s)", name, k)
			//log.Print(h.LongString(""))
			b.Fatalf("%s: Failed to BenchHamt32_T2F.Get(%s)", name, k)
		}
		if val != v {
			log.Printf("%s: retrieved val,%d != expected v,%d for k=%s", name, val, v, k)
			b.Fatalf("%s: retrieved val,%d != expected v,%d for k=%s", name, val, v, k)
		}
	}
}

var BenchHamt32_F2T hamt32.Hamt

func BenchmarkHamt32_F2T_Get(b *testing.B) {
	var name = "BenchmarkHamt32_F2T_Get"
	name += ":transient:" + hamt32.TableOptionName[TableOption]
	log.Print(name)

	if BenchHamt32_F2T == nil {
		var err error
		BenchHamt32_F2T, err = buildHamt32(name, KVS, true, TableOption)
		if err != nil {
			log.Printf("%s: failed buildHamt32(%q, KVS#%d, %t, %s) => %s", name,
				name, len(KVS), false, hamt32.TableOptionName[TableOption], err)
			b.Fatalf("%s: failed buildHamt32(%q, KVS#%d, %t, %s) => %s", name,
				name, len(KVS), false, hamt32.TableOptionName[TableOption], err)
		}
		BenchHamt32_F2T = BenchHamt32_F2T.ToTransient()
	}

	log.Printf("%s: Functional-to-Transient; b.N=%d;", name, b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = i % len(KVS)
		var k = KVS[j].Key
		var v = KVS[j].Val

		var val, found = BenchHamt32_F2T.Get(k)
		if !found {
			log.Printf("%s: Failed to BenchHamt32_F2T.Get(%s)", name, k)
			//log.Print(h.LongString(""))
			b.Fatalf("%s: Failed to BenchHamt32_F2T.Get(%s)", name, k)
		}
		if val != v {
			log.Printf("%s: retrieved val,%d != expected v,%d for k=%s", name, val, v, k)
			b.Fatalf("%s: retrieved val,%d != expected v,%d for k=%s", name, val, v, k)
		}
	}
}

func BenchmarkHamt32Put(b *testing.B) {
	var name = "BenchmarkHamt32Put"
	if Functional {
		name += ":functional:" + hamt32.TableOptionName[TableOption]
	} else {
		name += ":transient:" + hamt32.TableOptionName[TableOption]
	}
	log.Print(name)

	if b.N+InitHamtNumKvsForPut > len(KVS) {
		log.Printf("%s: Can't run: b.N+num > len(KVS)", name)
		b.Fatalf("%s: Can't run: b.N+num > len(KVS)", name)
	}

	var kvs = KVS[:InitHamtNumKvsForPut]

	var h, err = buildHamt32(name, kvs, Functional, TableOption)
	if err != nil {
		log.Printf("%s: failed buildHamt32(%q, KVS[:%d], %t, %s) => %s", name,
			name, InitHamtNumKvsForPut, Functional,
			hamt32.TableOptionName[TableOption], err)
		b.Fatalf("%s: failed buildHamt32(%q, KVS[:%d], %t, %s) => %s", name,
			name, InitHamtNumKvsForPut, Functional,
			hamt32.TableOptionName[TableOption], err)
	}

	log.Printf("%s: Functional=%t; b.N=%d;", name, Functional, b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var k = KVS[InitHamtNumKvsForPut+i].Key
		var v = KVS[InitHamtNumKvsForPut+i].Val

		var added bool
		h, added = h.Put(k, v)
		if !added {
			log.Printf("%s: failed to h.Put(%s, %d)", name, k, v)
			//log.Print(h.LongString(""))
			b.Fatalf("%s: failed to h.Put(%s, %d)", name, k, v)
		}
	}
}

func BenchmarkHamt32_T2F_Put(b *testing.B) {
	var name = "BenchmarkHamt32Put_T2F"
	name += ":functional:" + hamt32.TableOptionName[TableOption]
	log.Print(name)

	var InitHamtNumKvsForPut int //= 1000000 // 1 million; allows b.N=3,000,000
	if b.N+InitHamtNumKvsForPut > len(KVS) {
		log.Printf("%s: Can't run: b.N+num > len(KVS)", name)
		b.Fatalf("%s: Can't run: b.N+num > len(KVS)", name)
	}

	var kvs = KVS[:InitHamtNumKvsForPut]

	var h, err = buildHamt32(name, kvs, false, TableOption)
	if err != nil {
		log.Printf("%s: failed buildHamt32(%q, KVS[:%d], %t, %s) => %s", name,
			name, InitHamtNumKvsForPut, Functional,
			hamt32.TableOptionName[TableOption], err)
		b.Fatalf("%s: failed buildHamt32(%q, KVS[:%d], %t, %s) => %s", name,
			name, InitHamtNumKvsForPut, Functional,
			hamt32.TableOptionName[TableOption], err)
	}
	h = h.ToFunctional()

	log.Printf("%s: Transient-to-Functional; b.N=%d;", name, b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var k = KVS[InitHamtNumKvsForPut+i].Key
		var v = KVS[InitHamtNumKvsForPut+i].Val

		var added bool
		h, added = h.Put(k, v)
		if !added {
			log.Printf("%s: failed to h.Put(%s, %d)", name, k, v)
			//log.Print(h.LongString(""))
			b.Fatalf("%s: failed to h.Put(%s, %d)", name, k, v)
		}
	}
}

func BenchmarkHamt32Del(b *testing.B) {
	var name = "BenchmarkHamt32Del"
	if Functional {
		name += ":functional:" + hamt32.TableOptionName[TableOption]
	} else {
		name += ":transient:" + hamt32.TableOptionName[TableOption]
	}
	log.Print(name)

	var h, err = buildHamt32(name, KVS, Functional, TableOption)
	if err != nil {
		log.Printf("%s: failed buildHamt32(%q, KVS:%d, %t, %s) => %s", name,
			name, len(KVS), Functional,
			hamt32.TableOptionName[TableOption], err)
		b.Fatalf("%s: failed buildHamt32(%q, KVS:%d, %t, %s) => %s", name,
			name, len(KVS), Functional,
			hamt32.TableOptionName[TableOption], err)
	}

	log.Printf("%s: Functional=%t; b.N=%d;", name, Functional, b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var k = KVS[i].Key
		var v = KVS[i].Val

		var deleted bool
		var val interface{}
		h, val, deleted = h.Del(k)
		if !deleted {
			log.Printf("%s: failed to h.Del(%s)", name, k)
			//log.Print(h.LongString(""))
			b.Fatalf("%s: failed to h.Del(%s)", name, k)
		}
		if val != v {
			log.Printf("%s: failed val,%d != v,%d", name, val, v)
			b.Fatalf("%s: failed val,%d != v,%d", name, val, v)
		}
	}
}