package canvas

// Script is the templ-charts Canvas replay layer: a small, dependency-free
// vanilla-JS module (a sibling of charts/interact's Script) that renders the
// server-recorded draw-lists into <canvas> elements.
//
// For each `canvas.tc-canvas` on the page it reads the paired
// `script[type=application/json]` (referenced by the canvas's data-tc-canvas id)
// as an array of `[opcode, args…]` ops, sizes the backing store to the canvas's
// displayed size × devicePixelRatio, maps the chart coordinate space
// (data-tc-w/h) onto it, and replays the ops into the 2D context. It re-renders
// on window resize, on the canvas's own box resizing (ResizeObserver — the
// wrapper scales down in narrow containers), and on devicePixelRatio changes
// (moving a window between displays), so the raster stays crisp and never
// overflows its container — this is all client-side, with no server round-trip
// (unlike charts/interact's opt-in data-tc-observe re-fetch). Path ops are drawn
// through Path2D, so SVG "d" strings emitted by the existing shape builders
// render unchanged. Interactivity (hover/tooltip/nearest-point) is layered
// separately by charts/interact over a transparent hit-surface; this module
// only paints.
//
// Loaded once per page via CanvasScriptTag (or by inlining Script). Safe on
// pages with no canvases — it simply finds none.
const Script = `(function () {
  'use strict';
  var TAU = Math.PI * 2;

  // replay walks the op list, mapping each [opcode, args…] to a 2D-context call.
  function replay(ctx, ops) {
    for (var i = 0; i < ops.length; i++) {
      var op = ops[i];
      switch (op[0]) {
        case 'fs': ctx.fillStyle = op[1]; break;
        case 'ss': ctx.strokeStyle = op[1]; break;
        case 'lw': ctx.lineWidth = op[1]; break;
        case 'ga': ctx.globalAlpha = op[1]; break;
        case 'ft': ctx.font = op[1]; break;
        case 'ta': ctx.textAlign = op[1]; break;
        case 'tb': ctx.textBaseline = op[1]; break;
        case 'sv': ctx.save(); break;
        case 'rs': ctx.restore(); break;
        case 'tr': ctx.translate(op[1], op[2]); break;
        case 'fr': ctx.fillRect(op[1], op[2], op[3], op[4]); break;
        case 'sr': ctx.strokeRect(op[1], op[2], op[3], op[4]); break;
        case 'fc': ctx.beginPath(); ctx.arc(op[1], op[2], op[3], 0, TAU); ctx.fill(); break;
        case 'sc': ctx.beginPath(); ctx.arc(op[1], op[2], op[3], 0, TAU); ctx.stroke(); break;
        case 'fx': ctx.fillText(op[1], op[2], op[3]); break;
        case 'ln': ctx.beginPath(); ctx.moveTo(op[1], op[2]); ctx.lineTo(op[3], op[4]); ctx.stroke(); break;
        case 'pf': ctx.fill(new Path2D(op[1])); break;
        case 'sp': ctx.stroke(new Path2D(op[1])); break;
      }
    }
  }

  // render sizes the backing store to displayed-size × devicePixelRatio and
  // maps the chart coordinate space (data-tc-w/h — what the draw-list was
  // recorded in) onto it, so the chart scales down with its container (the
  // wrapper is max-width:100% with an aspect-ratio) while staying crisp.
  function render(canvas, ops) {
    var chartW = parseFloat(canvas.getAttribute('data-tc-w')) || canvas.clientWidth || 0;
    var chartH = parseFloat(canvas.getAttribute('data-tc-h')) || canvas.clientHeight || 0;
    if (!chartW || !chartH) return;
    var rect = canvas.getBoundingClientRect();
    var dispW = rect.width || chartW;
    var dispH = rect.height || chartH;
    var dpr = window.devicePixelRatio || 1;
    canvas.width = Math.round(dispW * dpr);
    canvas.height = Math.round(dispH * dpr);
    var ctx = canvas.getContext('2d');
    if (!ctx) return;
    ctx.setTransform(dpr * dispW / chartW, 0, 0, dpr * dispH / chartH, 0, 0);
    ctx.clearRect(0, 0, chartW, chartH);
    replay(ctx, ops);
  }

  function opsFor(canvas) {
    var ref = canvas.getAttribute('data-tc-canvas');
    var el = ref && document.getElementById(ref);
    if (!el) return null;
    try { return JSON.parse(el.textContent); } catch (e) { return null; }
  }

  var painted = [];
  // ro repaints a canvas when its own box resizes (container-driven size
  // changes the window resize event misses). One shared observer; entries map
  // back to their ops via the painted list.
  var ro = typeof ResizeObserver !== 'undefined' ? new ResizeObserver(onResize) : null;
  function scan() {
    if (ro) ro.disconnect();
    painted = [];
    var canvases = document.querySelectorAll('canvas.tc-canvas[data-tc-canvas]');
    for (var i = 0; i < canvases.length; i++) {
      var c = canvases[i];
      var ops = opsFor(c);
      if (!ops) continue;
      painted.push({ canvas: c, ops: ops });
      render(c, ops);
      if (ro) ro.observe(c);
    }
  }

  function repaint() {
    for (var i = 0; i < painted.length; i++) {
      if (painted[i].canvas.isConnected !== false) render(painted[i].canvas, painted[i].ops);
    }
  }

  var raf = 0;
  function onResize() {
    if (raf) return;
    raf = window.requestAnimationFrame(function () { raf = 0; repaint(); });
  }

  function init() {
    scan();
    window.addEventListener('resize', onResize);
    // Re-scan after HTMX swaps so canvases inserted dynamically get painted.
    document.body && document.body.addEventListener('htmx:afterSwap', scan);
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();`
