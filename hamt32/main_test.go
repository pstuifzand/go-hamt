package hamt32_test

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/lleo/go-hamt/hamt32"
	"github.com/lleo/stringutil"
	"github.com/pkg/errors"
)

type StrVal struct {
	Str string
	Val interface{}
}

// 4 million & change
var InitHamtNumSvsForPut = 1024 * 1024
var InitHamtNumSvs = (2 * 1024 * 1024) + InitHamtNumSvsForPut
var numSvs = InitHamtNumSvs + (4 * 1024)
var TwoKK = 2 * 1024 * 1024
var SVS []StrVal

var Functional bool
var TableOption int

var Hamt32 hamt32.Hamt

var Inc = stringutil.Lower.Inc

var StartTime = make(map[string]time.Time)
var RunTime = make(map[string]time.Duration)

func TestMain(m *testing.M) {
	var fixedonly, sparseonly, hybrid, all bool
	flag.BoolVar(&fixedonly, "F", false,
		"Use fixed tables only and exclude C and H Options.")
	flag.BoolVar(&sparseonly, "S", false,
		"Use sparse tables only and exclude F and H Options.")
	flag.BoolVar(&hybrid, "H", false,
		"Use sparse tables initially and exclude F and S Options.")
	flag.BoolVar(&all, "A", false,
		"Run all Tests w/ Options set to FixedTables, SparseTables, and HybridTables")

	var functional, transient, both bool
	flag.BoolVar(&functional, "f", false,
		"Run Tests against HamtFunctional struct; excludes transient option")
	flag.BoolVar(&transient, "t", false,
		"Run Tests against HamtFunctional struct; excludes functional option")
	flag.BoolVar(&both, "b", false,
		"Run Tests against both transient and functional Hamt types.")

	flag.Parse()

	// If all flag set, ignore fixedonly, sparseonly, and hybrid.
	if !all {

		// only one flag may be set between fixedonly, sparseonly, and hybrid
		if (fixedonly && (sparseonly || hybrid)) ||
			(sparseonly && (fixedonly || hybrid)) ||
			(hybrid && (sparseonly || fixedonly)) {
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	// If no flags given, run all tests.
	if !(all || fixedonly || sparseonly || hybrid) {
		all = true
	}

	if !both {
		if functional && transient {
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	if !(both || functional || transient) {
		both = true
	}

	log.SetFlags(log.Lshortfile)

	var logfn = fmt.Sprintf("test-%d.log", hamt32.IndexBits)
	var logfile, err = os.Create(logfn)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to os.Create(\"test.log\")"))
	}
	defer logfile.Close()

	log.SetOutput(logfile)

	log.Println("TestMain: and so it begins...")

	SVS = buildStrVals("TestMain", numSvs)

	log.Printf("TestMain: HashSize=%d\n", hamt32.HashSize)
	fmt.Printf("TestMain: HashSize=%d\n", hamt32.HashSize)
	log.Printf("TestMain: IndexBits=%d\n", hamt32.IndexBits)
	fmt.Printf("TestMain: IndexBits=%d\n", hamt32.IndexBits)
	log.Printf("TestMain: IndexLimit=%d\n", hamt32.IndexLimit)
	fmt.Printf("TestMain: IndexLimit=%d\n", hamt32.IndexLimit)
	log.Printf("TestMain: DepthLimit=%d\n", hamt32.DepthLimit)
	fmt.Printf("TestMain: DepthLimit=%d\n", hamt32.DepthLimit)
	log.Printf("TestMain: SizeofFixedTable=%d\n", hamt32.SizeofFixedTable)
	fmt.Printf("TestMain: SizeofFixedTable=%d\n", hamt32.SizeofFixedTable)
	log.Printf("TestMain: SizeofSparseTable=%d\n", hamt32.SizeofSparseTable)
	fmt.Printf("TestMain: SizeofSparseTable=%d\n", hamt32.SizeofSparseTable)
	log.Printf("TestMain: BitmapSize=%d\n", hamt32.BitmapSize)
	fmt.Printf("TestMain: BitmapSize=%d\n", hamt32.BitmapSize)

	// execute
	var xit int
	if all {
		if both {
			Functional = false
			log.Printf("TestMain: Functional=%t;\n", Functional)
			fmt.Printf("TestMain: Functional=%t;\n", Functional)

			xit = executeAll(m)
			if xit != 0 {
				log.Printf("\n", RunTimes())
				os.Exit(xit)
			}

			Hamt32 = nil

			Functional = true
			log.Printf("TestMain: Functional=%t;\n", Functional)
			fmt.Printf("TestMain: Functional=%t;\n", Functional)

			xit = executeAll(m)
		} else if functional {
			Functional = true
			log.Printf("TestMain: Functional=%t;\n", Functional)
			fmt.Printf("TestMain: Functional=%t;\n", Functional)

			xit = executeAll(m)
		} else if transient {
			Functional = false
			log.Printf("TestMain: Functional=%t;\n", Functional)
			fmt.Printf("TestMain: Functional=%t;\n", Functional)

			xit = executeAll(m)
		}
	} else {
		if hybrid {
			TableOption = hamt32.HybridTables
		} else if fixedonly {
			TableOption = hamt32.FixedTables
		} else /* if sparseonly */ {
			TableOption = hamt32.SparseTables
		}

		if both {
			Functional = false

			log.Printf("TestMain: Functional=%t;\n", Functional)
			fmt.Printf("TestMain: Functional=%t;\n", Functional)
			log.Printf("TestMain: TableOption=%s;\n",
				hamt32.TableOptionName[TableOption])
			fmt.Printf("TestMain: TableOption=%s;\n",
				hamt32.TableOptionName[TableOption])

			xit = m.Run()
			if xit != 0 {
				log.Printf("\n", RunTimes())
				os.Exit(xit)
			}

			Hamt32 = nil
			Functional = true

			log.Printf("TestMain: Functional=%t;\n", Functional)
			fmt.Printf("TestMain: Functional=%t;\n", Functional)
			log.Printf("TestMain: TableOption=%s;\n",
				hamt32.TableOptionName[TableOption])
			fmt.Printf("TestMain: TableOption=%s;\n",
				hamt32.TableOptionName[TableOption])

			xit = m.Run()
		} else {
			if functional {
				Functional = true
			} else /* if transient */ {
				Functional = false
			}

			log.Printf("TestMain: Functional=%t;\n", Functional)
			fmt.Printf("TestMain: Functional=%t;\n", Functional)
			log.Printf("TestMain: TableOption=%s;\n",
				hamt32.TableOptionName[TableOption])
			fmt.Printf("TestMain: TableOption=%s;\n",
				hamt32.TableOptionName[TableOption])
			xit = m.Run()
		}
	}

	log.Println("\n", RunTimes())
	os.Exit(xit)
}

func executeAll(m *testing.M) int {
	TableOption = hamt32.FixedTables

	log.Printf("TestMain: TableOption=%s;\n",
		hamt32.TableOptionName[TableOption])
	fmt.Printf("TestMain: TableOption=%s;\n",
		hamt32.TableOptionName[TableOption])

	var xit = m.Run()
	if xit != 0 {
		log.Println("\n", RunTimes())
		os.Exit(1)
	}

	Hamt32 = nil
	TableOption = hamt32.SparseTables

	log.Printf("TestMain: TableOption=%s;\n",
		hamt32.TableOptionName[TableOption])
	fmt.Printf("TestMain: TableOption=%s;\n",
		hamt32.TableOptionName[TableOption])

	xit = m.Run()
	if xit != 0 {
		log.Println("\n", RunTimes())
		os.Exit(1)
	}

	Hamt32 = nil
	TableOption = hamt32.HybridTables

	log.Printf("TestMain: TableOption=%s;\n",
		hamt32.TableOptionName[TableOption])
	fmt.Printf("TestMain: TableOption=%s;\n",
		hamt32.TableOptionName[TableOption])

	xit = m.Run()

	return xit
}

func buildStrVals(prefix string, num int) []StrVal {
	var name = fmt.Sprintf("%s-buildStrVals-%d", prefix, num)
	StartTime[name] = time.Now()

	var svs = make([]StrVal, num)
	var s = "aaa"

	for i := 0; i < num; i++ {
		svs[i] = StrVal{s, i}
		s = Inc(s)
	}

	RunTime[name] = time.Since(StartTime[name])
	return svs
}

func buildHamt32(
	prefix string,
	svs []StrVal,
	functional bool,
	opt int,
) (hamt32.Hamt, error) {
	var name = fmt.Sprintf("%s-buildHamt32-%d", prefix, len(svs))

	StartTime[name] = time.Now()
	var h = hamt32.New(functional, opt)
	for _, sv := range svs {
		var s = sv.Str
		var v = sv.Val

		var inserted bool
		h, inserted = h.Put([]byte(s), v)
		if !inserted {
			return nil, fmt.Errorf("failed to Put(%q, %v)", s, v)
		}
	}
	RunTime[name] = time.Since(StartTime[name])

	return h, nil
}

func RunTimes() string {
	// Grab list of keys from RunTime map; MAJOR un-feature of Go!
	var ks = make([]string, len(RunTime))
	var i int = 0
	for k := range RunTime {
		ks[i] = k
		i++
	}
	sort.Strings(ks)

	var s = ""

	s += "Key                                                Val\n"
	s += "==================================================+==========\n"

	var tot time.Duration
	for _, k := range ks {
		v := RunTime[k]
		s += fmt.Sprintf("%-50s %s\n", k, v)
		tot += v
	}
	s += fmt.Sprintf("%50s %s\n", "TOTAL", tot)

	return s
}
