package main

const dashboardPageHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Configra — Control Dashboard</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500;600&display=swap" rel="stylesheet">
  <style>
    :root {
      --bg-space: #F8F9FC;
      --card-bg: #FFFFFF;
      --border-color: #E2E8F0;
      --border-focus: #94A3B8;
      --text-main: #0F172A;
      --text-muted: #64748B;
      --text-light: #94A3B8;
      --accent-primary: #0F172A;
      --accent-hover: #1E293B;
      --accent-indigo: #6366F1;
      --accent-indigo-light: #EEF2FF;
      --success-bg: #ECFDF5;
      --success-text: #059669;
      --success-border: #A7F3D0;
      --error-bg: #FEF2F2;
      --error-text: #DC2626;
      --error-border: #FECACA;
      --shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.03);
      --shadow-card: 0 4px 12px -2px rgba(15, 23, 42, 0.04), 0 2px 4px -1px rgba(15, 23, 42, 0.02);
      --radius-sm: 6px;
      --radius-md: 10px;
      --radius-lg: 14px;
    }

    * {
      box-sizing: border-box;
      margin: 0;
      padding: 0;
    }

    body {
      font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
      background-color: var(--bg-space);
      background-image: 
        radial-gradient(circle at 15% 15%, rgba(99, 102, 241, 0.03) 0%, transparent 40%),
        radial-gradient(circle at 85% 85%, rgba(14, 165, 233, 0.03) 0%, transparent 40%);
      color: var(--text-main);
      min-height: 100vh;
      display: flex;
      flex-direction: column;
      -webkit-font-smoothing: antialiased;
    }

    /* Top Navigation Bar */
    header {
      background: rgba(255, 255, 255, 0.85);
      backdrop-filter: blur(12px);
      border-bottom: 1px solid var(--border-color);
      position: sticky;
      top: 0;
      z-index: 100;
    }

    .nav-container {
      max-width: 1200px;
      margin: 0 auto;
      padding: 0.85rem 1.5rem;
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 1.5rem;
    }

    .brand {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      text-decoration: none;
      color: var(--text-main);
    }

    .logo-badge {
      width: 32px;
      height: 32px;
      background: var(--accent-primary);
      color: white;
      border-radius: 8px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-weight: 700;
      font-size: 0.95rem;
      box-shadow: 0 2px 6px rgba(15, 23, 42, 0.15);
      transition: transform 0.2s cubic-bezier(0.16, 1, 0.3, 1);
    }

    .brand:hover .logo-badge {
      transform: scale(1.05);
    }

    .brand-title {
      font-weight: 700;
      font-size: 1.1rem;
      letter-spacing: -0.02em;
    }

    .brand-tag {
      font-size: 0.75rem;
      background: #F1F5F9;
      color: var(--text-muted);
      padding: 0.2rem 0.5rem;
      border-radius: 12px;
      font-weight: 500;
    }

    .header-controls {
      display: flex;
      align-items: center;
      gap: 1rem;
      flex: 1;
      max-width: 600px;
      justify-content: flex-end;
    }

    .api-key-box {
      display: flex;
      align-items: center;
      background: #F1F5F9;
      border: 1px solid var(--border-color);
      border-radius: var(--radius-sm);
      padding: 0.35rem 0.75rem;
      gap: 0.5rem;
      transition: border-color 0.15s ease, background 0.15s ease;
      width: 100%;
      max-width: 340px;
    }

    .api-key-box:focus-within {
      border-color: var(--border-focus);
      background: #FFFFFF;
      box-shadow: 0 0 0 3px rgba(148, 163, 184, 0.15);
    }

    .api-key-box svg {
      color: var(--text-muted);
      flex-shrink: 0;
    }

    .api-key-box input {
      border: none;
      background: transparent;
      outline: none;
      font-family: 'JetBrains Mono', monospace;
      font-size: 0.825rem;
      color: var(--text-main);
      width: 100%;
    }

    .nav-links {
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }

    .nav-link {
      font-size: 0.85rem;
      font-weight: 500;
      color: var(--text-muted);
      text-decoration: none;
      padding: 0.4rem 0.75rem;
      border-radius: var(--radius-sm);
      transition: all 0.15s ease;
    }

    .nav-link:hover {
      color: var(--text-main);
      background: rgba(0, 0, 0, 0.03);
    }

    /* Main Container */
    main {
      max-width: 1200px;
      margin: 2rem auto;
      padding: 0 1.5rem;
      width: 100%;
      flex: 1;
    }

    .dashboard-grid {
      display: grid;
      grid-template-columns: 240px 1fr;
      gap: 2rem;
      align-items: start;
    }

    @media (max-width: 860px) {
      .dashboard-grid {
        grid-template-columns: 1fr;
      }
    }

    /* Sidebar Navigation */
    .sidebar {
      background: var(--card-bg);
      border: 1px solid var(--border-color);
      border-radius: var(--radius-md);
      padding: 0.75rem;
      box-shadow: var(--shadow-card);
      position: sticky;
      top: 5rem;
    }

    .sidebar-heading {
      font-size: 0.7rem;
      font-weight: 700;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      color: var(--text-light);
      padding: 0.5rem 0.75rem 0.25rem 0.75rem;
    }

    .tab-btn {
      display: flex;
      align-items: center;
      gap: 0.65rem;
      width: 100%;
      padding: 0.65rem 0.85rem;
      border: none;
      background: transparent;
      color: var(--text-muted);
      font-family: inherit;
      font-size: 0.875rem;
      font-weight: 500;
      border-radius: var(--radius-sm);
      cursor: pointer;
      text-align: left;
      transition: all 0.15s ease;
      margin-bottom: 0.25rem;
    }

    .tab-btn:hover {
      background: #F1F5F9;
      color: var(--text-main);
    }

    .tab-btn.active {
      background: var(--accent-primary);
      color: #FFFFFF;
      font-weight: 600;
      box-shadow: var(--shadow-sm);
    }

    .tab-btn svg {
      width: 18px;
      height: 18px;
      flex-shrink: 0;
    }

    .env-selector {
      margin-top: 1rem;
      padding-top: 1rem;
      border-top: 1px solid var(--border-color);
    }

    .env-label {
      font-size: 0.75rem;
      font-weight: 600;
      color: var(--text-muted);
      margin-bottom: 0.5rem;
      display: block;
    }

    .select-input {
      width: 100%;
      padding: 0.5rem;
      border: 1px solid var(--border-color);
      border-radius: var(--radius-sm);
      background: #FFFFFF;
      font-family: inherit;
      font-size: 0.825rem;
      color: var(--text-main);
      outline: none;
      cursor: pointer;
      transition: border-color 0.15s ease;
    }

    .select-input:focus {
      border-color: var(--border-focus);
    }

    /* Content Panes */
    .content-area {
      min-width: 0;
    }

    .tab-pane {
      display: none;
      animation: fadeIn 0.2s cubic-bezier(0.16, 1, 0.3, 1);
    }

    .tab-pane.active {
      display: block;
    }

    @keyframes fadeIn {
      from { opacity: 0; transform: translateY(4px); }
      to { opacity: 1; transform: translateY(0); }
    }

    /* Card Components */
    .card {
      background: var(--card-bg);
      border: 1px solid var(--border-color);
      border-radius: var(--radius-lg);
      padding: 1.5rem;
      box-shadow: var(--shadow-card);
      margin-bottom: 1.5rem;
      transition: border-color 0.2s ease, box-shadow 0.2s ease;
    }

    .card-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 1.25rem;
      padding-bottom: 0.85rem;
      border-bottom: 1px solid #F1F5F9;
    }

    .card-title {
      font-size: 1.05rem;
      font-weight: 600;
      color: var(--text-main);
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }

    .card-subtitle {
      font-size: 0.825rem;
      color: var(--text-muted);
      margin-top: 0.25rem;
    }

    /* Form Controls */
    .form-group {
      margin-bottom: 1.25rem;
    }

    .form-label {
      display: block;
      font-size: 0.825rem;
      font-weight: 600;
      color: var(--text-main);
      margin-bottom: 0.4rem;
    }

    .form-hint {
      font-size: 0.75rem;
      color: var(--text-muted);
      margin-top: 0.25rem;
    }

    .text-input, .textarea-input {
      width: 100%;
      padding: 0.6rem 0.85rem;
      border: 1px solid var(--border-color);
      border-radius: var(--radius-sm);
      background: #FAFAFC;
      font-family: inherit;
      font-size: 0.875rem;
      color: var(--text-main);
      outline: none;
      transition: all 0.15s ease;
    }

    .text-input:focus, .textarea-input:focus {
      background: #FFFFFF;
      border-color: var(--text-main);
      box-shadow: 0 0 0 3px rgba(15, 23, 42, 0.06);
    }

    .textarea-input {
      font-family: 'JetBrains Mono', monospace;
      font-size: 0.825rem;
      min-height: 120px;
      resize: vertical;
      line-height: 1.5;
    }

    /* Action Buttons */
    .btn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      gap: 0.5rem;
      padding: 0.55rem 1.1rem;
      font-family: inherit;
      font-size: 0.85rem;
      font-weight: 600;
      border-radius: var(--radius-sm);
      cursor: pointer;
      border: 1px solid transparent;
      transition: all 0.15s cubic-bezier(0.16, 1, 0.3, 1);
    }

    .btn:active {
      transform: scale(0.98);
    }

    .btn-primary {
      background: var(--accent-primary);
      color: #FFFFFF;
    }

    .btn-primary:hover {
      background: var(--accent-hover);
      box-shadow: 0 2px 8px rgba(15, 23, 42, 0.15);
    }

    .btn-secondary {
      background: #FFFFFF;
      border-color: var(--border-color);
      color: var(--text-main);
    }

    .btn-secondary:hover {
      background: #F8FAFC;
      border-color: var(--border-focus);
    }

    .btn-sm {
      padding: 0.35rem 0.75rem;
      font-size: 0.775rem;
    }

    /* JSON Display / Code Block */
    .code-block {
      background: #0F172A;
      color: #F8FAFC;
      padding: 1rem;
      border-radius: var(--radius-sm);
      font-family: 'JetBrains Mono', monospace;
      font-size: 0.825rem;
      overflow-x: auto;
      line-height: 1.5;
      white-space: pre-wrap;
      word-break: break-word;
    }

    /* Output Alert Boxes */
    .status-alert {
      padding: 0.85rem 1rem;
      border-radius: var(--radius-sm);
      font-size: 0.85rem;
      display: flex;
      align-items: flex-start;
      gap: 0.75rem;
      margin-top: 1rem;
      animation: fadeIn 0.2s ease;
    }

    .status-alert.success {
      background: var(--success-bg);
      color: var(--success-text);
      border: 1px solid var(--success-border);
    }

    .status-alert.error {
      background: var(--error-bg);
      color: var(--error-text);
      border: 1px solid var(--error-border);
    }

    .status-alert svg {
      width: 18px;
      height: 18px;
      flex-shrink: 0;
      margin-top: 0.1rem;
    }

    /* Config Item Result Header */
    .result-meta {
      display: flex;
      gap: 1.5rem;
      margin-bottom: 0.75rem;
      font-size: 0.8rem;
    }

    .badge {
      display: inline-flex;
      align-items: center;
      padding: 0.2rem 0.55rem;
      border-radius: 12px;
      font-size: 0.75rem;
      font-weight: 600;
      text-transform: uppercase;
      letter-spacing: 0.03em;
    }

    .badge-indigo {
      background: var(--accent-indigo-light);
      color: var(--accent-indigo);
    }

    .badge-gray {
      background: #F1F5F9;
      color: var(--text-muted);
    }

    /* Quick Stats Bar */
    .stats-row {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
      gap: 1rem;
      margin-bottom: 1.5rem;
    }

    .stat-card {
      background: var(--card-bg);
      border: 1px solid var(--border-color);
      border-radius: var(--radius-md);
      padding: 1rem 1.25rem;
      box-shadow: var(--shadow-card);
    }

    .stat-label {
      font-size: 0.75rem;
      color: var(--text-muted);
      font-weight: 500;
    }

    .stat-value {
      font-size: 1.35rem;
      font-weight: 700;
      color: var(--text-main);
      margin-top: 0.2rem;
      letter-spacing: -0.02em;
    }

    footer {
      border-top: 1px solid var(--border-color);
      padding: 1.5rem;
      text-align: center;
      font-size: 0.8rem;
      color: var(--text-muted);
      margin-top: auto;
      background: #FFFFFF;
    }
  </style>
