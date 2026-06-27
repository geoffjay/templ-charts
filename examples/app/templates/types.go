package templates

// LayoutProps is the input to the Layout templ component: the page title and
// the active nav section (for highlighting).
type LayoutProps struct {
	Title string
	Nav   string
}

// css is the hand-written minimal CSS for the demo app (~80 lines). Kept
// inline so the app has zero external CSS deps.
const css = `
:root { --fg:#222; --muted:#666; --border:#e2e6ec; --bg:#fff; --card:#fff; --accent:#3b82f6; }
* { box-sizing: border-box; }
body { margin:0; font: 15px/1.5 system-ui,-apple-system,Segoe UI,Roboto,sans-serif; color:var(--fg); background:#f7f8fa; }
a { color:var(--accent); text-decoration:none; }
a:hover { text-decoration:underline; }
header { background:var(--bg); border-bottom:1px solid var(--border); padding:14px 24px; }
header h1 { margin:0; font-size:18px; font-weight:600; display:inline-block; }
header h1 a { color:var(--fg); }
nav { display:inline-block; margin-left:24px; }
nav a { margin-right:16px; color:var(--muted); }
nav a.active { color:var(--fg); font-weight:600; }
main { max-width:1200px; margin:0 auto; padding:24px; }
.intro { color:var(--muted); margin:0 0 24px; }
.grid { display:grid; grid-template-columns:repeat(auto-fill,minmax(360px,1fr)); gap:20px; }
.card { background:var(--card); border:1px solid var(--border); border-radius:8px; padding:18px; }
.card h2 { margin:0 0 4px; font-size:16px; }
.card p { margin:0 0 14px; color:var(--muted); font-size:13px; }
.chart { width:100%; height:auto; position:relative; }
.chart svg { width:100%; height:auto; display:block; }
.tooltip { position:absolute; pointer-events:none; z-index:10; background:rgba(255,255,255,0.97); border:1px solid var(--border); border-radius:4px; padding:5px 9px; font-size:12px; box-shadow:0 1px 3px rgba(0,0,0,0.12); }
.list { list-style:none; padding:0; margin:0; }
.list li { padding:10px 14px; border:1px solid var(--border); border-radius:6px; margin-bottom:10px; background:var(--card); }
.list li a { font-weight:600; }
.list li p { margin:4px 0 0; color:var(--muted); font-size:13px; }
footer { text-align:center; color:var(--muted); font-size:12px; padding:30px; }
code { background:#eef; padding:1px 4px; border-radius:3px; font-size:13px; }
`
