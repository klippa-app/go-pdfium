// Command corpus renders the first page of every PDF in a directory with each
// WebAssembly runtime go-pdfium supports and reports timings and pixel
// equality. All runtimes run in one process and are interleaved per document,
// so that machine state drift affects them equally.
//
//	go run ./experimental/benchmark/cmd/corpus -dir /path/to/pdfs [-runtimes wazero,wazy,wago] [-dpi 100] [-limit N] [-csv out.csv]
//
// A runtime that can not run the module on this machine is reported and left
// out.
package main

import (
	"crypto/sha256"
	"encoding/csv"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/experimental/wago"
	"github.com/klippa-app/go-pdfium/experimental/wazy"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/webassembly"

	wazyrt "github.com/samyfodil/wazy"
	wazyapi "github.com/samyfodil/wazy/api"
	"github.com/tetratelabs/wazero"
	wazeroapi "github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/experimental"
	wagort "github.com/wago-org/wago"
)

type runtimeSpec struct {
	name string
	// init creates the pool. interruptible enables the runtime's
	// close-on-context-done mode, which is what makes a Kill() of a stuck
	// document possible but costs execution speed.
	init func(interruptible bool) (pdfium.Pool, error)
}

var specs = []runtimeSpec{
	{"wazero", func(interruptible bool) (pdfium.Pool, error) {
		return webassembly.Init(webassembly.Config{MinIdle: 1, MaxIdle: 1, MaxTotal: 1, ReuseWorkers: true,
			RuntimeConfig: wazero.NewRuntimeConfig().WithCoreFeatures(wazeroapi.CoreFeaturesV2 | experimental.CoreFeaturesExceptionHandling).WithCloseOnContextDone(interruptible)})
	}},
	{"wazy", func(interruptible bool) (pdfium.Pool, error) {
		return wazy.Init(wazy.Config{MinIdle: 1, MaxIdle: 1, MaxTotal: 1, ReuseWorkers: true,
			RuntimeConfig: wazyrt.NewRuntimeConfig().WithCoreFeatures(wazyapi.CoreFeaturesV2 | wazyapi.CoreFeatureExceptionHandling).WithCloseOnContextDone(interruptible)})
	}},
	{"wago", func(interruptible bool) (pdfium.Pool, error) {
		return wago.Init(wago.Config{MinIdle: 1, MaxIdle: 1, MaxTotal: 1, ReuseWorkers: true, InterruptibleCalls: interruptible,
			RuntimeConfig: wagort.NewRuntimeConfig()})
	}},
}

type runtimeState struct {
	spec     runtimeSpec
	pool     pdfium.Pool
	instance pdfium.Pdfium
	compile  time.Duration
	times    []float64 // ms per successfully rendered document
	failed   int
	timedOut int
}

// result of one document on one runtime.
type result struct {
	ms   float64
	hash string
	err  error
}

func (s *runtimeState) render(data []byte, dpi int, timeout time.Duration) result {
	type out struct {
		r result
	}
	done := make(chan out, 1)
	start := time.Now()
	go func() {
		doc, err := s.instance.OpenDocument(&requests.OpenDocument{File: &data})
		if err != nil {
			done <- out{result{err: err}}
			return
		}
		rendered, err := s.instance.RenderPageInDPI(&requests.RenderPageInDPI{
			Page: requests.Page{ByIndex: &requests.PageByIndex{Document: doc.Document, Index: 0}},
			DPI:  dpi,
		})
		var hash string
		if err == nil {
			sum := sha256.Sum256(rendered.Result.Image.Pix)
			hash = fmt.Sprintf("%x", sum[:8])
			rendered.Cleanup()
		}
		s.instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{Document: doc.Document})
		done <- out{result{ms: float64(time.Since(start)) / 1e6, hash: hash, err: err}}
	}()
	select {
	case o := <-done:
		return o.r
	case <-time.After(timeout):
		// Kill the stuck instance and start over with a fresh one.
		s.instance.Kill()
		<-done
		instance, err := s.pool.GetInstance(time.Minute)
		if err != nil {
			panic(err)
		}
		s.instance = instance
		return result{err: fmt.Errorf("timeout after %s", timeout)}
	}
}

// readSmokeDocument returns the first PDF of the directory, used to check that
// a runtime works at all.
func readSmokeDocument(dir string) []byte {
	files, _ := filepath.Glob(filepath.Join(dir, "*.pdf"))
	sort.Strings(files)
	if len(files) == 0 {
		return nil
	}
	data, _ := os.ReadFile(files[0])
	return data
}

func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return math.NaN()
	}
	i := int(math.Ceil(p/100*float64(len(sorted)))) - 1
	return sorted[max(0, min(i, len(sorted)-1))]
}