</head>
<body>

  <!-- Top Navbar -->
  <header>
    <div class="nav-container">
      <a href="/" class="brand">
        <div class="logo-badge">C</div>
        <div>
          <span class="brand-title">Configra</span>
          <span class="brand-tag">v1.0</span>
        </div>
      </a>

      <div class="header-controls">
        <div class="api-key-box">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
          <input type="password" id="globalApiKey" placeholder="X-API-Key (e.g. test_key_123)" value="test_key_123">
        </div>
        <div class="nav-links">
          <a href="/" class="nav-link">Docs & Landing</a>
          <a href="/health" target="_blank" class="nav-link">Health API</a>
        </div>
      </div>
    </div>
  </header>

  <!-- Main Content -->
  <main>
    <!-- Stats Banner -->
    <div class="stats-row">
      <div class="stat-card">
        <div class="stat-label">System Status</div>
        <div class="stat-value" style="color: #059669; font-size: 1.1rem; display:flex; align-items:center; gap:0.4rem; margin-top:0.3rem;">
          <span style="width:8px; height:8px; border-radius:50%; background:#059669; display:inline-block;"></span>
          Operational
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-label">Selected Environment</div>
        <div class="stat-value" id="statEnvDisplay" style="font-size: 1.1rem; margin-top:0.3rem;">Development (1)</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">API Scope</div>
        <div class="stat-value" style="font-size: 1.1rem; margin-top:0.3rem;">Project Auth</div>
      </div>
    </div>

    <div class="dashboard-grid">
      <!-- Sidebar Nav -->
      <aside class="sidebar">
        <div class="sidebar-heading">Navigation</div>
        <button class="tab-btn active" onclick="switchTab('explorerTab')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path></svg>
          Config Manager
        </button>
        <button class="tab-btn" onclick="switchTab('rollbackTab')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="1 4 1 10 7 10"></polyline><path d="M3.51 15a9 9 0 1 0 2.13-9.36L1 10"></path></svg>
          Rollback Engine
        </button>
        <button class="tab-btn" onclick="switchTab('validatorTab')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><polyline points="22 4 12 14.01 9 11.01"></polyline></svg>
          Schema Validator
        </button>
        <button class="tab-btn" onclick="switchTab('fetchTab')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          External Fetch
        </button>

        <div class="env-selector">
          <label class="env-label">Environment Scope</label>
          <select class="select-input" id="globalEnvId" onchange="updateEnvStat()">
            <option value="1" selected>Env 1 — Development</option>
            <option value="2">Env 2 — Staging</option>
            <option value="3">Env 3 — Production</option>
          </select>
        </div>
      </aside>

      <!-- Content Area -->
      <section class="content-area">
        
        <!-- Tab 1: Config Manager -->
        <div id="explorerTab" class="tab-pane active">
          <!-- Fetch Existing Config -->
          <div class="card">
            <div class="card-header">
              <div>
                <h2 class="card-title">
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
                  Get Active Configuration
                </h2>
                <p class="card-subtitle">Fetch latest dynamic configuration by key and environment</p>
              </div>
            </div>

            <div style="display: flex; gap: 0.75rem; align-items: flex-end;">
              <div style="flex: 1;">
                <label class="form-label">Config Key Name</label>
                <input type="text" id="fetchKeyInput" class="text-input" value="feature_flags" placeholder="e.g. feature_flags or rate_limits">
              </div>
              <button class="btn btn-primary" onclick="fetchConfig()">
                Get Config
              </button>
            </div>

            <div class="form-group" style="margin-top:1rem;">
              <label class="form-label">Version history</label>
              <button class="btn btn-secondary btn-sm" onclick="loadHistory()">Browse versions</button>
              <select id="historyVersion" class="select-input" onchange="showVersionDiff()" style="margin-top:.5rem"><option value="">Select a version</option></select>
              <button class="btn btn-secondary btn-sm" onclick="showVersionDiff()">Compare with current</button>
              <button class="btn btn-secondary btn-sm" onclick="rollbackSelected()">Rollback selected version</button>
              <pre class="code-block" id="historyResult" style="margin-top:.5rem;display:none"></pre>
            </div>

            <div id="fetchResultArea" style="margin-top: 1.25rem; display: none;">
              <div class="result-meta">
                <span class="badge badge-indigo" id="fetchVerBadge">Version: --</span>
                <span class="badge badge-gray" id="fetchKeyBadge">Key: --</span>
                <span class="badge badge-gray" id="fetchTimeBadge">Updated: --</span>
              </div>
              <pre class="code-block" id="fetchResultCode"></pre>
            </div>
          </div>

          <!-- Create / Push Config -->
          <div class="card">
            <div class="card-header">
              <div>
                <h2 class="card-title">
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
                  Push New Version
                </h2>
                <p class="card-subtitle">Publish new key-value data with runtime schema constraints</p>
              </div>
            </div>

            <div class="form-group">
              <label class="form-label">Key Name</label>
              <input type="text" id="pushKey" class="text-input" value="feature_flags" placeholder="e.g. feature_flags">
            </div>

            <div class="form-group">
              <label class="form-label">Data Payload (JSON)</label>
              <textarea id="pushData" class="textarea-input" placeholder='{\n  "dark_mode": true,\n  "max_items": 50\n}'>{"dark_mode": true, "max_items": 50}</textarea>
            </div>

            <div class="form-group">
              <label class="form-label">Schema Definitions (JSON)</label>
              <textarea id="pushSchema" class="textarea-input" placeholder='{\n  "dark_mode": "bool",\n  "max_items": "int"\n}'>{"dark_mode": "bool", "max_items": "int"}</textarea>
            </div>

            <div style="display:flex; justify-content:space-between; align-items:center;">
              <button class="btn btn-secondary btn-sm" onclick="formatPushJson()">Format JSON</button>
              <button class="btn btn-primary" onclick="pushConfig()">Push Configuration</button>
            </div>

            <div id="pushStatus"></div>
          </div>
        </div>

        <!-- Tab 2: Rollback Engine -->
        <div id="rollbackTab" class="tab-pane">
          <div class="card">
            <div class="card-header">
              <div>
                <h2 class="card-title">
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="1 4 1 10 7 10"></polyline><path d="M3.51 15a9 9 0 1 0 2.13-9.36L1 10"></path></svg>
                  Immutable Rollback Engine
                </h2>
                <p class="card-subtitle">Instantly restore any past version. Configra creates a new version preserving audit history.</p>
              </div>
            </div>

            <div class="form-group">
              <label class="form-label">Config Key</label>
              <input type="text" id="rollbackKey" class="text-input" value="feature_flags" placeholder="e.g. feature_flags">
            </div>

            <div class="form-group">
              <label class="form-label">Target Version Number</label>
              <input type="number" id="rollbackVersion" class="text-input" value="1" min="1" placeholder="e.g. 1">
              <p class="form-hint">Specifying version 1 will fetch version 1's content and publish it as the newest version.</p>
            </div>

            <button class="btn btn-primary" onclick="executeRollback()">Execute Rollback</button>

            <div id="rollbackStatus"></div>
          </div>
        </div>

        <!-- Tab 3: Schema Validator -->
        <div id="validatorTab" class="tab-pane">
          <div class="card">
            <div class="card-header">
              <div>
                <h2 class="card-title">
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><polyline points="22 4 12 14.01 9 11.01"></polyline></svg>
                  Schema Validation Playground
                </h2>
                <p class="card-subtitle">Test strict type checking without saving to the database</p>
              </div>
            </div>

            <div class="form-group">
              <label class="form-label">Test Config (JSON)</label>
              <textarea id="valConfig" class="textarea-input">{\n  "service_name": "configra",\n  "port": 8080,\n  "enabled": true\n}</textarea>
            </div>

            <div class="form-group">
              <label class="form-label">Schema Types (string, int, bool, float, json)</label>
              <textarea id="valSchema" class="textarea-input">{\n  "service_name": "string",\n  "port": "int",\n  "enabled": "bool"\n}</textarea>
            </div>

            <button class="btn btn-primary" onclick="validateSchema()">Run Type Checker</button>

            <div id="valStatus"></div>
          </div>
        </div>

        <!-- Tab 4: External Fetch Proxy -->
        <div id="fetchTab" class="tab-pane">
          <div class="card">
            <div class="card-header">
              <div>
                <h2 class="card-title">
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
                  External Fetch Source
                </h2>
                <p class="card-subtitle">Fetch remote configuration payloads safely via Configra proxy</p>
              </div>
            </div>

            <div class="form-group">
              <label class="form-label">Remote JSON Endpoint URL</label>
              <input type="text" id="extUrlInput" class="text-input" value="https://httpbin.org/json" placeholder="https://api.example.com/config.json">
            </div>

            <button class="btn btn-primary" onclick="fetchExternalSource()">Proxy Remote Fetch</button>

            <div id="extStatus"></div>
          </div>
        </div>

      </section>
    </div>
  </main>

  <footer>
    Configra Light Minimalist Control Center • Built for fast dynamic configuration management
  </footer>

  <script>
    // Tab switching logic
    function switchTab(tabId) {
      document.querySelectorAll('.tab-btn').forEach(btn => btn.classList.remove('active'));
      document.querySelectorAll('.tab-pane').forEach(pane => pane.classList.remove('active'));
      
      event.currentTarget.classList.add('active');
      document.getElementById(tabId).classList.add('active');
    }

    function updateEnvStat() {
      const select = document.getElementById('globalEnvId');
      const text = select.options[select.selectedIndex].text;
      document.getElementById('statEnvDisplay').innerText = text;
    }

    function getApiKey() {
      return document.getElementById('globalApiKey').value.trim();
    }

    function getEnvId() {
      return parseInt(document.getElementById('globalEnvId').value, 10);
    }

    // 1. Fetch Config
    async function fetchConfig() {
      const key = document.getElementById('fetchKeyInput').value.trim();
      const apiKey = getApiKey();
      const envId = getEnvId();
      const resArea = document.getElementById('fetchResultArea');
      const codeElem = document.getElementById('fetchResultCode');

      if (!key) return alert("Please specify a key name.");
      if (!apiKey) return alert("Please enter your X-API-Key in the top navigation bar.");

      codeElem.innerText = "Loading...";
      resArea.style.display = "block";

      try {
        const resp = await fetch('/v1/configs/' + encodeURIComponent(key) + '?env_id=' + envId, {
          headers: { 'X-API-Key': apiKey }
        });
        const data = await resp.json();

        if (!resp.ok) {
          codeElem.innerText = JSON.stringify(data, null, 2);
          document.getElementById('fetchVerBadge').innerText = "Status: " + resp.status;
          document.getElementById('fetchKeyBadge').innerText = "Key: " + key;
          document.getElementById('fetchTimeBadge').innerText = "Error";
          return;
        }

        document.getElementById('fetchVerBadge').innerText = "Version " + data.version;
        document.getElementById('fetchKeyBadge').innerText = "Key: " + data.key;
        document.getElementById('fetchTimeBadge').innerText = data.updated_at ? new Date(data.updated_at).toLocaleTimeString() : "Just now";
        codeElem.innerText = JSON.stringify(data, null, 2);
      } catch (err) {
        codeElem.innerText = "Error: " + err.message;
      }
    }

    async function loadHistory() {
      const key=document.getElementById('fetchKeyInput').value.trim(), apiKey=getApiKey(), envId=getEnvId();
      const select=document.getElementById('historyVersion'), result=document.getElementById('historyResult');
      const resp=await fetch('/v1/configs/'+encodeURIComponent(key)+'/versions?env_id='+envId,{headers:{'X-API-Key':apiKey}}); const versions=await resp.json();
      if(!resp.ok){result.style.display='block';result.innerText=versions.error||'Unable to load history';return;}
      select.innerHTML='<option value="">Select a version</option>'+versions.map(v=>'<option value="'+v.version+'">v'+v.version+' · '+new Date(v.created_at).toLocaleString()+' · author '+(v.author_id||'unknown')+'</option>').join('');
      result.style.display='block'; result.innerText=versions.length+' versions loaded. Select one and compare or roll back.';
    }
    async function showVersionDiff() {
      const key=document.getElementById('fetchKeyInput').value.trim(), selected=document.getElementById('historyVersion').value, result=document.getElementById('historyResult'); if(!selected)return;
      const apiKey=getApiKey(), envId=getEnvId();
      const latestResp=await fetch('/v1/configs/'+encodeURIComponent(key)+'?env_id='+envId,{headers:{'X-API-Key':apiKey}}); const latest=await latestResp.json();
      const resp=await fetch('/v1/configs/'+encodeURIComponent(key)+'/diff?env_id='+envId+'&from='+selected+'&to='+latest.version,{headers:{'X-API-Key':apiKey}}); const diff=await resp.json(); result.style.display='block'; result.innerText=JSON.stringify(diff,null,2);
    }
    async function rollbackSelected() {
      const key=document.getElementById('fetchKeyInput').value.trim(), version=document.getElementById('historyVersion').value; if(!version)return alert('Select a version first.');
      if(!confirm('Roll back '+key+' to version '+version+'? This publishes a new version with that content.'))return;
      const resp=await fetch('/v1/configs/'+encodeURIComponent(key)+'/rollback',{method:'POST',headers:{'Content-Type':'application/json','X-API-Key':getApiKey()},body:JSON.stringify({env_id:getEnvId(),target_version:Number(version)})}); const body=await resp.json();
      document.getElementById('historyResult').style.display='block'; document.getElementById('historyResult').innerText=JSON.stringify(body,null,2);
    }

    // 2. Format JSON helper
    function formatPushJson() {
      try {
        const d = JSON.parse(document.getElementById('pushData').value);
        document.getElementById('pushData').value = JSON.stringify(d, null, 2);
        const s = JSON.parse(document.getElementById('pushSchema').value);
        document.getElementById('pushSchema').value = JSON.stringify(s, null, 2);
      } catch(e) {
        alert("Invalid JSON format: " + e.message);
      }
    }

    // 3. Push Config
    async function pushConfig() {
      const key = document.getElementById('pushKey').value.trim();
      const apiKey = getApiKey();
      const envId = getEnvId();
      const statusContainer = document.getElementById('pushStatus');

      try {
        const data = JSON.parse(document.getElementById('pushData').value);
        const schema = JSON.parse(document.getElementById('pushSchema').value);

        statusContainer.innerHTML = '<div class="status-alert">Publishing new version...</div>';

        const resp = await fetch('/v1/configs', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'X-API-Key': apiKey
          },
          body: JSON.stringify({
            env_id: envId,
            key: key,
            data: data,
            schema: schema
          })
        });

        const resData = await resp.json();

        if (!resp.ok) {
          statusContainer.innerHTML = '<div class="status-alert error"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><line x1="15" y1="9" x2="9" y2="15"></line><line x1="9" y1="9" x2="15" y2="15"></line></svg><div><strong>Error:</strong> ' + (resData.error || 'Failed to push config') + '</div></div>';
        } else {
          statusContainer.innerHTML = '<div class="status-alert success"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="20 6 9 17 4 12"></polyline></svg><div><strong>Success!</strong> Created Version ' + resData.version + ' for key "' + resData.key + '"</div></div>';
          document.getElementById('fetchKeyInput').value = key;
        }
      } catch (err) {
        statusContainer.innerHTML = '<div class="status-alert error"><div><strong>JSON Error:</strong> ' + err.message + '</div></div>';
      }
    }

    // 4. Rollback
    async function executeRollback() {
      const key = document.getElementById('rollbackKey').value.trim();
      const ver = parseInt(document.getElementById('rollbackVersion').value, 10);
      const apiKey = getApiKey();
      const envId = getEnvId();
      const container = document.getElementById('rollbackStatus');

      if (!key || !ver) return alert("Please specify key and target version.");

      container.innerHTML = '<div class="status-alert">Executing rollback...</div>';

      try {
        const resp = await fetch('/v1/rollback', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'X-API-Key': apiKey
          },
          body: JSON.stringify({
            env_id: envId,
            key: key,
            target_version: ver
          })
        });

        const resData = await resp.json();

        if (!resp.ok) {
          container.innerHTML = '<div class="status-alert error"><div><strong>Rollback Failed:</strong> ' + (resData.error || 'Failed') + '</div></div>';
        } else {
          container.innerHTML = '<div class="status-alert success"><div><strong>Rollback Successful!</strong> Restored version ' + ver + ' contents into new version <strong>' + resData.version + '</strong>.</div></div>';
        }
      } catch (err) {
        container.innerHTML = '<div class="status-alert error"><div><strong>Error:</strong> ' + err.message + '</div></div>';
      }
    }

    // 5. Validate Schema
    async function validateSchema() {
      const container = document.getElementById('valStatus');
      try {
        const cfg = JSON.parse(document.getElementById('valConfig').value);
        const sch = JSON.parse(document.getElementById('valSchema').value);

        container.innerHTML = '<div class="status-alert">Validating schema...</div>';

        const resp = await fetch('/v1/validate', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ config: cfg, schema: sch })
        });

        const resData = await resp.json();

        if (!resp.ok) {
          container.innerHTML = '<div class="status-alert error"><div><strong>Validation Failed:</strong> ' + (resData.error || 'Invalid schema') + '</div></div>';
        } else {
          container.innerHTML = '<div class="status-alert success"><div><strong>Validation Passed:</strong> Configuration strictly matches schema definitions.</div></div>';
        }
      } catch (err) {
        container.innerHTML = '<div class="status-alert error"><div><strong>JSON Error:</strong> ' + err.message + '</div></div>';
      }
    }

    // 6. Proxy Fetch Source
    async function fetchExternalSource() {
      const url = document.getElementById('extUrlInput').value.trim();
      const apiKey = getApiKey();
      const container = document.getElementById('extStatus');

      if (!url) return alert("Please specify a URL.");

      container.innerHTML = '<div class="status-alert">Fetching remote payload...</div>';

      try {
        const resp = await fetch('/fetch', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'X-API-Key': apiKey
          },
          body: JSON.stringify({ url: url })
        });

        const resData = await resp.json();

        if (!resp.ok) {
          container.innerHTML = '<div class="status-alert error"><div><strong>Fetch Failed:</strong> ' + (resData.error || 'Failed') + '</div></div>';
        } else {
          container.innerHTML = '<div class="status-alert success"><div><strong>Fetched Payload Successfully:</strong></div></div><pre class="code-block" style="margin-top:0.5rem">' + JSON.stringify(resData, null, 2) + '</pre>';
        }
      } catch (err) {
        container.innerHTML = '<div class="status-alert error"><div><strong>Error:</strong> ' + err.message + '</div></div>';
      }
    }
  </script>
</body>
</html>`
