(function() {
  var cfg = window.__SCROBBLET__;
  var card      = document.getElementById('now-playing-card');
  var bar       = document.getElementById('np-bar');
  var timeEl    = document.getElementById('np-time');
  var titleEl   = document.getElementById('np-title');
  var metaEl    = document.getElementById('np-meta');
  var statusEl  = document.getElementById('np-status');
  var labelEl   = document.getElementById('np-label-text');

  var rafId = null;
  var pos = 0, dur = 0, playing = false, animStart = 0;

  function fmt(ms) {
    var s = Math.floor(ms / 1000), m = Math.floor(s / 60);
    s = s % 60;
    return m + ':' + (s < 10 ? '0' : '') + s;
  }

  function tick() {
    var current = playing ? Math.min(pos + (Date.now() - animStart), dur) : pos;
    bar.style.width = (dur > 0 ? (current / dur * 100) : 0) + '%';
    timeEl.textContent = fmt(current) + ' / ' + fmt(dur);
    if (playing && current < dur) {
      rafId = requestAnimationFrame(tick);
    }
  }

  function startAnimation(posMs, durMs, isPlaying) {
    if (rafId !== null) { cancelAnimationFrame(rafId); rafId = null; }
    pos = posMs;
    dur = durMs;
    playing = isPlaying;
    animStart = Date.now();
    tick();
  }

  function applyTrackChange(ev) {
    if (!ev.title) {
      card.style.display = 'none';
      if (rafId !== null) { cancelAnimationFrame(rafId); rafId = null; }
      return;
    }
    titleEl.textContent = ev.title;
    metaEl.textContent  = ev.artist;
    statusEl.className  = 'np-status ' + (ev.isPlaying ? 'np-status-playing' : 'np-status-paused');
    labelEl.textContent = ev.isPlaying ? 'Now Playing' : 'Paused';
    card.style.display  = '';
    startAnimation(ev.positionMs, ev.durationMs, ev.isPlaying);
  }

  if (cfg.nowPlaying) {
    startAnimation(cfg.nowPlaying.positionMs, cfg.nowPlaying.durationMs, cfg.nowPlaying.isPlaying);
  }

  var targets = cfg.targets;
  var svg = document.getElementById('graph-svg');

  var NS = 'http://www.w3.org/2000/svg';
  var DURATION = 900;

  function findTargetIndex(label) {
    for (var i = 0; i < targets.length; i++) {
      if (targets[i].label === label) return targets[i].index;
    }
    return -1;
  }

  function flashTarget(index, success) {
    var rect = document.getElementById('target-rect-' + index);
    if (!rect) return;
    var color = success ? '#1db954' : '#e74c3c';
    rect.setAttribute('stroke', color);
    rect.setAttribute('stroke-width', '3');
    setTimeout(function() {
      rect.setAttribute('stroke', '#2a2a4a');
      rect.setAttribute('stroke-width', '1.5');
    }, 1200);
  }

  function animateDot(targetIndex, label, success) {
    var path = document.getElementById('edge-' + targetIndex);
    if (!path) return;
    var totalLen = path.getTotalLength();

    var g = document.createElementNS(NS, 'g');

    var circle = document.createElementNS(NS, 'circle');
    circle.setAttribute('r', '7');
    circle.setAttribute('fill', '#1db954');
    circle.setAttribute('opacity', '0.95');

    var txt = document.createElementNS(NS, 'text');
    txt.textContent = label;
    txt.setAttribute('fill', '#fff');
    txt.setAttribute('font-size', '10');
    txt.setAttribute('font-family', 'sans-serif');
    txt.setAttribute('font-weight', '600');
    txt.setAttribute('dy', '-10');
    txt.setAttribute('text-anchor', 'middle');

    g.appendChild(circle);
    g.appendChild(txt);
    svg.appendChild(g);

    var startTime = null;

    function step(now) {
      if (!startTime) startTime = now;
      var t = Math.min((now - startTime) / DURATION, 1);
      var eased = t < 0.5 ? 2 * t * t : -1 + (4 - 2 * t) * t;
      var pt = path.getPointAtLength(eased * totalLen);
      circle.setAttribute('cx', pt.x);
      circle.setAttribute('cy', pt.y);
      txt.setAttribute('x', pt.x);
      txt.setAttribute('y', pt.y);
      if (t < 1) {
        requestAnimationFrame(step);
      } else {
        g.remove();
        flashTarget(targetIndex, success);
      }
    }

    requestAnimationFrame(step);
  }

  var es = new EventSource('/api/events');
  es.onmessage = function(e) {
    try {
      var ev = JSON.parse(e.data);
      if (ev.type === 'track_change') {
        applyTrackChange(ev);
        return;
      }
      if (ev.type !== 'scrobble') return;
      var idx = findTargetIndex(ev.target);
      if (idx < 0) return;
      animateDot(idx, ev.title, ev.success);
    } catch(err) {}
  };
  es.onerror = function() {
    setTimeout(function() {
      es.close();
      es = new EventSource('/api/events');
    }, 3000);
  };
})();
