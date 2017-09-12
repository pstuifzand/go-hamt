package hamt64_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"log"
	"testing"
	"time"

	"github.com/lleo/go-hamt/hamt64"
	"github.com/pkg/errors"
)

func TestBuild64(t *testing.T) {
	var name = "TestBuild64"
	if Functional {
		name += ":functional:" + hamt64.TableOptionName[TableOption]
	} else {
		name += ":transient:" + hamt64.TableOptionName[TableOption]
	}

	var h = hamt64.New(Functional, TableOption)

	for _, kv := range KVS[:30] {
		var k = kv.Key
		var v = kv.Val

		var inserted bool
		h, inserted = h.Put(k, v)
		if !inserted {
			log.Printf("%s: failed to insert s=%q, v=%d", name, string(k), v)
			t.Fatalf("%s: failed to insert s=%q, v=%d", name, string(k), v)
		}

		//log.Print(h.LongString(""))
	}
}

func TestHamt64Put(t *testing.T) {
	runTestHamt64Put(t, KVS, Functional, TableOption)
}

func runTestHamt64Put(
	t *testing.T,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "TestHamt64Put"
	if functional {
		name += ":functional:" + hamt64.TableOptionName[tblOpt]
	} else {
		name += ":transient:" + hamt64.TableOptionName[tblOpt]
	}

	StartTime[name] = time.Now()
	Hamt64 = hamt64.New(functional, tblOpt)
	for _, kv := range kvs {
		var k = kv.Key
		var v = kv.Val

		var inserted bool
		Hamt64, inserted = Hamt64.Put(k, v)
		if !inserted {
			log.Printf("%s: failed to Hamt64.Put(%q, %v)", name, string(k), v)
			t.Fatalf("%s: failed to Hamt64.Put(%q, %v)", name, string(k), v)
		}

		var val, found = Hamt64.Get(k)
		if !found {
			log.Printf("%s: failed to Hamt64.Get(%q)", name, string(k))
			//log.Print(Hamt64.LongString(""))
			t.Fatalf("%s: failed to Hamt64.Get(%q)", name, string(k))
		}
		if val != v {
			log.Printf("%s: returned val,%d != expected v,%d for s=%q", name, val, v, string(k))
			t.Fatalf("%s: returned val,%d != expected v,%d for s=%q", name, val, v, string(k))
		}
	}
	RunTime[name] = time.Since(StartTime[name])

	StartTime["Hamt64.Stats()"] = time.Now()
	var stats = Hamt64.Stats()
	RunTime["Hamt64.Stats()"] = time.Since(StartTime["Hamt64.Stats()"])
	log.Printf("%s: stats=%+v;\n", name, stats)
}

func TestHamt64IterFunc(t *testing.T) {
	runTestHamt64IterFunc(t, KVS, Functional, TableOption)
}

func runTestHamt64IterFunc(
	t *testing.T,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "TestHamt64IterFunc"
	if functional {
		name += ":functional:" + hamt64.TableOptionName[tblOpt]
	} else {
		name += ":transient:" + hamt64.TableOptionName[tblOpt]
	}

	if Hamt64 == nil {
		var err error
		Hamt64, err = buildHamt64(name, kvs, functional, tblOpt)
		if err != nil {
			log.Printf("%s: failed buildHamt64(%q, kvs#%d, %t, %s) => %s", name,
				name, len(kvs), functional,
				hamt64.TableOptionName[tblOpt], err)
			t.Fatalf("%s: failed buildHamt64(%q, kvs#%d, %t, %s) => %s", name,
				name, len(kvs), functional,
				hamt64.TableOptionName[tblOpt], err)
		}

		StartTime["Hamt64.Stats()"] = time.Now()
		var stats = Hamt64.Stats()
		RunTime["Hamt64.Stats()"] = time.Since(StartTime["Hamt64.Stats()"])
		log.Printf("%s: stats=%+v;\n", name, stats)
	}

	StartTime[name] = time.Now()

	var i int
	var next = Hamt64.Iter()
	for kv, ok := next(); ok; kv, ok = next() {
		var val, ok = Hamt64.Get(kv.Key)
		if !ok {
			t.Fatalf("failed to lookup %s in Hamt64", kv.Key)
		}

		if val != kv.Val {
			t.Fatalf("val,%v != kv.Val,%v\n", val, kv.Val)
		}

		i++
	}

	if len(kvs) != i {
		t.Fatalf("Expected len(kvs),%d go i,%d; Hamt64.Nentries()=%d;",
			len(kvs), i, Hamt64.Nentries())
	}

	RunTime[name] = time.Since(StartTime[name])
}

