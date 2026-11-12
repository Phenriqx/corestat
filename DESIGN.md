# Corestat interface design

## Status and direction

This document specifies the intended interface for Corestat, a Linux desktop system monitor built with Wails and Svelte. It is a design reference, not a description of the current implementation. All interface text is in English.

**North star:** Make changing system load immediately legible, with CPU and memory as the main story and the rest of the machine one selection away.

Today, the Svelte app renders RAM, CPU, disk, host, and process monitors sequentially on one page. The design below replaces that visual structure with a focused workspace. It does not authorize backend contract, dependency, or process-control changes; follow `AGENTS.md` for those decisions.

## References and choices

- **Provided image 3:** Primary reference for the dark desktop workspace, restrained green accent, left navigation, metric cards, and prominent chart. Borrow its spatial relationships, not its simulated-data controls, project links, or copy.
- **Provided images 1 and 2:** References for dashboard hierarchy: a concise summary row, one dominant graphic, and secondary detail. Their business content and blue/orange branding do not belong in Corestat.
- **[Charts Supply](https://charts.supply/):** Reference for clear charts with legible axes, legends, series, and tooltips. It is a visual reference, not a required Framer component or dependency.

The visual direction should feel like a precise Linux monitoring instrument: calm, dense enough to be useful, and focused on real readings. The distinctive moment is seeing the chart and large current reading change together as the machine's load changes.

## User task and navigation

The first view answers: **How busy are my CPU and memory right now?** The user then selects a resource in the sidebar to inspect its live chart and details. The selected resource is the main visual focus. Passive monitoring has no prominent action button. The Overview should fit in a normal desktop window without page scrolling; compact its sections rather than clipping content.

Use a persistent left sidebar in a normal desktop window with these destinations:

1. **Overview** — default landing view; CPU and memory visible together.
2. **CPU** — CPU activity and per-core detail.
3. **Memory** — RAM use and breakdown.
4. **Disks** — mount-point usage.
5. **Processes** — process list and hierarchy.

Put Corestat branding at the top. Mark the active destination with a green accent, distinct surface, and text label; selection cannot depend on color alone. Place host name and essential machine context in a quiet sidebar footer or compact Overview detail. Do not add the references' search, export, account, simulation, pricing, or project-navigation controls without a real product requirement.

The content region starts with a compact heading, a short description, and live or last-updated status. Summary readings come next, followed by the dominant chart and resource-specific details. Preserve the reading order **heading → readings → chart → details**.

### Main window structure

```text
┌───────────────────┬─────────────────────────────────────────────────────────┐
│ Corestat          │ Overview                                Live · updated … │
│                   │ CPU and memory at a glance                           │
│ Overview          │ ┌──────────────────┐ ┌──────────────────┐            │
│ CPU               │ │ CPU usage        │ │ Memory used      │            │
│ Memory            │ └──────────────────┘ └──────────────────┘            │
│ Disks             │ ┌──────────────────────────────────────┐ ┌─────────┐  │
│ Processes         │ │ Live CPU activity                    │ │ Memory  │  │
│                   │ │ time axis · value axis · values      │ │ summary │  │
│ Host context      │ └──────────────────────────────────────┘ └─────────┘  │
│                   │ Relevant details or a short process preview           │
└───────────────────┴─────────────────────────────────────────────────────────┘
```

Overview gives CPU the larger chart and memory a substantial adjacent summary or smaller trend. Selecting CPU, Memory, or Disks changes both the main chart and its supporting details. Processes is list led: scanning and inspecting processes is more useful there than forcing a chart into that view.

## View content

| View | Main visual | Supporting information |
| --- | --- | --- |
| Overview | Live CPU usage trend, with memory still visible alongside it | Current CPU %, used RAM %, and compact host context. |
| CPU | Large CPU usage time series | Show every logical core separately with its current usage and frequency where available, inspired by btop's compact per-core display. Show temperature as an overall sensor reading unless the backend can identify a reliable sensor-to-core mapping. |
| Memory | Large used-memory time series with a labeled percentage or capacity scale | Total, used, free, and cached memory. Never stack overlapping categories as though they sum to total. |
| Disks | Usage trend for the selected mount point | Mount-point list with used, free, total, and used %. A nearly flat trend is valid; do not exaggerate its scale. |
| Processes | Scan-friendly table or tree, initially sorted by CPU consumption descending | Name, PID, CPU, memory, threads, children, and an inspectable selected-process detail. Other sorting and filtering options can come later. |

Place process termination controls only in the selected process detail, visibly separate from passive metrics. Distinguish graceful termination from force kill by label and confirmation. Changes to the actual process action require consultation under `AGENTS.md`.

## Live charts

- Plot actual readings collected while the app runs. Never fill gaps with simulated values or invent historical trends.
- Keep the latest 60 seconds of samples in memory and label the visible span accurately. Closing and reopening the app starts a new chart history.
- Show the newest reading numerically beside the chart. On pointer hover, place the hovered reading's percentage above its point; keep a keyboard-accessible way to inspect exact values. Use a time axis and labeled units. Percentage axes span 0–100%; capacity axes identify their units and range.
- When multiple series appear, name them in a legend and distinguish them by line style or labels as well as color.
- Treat unavailable data differently from zero. On a polling failure after successful readings, retain the last valid chart for context, mark it stale, show the last update time, and offer retry.
- Keep updates smooth and legible. Reduce polling or rendering for hidden views when practical.

## Visual language

Use a near-black, slightly green-tinted foundation with restrained bright green for active navigation, key readings, and the primary data series. Avoid turning every panel into a neon accent. Suggested starting tokens, to refine against actual screenshots:

| Role | Starting value | Use |
| --- | --- | --- |
| Canvas | `#0D1110` | Main background |
| Sidebar | `#111715` | Navigation background |
| Raised surface | `#18201C` | Cards and chart panels |
| Border | `#2A3830` | Quiet separation |
| Primary text | `#F1F6F2` | Titles and readings |
| Secondary text | `#AABDB0` | Labels and context |
| Accent | `#B8F36B` | Active item, focus, primary chart line |
| Secondary series | `#6CBF95` | Another series when necessary |
| Error | `#FF8A7A` | Failures and destructive-action emphasis |

Use a code-oriented typeface such as **Source Code Pro** or **Roboto Mono**, with a local monospace fallback so the desktop UI remains usable offline. Use tabular figures for changing metrics and aligned columns. Large readings should be more prominent than labels; section headings should support the data instead of resembling marketing headlines. In the Memory Used summary card, show the used GiB value first and the percentage below it. Start around 14–16 px for normal desktop text and avoid tiny chart labels. Use moderate radii, thin borders, and spacing to group content; keep shadows subtle. Use one coherent icon style only where icons improve scanning.

Charts use fine, low-contrast grid lines, a clear accent line, restrained fills, and readable axis labels. Give the active series more visual weight than comparison series. Keep text, chart lines, selection, and focus legible at normal and enlarged display scales. Health and warning states must not rely on hue alone.

## Window size and input

This is a resizable Linux desktop window, not a mobile web page. At the smallest practical desktop width, compact the sidebar or turn it into a labeled navigation control, then stack chart and details. At medium and wide widths, keep the chart dominant and place details beside or below it as space permits. Avoid horizontal page scrolling; the process table may use its own controlled overflow when necessary. Check narrow, medium, and wide desktop widths and increased text scaling.

All destinations, chart selectors, process rows, and actions must work by keyboard and show visible focus. Use semantic controls and an appropriate table/tree pattern. Respect reduced-motion settings; live chart transitions should be quiet and should not obscure value changes.

## States

- **Initial load:** Reserve the layout and show that readings are being collected; no fake values.
- **Insufficient samples:** Show the current reading immediately while the chart grows from its first samples.
- **Live:** Show latest values, an accurate update indicator, and the rolling chart.
- **Unavailable metric:** Explain that a sensor or mount point could not be read while leaving other metrics usable.
- **Polling error or stale data:** Keep the last valid reading visibly dated and offer retry.
- **Empty result:** Explain an empty process or disk list instead of leaving a blank panel.
- **Process action:** Show the target and confirmation, then progress and outcome without losing the user's place.

## Implementation boundary

This is a design specification, not a demand to add a chart library or change the Go/Wails API now. Before implementation, reconcile available backend readings with each chart and consult the maintainer before significant backend, dependency, or process-action changes. Use `AGENTS.md` for commands, tests, and approvals.
