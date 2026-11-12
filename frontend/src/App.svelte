<script>
  import { onMount, tick } from "svelte";
  import TrendChart from "./TrendChart.svelte";
  import {
    GetCPU, GetRAM, GetDiskUsage, GetHostInfo, GetProcesses,
    SigTerminateProcess, SigKillProcess,
  } from "../wailsjs/go/main/App.js";

  const destinations = ["Overview", "CPU", "Memory", "Disks", "Processes"];
  const descriptions = {
    Overview: "CPU and memory at a glance",
    CPU: "Live processor activity and core readings",
    Memory: "Live memory use and capacity",
    Disks: "Space used across mounted volumes",
    Processes: "Inspect running processes",
  };
  const HISTORY_MS = 60000;

  let active = "Overview";
  let menuOpen = false;
  let cpu = null;
  let memory = null;
  let disks = {};
  let host = null;
  let processes = [];
  let selectedMount = "";
  let selectedPid = null;
  let processQuery = "";
  let cpuHistory = [];
  let memoryHistory = [];
  let diskHistory = {};
  let errors = {};
  let updated = {};
  let pendingAction = null;
  let actionBusy = false;
  let actionMessage = "";
  let actionError = "";
  let cancelActionButton;
  let priorFocus;
  const inFlight = {};

  $: cpuPercent = cpu?.CPUPercent?.length
    ? cpu.CPUPercent.reduce((sum, value) => sum + Number(value || 0), 0) / cpu.CPUPercent.length
    : null;
  $: memoryPercent = Number.isFinite(memory?.["Percent used"]) ? memory["Percent used"] : null;
  $: mountNames = Object.keys(disks).sort();
  $: if (mountNames.length && !mountNames.includes(selectedMount)) selectedMount = mountNames[0];
  $: selectedDisk = disks[selectedMount] || null;
  $: selectedProcess = processes.find((process) => process.PID === selectedPid) || null;
  $: visibleProcesses = processes
    .filter((process) => `${process.Name || ""} ${process.PID || ""}`.toLowerCase().includes(processQuery.toLowerCase().trim()))
    .sort((a, b) => (Number(b.CPUPercent) || 0) - (Number(a.CPUPercent) || 0) || a.PID - b.PID);
  $: currentResource = active === "Overview" ? "cpu" : active.toLowerCase();
  $: currentError = errors[currentResource];
  $: currentUpdated = updated[currentResource];

  function addReading(history, value, time = Date.now()) {
    if (!Number.isFinite(value)) return history;
    return [...history, { time, value }].filter((reading) => reading.time >= time - HISTORY_MS);
  }

  function message(error) {
    return error instanceof Error ? error.message : String(error);
  }

  async function fetchCpu() {
    if (inFlight.cpu) return;
    inFlight.cpu = true;
    try {
      const result = await GetCPU();
      if (!result) throw new Error("No CPU readings received.");
      cpu = result;
      const values = Array.isArray(result.CPUPercent) ? result.CPUPercent : [];
      const average = values.length ? values.reduce((sum, value) => sum + Number(value || 0), 0) / values.length : null;
      const time = Date.now();
      cpuHistory = addReading(cpuHistory, average, time);
      updated = { ...updated, cpu: time };
      errors = { ...errors, cpu: null };
    } catch (error) {
      errors = { ...errors, cpu: message(error) };
    } finally {
      inFlight.cpu = false;
    }
  }

  async function fetchMemory() {
    if (inFlight.memory) return;
    inFlight.memory = true;
    try {
      const result = await GetRAM();
      if (!result) throw new Error("No memory readings received.");
      memory = result;
      const time = Date.now();
      memoryHistory = addReading(memoryHistory, Number(result["Percent used"]), time);
      updated = { ...updated, memory: time };
      errors = { ...errors, memory: null };
    } catch (error) {
      errors = { ...errors, memory: message(error) };
    } finally {
      inFlight.memory = false;
    }
  }

  async function fetchDisks() {
    if (inFlight.disks) return;
    inFlight.disks = true;
    try {
      const result = await GetDiskUsage();
      if (!result) throw new Error("No disk readings received.");
      disks = result;
      const time = Date.now();
      const next = { ...diskHistory };
      for (const [mount, usage] of Object.entries(result)) {
        next[mount] = addReading(next[mount] || [], Number(usage.usedPercent), time);
      }
      diskHistory = next;
      updated = { ...updated, disks: time };
      errors = { ...errors, disks: null };
    } catch (error) {
      errors = { ...errors, disks: message(error) };
    } finally {
      inFlight.disks = false;
    }
  }

  async function fetchHost() {
    if (inFlight.host) return;
    inFlight.host = true;
    try {
      host = await GetHostInfo();
    } catch (error) {
      console.error("Host information unavailable", error);
    } finally {
      inFlight.host = false;
    }
  }

  async function fetchProcesses() {
    if (inFlight.processes) return;
    inFlight.processes = true;
    try {
      const result = await GetProcesses();
      processes = Array.isArray(result) ? result : [];
      updated = { ...updated, processes: Date.now() };
      errors = { ...errors, processes: null };
    } catch (error) {
      errors = { ...errors, processes: message(error) };
    } finally {
      inFlight.processes = false;
    }
  }

  function pollActive() {
    if (active === "Overview") {
      fetchCpu();
      fetchMemory();
    } else if (active === "CPU") fetchCpu();
    else if (active === "Memory") fetchMemory();
    else if (active === "Disks") fetchDisks();
    else if (active === "Processes") fetchProcesses();
  }

  function selectView(view) {
    active = view;
    menuOpen = false;
    actionMessage = "";
    pollActive();
  }

  function retry() {
    pollActive();
    if (active === "Overview") fetchHost();
  }

  function formatPercent(value) {
    return Number.isFinite(value) ? `${value.toFixed(1)}%` : "—";
  }

  function formatBytes(value) {
    if (!Number.isFinite(value)) return "—";
    return `${(value / (1024 ** 3)).toFixed(1)} GiB`;
  }

  function formatProcessBytes(value) {
    if (!Number.isFinite(value)) return "—";
    return value >= 1024 ** 3 ? formatBytes(value) : `${(value / (1024 ** 2)).toFixed(0)} MiB`;
  }

  function frequency() {
    const values = cpu?.CPUFrequency?.filter((value) => Number(value) > 0) || [];
    return values.length
      ? `${(values.reduce((sum, value) => sum + Number(value), 0) / values.length / 1000).toFixed(2)} GHz`
      : "Unavailable";
  }

  function coreFrequency(index) {
    const value = Number(cpu?.CPUFrequency?.[index]);
    return value > 0 ? `${(value / 1000).toFixed(2)} GHz` : "Frequency unavailable";
  }

  function temperature() {
    const values = Object.values(cpu?.CPUTemperature || {}).map(Number).filter(Number.isFinite);
    return values.length ? `${Math.max(...values).toFixed(0)} °C` : "Unavailable";
  }

  function lastUpdated(time) {
    return time ? new Date(time).toLocaleTimeString() : "Waiting for readings";
  }

  async function openAction(kind) {
    if (!selectedProcess) return;
    priorFocus = document.activeElement;
    pendingAction = { kind, pid: selectedProcess.PID, name: selectedProcess.Name };
    actionError = "";
    await tick();
    cancelActionButton?.focus();
  }

  function closeAction() {
    pendingAction = null;
    tick().then(() => priorFocus?.focus());
  }

  async function runAction() {
    if (!pendingAction || actionBusy) return;
    actionBusy = true;
    actionError = "";
    try {
      if (pendingAction.kind === "terminate") await SigTerminateProcess(pendingAction.pid);
      else await SigKillProcess(pendingAction.pid);
      actionMessage = `${pendingAction.kind === "terminate" ? "Termination" : "Force kill"} requested for ${pendingAction.name} (PID ${pendingAction.pid}).`;
      closeAction();
      await fetchProcesses();
    } catch (error) {
      actionError = message(error);
    } finally {
      actionBusy = false;
    }
  }

  function onKeydown(event) {
    if (event.key === "Escape" && pendingAction && !actionBusy) closeAction();
    if (event.key === "Tab" && pendingAction) {
      const buttons = [...document.querySelectorAll(".confirm-dialog button:not(:disabled)")];
      if (!buttons.length) return;
      const first = buttons[0];
      const last = buttons[buttons.length - 1];
      if (event.shiftKey && document.activeElement === first) {
        event.preventDefault();
        last.focus();
      } else if (!event.shiftKey && document.activeElement === last) {
        event.preventDefault();
        first.focus();
      }
    }
  }

  onMount(() => {
    pollActive();
    fetchHost();
    const interval = setInterval(pollActive, 3000);
    return () => clearInterval(interval);
  });