func TestHamt64IterChan(t *testing.T) {
	runTestHamt64IterChan(t, KVS, Functional, TableOption)
}

func runTestHamt64IterChan(
	t *testing.T,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "TestHamt64IterChan"
	if functional {
		name += ":functional:" + hamt64.TableOptionName[tblOpt]
	} else {
		name += ":transient:" + hamt64.TableOptionName[tblOpt]
	}

	if Hamt64 == nil {
		var err error
		Hamt64, err = buildHamt64(name, kvs, functional, tblOpt)
		if err != nil {
			log.Printf("%s: failed buildHamt64(%q, kvs#%d, %t, %s) => %s", name,
				name, len(kvs), functional,
				hamt64.TableOptionName[tblOpt], err)
			t.Fatalf("%s: failed buildHamt64(%q, kvs#%d, %t, %s) => %s", name,
				name, len(kvs), functional,
				hamt64.TableOptionName[tblOpt], err)
		}

		StartTime["Hamt64.Stats()"] = time.Now()
		var stats = Hamt64.Stats()
		RunTime["Hamt64.Stats()"] = time.Since(StartTime["Hamt64.Stats()"])
		log.Printf("%s: stats=%+v;\n", name, stats)
	}

	StartTime[name] = time.Now()

	var i int
	for kv := range Hamt64.IterChan(0, nil) {
		var val, ok = Hamt64.Get(kv.Key)
		if !ok {
			t.Fatalf("failed to lookup %s in Hamt64", kv.Key)
		}

		if val != kv.Val {
			t.Fatalf("val,%v != kv.Val,%v\n", val, kv.Val)
		}

		i++
	}

	if len(kvs) != i {
		t.Fatalf("Expected len(kvs),%d go i,%d", len(kvs), i)
	}

	RunTime[name] = time.Since(StartTime[name])
}

func TestHamt64IterChanContext(t *testing.T) {
	runTestHamt64IterChanContext(t, KVS, Functional, TableOption)
}

func runTestHamt64IterChanContext(
	t *testing.T,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "TestHamt64IterChanContext"
	if functional {
		name += ":functional:" + hamt64.TableOptionName[tblOpt]
	} else {
		name += ":transient:" + hamt64.TableOptionName[tblOpt]
	}

	if Hamt64 == nil {
		var err error
		Hamt64, err = buildHamt64(name, kvs, functional, tblOpt)
		if err != nil {
			log.Printf("%s: failed buildHamt64(%q, kvs#%d, %t, %s) => %s", name,
				name, len(kvs), functional,
				hamt64.TableOptionName[tblOpt], err)
			t.Fatalf("%s: failed buildHamt64(%q, kvs#%d, %t, %s) => %s", name,
				name, len(kvs), functional,
				hamt64.TableOptionName[tblOpt], err)
		}

		StartTime["Hamt64.Stats()"] = time.Now()
		var stats = Hamt64.Stats()
		RunTime["Hamt64.Stats()"] = time.Since(StartTime["Hamt64.Stats()"])
		log.Printf("%s: stats=%+v;\n", name, stats)
	}

	StartTime[name] = time.Now()

	var i int
	var stopKey = kvs[0].Key // "aaa" but key from iter are random
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	var iterChan = Hamt64.IterChan(0, ctx)
	for kv := range iterChan {
		var val, ok = Hamt64.Get(kv.Key)
		if !ok {
			t.Fatalf("failed to lookup %s in Hamt64", kv.Key)
		}

		if val != kv.Val {
			t.Fatalf("val,%v != kv.Val,%v\n", val, kv.Val)
		}

		i++

		if bytes.Equal(kv.Key, stopKey) {
			break
		}
	}

	log.Printf("%s: stopped after %d iterations", name, i)

	RunTime[name] = time.Since(StartTime[name])
}

