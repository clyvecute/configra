package main

// landingPageHTML is the portfolio-grade HTML landing page served at /.
// It is embedded directly so the binary has no static file dependencies.
const landingPageHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Configra — Reliable Configuration, Without the Complexity</title>
  <meta name="description" content="Configra is a cloud-native configuration management and feature flag service with immutable versioning, schema validation, and atomic rollbacks." />
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700;800&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet" />
  <style>
    *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
    :root {
      --bg:#0a0d14; --bg-card:#0f1320; --border:rgba(99,120,255,0.18);
      --accent:#6378ff; --accent2:#a78bfa; --green:#34d399;
      --text:#e2e8f0; --muted:#64748b; --radius:14px;
    }
    html { scroll-behavior: smooth; }
    body { background:var(--bg); color:var(--text); font-family:'Inter',system-ui,sans-serif; line-height:1.6; min-height:100vh; overflow-x:hidden; }
    body::before { content:''; position:fixed; inset:0; background-image:linear-gradient(rgba(99,120,255,0.04) 1px,transparent 1px),linear-gradient(90deg,rgba(99,120,255,0.04) 1px,transparent 1px); background-size:48px 48px; pointer-events:none; z-index:0; }
    .orb { position:fixed; border-radius:50%; filter:blur(120px); opacity:0.25; pointer-events:none; z-index:0; animation:drift 18s ease-in-out infinite alternate; }
    .orb-1 { width:600px; height:600px; background:#6378ff; top:-200px; left:-200px; }
    .orb-2 { width:500px; height:500px; background:#a78bfa; bottom:-200px; right:-150px; animation-delay:-9s; }
    @keyframes drift { from{transform:translate(0,0)} to{transform:translate(40px,30px)} }
    .container { max-width:1100px; margin:0 auto; padding:0 24px; position:relative; z-index:1; }
    nav { display:flex; align-items:center; justify-content:space-between; padding:20px 0; }
    .logo { font-size:1.35rem; font-weight:700; background:linear-gradient(135deg,var(--accent),var(--accent2)); -webkit-background-clip:text; -webkit-text-fill-color:transparent; background-clip:text; letter-spacing:-0.5px; }
    .nav-badge { background:rgba(99,120,255,0.15); border:1px solid var(--border); color:var(--accent2); font-size:0.75rem; font-weight:500; padding:4px 12px; border-radius:99px; font-family:'JetBrains Mono',monospace; }
    .hero { text-align:center; padding:90px 0 60px; }
    .hero-eyebrow { display:inline-flex; align-items:center; gap:8px; background:rgba(52,211,153,0.1); border:1px solid rgba(52,211,153,0.25); color:var(--green); font-size:0.8rem; font-weight:600; padding:5px 14px; border-radius:99px; letter-spacing:0.04em; text-transform:uppercase; margin-bottom:28px; }
    .hero-eyebrow::before { content:''; display:block; width:7px; height:7px; border-radius:50%; background:var(--green); animation:pulse 2s infinite; }
    @keyframes pulse { 0%,100%{opacity:1} 50%{opacity:0.4} }
    h1 { font-size:clamp(2.4rem,5vw,3.8rem); font-weight:800; line-height:1.12; letter-spacing:-1.5px; margin-bottom:20px; }
    h1 span { background:linear-gradient(135deg,var(--accent) 0%,var(--accent2) 100%); -webkit-background-clip:text; -webkit-text-fill-color:transparent; background-clip:text; }
    .hero-sub { font-size:1.15rem; color:var(--muted); max-width:580px; margin:0 auto 40px; }
    .hero-actions { display:flex; gap:14px; justify-content:center; flex-wrap:wrap; }
    .btn { display:inline-flex; align-items:center; gap:8px; padding:13px 28px; border-radius:10px; font-size:0.95rem; font-weight:600; text-decoration:none; transition:all 0.2s; border:none; cursor:pointer; }
    .btn-primary { background:linear-gradient(135deg,var(--accent),var(--accent2)); color:#fff; box-shadow:0 0 30px rgba(99,120,255,0.35); }
    .btn-primary:hover { transform:translateY(-2px); box-shadow:0 0 45px rgba(99,120,255,0.55); }
    .btn-ghost { background:rgba(255,255,255,0.05); color:var(--text); border:1px solid var(--border); }
    .btn-ghost:hover { background:rgba(255,255,255,0.1); transform:translateY(-2px); }
    .terminal-wrapper { max-width:680px; margin:64px auto 0; border-radius:var(--radius); border:1px solid var(--border); background:#080b12; box-shadow:0 24px 80px rgba(0,0,0,0.6),0 0 0 1px rgba(99,120,255,0.1); overflow:hidden; }
    .terminal-bar { display:flex; align-items:center; gap:7px; padding:12px 16px; background:rgba(255,255,255,0.04); border-bottom:1px solid var(--border); }
    .terminal-bar span { width:12px; height:12px; border-radius:50%; }
    .t-red{background:#ff5f57} .t-amber{background:#febc2e} .t-green{background:#28c840}
    .terminal-title { margin-left:8px; font-size:0.75rem; color:var(--muted); font-family:'JetBrains Mono',monospace; }
    .terminal-body { padding:22px 24px; font-family:'JetBrains Mono',monospace; font-size:0.82rem; line-height:2; }
    .t-prompt{color:var(--accent)} .t-cmd{color:var(--text)} .t-ok{color:var(--green)} .t-info{color:var(--muted)} .t-out{color:var(--accent2)}
    .t-cursor { display:inline-block; width:9px; height:1.1em; background:var(--accent); vertical-align:text-bottom; animation:blink 1s step-end infinite; }
    @keyframes blink { 0%,100%{opacity:1} 50%{opacity:0} }
    section { padding:80px 0; }
    .section-label { text-align:center; font-size:0.75rem; font-weight:600; letter-spacing:0.1em; text-transform:uppercase; color:var(--accent); margin-bottom:12px; }
    h2 { text-align:center; font-size:clamp(1.7rem,3vw,2.4rem); font-weight:700; letter-spacing:-0.8px; margin-bottom:12px; }
    .section-sub { text-align:center; color:var(--muted); margin-bottom:52px; font-size:1rem; }
    .features-grid { display:grid; grid-template-columns:repeat(auto-fit,minmax(300px,1fr)); gap:20px; }
    .card { background:var(--bg-card); border:1px solid var(--border); border-radius:var(--radius); padding:28px; transition:all 0.25s; position:relative; overflow:hidden; }
    .card::before { content:''; position:absolute; inset:0; background:radial-gradient(circle at top left,rgba(99,120,255,0.06),transparent 60%); pointer-events:none; }
    .card:hover { border-color:rgba(99,120,255,0.45); transform:translateY(-4px); box-shadow:0 16px 48px rgba(0,0,0,0.4); }
    .card-icon { width:44px; height:44px; border-radius:10px; display:flex; align-items:center; justify-content:center; font-size:1.3rem; margin-bottom:16px; }
    .card h3 { font-size:1rem; font-weight:600; margin-bottom:8px; }
    .card p { font-size:0.88rem; color:var(--muted); line-height:1.65; }
    .api-table-wrapper { border:1px solid var(--border); border-radius:var(--radius); overflow:hidden; }
    table { width:100%; border-collapse:collapse; }
    thead { background:rgba(99,120,255,0.08); }
    th { padding:14px 20px; text-align:left; font-size:0.78rem; font-weight:600; letter-spacing:0.06em; text-transform:uppercase; color:var(--muted); }
    td { padding:14px 20px; font-size:0.88rem; border-top:1px solid var(--border); }
    tr:hover td { background:rgba(255,255,255,0.02); }
    .badge { display:inline-block; font-family:'JetBrains Mono',monospace; font-size:0.72rem; font-weight:600; padding:3px 10px; border-radius:6px; }
    .badge-get { background:rgba(52,211,153,0.15); color:var(--green); }
    .badge-post { background:rgba(99,120,255,0.15); color:var(--accent2); }
    .endpoint { font-family:'JetBrains Mono',monospace; font-size:0.82rem; color:var(--text); }
    .endpoint-auth { font-size:0.75rem; color:var(--muted); }
    .stack-pills { display:flex; flex-wrap:wrap; gap:12px; justify-content:center; }
    .pill { background:var(--bg-card); border:1px solid var(--border); border-radius:99px; padding:8px 20px; font-size:0.85rem; font-weight:500; transition:all 0.2s; }
    .pill:hover { border-color:var(--accent); color:var(--accent2); }
    footer { border-top:1px solid var(--border); padding:32px 0; display:flex; align-items:center; justify-content:space-between; flex-wrap:wrap; gap:12px; font-size:0.82rem; color:var(--muted); }
    footer a { color:var(--muted); text-decoration:none; transition:color 0.2s; }
    footer a:hover { color:var(--accent2); }
  </style>
</head>
<body>
  <div class="orb orb-1"></div>
  <div class="orb orb-2"></div>
  <div class="container">
    <nav>
      <div class="logo">Configra</div>
      <span class="nav-badge">v1 API</span>
    </nav>
    <section class="hero">
      <div class="hero-eyebrow">Live &amp; Running</div>
      <h1>Reliable Configuration,<br /><span>Without the Complexity.</span></h1>
      <p class="hero-sub">Cloud-native config management with immutable versioning, strict schema validation, and zero-downtime rollbacks &mdash; built on Twelve-Factor principles.</p>
      <div class="hero-actions">
        <a href="/dashboard" class="btn btn-primary" id="hero-dashboard-btn">&#9889; Control Dashboard</a>
        <a href="/health" class="btn btn-ghost" id="hero-health-btn">&#10003; Health Check</a>
        <a href="https://github.com/clyvecute/configra" class="btn btn-ghost" id="hero-github-btn" target="_blank" rel="noopener">&#128279; View on GitHub</a>
      </div>
      <div class="terminal-wrapper">
        <div class="terminal-bar">
          <span class="t-red"></span><span class="t-amber"></span><span class="t-green"></span>
          <span class="terminal-title">configra cli</span>
        </div>
        <div class="terminal-body">
          <div><span class="t-prompt">$</span> <span class="t-cmd">configra validate -schema schema.json -config config.json</span></div>
          <div><span class="t-ok">&#10003; Configuration is VALID.</span></div>
          <div>&nbsp;</div>
          <div><span class="t-prompt">$</span> <span class="t-cmd">configra push -file config.json -project 42 -key feature_flags -env 1</span></div>
          <div><span class="t-info">&nbsp;&nbsp;Validating locally...</span></div>
          <div><span class="t-ok">&nbsp;&nbsp;Successfully pushed config to server!</span></div>
          <div>&nbsp;</div>
          <div><span class="t-prompt">$</span> <span class="t-cmd">configra fetch -project 42 -key feature_flags -env 1</span></div>
          <div><span class="t-out">&nbsp;&nbsp;{ "dark_mode": true, "max_upload_mb": 25, "tier": "pro" }</span></div>
          <div>&nbsp;</div>
          <div><span class="t-prompt">$</span> <span class="t-cursor"></span></div>
        </div>
      </div>
    </section>
    <section id="features">
      <div class="section-label">Why Configra</div>
      <h2>Everything configs need. Nothing they don&rsquo;t.</h2>
      <p class="section-sub">Stop worrying about bad deploys caused by misconfiguration.</p>
      <div class="features-grid">
        <div class="card">
          <div class="card-icon" style="background:rgba(99,120,255,0.15)">&#128274;</div>
          <h3>Immutable History</h3>
          <p>Every change creates a new version. See who changed what and when. No blind overwrites, ever.</p>
        </div>
        <div class="card">
          <div class="card-icon" style="background:rgba(167,139,250,0.15)">&#9989;</div>
          <h3>Schema Enforcement</h3>
          <p>Configs validated against strict typed schemas &mdash; strings, ints, enums, booleans &mdash; before acceptance.</p>
        </div>
        <div class="card">
          <div class="card-icon" style="background:rgba(52,211,153,0.15)">&#9889;</div>
          <h3>Zero-Downtime Rollbacks</h3>
          <p>Revert to any previous known-good version in one CLI command. History is preserved, not overwritten.</p>
        </div>
        <div class="card">
          <div class="card-icon" style="background:rgba(251,191,36,0.15)">&#128736;</div>
          <h3>API Key Auth</h3>
          <p>Projects isolated by API key. Middleware enforces project scope &mdash; no IDOR possible by design.</p>
        </div>
        <div class="card">
          <div class="card-icon" style="background:rgba(239,68,68,0.15)">&#9729;</div>
          <h3>Cloud Native</h3>
          <p>Stateless Go binary runs identically on Docker Compose locally and Google Cloud Run in production.</p>
        </div>
        <div class="card">
          <div class="card-icon" style="background:rgba(99,120,255,0.15)">&#128196;</div>
          <h3>Auto-Migration</h3>
          <p>The service manages its own database schema on startup. No manual migration scripts to run.</p>
        </div>
      </div>
    </section>
    <section id="api">
      <div class="section-label">API Reference</div>
      <h2>Clean, predictable endpoints.</h2>
      <p class="section-sub">All responses are JSON. Protected routes require <code style="color:var(--accent2);font-family:'JetBrains Mono',monospace">X-API-Key</code> header.</p>
      <div class="api-table-wrapper">
        <table>
          <thead><tr><th>Method</th><th>Endpoint</th><th>Description</th></tr></thead>
          <tbody>
            <tr><td><span class="badge badge-get">GET</span></td><td><span class="endpoint">/health</span></td><td>Service health check &mdash; no auth required</td></tr>
            <tr><td><span class="badge badge-post">POST</span></td><td><span class="endpoint">/v1/validate</span><br/><span class="endpoint-auth">No auth required</span></td><td>Dry-run schema validation of a config payload</td></tr>
            <tr><td><span class="badge badge-post">POST</span></td><td><span class="endpoint">/v1/configs</span><br/><span class="endpoint-auth">&#128274; API Key</span></td><td>Create or update a configuration version</td></tr>
            <tr><td><span class="badge badge-get">GET</span></td><td><span class="endpoint">/v1/configs?key=&amp;env_id=</span><br/><span class="endpoint-auth">&#128274; API Key</span></td><td>Fetch the latest config for a project/env/key</td></tr>
            <tr><td><span class="badge badge-post">POST</span></td><td><span class="endpoint">/v1/rollback</span><br/><span class="endpoint-auth">&#128274; API Key</span></td><td>Roll back to a specific previous version</td></tr>
          </tbody>
        </table>
      </div>
    </section>
    <section id="stack" style="padding-top:0">
      <div class="section-label">Technology</div>
      <h2>Built on a modern, maintainable stack.</h2>
      <p class="section-sub">No unnecessary dependencies. Just solid engineering choices.</p>
      <div class="stack-pills">
        <div class="pill">Go 1.21+</div><div class="pill">PostgreSQL 15</div><div class="pill">Docker</div>
        <div class="pill">Terraform</div><div class="pill">Google Cloud Run</div><div class="pill">GitHub Actions</div>
        <div class="pill">Clean Architecture</div><div class="pill">Twelve-Factor</div>
      </div>
    </section>
    <footer>
      <span>Configra &mdash; MIT License</span>
      <span><a href="/health" id="footer-health">Status</a> &middot; <a href="https://github.com/clyvecute/configra" target="_blank" rel="noopener" id="footer-github">GitHub</a></span>
    </footer>
  </div>
</body>
</html>`
