// Package golden provides a small snapshot-test helper for the templ-charts
// repo. A golden file holds the expected output (e.g. an SVG string or an SVG
// path-data string). Tests call Assert with the actual output; on mismatch
// the failure includes a diff-style message and the test fails.
//
// To (re)generate the golden files after an intentional change, run:
//
//	go test ./... -update
//
// (any package importing internal/golden honors the -update flag). The golden
// files live under each package's testdata/golden/ directory and are committed
// to the repo so CI catches unintended output drift.
package golden

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// update is the package-wide flag controlling golden-file regeneration. It
// is registered once via init(); importing this package is enough to make
// `-update` available on `go test`.
var update = flag.Bool("update", false, "regenerate golden files")

// Assert compares `got` against the golden file named `name` (without
// extension) inside the calling package's testdata/golden/ directory. On
// mismatch it fails the test with a unified-diff-ish message. When `-update`
// is passed, the golden file is (re)written and the test is treated as
// passing (no failure).
func Assert(t *testing.T, name, got string) {
	t.Helper()
	path := goldenPath(name)
	if *update {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("golden: mkdir %s: %v", filepath.Dir(path), err)
		}
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatalf("golden: write %s: %v", path, err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("golden: read %s: %v\n(if this is a new snapshot, run `go test ./... -update`", path, err)
	}
	if string(want) != got {
		t.Fatalf("golden mismatch for %s:\n--- want (%s)\n+++ got\n%s",
			name, path, diff(string(want), got))
	}
}

// goldenPath returns the absolute path to the golden file for `name`,
// relative to the calling test's source file. The calling file is resolved
// via runtime.Caller so each package's golden files live alongside its tests.
func goldenPath(name string) string {
	_, file, _, _ := runtime.Caller(2) // 2: skip Assert + goldenPath
	dir := filepath.Join(filepath.Dir(file), "testdata", "golden")
	return filepath.Join(dir, name+".txt")
}

// diff produces a small unified-diff-style string showing the first hunk
// where want and got diverge. It is intentionally compact (no full diff
// library) — good enough to point a human at the change.
func diff(want, got string) string {
	wl := strings.Split(want, "\n")
	gl := strings.Split(got, "\n")
	var b strings.Builder
	n := len(wl)
	if len(gl) > n {
		n = len(gl)
	}
	for i := 0; i < n; i++ {
		var w, g string
		if i < len(wl) {
			w = wl[i]
		}
		if i < len(gl) {
			g = gl[i]
		}
		if w != g {
			b.WriteString("  ")
			b.WriteString(w)
			b.WriteString("\n+ ")
			b.WriteString(g)
			b.WriteString("\n")
		}
	}
	return b.String()
}