func TestHamt64Get(t *testing.T) {
	runTestHamt64Get(t, KVS, Functional, TableOption)
}

func runTestHamt64Get(
	t *testing.T,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "TestHamt64Get"
	if functional {
		name += ":functional:" + hamt64.TableOptionName[tblOpt]
	} else {
		name += ":transient:" + hamt64.TableOptionName[tblOpt]
	}

	if Hamt64 == nil {
		var err error
		Hamt64, err = buildHamt64(name, kvs, functional, tblOpt)
		if err != nil {
			log.Printf("%s: failed buildHamt64(%q, kvs#%d, %t, %s) => %s", name,
				name, len(kvs), functional,
				hamt64.TableOptionName[tblOpt], err)
			t.Fatalf("%s: failed buildHamt64(%q, kvs#%d, %t, %s) => %s", name,
				name, len(kvs), functional,
				hamt64.TableOptionName[tblOpt], err)
		}

		StartTime["Hamt64.Stats()"] = time.Now()
		var stats = Hamt64.Stats()
		RunTime["Hamt64.Stats()"] = time.Since(StartTime["Hamt64.Stats()"])
		log.Printf("%s: stats=%+v;\n", name, stats)
	}

	StartTime[name] = time.Now()
	for _, kv := range kvs {
		var k = kv.Key
		var v = kv.Val

		var val, found = Hamt64.Get(k)
		if !found {
			log.Printf("%s: Failed to Hamt64.Get(%q)", name, string(k))
			log.Print(Hamt64.LongString(""))
			t.Fatalf("%s: Failed to Hamt64.Get(%q)", name, string(k))
		}
		if val != v {
			log.Printf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, string(k))
			t.Fatalf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, string(k))
		}
	}
	RunTime[name] = time.Since(StartTime[name])
}

func TestHamt64Del(t *testing.T) {
	runTestHamt64Del(t, KVS, Functional, TableOption)
}

func runTestHamt64Del(
	t *testing.T,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "TestHamt64Del"
	if functional {
		name += ":functional:" + hamt64.TableOptionName[tblOpt]
	} else {
		name += ":transient:" + hamt64.TableOptionName[tblOpt]
	}

	if Hamt64 == nil {
		var err error
		Hamt64, err = buildHamt64(name, kvs, functional, tblOpt)
		if err != nil {
			log.Printf("%s: failed buildHamt64(%q, kvs#%d, %t, %s) => %s", name,
				name, len(kvs), functional,
				hamt64.TableOptionName[tblOpt], err)
			t.Fatalf("%s: failed buildHamt64(%q, kvs#%d, %t, %s) => %s", name,
				name, len(kvs), functional,
				hamt64.TableOptionName[tblOpt], err)
		}

		StartTime["Hamt64.Stats()"] = time.Now()
		var stats = Hamt64.Stats()
		RunTime["Hamt64.Stats()"] = time.Since(StartTime["Hamt64.Stats()"])
		log.Printf("%s: stats=%+v;\n", name, stats)
	}

	StartTime[name] = time.Now()
	for _, kv := range kvs {
		var k = kv.Key
		var v = kv.Val

		var val interface{}
		var deleted bool
		Hamt64, val, deleted = Hamt64.Del(k)
		if !deleted {
			log.Printf("%s: Failed to Hamt64.Del(%q)", name, string(k))
			log.Print(Hamt64.LongString(""))
			t.Fatalf("%s: Failed to Hamt64.Del(%q)", name, string(k))
		}
		if val != v {
			log.Printf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, string(k))
			t.Fatalf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, string(k))
		}
	}
	RunTime[name] = time.Since(StartTime[name])
}

func BenchmarkHamt64Get(b *testing.B) {
	runBenchmarkHamt64Get(b, KVS, Functional, TableOption)
}