</script>

<svelte:window on:keydown={onKeydown} />

<div class="shell">
  <aside class:open={menuOpen} class="sidebar">
    <div class="brand"><span class="brand-mark" aria-hidden="true"><i></i><i></i><i></i><i></i></span><span>corestat<span class="brand-dot">.</span></span></div>
    <div class="brand-caption">SYSTEM MONITOR / LINUX</div>
    <nav aria-label="Main navigation">
      <div class="nav-label">MONITOR</div>
      {#each destinations as destination}
        <button type="button" class:active={active === destination} aria-current={active === destination ? "page" : undefined} on:click={() => selectView(destination)}>
          <span class="nav-symbol" aria-hidden="true">{destination === "Overview" ? "▦" : destination === "CPU" ? "◫" : destination === "Memory" ? "▥" : destination === "Disks" ? "▤" : "≡"}</span>
          <span>{destination}</span>
        </button>
      {/each}
    </nav>
    <div class="sidebar-footer">
      <span class="footer-dot" aria-hidden="true"></span>
      <span>{host?.MajorInfo?.hostname || "Local machine"}</span>
      <small>{host?.MajorInfo?.platform || "Linux"}</small>
    </div>
  </aside>

  <div class="workspace">
    <header class="topbar">
      <button class="menu-toggle" type="button" aria-label="Toggle navigation" aria-expanded={menuOpen} on:click={() => menuOpen = !menuOpen}>☰</button>
      <span class="breadcrumb">CORESTAT <span>/</span> {active.toUpperCase()}</span>
      <span class="topbar-context">LOCAL SYSTEM · LIVE MONITORING</span>
    </header>

    <main class="content" class:overview-content={active === "Overview"}>
      <div class="page-heading">
        <div><p class="eyebrow">SYSTEM INSIGHTS</p><h1>{active}</h1><p class="subtitle">{descriptions[active]}</p></div>
        <div class:stale={Boolean(currentError)} class="live-status"><span class="status-dot" aria-hidden="true"></span><span>{currentError ? "Data unavailable" : currentUpdated ? "Live readings" : "Connecting"}</span><small>{lastUpdated(currentUpdated)}</small></div>
      </div>

      {#if currentError}
        <div class="error-banner" role="alert"><div><strong>{currentUpdated ? "Readings are stale" : "Unable to load readings"}</strong><span>{currentError}</span></div><button type="button" on:click={retry}>Retry</button></div>
      {/if}
      {#if active === "Overview" && errors.memory}
        <div class="error-banner" role="alert"><div><strong>Memory readings are unavailable</strong><span>{errors.memory}</span></div><button type="button" on:click={retry}>Retry</button></div>
      {/if}

      {#if active === "Overview"}
        <section class="metric-grid" aria-label="Current system readings">
          <div class="metric-card"><span class="metric-label">CPU UTILIZATION</span><strong>{formatPercent(cpuPercent)}</strong><small>{cpu?.CPUCores ? `${cpu.CPUCores} logical cores` : "Waiting for CPU"}</small><span class="metric-accent"></span></div>
          <div class="metric-card"><span class="metric-label">MEMORY USED</span><strong>{formatBytes(memory?.["Used memory"])}</strong><small>{memory ? `${formatPercent(memoryPercent)} of ${formatBytes(memory["Total memory"])} total` : "Waiting for memory"}</small><span class="metric-accent secondary"></span></div>
        </section>
        <section class="overview-grid">
          <div class="panel chart-panel"><div class="panel-heading"><div><p class="eyebrow">LIVE METRIC</p><h2>CPU activity</h2><p>Average across logical cores · last 60 seconds</p></div><span class="legend"><i></i>CPU total</span></div><TrendChart samples={cpuHistory} title="CPU activity" compact /></div>
          <div class="panel memory-panel"><div class="panel-heading"><div><p class="eyebrow">ALSO IN VIEW</p><h2>Memory</h2><p>Current physical memory</p></div></div><div class="memory-figure">{formatPercent(memoryPercent)}</div><div class="memory-track"><span style={`width: ${Math.min(100, Math.max(0, memoryPercent || 0))}%`}></span></div><div class="memory-breakdown"><span>Used</span><strong>{formatBytes(memory?.["Used memory"])}</strong><span>Available</span><strong>{formatBytes(memory?.["Free memory"])}</strong><span>Total</span><strong>{formatBytes(memory?.["Total memory"])}</strong></div><button class="text-link" type="button" on:click={() => selectView("Memory")}>Explore memory <span aria-hidden="true">→</span></button></div>
        </section>
      {:else if active === "CPU"}
        <section class="metric-grid three" aria-label="CPU readings"><div class="metric-card"><span class="metric-label">TOTAL USAGE</span><strong>{formatPercent(cpuPercent)}</strong><small>Average of logical cores</small><span class="metric-accent"></span></div><div class="metric-card"><span class="metric-label">FREQUENCY</span><strong>{frequency()}</strong><small>Average available reading</small></div><div class="metric-card"><span class="metric-label">TEMPERATURE</span><strong>{temperature()}</strong><small>Highest available sensor</small></div></section>
        <section class="panel chart-panel"><div class="panel-heading"><div><p class="eyebrow">PROCESSOR</p><h2>CPU usage over time</h2><p>Average across logical cores · last 60 seconds</p></div><span class="legend"><i></i>CPU total</span></div><TrendChart samples={cpuHistory} title="CPU usage" /></section>
        <section class="panel detail-panel"><div class="panel-heading"><div><h2>Logical cores</h2><p>{cpu?.CPUInfo?.[0]?.modelName || "Per-core readings"}</p></div></div>{#if cpu?.CPUPercent?.length}<div class="core-grid">{#each cpu.CPUPercent as value, index}<div class="core-card"><div class="core-card-top"><span>CORE {index}</span><strong>{formatPercent(Number(value))}</strong></div><div class="core-track"><span style={`width: ${Math.min(100, Math.max(0, Number(value) || 0))}%`}></span></div><small>{coreFrequency(index)}</small></div>{/each}</div>{:else}<p class="empty-copy">Per-core readings are unavailable.</p>{/if}</section>
      {:else if active === "Memory"}
        <section class="metric-grid three" aria-label="Memory readings"><div class="metric-card"><span class="metric-label">MEMORY USED</span><strong>{formatBytes(memory?.["Used memory"])}</strong><small>{formatPercent(memoryPercent)} of total</small><span class="metric-accent"></span></div><div class="metric-card"><span class="metric-label">TOTAL CAPACITY</span><strong>{formatBytes(memory?.["Total memory"])}</strong><small>Physical memory</small></div><div class="metric-card"><span class="metric-label">FREE MEMORY</span><strong>{formatBytes(memory?.["Free memory"])}</strong><small>Current reading</small></div></section>
        <section class="panel chart-panel"><div class="panel-heading"><div><p class="eyebrow">PHYSICAL MEMORY</p><h2>Memory use over time</h2><p>Used share of total · last 60 seconds</p></div><span class="legend"><i></i>Used memory</span></div><TrendChart samples={memoryHistory} title="Memory usage" /></section>
        <section class="panel detail-panel"><div class="panel-heading"><div><h2>Memory breakdown</h2><p>Current values from the system</p></div></div><div class="data-list"><div><span>Used</span><strong>{formatBytes(memory?.["Used memory"])}</strong></div><div><span>Cached</span><strong>{formatBytes(memory?.["Cached memory"])}</strong></div><div><span>Free</span><strong>{formatBytes(memory?.["Free memory"])}</strong></div><div><span>Total</span><strong>{formatBytes(memory?.["Total memory"])}</strong></div></div></section>
      {:else if active === "Disks"}
        {#if mountNames.length}
          <div class="disk-layout"><section class="panel mount-panel"><div class="panel-heading"><div><p class="eyebrow">MOUNTED VOLUMES</p><h2>Disks</h2><p>Select a mount point</p></div></div><div class="mount-list">{#each mountNames as mount}<button type="button" class:selected={mount === selectedMount} aria-pressed={mount === selectedMount} on:click={() => selectedMount = mount}><span class="mount-name" title={mount}>{mount}</span><strong>{formatPercent(disks[mount]?.usedPercent)}</strong><span class="mount-track"><i style={`width: ${Math.min(100, Math.max(0, Number(disks[mount]?.usedPercent) || 0))}%`}></i></span></button>{/each}</div></section><section class="panel chart-panel"><div class="panel-heading"><div><p class="eyebrow">SELECTED VOLUME</p><h2 title={selectedMount}>{selectedMount}</h2><p>Used space · last 60 seconds</p></div><span class="legend"><i></i>Used</span></div><TrendChart samples={diskHistory[selectedMount] || []} title={`Disk usage for ${selectedMount}`} /><div class="data-list compact"><div><span>Used</span><strong>{formatBytes(selectedDisk?.used)}</strong></div><div><span>Free</span><strong>{formatBytes(selectedDisk?.free)}</strong></div><div><span>Total</span><strong>{formatBytes(selectedDisk?.total)}</strong></div></div></section></div>
        {:else if !errors.disks}<div class="panel empty-panel"><h2>{updated.disks ? "No mounted disks found" : "Looking for mounted disks"}</h2><p>{updated.disks ? "No readable volumes are currently available." : "Collecting disk readings…"}</p></div>{/if}
      {:else if active === "Processes"}
        {#if actionMessage}<div class="success-banner" role="status">{actionMessage}</div>{/if}
        <section class="panel process-panel"><div class="panel-heading process-heading"><div><p class="eyebrow">RUNNING NOW</p><h2>Processes</h2><p>{processes.length} processes returned by the system</p></div><label class="search-label">Search processes<input type="search" bind:value={processQuery} placeholder="Name or PID" /></label></div><div class="table-scroll"><table><thead><tr><th scope="col">Process</th><th scope="col">PID</th><th scope="col">CPU</th><th scope="col">Memory</th><th scope="col">Threads</th></tr></thead><tbody>{#each visibleProcesses as process (process.PID)}<tr class:selected={selectedPid === process.PID}><td><button class="process-name" type="button" on:click={() => selectedPid = process.PID}>{process.Name || "Unnamed process"}</button></td><td>{process.PID}</td><td>{formatPercent(Number(process.CPUPercent))}</td><td>{formatProcessBytes(Number(process.MemoryInfo?.rss))}</td><td>{process.Threads ?? "—"}</td></tr>{/each}</tbody></table></div>{#if !visibleProcesses.length}<div class="empty-table">{processQuery ? "No processes match your search." : updated.processes ? "No processes were returned." : "Loading processes…"}</div>{/if}</section>
        {#if selectedProcess}<section class="panel detail-panel process-detail"><div class="panel-heading"><div><p class="eyebrow">SELECTED PROCESS</p><h2>{selectedProcess.Name || "Unnamed process"}</h2><p>PID {selectedProcess.PID}{selectedProcess.HostUser ? ` · ${selectedProcess.HostUser}` : ""}</p></div><button class="quiet-button" type="button" on:click={() => selectedPid = null}>Close details</button></div><div class="data-list"><div><span>CPU</span><strong>{formatPercent(Number(selectedProcess.CPUPercent))}</strong></div><div><span>Memory</span><strong>{formatProcessBytes(Number(selectedProcess.MemoryInfo?.rss))}</strong></div><div><span>Threads</span><strong>{selectedProcess.Threads ?? "—"}</strong></div><div><span>Working directory</span><strong class="path-value" title={selectedProcess.Cwd || ""}>{selectedProcess.Cwd || "Unavailable"}</strong></div></div>{#if selectedProcess.Children?.length}<div class="children"><h3>Child processes</h3>{#each selectedProcess.Children as child}<div><span>{child.Name || "Unnamed process"}</span><span>PID {child.PID}</span></div>{/each}</div>{/if}<div class="process-actions"><button class="quiet-button" type="button" on:click={() => openAction("terminate")}>Terminate process</button><button class="danger-button" type="button" on:click={() => openAction("kill")}>Force kill…</button></div></section>{/if}
      {/if}
    </main>
  </div>
</div>

{#if pendingAction}
  <div class="dialog-scrim"><div class="confirm-dialog" role="dialog" aria-modal="true" aria-labelledby="confirm-title" aria-describedby="confirm-description"><p class="eyebrow">PROCESS ACTION</p><h2 id="confirm-title">{pendingAction.kind === "terminate" ? "Terminate process?" : "Force kill process?"}</h2><p id="confirm-description">{pendingAction.kind === "terminate" ? "Request a graceful termination" : "Immediately request a force kill"} for <strong>{pendingAction.name}</strong> (PID {pendingAction.pid}).</p>{#if actionError}<p class="dialog-error" role="alert">{actionError}</p>{/if}<div class="dialog-actions"><button bind:this={cancelActionButton} class="quiet-button" type="button" disabled={actionBusy} on:click={closeAction}>Cancel</button><button class:danger-button={pendingAction.kind === "kill"} class="confirm-button" type="button" disabled={actionBusy} on:click={runAction}>{actionBusy ? "Requesting…" : pendingAction.kind === "kill" ? "Force kill" : "Terminate"}</button></div></div></div>
{/if}
