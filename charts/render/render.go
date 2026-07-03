// Package render provides small convenience helpers for rendering a chart
// (or any templ.Component) to a string or an io.Writer.
//
// Every chart in this library is a templ.Component, rendered with the standard
// two-line dance:
//
//	var b strings.Builder
//	if err := bar.Bar(props).Render(context.Background(), &b); err != nil { ... }
//	html := b.String()
//
// These helpers collapse that into a single call:
//
//	svg, err := render.String(bar.Bar(props))
//
// The functions are deliberately generic over templ.Component — they carry no
// chart-specific knowledge — so they work for every chart type and for any
// composed component. String/To use a background context; the *Ctx variants
// take an explicit context.Context for cancellation.
package render

import (
	"context"
	"io"
	"strings"

	"github.com/a-h/templ"
)

// String renders c to a string using a background context.
func String(c templ.Component) (string, error) {
	return StringCtx(context.Background(), c)
}

// To renders c into w using a background context.
func To(w io.Writer, c templ.Component) error {
	return ToCtx(context.Background(), w, c)
}

// StringCtx renders c to a string using the supplied context.
func StringCtx(ctx context.Context, c templ.Component) (string, error) {
	var b strings.Builder
	if err := c.Render(ctx, &b); err != nil {
		return "", err
	}
	return b.String(), nil
}

// ToCtx renders c into w using the supplied context. It is a thin wrapper over
// c.Render provided for symmetry with StringCtx.
func ToCtx(ctx context.Context, w io.Writer, c templ.Component) error {
	return c.Render(ctx, w)
}