func runBenchmarkHamt64Get(
	b *testing.B,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "BenchmarkHamt64Get"
	if functional {
		name += ":functional:" + hamt64.TableOptionName[tblOpt]
	} else {
		name += ":transient:" + hamt64.TableOptionName[tblOpt]
	}

	var BenchHamt64Get, err = buildHamt64(name, kvs, functional, tblOpt)
	if err != nil {
		log.Printf("%s: failed buildHamt64(%q, kvs#%d, %t, %s) => %s", name,
			name, len(kvs), false, hamt64.TableOptionName[tblOpt], err)
		b.Fatalf("%s: failed buildHamt64(%q, kvs#%d, %t, %s) => %s", name,
			name, len(kvs), false, hamt64.TableOptionName[tblOpt], err)
	}

	log.Printf("%s: b.N=%d", name, b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = i % len(kvs)
		var k = kvs[j].Key
		var v = kvs[j].Val

		var val, found = BenchHamt64Get.Get(k)
		if !found {
			log.Printf("%s: Failed to h.Get(%q)", name, string(k))
			//log.Print(h.LongString(""))
			b.Fatalf("%s: Failed to h.Get(%q)", name, string(k))
		}
		if val != v {
			log.Printf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, string(k))
			b.Fatalf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, string(k))
		}
	}
}

func BenchmarkHamt64Put(b *testing.B) {
	runBenchmarkHamt64Put(b, KVS, Functional, TableOption)
}

func runBenchmarkHamt64Put(
	b *testing.B,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "BenchmarkHamt64Put"
	if functional {
		name += ":functional:" + hamt64.TableOptionName[tblOpt]
	} else {
		name += ":transient:" + hamt64.TableOptionName[tblOpt]
	}

	if b.N+InitHamtNumKvsForPut > len(kvs) {
		log.Printf("%s: Can't run: b.N+num > len(kvs)", name)
		b.Fatalf("%s: Can't run: b.N+num > len(kvs)", name)
	}

	var initKvs = kvs[:InitHamtNumKvsForPut]

	var h, err = buildHamt64(name, initKvs, functional, tblOpt)
	if err != nil {
		log.Printf("%s: failed buildHamt64(%q, kvs[:%d], %t, %s) => %s", name,
			name, InitHamtNumKvsForPut, functional,
			hamt64.TableOptionName[tblOpt], err)
		b.Fatalf("%s: failed buildHamt64(%q, kvs[:%d], %t, %s) => %s", name,
			name, InitHamtNumKvsForPut, functional,
			hamt64.TableOptionName[tblOpt], err)
	}

	log.Printf("%s: b.N=%d;", name, b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var k = kvs[InitHamtNumKvsForPut+i].Key
		var v = kvs[InitHamtNumKvsForPut+i].Val

		var added bool
		h, added = h.Put(k, v)
		if !added {
			log.Printf("%s: failed to h.Put(%q, %d)", name, string(k), v)
			//log.Print(h.LongString(""))
			b.Fatalf("%s: failed to h.Put(%q, %d)", name, string(k), v)
		}
	}
}

func BenchmarkHamt64Del(b *testing.B) {
	runBenchmarkHamt64Del(b, KVS, Functional, TableOption)
}

func runBenchmarkHamt64Del(
	b *testing.B,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "BenchmarkHamt64Del"
	if functional {
		name += ":functional:" + hamt64.TableOptionName[tblOpt]
	} else {
		name += ":transient:" + hamt64.TableOptionName[tblOpt]
	}

	var h, err = buildHamt64(name, kvs, functional, tblOpt)
	if err != nil {
		log.Printf("%s: failed buildHamt64(%q, kvs:%d, %t, %s) => %s", name,
			name, len(kvs), functional,
			hamt64.TableOptionName[tblOpt], err)
		b.Fatalf("%s: failed buildHamt64(%q, kvs:%d, %t, %s) => %s", name,
			name, len(kvs), functional,
			hamt64.TableOptionName[tblOpt], err)
	}

	log.Printf("%s: b.N=%d;", name, b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var k = kvs[i].Key
		var v = kvs[i].Val

		var deleted bool
		var val interface{}
		h, val, deleted = h.Del(k)
		if !deleted {
			log.Printf("%s: failed to h.Del(%q)", name, string(k))
			//log.Print(h.LongString(""))
			b.Fatalf("%s: failed to h.Del(%q)", name, string(k))
		}
		if val != v {
			log.Printf("%s: failed val,%d != v,%d", name, val, v)
			b.Fatalf("%s: failed val,%d != v,%d", name, val, v)
		}
	}
}

func BenchmarkHamt64IterFunc(b *testing.B) {
	runBenchmarkHamt64IterFunc(b, KVS, Functional, TableOption)
}

