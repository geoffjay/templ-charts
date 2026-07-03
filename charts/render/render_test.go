package render

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func TestString(t *testing.T) {
	got, err := String(templ.Raw("<svg></svg>"))
	if err != nil {
		t.Fatalf("String: unexpected error: %v", err)
	}
	if got != "<svg></svg>" {
		t.Fatalf("String: got %q, want %q", got, "<svg></svg>")
	}
}

func TestTo(t *testing.T) {
	var b strings.Builder
	if err := To(&b, templ.Raw("<g/>")); err != nil {
		t.Fatalf("To: unexpected error: %v", err)
	}
	if b.String() != "<g/>" {
		t.Fatalf("To: got %q, want %q", b.String(), "<g/>")
	}
}

func TestStringCtx_PropagatesContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), ctxKey{}, "v")
	var seen any
	c := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		seen = ctx.Value(ctxKey{})
		return nil
	})
	if _, err := StringCtx(ctx, c); err != nil {
		t.Fatalf("StringCtx: unexpected error: %v", err)
	}
	if seen != "v" {
		t.Fatalf("StringCtx: context not propagated, got %v", seen)
	}
}

func TestString_ReturnsRenderError(t *testing.T) {
	want := errors.New("boom")
	c := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return want
	})
	got, err := String(c)
	if !errors.Is(err, want) {
		t.Fatalf("String: got err %v, want %v", err, want)
	}
	if got != "" {
		t.Fatalf("String: got %q on error, want empty", got)
	}
}

type ctxKey struct{}
