package interact

// Script is the templ-charts client-side interactivity layer: a small,
// dependency-free vanilla-JS module that handles the *ephemeral* hover
// interactions the v2 plan moves off the server (docs/PLAN-v2.md §5):
//
//   - Element tooltips: any SVG element carrying a data-tc-tooltip attribute
//     (or whose ancestor does) inside a .tc-chart container shows a tooltip
//     with that attribute's HTML, positioned at the cursor and tracking it.
//   - Mesh hover: an element carrying data-tc-mesh (a JSON array of
//     {x,y,html} points in the chart's inner SVG coordinates) resolves the
//     nearest point on mousemove, shows its tooltip, and draws a crosshair —
//     replacing the line chart's per-mousemove server round-trip.
//
// It also provides an opt-in ResizeObserver (data-tc-observe): elements that
// carry that attribute are re-fetched on container resize for a pixel-accurate
// re-render (tick/label density on axes-heavy charts), with the new content-box
// width/height appended as ?w=&h=. This is the only part that round-trips, and
// only for elements that explicitly opt in — cosmetic fluid scaling is already
// covered by the Responsive viewBox, so pages that never set data-tc-observe
// pay nothing.
//
// State-changing interactions (series toggle, active-arc) stay on the server
// via HTMX; the hover layer never round-trips. It reuses the .tc-chart /
// .tc-chart-tooltip container conventions emitted by htmx.Mount, so charts
// embedded with zero JS still render statically. Loaded once per page via
// ScriptTag (or by inlining Script).
const Script = `(function () {
  'use strict';

  // tooltipFor returns the .tc-chart-tooltip associated with a chart container
  // (the sibling emitted by htmx.Mount, or any tooltip within the same parent).
  // For static embeds with no pre-rendered tooltip, one is created lazily and
  // appended inside the chart container (made position:relative if needed) so
  // client hover works with zero extra markup.
  function tooltipFor(chart) {
    var t = chart.nextElementSibling;
    if (t && t.classList && t.classList.contains('tc-chart-tooltip')) return t;
    var p = chart.parentElement;
    if (p) {
      var found = p.querySelector('.tc-chart-tooltip');
      if (found) return found;
    }
    if (chart._tcTooltip && chart._tcTooltip.parentNode) return chart._tcTooltip;
    if (getComputedStyle(chart).position === 'static') chart.style.position = 'relative';
    var div = document.createElement('div');
    div.className = 'tc-chart-tooltip';
    div.style.position = 'absolute';
    div.style.pointerEvents = 'none';
    div.style.display = 'none';
    chart.appendChild(div);
    chart._tcTooltip = div;
    return div;
  }

  // position places the tooltip near the cursor, flipping sides to stay inside
  // the chart bounds. Mirrors the original inline demo script's logic.
  function position(t, chart, cx, cy) {
    var op = t.offsetParent || document.body;
    var opr = op.getBoundingClientRect();
    var chr = chart.getBoundingClientRect();
    var left = cx - opr.left + 12;
    var top = cy - opr.top + 12;
    if (cx + 12 + t.offsetWidth > chr.right) left = cx - opr.left - t.offsetWidth - 12;
    if (cy + 12 + t.offsetHeight > chr.bottom) top = cy - opr.top - t.offsetHeight - 12;
    t.style.left = left + 'px';
    t.style.top = top + 'px';
  }

  function show(chart, html, cx, cy) {
    var t = tooltipFor(chart);
    if (!t) return;
    if (t.innerHTML !== html) t.innerHTML = html;
    t.style.display = 'block';
    position(t, chart, cx, cy);
  }

  function hide(chart) {
    var t = tooltipFor(chart);
    if (t) { t.style.display = 'none'; t.innerHTML = ''; }
    hideCrosshair(chart);
  }

  // --- crosshair (mesh hover) ---------------------------------------------

  var SVGNS = 'http://www.w3.org/2000/svg';

  // crosshairFor lazily creates a <g class="tc-crosshair"> as a sibling of the
  // mesh element, so it shares the mesh's (margin-translated) coordinate space.
  function crosshairFor(mesh) {
    var parent = mesh.parentNode;
    if (!parent) return null;
    var g = parent._tcCrosshair;
    if (g && g.parentNode === parent) return g;
    g = document.createElementNS(SVGNS, 'g');
    g.setAttribute('class', 'tc-crosshair');
    g.style.pointerEvents = 'none';
    for (var i = 0; i < 2; i++) {
      var l = document.createElementNS(SVGNS, 'line');
      l.setAttribute('stroke', '#000');
      l.setAttribute('stroke-width', '1');
      l.setAttribute('stroke-opacity', '0.35');
      l.setAttribute('stroke-dasharray', '6 6');
      g.appendChild(l);
    }
    parent.appendChild(g);
    parent._tcCrosshair = g;
    return g;
  }

  function hideCrosshair(chart) {
    if (!chart) return;
    var gs = chart.querySelectorAll('.tc-crosshair');
    for (var i = 0; i < gs.length; i++) gs[i].style.display = 'none';
  }

  // localPoint maps a screen-space event to the element's local SVG user
  // coordinates via the inverse screen CTM (correct under viewBox scaling).
  function localPoint(el, e) {
    var svg = el.ownerSVGElement;
    var m = el.getScreenCTM && el.getScreenCTM();
    if (!svg || !svg.createSVGPoint || !m) return { x: e.offsetX, y: e.offsetY };
    var p = svg.createSVGPoint();
    p.x = e.clientX;
    p.y = e.clientY;
    return p.matrixTransform(m.inverse());
  }

  function meshPoints(mesh) {
    if (mesh._tcPts) return mesh._tcPts;
    try {
      mesh._tcPts = JSON.parse(mesh.getAttribute('data-tc-mesh') || '[]');
    } catch (err) {
      mesh._tcPts = [];
    }
    return mesh._tcPts;
  }

  // meshRadius2 returns the squared detection radius for a mesh element, or
  // Infinity when unset (unlimited). Points farther than this from the cursor
  // are not detected — matching nivo's voronoi-mesh detectionRadius.
  function meshRadius2(mesh) {
    if (mesh._tcR2 === undefined) {
      var r = parseFloat(mesh.getAttribute('data-tc-mesh-radius'));
      mesh._tcR2 = (r > 0) ? r * r : Infinity;
    }
    return mesh._tcR2;
  }

  function handleMesh(chart, mesh, e) {
    var pts = meshPoints(mesh);
    if (!pts.length) { hide(chart); return; }
    var lp = localPoint(mesh, e);
    var best = null, bestD = Infinity;
    for (var i = 0; i < pts.length; i++) {
      var dx = pts[i].x - lp.x, dy = pts[i].y - lp.y, d = dx * dx + dy * dy;
      if (d < bestD) { bestD = d; best = pts[i]; }
    }
    if (!best || bestD > meshRadius2(mesh)) { hide(chart); return; }
    show(chart, best.html || '', e.clientX, e.clientY);
    if (mesh.getBBox) {
      var g = crosshairFor(mesh);
      if (g) {
        var bb = mesh.getBBox();
        var ln = g.childNodes;
        ln[0].setAttribute('x1', best.x); ln[0].setAttribute('x2', best.x);
        ln[0].setAttribute('y1', bb.y); ln[0].setAttribute('y2', bb.y + bb.height);
        ln[1].setAttribute('x1', bb.x); ln[1].setAttribute('x2', bb.x + bb.width);
        ln[1].setAttribute('y1', best.y); ln[1].setAttribute('y2', best.y);
        g.style.display = '';
      }
    }
  }

  // isClientChart reports whether a chart participates in the client hover
  // layer (it contains elements with data-tc-tooltip / data-tc-mesh). Charts
  // that drive their tooltip through the server (htmx) hover path instead —
  // e.g. bar and pie — have neither, so this script leaves them entirely alone
  // and never clobbers their htmx-populated tooltip. Cached per container.
  function isClientChart(chart) {
    if (chart._tcClient === undefined) {
      chart._tcClient = !!chart.querySelector('[data-tc-tooltip],[data-tc-mesh]');
    }
    return chart._tcClient;
  }

  // --- delegation ----------------------------------------------------------

  document.addEventListener('mousemove', function (e) {
    var chart = e.target.closest ? e.target.closest('.tc-chart') : null;
    if (!chart || !isClientChart(chart)) return;
    var el = e.target.closest('[data-tc-tooltip]');
    if (el && chart.contains(el)) {
      hideCrosshair(chart);
      show(chart, el.getAttribute('data-tc-tooltip'), e.clientX, e.clientY);
      return;
    }
    var mesh = e.target.closest('[data-tc-mesh]');
    if (mesh && chart.contains(mesh)) { handleMesh(chart, mesh, e); return; }
    hide(chart);
  });

  // Hide when leaving a chart entirely (capture phase: mouseleave doesn't
  // bubble). Only for client-managed charts — htmx charts handle their own
  // leave reset.
  document.addEventListener('mouseout', function (e) {
    var chart = e.target.closest ? e.target.closest('.tc-chart') : null;
    if (!chart || !isClientChart(chart)) return;
    var to = e.relatedTarget;
    if (!to || !chart.contains(to)) hide(chart);
  }, true);

  // --- opt-in ResizeObserver re-fetch --------------------------------------
  // Elements carrying data-tc-observe="<url>" are re-fetched on resize for a
  // pixel-accurate re-render; the new content-box w/h are appended to the URL.
  // Uses htmx when present (preserving hx-* swap semantics), else fetch +
  // innerHTML. Purely opt-in — pages without data-tc-observe are untouched.
  function observeURL(el, w, h) {
    var base = el.getAttribute('data-tc-observe');
    if (!base) return null;
    return base + (base.indexOf('?') === -1 ? '?' : '&') + 'w=' + w + '&h=' + h;
  }

  function refetch(el) {
    var r = el.getBoundingClientRect();
    var w = Math.round(r.width), h = Math.round(r.height);
    if (!w || !h || (el._tcW === w && el._tcH === h)) return;
    el._tcW = w; el._tcH = h;
    var url = observeURL(el, w, h);
    if (!url) return;
    var sel = el.getAttribute('data-tc-observe-target');
    var tgt = sel ? document.querySelector(sel) : el;
    if (!tgt) return;
    if (window.htmx && typeof window.htmx.ajax === 'function') {
      window.htmx.ajax('GET', url, { target: tgt, swap: 'innerHTML' });
      return;
    }
    fetch(url).then(function (res) { return res.text(); }).then(function (html) {
      tgt.innerHTML = html;
    }).catch(function () {});
  }

  if (typeof ResizeObserver === 'function') {
    var ro = new ResizeObserver(function (entries) {
      for (var i = 0; i < entries.length; i++) {
        (function (el) {
          // The observer fires once on observe(); record the baseline size and
          // skip so we only re-fetch on an actual later change.
          var rr = el.getBoundingClientRect();
          if (el._tcW === undefined) { el._tcW = Math.round(rr.width); el._tcH = Math.round(rr.height); return; }
          clearTimeout(el._tcTimer);
          el._tcTimer = setTimeout(function () { refetch(el); }, 150);
        })(entries[i].target);
      }
    });
    var observeAll = function () {
      var els = document.querySelectorAll('[data-tc-observe]');
      for (var i = 0; i < els.length; i++) {
        if (!els[i]._tcObserved) { els[i]._tcObserved = true; ro.observe(els[i]); }
      }
    };
    observeAll();
    document.addEventListener('htmx:afterSwap', observeAll);
    if (document.readyState === 'loading') document.addEventListener('DOMContentLoaded', observeAll);
  }
})();`