func runBenchmarkHamt64IterFunc(
	b *testing.B,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "BenchmarkHamt64IterFunc"
	if Functional {
		name += ":functional:" + hamt64.TableOptionName[tblOpt]
	} else {
		name += ":transient:" + hamt64.TableOptionName[tblOpt]
	}

	var BenchHamt64Get, err = buildHamt64(name, kvs, functional, tblOpt)
	if err != nil {
		log.Printf("%s: failed buildHamt64(%q, kvs#%d, %t, %s) => %s", name,
			name, len(kvs), false, hamt64.TableOptionName[tblOpt], err)
		b.Fatalf("%s: failed buildHamt64(%q, kvs#%d, %t, %s) => %s", name,
			name, len(kvs), false, hamt64.TableOptionName[tblOpt], err)
	}

	log.Printf("%s: b.N=%d", name, b.N)
	b.ResetTimer()

	var next = BenchHamt64Get.Iter()
	var kv, ok = next()
	for i := 0; i < b.N; i++ {
		if len(kv.Key) < 0 {
			b.Fatal("stupid test to touch the kv")
		}

		if i >= b.N {
			break
		}

		kv, ok = next()
		if !ok {
			next = BenchHamt64Get.Iter()
			kv, ok = next()
		}
	}
}

func BenchmarkHamt64IterChan(b *testing.B) {
	runBenchmarkHamt64IterChan(b, KVS, Functional, TableOption)
}

func runBenchmarkHamt64IterChan(
	b *testing.B,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "BenchmarkHamt64IterChan"
	if Functional {
		name += ":functional:" + hamt64.TableOptionName[tblOpt]
	} else {
		name += ":transient:" + hamt64.TableOptionName[tblOpt]
	}

	var BenchHamt64Get, err = buildHamt64(name, kvs, functional, tblOpt)
	if err != nil {
		log.Printf("%s: failed buildHamt64(%q, kvs#%d, %t, %s) => %s", name,
			name, len(kvs), false, hamt64.TableOptionName[tblOpt], err)
		b.Fatalf("%s: failed buildHamt64(%q, kvs#%d, %t, %s) => %s", name,
			name, len(kvs), false, hamt64.TableOptionName[tblOpt], err)
	}

	log.Printf("%s: b.N=%d", name, b.N)
	b.ResetTimer()

	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	var nents = BenchHamt64Get.Nentries()
	var iterChan = BenchHamt64Get.IterChan(20, ctx)
ForLoop:
	for i := 0; i < b.N; i++ {
		//// This code is ~100ns slower than the code after this select
		//select {
		//case <-ctx.Done():
		//	break ForLoop
		//case kv, ok := <-iterChan:
		//	if !ok {
		//		// iterChan closed
		//		break ForLoop
		//	}
		//
		//	if len(kv.Key) < 0 {
		//		b.Fatal("stupid test to touch the kv")
		//	}
		//
		//	if uint(i) == nents {
		//		cancel()
		//		ctx, cancel = context.WithCancel(context.Background())
		//		defer cancel()
		//		iterChan = BenchHamt64Get.IterChan(20, ctx)
		//	}
		//}
		var kv = <-iterChan

		if len(kv.Key) < 0 {
			b.Fatal("stupid test to touch the kv")
			break ForLoop //meaningless...just to use the ForLoop label
		}

		if uint(i) == nents {
			cancel()
			ctx, cancel = context.WithCancel(context.Background())
			defer cancel()
			iterChan = BenchHamt64Get.IterChan(20, ctx)
		}
	}
}

func BenchmarkHamt64Stats(b *testing.B) {
	runBenchmarkHamt64Stats(b, KVS, Functional, TableOption)
}

func runBenchmarkHamt64Stats(
	b *testing.B,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "BenchmarkHamt64Stats"
	if Functional {
		name += ":functional:" + hamt64.TableOptionName[tblOpt]
	} else {
		name += ":transient:" + hamt64.TableOptionName[tblOpt]
	}

	var h, err = buildHamt64(name, kvs, functional, tblOpt)
	if err != nil {
		log.Printf("%s: failed buildHamt64(%q, kvs#%d, %t, %s) => %s", name,
			name, len(kvs), false, hamt64.TableOptionName[tblOpt], err)
		b.Fatalf("%s: failed buildHamt64(%q, kvs#%d, %t, %s) => %s", name,
			name, len(kvs), false, hamt64.TableOptionName[tblOpt], err)
	}

	log.Printf("%s: b.N=%d", name, b.N)
	b.ResetTimer()

	var stats *hamt64.Stats
	for i := 0; i < b.N; i++ {
		stats = h.Stats()
	}

	log.Printf("%s: stats=%+v;\n", name, stats)
}