func main() {
	dir := flag.String("dir", "", "directory with PDF files")
	names := flag.String("runtimes", "wazero,wazy,wago", "comma separated runtimes to compare")
	dpi := flag.Int("dpi", 100, "render DPI")
	limit := flag.Int("limit", 0, "only use the first N files (0 = all)")
	stride := flag.Int("stride", 1, "use every Nth file")
	timeout := flag.Duration("timeout", 60*time.Second, "per document timeout")
	csvPath := flag.String("csv", "", "write per document results to this CSV file")
	interruptible := flag.Bool("interruptible", false, "enable close-on-context-done so a stuck document can be killed (slower)")
	flag.Parse()
	if *dir == "" {
		flag.Usage()
		os.Exit(2)
	}

	files, err := filepath.Glob(filepath.Join(*dir, "*.pdf"))
	if err != nil || len(files) == 0 {
		fmt.Fprintln(os.Stderr, "no PDF files found:", err)
		os.Exit(1)
	}
	sort.Strings(files)
	var selected []string
	for i := 0; i < len(files); i += *stride {
		selected = append(selected, files[i])
	}
	if *limit > 0 && len(selected) > *limit {
		selected = selected[:*limit]
	}

	var states []*runtimeState
	for _, spec := range specs {
		if !strings.Contains(","+*names+",", ","+spec.name+",") {
			continue
		}
		start := time.Now()
		pool, err := spec.init(*interruptible)
		if err != nil {
			fmt.Printf("%s: not available: %v\n", spec.name, err)
			continue
		}
		instance, err := pool.GetInstance(time.Minute)
		if err != nil {
			panic(err)
		}
		compile := time.Since(start)
		// Check one render up front, a runtime that miscompiles the module on
		// this machine would otherwise fail on every document.
		state := &runtimeState{spec: spec, pool: pool, instance: instance, compile: compile}
		if smoke := state.render(readSmokeDocument(*dir), *dpi, *timeout); smoke.err != nil {
			fmt.Printf("%s: can not run the module on this machine: %v\n", spec.name, smoke.err)
			instance.Close()
			pool.Close()
			continue
		}
		states = append(states, state)
	}
	if len(states) == 0 {
		fmt.Println("no runtime available")
		os.Exit(1)
	}

	var writer *csv.Writer
	if *csvPath != "" {
		f, err := os.Create(*csvPath)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		writer = csv.NewWriter(f)
		defer writer.Flush()
		header := []string{"file", "bytes"}
		for _, s := range states {
			header = append(header, s.spec.name+"_ms", s.spec.name+"_hash", s.spec.name+"_error")
		}
		writer.Write(header)
	}

	// Per document ratios against the first runtime, for the geomean.
	logRatios := make([][]float64, len(states))
	mismatches := make([]int, len(states))
	compared := make([]int, len(states))
	mismatchFiles := make([][]string, len(states))
	runStart := time.Now()
	for n, file := range selected {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		// Rotate the order per document so that no runtime always runs first
		// on a cold cache.
		results := make([]result, len(states))
		for k := range states {
			i := (k + n) % len(states)
			s := states[i]
			results[i] = s.render(data, *dpi, *timeout)
			switch {
			case results[i].err != nil && strings.HasPrefix(results[i].err.Error(), "timeout"):
				s.timedOut++
			case results[i].err != nil:
				s.failed++
			default:
				s.times = append(s.times, results[i].ms)
			}
		}
		for i := 1; i < len(states); i++ {
			if results[0].err == nil && results[i].err == nil {
				logRatios[i] = append(logRatios[i], math.Log(results[i].ms/results[0].ms))
				compared[i]++
				if results[i].hash != results[0].hash {
					mismatches[i]++
					if len(mismatchFiles[i]) < 10 {
						mismatchFiles[i] = append(mismatchFiles[i], filepath.Base(file))
					}
				}
			}
		}
		if writer != nil {
			row := []string{filepath.Base(file), strconv.Itoa(len(data))}
			for _, r := range results {
				errText := ""
				if r.err != nil {
					errText = r.err.Error()
				}
				row = append(row, strconv.FormatFloat(r.ms, 'f', 3, 64), r.hash, errText)
			}
			writer.Write(row)
		}
		if (n+1)%250 == 0 {
			fmt.Fprintf(os.Stderr, "%d/%d documents, %s elapsed\n", n+1, len(selected), time.Since(runStart).Round(time.Second))
		}
	}

	fmt.Printf("documents: %d, dpi: %d, page: first, interruptible: %v\n\n", len(selected), *dpi, *interruptible)
	fmt.Println("| Runtime | Compile + first worker | Rendered | Failed | Timed out | Total render time | Mean | Median | p90 | p99 | Geomean ratio vs " + states[0].spec.name + " | Pixel mismatches vs " + states[0].spec.name + " |")
	fmt.Println("|---|---|---|---|---|---|---|---|---|---|---|---|")
	for i, s := range states {
		sorted := append([]float64(nil), s.times...)
		sort.Float64s(sorted)
		total := 0.0
		for _, t := range sorted {
			total += t
		}
		ratio, mismatch := "1.00x", "reference"
		if i > 0 {
			sum := 0.0
			for _, l := range logRatios[i] {
				sum += l
			}
			if len(logRatios[i]) > 0 {
				ratio = fmt.Sprintf("%.3fx", math.Exp(sum/float64(len(logRatios[i]))))
			}
			mismatch = fmt.Sprintf("%d of %d", mismatches[i], compared[i])
		}
		fmt.Printf("| %s | %s | %d | %d | %d | %.1f s | %.2f ms | %.2f ms | %.1f ms | %.0f ms | %s | %s |\n",
			s.spec.name, s.compile.Round(time.Millisecond), len(sorted), s.failed, s.timedOut, total/1000,
			total/float64(max(len(sorted), 1)), percentile(sorted, 50), percentile(sorted, 90), percentile(sorted, 99), ratio, mismatch)
	}
	for i := 1; i < len(states); i++ {
		if len(mismatchFiles[i]) > 0 {
			fmt.Printf("\n%s pixel mismatches (first %d): %s\n", states[i].spec.name, len(mismatchFiles[i]), strings.Join(mismatchFiles[i], ", "))
		}
	}
	for _, s := range states {
		s.instance.Close()
		s.pool.Close()
	}
}
