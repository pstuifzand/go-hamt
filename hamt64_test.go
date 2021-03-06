package hamt_test

import (
	"log"
	"testing"
	"time"

	"github.com/pstuifzand/go-hamt/hamt64"
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
		var k = hamt64.StringKey(kv.Key)
		var v = kv.Val

		var inserted bool
		h, inserted = h.Put(k, v)
		if !inserted {
			log.Printf("%s: failed to insert s=%q, v=%d", name, k, v)
			t.Fatalf("%s: failed to insert s=%q, v=%d", name, k, v)
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
		var k = hamt64.StringKey(kv.Key)
		var v = kv.Val

		var inserted bool
		Hamt64, inserted = Hamt64.Put(k, v)
		if !inserted {
			log.Printf("%s: failed to Hamt64.Put(%q, %v)", name, k, v)
			t.Fatalf("%s: failed to Hamt64.Put(%q, %v)", name, k, v)
		}

		var val, found = Hamt64.Get(k)
		if !found {
			log.Printf("%s: failed to Hamt64.Get(%q)", name, k)
			//log.Print(Hamt64.LongString(""))
			t.Fatalf("%s: failed to Hamt64.Get(%q)", name, k)
		}
		if val != v {
			log.Printf("%s: returned val,%d != expected v,%d for s=%q", name, val, v, k)
			t.Fatalf("%s: returned val,%d != expected v,%d for s=%q", name, val, v, k)
		}
	}
	RunTime[name] = time.Since(StartTime[name])

	StartTime["Hamt64.Stats()"] = time.Now()
	var stats = Hamt64.Stats()
	RunTime["Hamt64.Stats()"] = time.Since(StartTime["Hamt64.Stats()"])
	log.Printf("%s: stats=%+v;\n", name, stats)
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
		var k = hamt64.StringKey(kv.Key)
		var v = kv.Val

		var val, found = Hamt64.Get(k)
		if !found {
			log.Printf("%s: Failed to Hamt64.Get(%q)", name, k)
			log.Print(Hamt64.LongString(""))
			t.Fatalf("%s: Failed to Hamt64.Get(%q)", name, k)
		}
		if val != v {
			log.Printf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, k)
			t.Fatalf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, k)
		}
	}
	RunTime[name] = time.Since(StartTime[name])
}

func TestHamt64Range(t *testing.T) {
	runTestHamt64Range(t, KVS, Functional, TableOption)
}

func runTestHamt64Range(
	t *testing.T,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "TestHamt64Range"
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

	var kvMap = make(map[string]int, len(KVS))
	for _, kv := range KVS {
		kvMap[kv.Key] = kv.Val
	}

	// Reconstruct KVS as tmpKVS
	var totalKvs int
	var visitKeyVal = func(k hamt64.KeyI, v interface{}) bool {
		var s = string(k.(hamt64.StringKey))
		var i = v.(int)
		var expected_i, found = kvMap[s]

		if !found {
			t.Fatalf("%s: Range(visitKeyVal) KeyI.(StringKey),%q not in kvMap",
				name, s)
		}

		if expected_i != i {
			t.Fatalf("%s: Range(visitKeyVal) for KeyI.(StringKey),%q found i,%d != expected_i,%d", name, s, i, expected_i)
		}

		totalKvs++
		return true
	}
	Hamt64.Range(visitKeyVal)

	if totalKvs != len(KVS) {
		t.Fatalf("%s: Range(visitKeyVal) found totalKvs,%d != len(KVS),%d",
			name, totalKvs, len(KVS))
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
		var k = hamt64.StringKey(kv.Key)
		var v = kv.Val

		var val interface{}
		var deleted bool
		Hamt64, val, deleted = Hamt64.Del(k)
		if !deleted {
			log.Printf("%s: Failed to Hamt64.Del(%q)", name, k)
			log.Print(Hamt64.LongString(""))
			t.Fatalf("%s: Failed to Hamt64.Del(%q)", name, k)
		}
		if val != v {
			log.Printf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, k)
			t.Fatalf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, k)
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
		var k = hamt64.StringKey(kvs[j].Key)
		var v = kvs[j].Val

		var val, found = BenchHamt64Get.Get(k)
		if !found {
			log.Printf("%s: Failed to h.Get(%q)", name, k)
			//log.Print(h.LongString(""))
			b.Fatalf("%s: Failed to h.Get(%q)", name, k)
		}
		if val != v {
			log.Printf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, k)
			b.Fatalf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, k)
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
		var k = hamt64.StringKey(kvs[InitHamtNumKvsForPut+i].Key)
		var v = kvs[InitHamtNumKvsForPut+i].Val

		var added bool
		h, added = h.Put(k, v)
		if !added {
			log.Printf("%s: failed to h.Put(%q, %d)", name, k, v)
			//log.Print(h.LongString(""))
			b.Fatalf("%s: failed to h.Put(%q, %d)", name, k, v)
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
		var k = hamt64.StringKey(kvs[i].Key)
		var v = kvs[i].Val

		var deleted bool
		var val interface{}
		h, val, deleted = h.Del(k)
		if !deleted {
			log.Printf("%s: failed to h.Del(%q)", name, k)
			//log.Print(h.LongString(""))
			b.Fatalf("%s: failed to h.Del(%q)", name, k)
		}
		if val != v {
			log.Printf("%s: failed val,%d != v,%d", name, val, v)
			b.Fatalf("%s: failed val,%d != v,%d", name, val, v)
		}
	}
}

func BenchmarkHamt64Range(b *testing.B) {
	runBenchmarkHamt64Range(b, KVS, Functional, TableOption)
}

func runBenchmarkHamt64Range(
	b *testing.B,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "BenchmarkHamt64Range"
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

	var kvMap = make(map[string]int, len(KVS))
	for _, kv := range KVS {
		kvMap[kv.Key] = kv.Val
	}

	log.Printf("%s: b.N=%d", name, b.N)
	b.ResetTimer()

	var i int
	h.Range(func(k hamt64.KeyI, v interface{}) bool {
		var sk = string(k.(hamt64.StringKey))
		var iv = v.(int)

		if kvMap[sk] != iv {
			b.Fatalf("%s: for kvMap[%q],%d != i,%d", name, sk, kvMap[sk], iv)
		}

		i++
		if i >= b.N {
			return false //stop Range()
		}

		return true
	})
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
		var val, found = h.Get(hamt64.StringKey(kv.Key))
		if !found {
			b.Fatalf("Failed to find h.Get(%q)", kv.Key)
		}

		if val != kv.Val {
			b.Fatalf("Retrieved val,%d != kv.Val,%d", val, kv.Val)
		}
	}
}