func BenchmarkHamt64_GetN30(b *testing.B) {
	runBenchmarkHamt64GetN(b, KVS[:30], Functional, TableOption)
}

func BenchmarkHamt64_GetN1000(b *testing.B) {
	runBenchmarkHamt64GetN(b, KVS[:1000], Functional, TableOption)
}

func BenchmarkHamt64_GetN10000(b *testing.B) {
	runBenchmarkHamt64GetN(b, KVS[:10000], Functional, TableOption)
}

func runBenchmarkHamt64GetN(
	b *testing.B,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "runBenchmarkHamt64GetN"
	if Functional {
		name += ":functional:" + hamt64.TableOptionName[tblOpt]
	} else {
		name += ":transient:" + hamt64.TableOptionName[tblOpt]
	}

	log.Println(name, b.N)

	var h, err = buildHamt64(name, kvs, functional, tblOpt)
	if err != nil {
		log.Printf("%s: failed buildHamt64(%q, kvs#%d, %t, %s) => %s", name,
			name, len(kvs), false, hamt64.TableOptionName[tblOpt], err)
		b.Fatalf("%s: failed buildHamt64(%q, kvs#%d, %t, %s) => %s", name,
			name, len(kvs), false, hamt64.TableOptionName[tblOpt], err)
	}

	log.Printf("%s: b.N=%d", name, b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var kv = kvs[i%len(kvs)]
		var val, found = h.Get(kv.Key)
		if !found {
			b.Fatalf("Failed to find h.Get(%q)", kv.Key)
		}

		if val != kv.Val {
			b.Fatalf("Retrieved val,%d != kv.Val,%d", val, kv.Val)
		}
	}
}

func BenchmarkHamt64_CalcHash(b *testing.B) {
	runBenchmarkHamt64CalcHash(b, KVS)
}

func runBenchmarkHamt64CalcHash(b *testing.B, kvs []KeyVal) {
	for i := 0; i < b.N; i++ {
		hamt64.CalcHash(kvs[i%len(kvs)].Key)
	}
}

func BenchmarkHamt64_CalcHashFromStringInterface(b *testing.B) {
	runBenchmarkHamt64CalcHashFomStringInterface(b, KVS)
}

func copyByteSlice(key []byte) []byte {
	var k = make([]byte, len(key))
	copy(k, key)
	return k
}

func toByteSlice(v interface{}) ([]byte, error) {
	switch x := v.(type) {
	case string:
		return []byte(x), nil
	default:
		var buf bytes.Buffer
		var err = binary.Write(&buf, binary.LittleEndian, v)
		if err != nil {
			err = errors.Wrapf(err, "Failed to encode %q to bytes.Buffer", v)
			return nil, err
		}
		return copyByteSlice(buf.Bytes()), nil
	}
}

func runBenchmarkHamt64CalcHashFomStringInterface(
	b *testing.B,
	kvs []KeyVal,
) {
	var keys = make([]interface{}, len(kvs))
	for i, kv := range kvs {
		keys[i] = string(kv.Key)
	}

	log.Println("b.N = ", b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var key, err = toByteSlice(keys[i%len(keys)])
		if err != nil {
			log.Panic(err)
		}
		hamt64.CalcHash(key)
	}
}

func BenchmarkHamt64_stringToByteSlice(b *testing.B) {
	var keys = make([]string, len(KVS))
	for i, kv := range KVS {
		keys[i] = string(kv.Key)
	}

	log.Println("b.N = ", b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = []byte(keys[i%len(keys)])
		//bs := []byte(keys[i%len(keys)])
		//if len(bs) < 0 {
		//	panic("never")
		//}
	}
}

//func BenchmarkHamt64_noop(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		i = i / 1
//	}
//}
