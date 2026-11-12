<script>
  export let samples = [];
  export let title = "Live usage";
  export let unit = "%";
  export let color = "var(--accent)";
  export let maximum = 100;
  export let compact = false;

  let width = 760;
  const left = 44;
  const right = 14;
  const top = 18;
  const bottom = 31;
  const clipId = `chart-${Math.random().toString(36).slice(2)}`;
  let activeIndex = -1;

  $: height = compact ? 210 : 280;
  $: plotWidth = Math.max(1, width - left - right);
  $: plotHeight = height - top - bottom;
  $: latestTime = samples.length ? samples[samples.length - 1].time : Date.now();
  $: visible = samples.filter((point) => point.time >= latestTime - 60000 && Number.isFinite(point.value));
  $: maxValue = Math.max(1, maximum);
  $: coordinates = visible.map((point) => ({
    x: left + ((point.time - (latestTime - 60000)) / 60000) * plotWidth,
    y: top + (1 - Math.min(maxValue, Math.max(0, point.value)) / maxValue) * plotHeight,
  }));
  $: line = coordinates.map((point) => `${point.x},${point.y}`).join(" ");
  $: area = coordinates.length > 1
    ? `${coordinates[0].x},${top + plotHeight} ${line} ${coordinates[coordinates.length - 1].x},${top + plotHeight}`
    : "";
  $: selected = activeIndex >= 0 && activeIndex < visible.length ? visible[activeIndex] : null;
  $: selectedPoint = activeIndex >= 0 && activeIndex < coordinates.length ? coordinates[activeIndex] : null;
  $: tooltipX = selectedPoint ? Math.min(width - right - 42, Math.max(left + 42, selectedPoint.x)) : 0;
  $: tooltipY = selectedPoint ? Math.max(top + 29, selectedPoint.y - 13) : 0;

  function move(event) {
    if (!coordinates.length) return;
    const rect = event.currentTarget.getBoundingClientRect();
    const x = ((event.clientX - rect.left) / rect.width) * width;
    let nearest = 0;
    for (let i = 1; i < coordinates.length; i += 1) {
      if (Math.abs(coordinates[i].x - x) < Math.abs(coordinates[nearest].x - x)) nearest = i;
    }
    activeIndex = nearest;
  }

  function displayValue(value) {
    return `${value.toFixed(1)}${unit === "%" ? "%" : ` ${unit}`}`;
  }
</script>

<div class="chart-wrap" class:compact bind:clientWidth={width}>
  {#if visible.length}
    <svg viewBox={`0 0 ${width} ${height}`} role="img" aria-label={`${title}. Latest reading ${displayValue(visible[visible.length - 1].value)}. Previous 60 seconds.`} on:mousemove={move} on:mouseleave={() => activeIndex = -1}>
      <defs>
        <clipPath id={clipId}><rect x={left} y={top} width={plotWidth} height={plotHeight} /></clipPath>
      </defs>
      {#each [0, 25, 50, 75, 100] as percent}
        {@const y = top + (1 - percent / 100) * plotHeight}
        <line class="grid-line" x1={left} x2={width - right} y1={y} y2={y} />
        <text class="axis-label" x={left - 10} y={y + 4} text-anchor="end">{unit === "%" ? `${percent}%` : `${(maxValue * percent / 100).toFixed(0)}`}</text>
      {/each}
      <line class="axis-line" x1={left} x2={width - right} y1={top + plotHeight} y2={top + plotHeight} />
      <text class="axis-label" x={left} y={height - 5}>60s ago</text>
      <text class="axis-label" x={left + plotWidth / 2} y={height - 5} text-anchor="middle">30s ago</text>
      <text class="axis-label" x={width - right} y={height - 5} text-anchor="end">Now</text>
      <g clip-path={`url(#${clipId})`}>
        {#if area}<polygon points={area} fill={color} opacity="0.09" />{/if}
        {#if coordinates.length > 1}<polyline points={line} fill="none" stroke={color} stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" />{/if}
        {#if coordinates.length === 1}<circle cx={coordinates[0].x} cy={coordinates[0].y} r="4" fill={color} />{/if}
        {#if selectedPoint}
          <line class="hover-line" x1={selectedPoint.x} x2={selectedPoint.x} y1={top} y2={top + plotHeight} />
          <circle cx={selectedPoint.x} cy={selectedPoint.y} r="5" fill={color} stroke="var(--surface)" stroke-width="2" />
        {/if}
      </g>
      {#if selected && selectedPoint}
        <g class="hover-tooltip" transform={`translate(${tooltipX}, ${tooltipY})`} aria-hidden="true">
          <rect x="-40" y="-29" width="80" height="25" rx="5" />
          <text x="0" y="-12" text-anchor="middle">{displayValue(selected.value)}</text>
        </g>
      {/if}
    </svg>
    <div class="chart-reading" aria-live="off">
      {#if selected}
        <span>{new Date(selected.time).toLocaleTimeString()}</span>
        <strong>{displayValue(selected.value)}</strong>
      {:else}
        <span>Latest reading</span>
        <strong>{displayValue(visible[visible.length - 1].value)}</strong>
      {/if}
    </div>
    <label class="chart-scrubber">Inspect chart readings
      <input type="range" min="0" max={visible.length - 1} value={activeIndex >= 0 ? activeIndex : visible.length - 1} on:input={(event) => activeIndex = Number(event.currentTarget.value)} aria-valuetext={selected ? `${new Date(selected.time).toLocaleTimeString()}, ${displayValue(selected.value)}` : displayValue(visible[visible.length - 1].value)} />
    </label>
  {:else}
    <div class="chart-empty">Collecting live readings…</div>
  {/if}
</div>

<style>
  .chart-wrap { min-width: 0; }
  svg { display: block; width: 100%; height: auto; max-height: 320px; overflow: visible; }
  .grid-line { stroke: var(--border); stroke-width: 1; stroke-dasharray: 3 6; }
  .axis-line { stroke: var(--border); stroke-width: 1; }
  .hover-line { stroke: var(--muted); stroke-width: 1; stroke-dasharray: 3 4; }
  .hover-tooltip { pointer-events: none; }
  .hover-tooltip rect { fill: var(--accent); }
  .hover-tooltip text { fill: var(--canvas); font-family: inherit; font-size: 12px; font-weight: 700; }
  .axis-label { fill: var(--muted); font-size: 12px; font-family: inherit; }
  .chart-reading { display: flex; justify-content: space-between; gap: 1rem; color: var(--muted); font-size: 0.82rem; padding: 0.3rem 0.1rem 0; }
  .chart-reading strong { color: var(--text); font-weight: 600; font-variant-numeric: tabular-nums; }
  .chart-scrubber { display: block; color: var(--muted); font-size: 0.78rem; margin-top: 0.75rem; }
  .chart-scrubber input { display: block; width: 100%; margin: 0.3rem 0 0; accent-color: var(--accent); }
  .chart-empty { min-height: 250px; display: grid; place-items: center; color: var(--muted); font-size: 0.9rem; }
  .chart-wrap.compact .chart-empty { min-height: 170px; }
</style>
