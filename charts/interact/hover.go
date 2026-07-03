package interact

import "strings"

// Hover-others highlighting — the shared, client-only "dim the rest, re-light
// the hovered element and everything connected to it" pattern first built for
// charts/chord and now used by chord, sankey and network. It emits a scoped
// CSS <style> block driven entirely by the `:has(...:hover)` selector: no JS,
// no server round-trip, and a graceful no-op where `:has` is unsupported (the
// resting rules still apply). Charts tag their elements with a per-instance
// scope class plus correspondence classes; this helper turns a declarative
// description of the resting/dim/relight rules into the matching CSS.
//
// The output is deliberately whitespace-free and deterministic so charts can
// golden-test it. A chart with N hover targets emits N `:has()` groups; each
// group dims the whole family and then re-lights (via higher selector
// specificity) the elements that belong to the hovered target.

// HoverRule is one CSS rule within the highlight block: one or more selector
// suffixes (each appended to the instance scope) sharing a declaration body.
// Multiple suffixes are comma-joined, with the scope prefix repeated for each —
// so a single HoverRule can re-light several element classes at once.
type HoverRule struct {
	// Sels are selector fragments appended after the scope, e.g. ".tc-sk-node".
	// When more than one is given they are emitted as a comma-separated selector
	// list, each carrying the full prefix.
	Sels []string
	// Body is the declaration text without braces, e.g. "fill-opacity:0.35".
	Body string
}

// HoverGroup is the set of rules that apply while a single target is hovered.
// Trigger is the selector fragment (appended to the scope) whose `:hover`
// activates the group; Rules are applied in order — dim rules first, then the
// higher-specificity relight rules.
type HoverGroup struct {
	Trigger string
	Rules   []HoverRule
}

// HoverHighlight describes a complete hover-others style block for one chart
// instance. Scope is the per-instance class shared by every participating
// element (e.g. ".tc-sk<cid>"), which keeps the rules from leaking across
// multiple interactive charts on one page. Resting rules apply unconditionally
// (the resting opacity and the `:has`-unsupported fallback); Groups carry the
// per-target dim/relight rules.
type HoverHighlight struct {
	Scope   string
	Resting []HoverRule
	Groups  []HoverGroup
}

// Style renders the <style> block. It returns "" when there is nothing to emit,
// so a chart can call it unconditionally and stay byte-stable when the block is
// empty.
func (h HoverHighlight) Style() string {
	if len(h.Resting) == 0 && len(h.Groups) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<style>`)
	for _, r := range h.Resting {
		writeRule(&b, h.Scope, r)
	}
	for _, g := range h.Groups {
		pre := `svg:has(` + h.Scope + g.Trigger + `:hover) ` + h.Scope
		for _, r := range g.Rules {
			writeRulePrefixed(&b, pre, r)
		}
	}
	b.WriteString(`</style>`)
	return b.String()
}

// writeRule emits a resting rule: "<scope><sel>,<scope><sel2>{body}".
func writeRule(b *strings.Builder, scope string, r HoverRule) {
	writeRulePrefixed(b, scope, r)
}

// writeRulePrefixed emits "<pre><sel1>,<pre><sel2>...{body}".
func writeRulePrefixed(b *strings.Builder, pre string, r HoverRule) {
	for i, sel := range r.Sels {
		if i > 0 {
			b.WriteString(`,`)
		}
		b.WriteString(pre)
		b.WriteString(sel)
	}
	b.WriteString(`{`)
	b.WriteString(r.Body)
	b.WriteString(`}`)
}
