// Machine mode dashboard: the left panel lists every machine grouped by its current mode
// (GET /api/v1/health/machines), with a count per mode that filters the list. Selecting a machine
// reads GET /api/v1/health/machines/{id}/modes and /cycles (all public, no session) and draws its
// mode periods as a timeline, one row per mode, with a "Cycles" row of its part program cycles,
// then the cycles (summary, cycle times, list) and the list of periods.
(function () {
  "use strict";

  const API_PATH = "/api/v1/health/machines/";
  const MACHINES_PATH = "/api/v1/health/machines";
  const CYCLES_SUFFIX = "/cycles";
  const CYCLE_HISTORY_HOURS = 24; // dto.MachineCycleHistoryHours
  const CYCLE_TABLE_LIMIT = 200; // newest cycles listed in the table
  const CYCLE_BARS_SHOWN = 40; // newest cycles the cycle time chart opens on
  const CYCLES_ROW = "Cycles";
  const LIST_REFRESH_MS = 30 * 1000;
  const NO_MODE = "none"; // group key of machines that have no mode yet
  const ALL_MODES_ROW = "All modes";
  const SECOND_MS = 1000;
  const MINUTE_MS = 60 * SECOND_MS;
  const HOUR_MS = 60 * MINUTE_MS;
  const DAY_MS = 24 * HOUR_MS;
  const ZOOM_STEP = 2;
  // Zoom level buttons: each shows the chosen span of time on the chart.
  const ZOOM_LEVELS = [
    { label: "30s", title: "30 seconds", spanMs: 30 * SECOND_MS },
    { label: "1m", title: "1 minute", spanMs: MINUTE_MS },
    { label: "2m", title: "2 minutes", spanMs: 2 * MINUTE_MS },
    { label: "5m", title: "5 minutes", spanMs: 5 * MINUTE_MS },
    { label: "10m", title: "10 minutes", spanMs: 10 * MINUTE_MS },
    { label: "30m", title: "30 minutes", spanMs: 30 * MINUTE_MS },
    { label: "1h", title: "1 hour", spanMs: HOUR_MS },
    { label: "3h", title: "3 hours", spanMs: 3 * HOUR_MS },
    { label: "8h", title: "8 hours", spanMs: 8 * HOUR_MS },
    { label: "24h", title: "24 hours", spanMs: DAY_MS },
    { label: "3d", title: "3 days", spanMs: 3 * DAY_MS },
    { label: "10d", title: "10 days", spanMs: 10 * DAY_MS },
  ];
  const DEFAULT_ZOOM_MS = 5 * MINUTE_MS; // a newly selected machine shows its latest 5 minutes (the "5m" level)
  const MIN_SPAN_MS = ZOOM_LEVELS[0].spanMs;
  const ZOOM_MATCH_TOLERANCE = 0.01; // a view within 1 % of a level's span highlights that level

  // machine_modes ids. Rows always shown come first, as on the dashboard; the other modes get
  // a row only when the machine was in them.
  const MODES = {
    1: { name: "Setup", color: "#3b82f6" },
    2: { name: "Loading / Unloading Time", color: "#1f8fa3" },
    3: { name: "In Production", color: "#2e9d4a" },
    4: { name: "Hold", color: "#6d28d9" },
    5: { name: "Idle", color: "#b7791f" },
    6: { name: "Downtime", color: "#c0262d" },
    7: { name: "Planned Downtime", color: "#a3921f" },
    8: { name: "Unfinished Cycle", color: "#db2777" },
    9: { name: "Offline", color: "#8a94a6" },
    10: { name: "Resume Cycle", color: "#0d9488" },
    11: { name: "End Cycle", color: "#64748b" },
    12: { name: "Rapid Mode", color: "#84cc16" },
    13: { name: "Leave Planned Downtime", color: "#ca8a04" },
  };
  const ALWAYS_SHOWN_MODES = [2, 3, 4, 5, 6, 9, 7];
  // Order of the mode groups in the machine list. The tiles of DASHBOARD_TILE_MODES always show;
  // other modes get a tile only while a machine is in them.
  const DASHBOARD_MODE_ORDER = [3, 4, 5, 6, 2, 1, 7, 8, 10, 11, 12, 13, 9];
  const DASHBOARD_TILE_MODES = [3, 4, 5, 6, 9];
  const UNKNOWN_MODE_COLOR = "#9ca3af";

  // machine_cycles.status and end_reason in words an operator understands.
  const CYCLE_STATUSES = {
    running: { name: "Running", color: "#f28c28" },
    completed: { name: "Completed", color: "#2e9d4a" },
    unfinished: { name: "Unfinished", color: "#c0262d" },
  };
  const CYCLE_END_REASONS = {
    part_count: "Part counted",
    program_loop: "Program started over",
    program_end: "Program finished",
    alarm_reset: "Stopped by an alarm",
    operator_reset: "Stopped by the operator",
    program_changed: "Program changed",
    offline_reset: "Machine went offline",
    collector_stop: "Server restarted",
    timeout: "Open too long, closed",
  };
  // The time of a cycle, split by machine mode (the cycle's *_sec fields). "Not split yet" is
  // the time of a running cycle the server has not written per mode yet (every 5 minutes).
  const CYCLE_TIME_PARTS = [
    { field: "productionSec", modeId: 3 },
    { field: "holdSec", modeId: 4 },
    { field: "downtimeSec", modeId: 6 },
    { field: "offlineSec", modeId: 9 },
  ];
  const NOT_SPLIT_NAME = "Not split yet";
  const NOT_SPLIT_COLOR = "#3a414c";
  const MONTHS = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];

  const elements = {
    refresh: document.getElementById("refresh"),
    updated: document.getElementById("updated"),
    tiles: document.getElementById("tiles"),
    search: document.getElementById("search"),
    machineList: document.getElementById("machine-list"),
    zoomLevels: document.getElementById("zoom-levels"),
    zoomIn: document.getElementById("zoom-in"),
    zoomOut: document.getElementById("zoom-out"),
    machineName: document.getElementById("machine-name"),
    programName: document.getElementById("program-name"),
    currentMode: document.getElementById("current-mode"),
    selectedPeriod: document.getElementById("selected-period"),
    status: document.getElementById("status"),
    chart: document.getElementById("chart"),
    totals: document.getElementById("totals"),
    periodCount: document.getElementById("period-count"),
    periodRows: document.getElementById("period-rows"),
    currentCycle: document.getElementById("current-cycle"),
    cycleStatus: document.getElementById("cycle-status"),
    cycleCards: document.getElementById("cycle-cards"),
    cycleSplit: document.getElementById("cycle-split"),
    cycleChart: document.getElementById("cycle-chart"),
    cycleCount: document.getElementById("cycle-count"),
    cycleRows: document.getElementById("cycle-rows"),
    cycleNote: document.getElementById("cycle-note"),
  };

  const chart = echarts.init(elements.chart, null, { renderer: "canvas" });
  const cycleChart = echarts.init(elements.cycleChart, null, { renderer: "canvas" });
  let modeHistory = null;
  let cycleHistory = null; // the selected machine's cycles, null when they could not be read
  let machines = [];
  let modeFilter = null; // a mode group key, or null for every machine
  let selectedMachineId = null;

  function modeName(modeId, fallback) {
    return (MODES[modeId] && MODES[modeId].name) || fallback || "Mode " + modeId;
  }

  function modeColor(modeId) {
    return (MODES[modeId] && MODES[modeId].color) || UNKNOWN_MODE_COLOR;
  }

  function formatTime(date) {
    return date.toLocaleTimeString("en-US", { hour: "numeric", minute: "2-digit", second: "2-digit" });
  }

  function formatDateTime(date) {
    return date.getDate() + " " + MONTHS[date.getMonth()] + " " + formatTime(date);
  }

  function formatDuration(totalSeconds) {
    let seconds = Math.max(0, Math.round(totalSeconds));
    const days = Math.floor(seconds / 86400);
    seconds %= 86400;
    const hours = Math.floor(seconds / 3600);
    seconds %= 3600;
    const minutes = Math.floor(seconds / 60);
    seconds %= 60;
    const parts = [];
    if (days) parts.push(days + "d");
    if (hours) parts.push(hours + "h");
    if (minutes) parts.push(minutes + "m");
    if (seconds || parts.length === 0) parts.push(seconds + "s");
    return parts.join(" ");
  }

  function setStatus(message, isError) {
    elements.status.textContent = message;
    elements.status.classList.toggle("error", Boolean(isError));
    elements.status.hidden = !message;
  }

  // periodEnd is when a period ends on the chart: the open period runs up to the response's "to".
  function periodEnd(period) {
    return new Date(period.endedAt || modeHistory.to);
  }

  function chartRows(periods) {
    const modeIds = ALWAYS_SHOWN_MODES.slice();
    periods.forEach(function (period) {
      if (modeIds.indexOf(period.modeId) === -1) modeIds.push(period.modeId);
    });
    return modeIds;
  }

  // Rows of the timeline: "All modes", "Cycles", then one per mode.
  const CYCLES_ROW_INDEX = 1;
  const FIRST_MODE_ROW = 2;

  // Each period is drawn twice: on the "All modes" row and on its own mode's row.
  function chartData(periods, modeIds) {
    const data = [];
    periods.forEach(function (period, index) {
      const start = new Date(period.startedAt).getTime();
      const end = periodEnd(period).getTime();
      const itemStyle = { color: modeColor(period.modeId) };
      data.push({ value: [0, start, end, index], itemStyle: itemStyle });
      data.push({ value: [modeIds.indexOf(period.modeId) + FIRST_MODE_ROW, start, end, index], itemStyle: itemStyle });
    });
    return data;
  }

  function cycleList() {
    return (cycleHistory && cycleHistory.cycles) || [];
  }

  // cycleEnd is when a cycle ends on the chart: a running cycle runs up to now.
  function cycleEnd(cycle) {
    return new Date(cycle.endedAt || modeHistory.to);
  }

  function cycleStatus(cycle) {
    return CYCLE_STATUSES[cycle.status] || { name: cycle.status, color: UNKNOWN_MODE_COLOR };
  }

  // Every other cycle is drawn a little lighter, so cycles that follow each other (a loop
  // program) still show where one ends and the next begins.
  function cycleData() {
    return cycleList().map(function (cycle, index) {
      return {
        value: [CYCLES_ROW_INDEX, new Date(cycle.startedAt).getTime(), cycleEnd(cycle).getTime(), index],
        itemStyle: { color: cycleStatus(cycle).color, opacity: index % 2 ? 0.7 : 1 },
      };
    });
  }

  // cycleNumber is a cycle's number in the shown list, oldest first ("Cycle 12").
  function cycleNumber(cycle) {
    return cycleList().indexOf(cycle) + 1;
  }

  function cycleJob(cycle) {
    return [cycle.jobId && "Job " + cycle.jobId, cycle.partName, cycle.processName].filter(Boolean).join(" · ");
  }

  function cycleEndReason(cycle) {
    return cycle.endReason ? CYCLE_END_REASONS[cycle.endReason] || cycle.endReason : "";
  }

  function cycleTooltip(cycle) {
    const end = cycle.endedAt ? formatDateTime(new Date(cycle.endedAt)) : "now (running)";
    const lines = [
      "<b>Cycle " + cycleNumber(cycle) + "</b> · " + escapeHtml(cycleStatus(cycle).name) +
        (cycle.endReason ? " (" + escapeHtml(cycleEndReason(cycle)) + ")" : ""),
      formatDateTime(new Date(cycle.startedAt)) + " — " + end,
      "Cycle time: <b>" + formatDuration(cycle.durationSec) + "</b>" +
        (cycle.standardCycleSec ? " (planned " + formatDuration(cycle.standardCycleSec) + ")" : ""),
    ];
    CYCLE_TIME_PARTS.forEach(function (part) {
      if (cycle[part.field]) lines.push(escapeHtml(modeName(part.modeId)) + ": " + formatDuration(cycle[part.field]));
    });
    if (cycle.programName) lines.push("Program: " + escapeHtml(cycle.programName));
    if (cycleJob(cycle)) lines.push(escapeHtml(cycleJob(cycle)));
    if (cycle.partsMade) lines.push("Parts: " + cycle.partsMade);
    return lines.join("<br>");
  }

  // renderCycleBar draws a cycle on the Cycles row, 1 px shorter so neighbours do not merge.
  function renderCycleBar(params, api) {
    const shape = renderBar(params, api);
    if (shape && shape.shape.width > 3) shape.shape.width -= 1;
    return shape;
  }

  function renderBar(params, api) {
    const row = api.value(0);
    const start = api.coord([api.value(1), row]);
    const end = api.coord([api.value(2), row]);
    const height = api.size([0, 1])[1] * 0.4;
    const shape = echarts.graphic.clipRectByRect(
      { x: start[0], y: start[1] - height / 2, width: Math.max(end[0] - start[0], 1), height: height },
      { x: params.coordSys.x, y: params.coordSys.y, width: params.coordSys.width, height: params.coordSys.height }
    );
    return shape && { type: "rect", transition: ["shape"], shape: shape, style: { fill: api.visual("color") } };
  }

  function periodTooltip(period) {
    const end = period.endedAt ? formatDateTime(new Date(period.endedAt)) : "now (running)";
    return (
      "<b>" + escapeHtml(modeName(period.modeId, period.modeName)) + "</b><br>" +
      formatDateTime(new Date(period.startedAt)) + " — " + end + "<br>" +
      "Duration: " + formatDuration(period.durationSec) +
      (period.programName ? "<br>Program: " + escapeHtml(period.programName) : "")
    );
  }

  function escapeHtml(text) {
    const span = document.createElement("span");
    span.textContent = text;
    return span.innerHTML;
  }

  function drawChart() {
    const periods = modeHistory.periods || [];
    const modeIds = chartRows(periods);
    const rowNames = [ALL_MODES_ROW, CYCLES_ROW].concat(modeIds.map(function (modeId) { return modeName(modeId); }));
    const from = new Date(modeHistory.from).getTime();
    const to = new Date(modeHistory.to).getTime();

    chart.setOption({
      backgroundColor: "transparent",
      animation: false,
      grid: { left: 190, right: 30, top: 40, bottom: 60 },
      tooltip: {
        trigger: "item",
        backgroundColor: "#1f242b",
        borderColor: "#2a2f37",
        textStyle: { color: "#d7dbe0" },
        formatter: function (params) {
          return params.seriesIndex === 1 ? cycleTooltip(cycleList()[params.value[3]]) : periodTooltip(periods[params.value[3]]);
        },
      },
      dataZoom: [
        { type: "inside", xAxisIndex: 0, filterMode: "weakFilter", minValueSpan: MIN_SPAN_MS },
        {
          type: "slider", xAxisIndex: 0, filterMode: "weakFilter", height: 18, bottom: 16,
          minValueSpan: MIN_SPAN_MS, borderColor: "#2a2f37", textStyle: { color: "#8b939e" },
          labelFormatter: function (value) { return formatDateTime(new Date(value)); },
        },
      ],
      xAxis: {
        type: "time",
        position: "top",
        min: from,
        max: to,
        axisLine: { show: false },
        axisTick: { show: false },
        splitLine: { show: true, lineStyle: { color: "#232830" } },
        axisLabel: {
          color: "#f28c28",
          hideOverlap: true,
          formatter: function (value) {
            const date = new Date(value);
            const span = currentSpan();
            return span > DAY_MS ? date.getDate() + " " + MONTHS[date.getMonth()] + "\n" + formatTime(date) : formatTime(date);
          },
        },
      },
      yAxis: {
        type: "category",
        inverse: true,
        data: rowNames,
        axisLine: { show: false },
        axisTick: { show: false },
        splitLine: { show: true, lineStyle: { color: "#232830" } },
        axisLabel: {
          color: "#d7dbe0",
          fontSize: 13,
          formatter: function (name) { return name === ALL_MODES_ROW ? "" : name; },
        },
      },
      series: [{
        type: "custom",
        renderItem: renderBar,
        encode: { x: [1, 2], y: 0 },
        data: chartData(periods, modeIds),
      }, {
        type: "custom",
        renderItem: renderCycleBar,
        encode: { x: [1, 2], y: 0 },
        data: cycleData(),
      }],
    }, true);

    showSpan(DEFAULT_ZOOM_MS, new Date(modeHistory.to).getTime());
  }

  // currentWindow is the time range the chart shows now, or null before a machine is drawn.
  function currentWindow() {
    const option = chart.getOption();
    if (!option || !option.dataZoom || !option.dataZoom[0] || !modeHistory) return null;
    const zoom = option.dataZoom[0];
    const from = new Date(modeHistory.from).getTime();
    const total = new Date(modeHistory.to).getTime() - from;
    return {
      start: zoom.startValue != null ? zoom.startValue : from + total * (zoom.start || 0) / 100,
      end: zoom.endValue != null ? zoom.endValue : from + total * (zoom.end == null ? 100 : zoom.end) / 100,
    };
  }

  function currentSpan() {
    const view = currentWindow();
    return view ? view.end - view.start : 0;
  }

  // showSpan shows spanMs of time ending at end, kept inside the history window.
  function showSpan(spanMs, end) {
    const from = new Date(modeHistory.from).getTime();
    const to = new Date(modeHistory.to).getTime();
    const span = Math.min(to - from, Math.max(MIN_SPAN_MS, spanMs));
    const newEnd = Math.min(to, Math.max(from + span, end));
    chart.dispatchAction({ type: "dataZoom", startValue: newEnd - span, endValue: newEnd });
  }

  // zoomToLevel shows a zoom level's span ending where the current view ends, so the levels
  // zoom into the part being looked at (at first the latest data).
  function zoomToLevel(spanMs) {
    if (!modeHistory) return;
    const view = currentWindow();
    showSpan(spanMs, view ? view.end : new Date(modeHistory.to).getTime());
  }

  // zoom multiplies the shown span by factor around the middle of the view (the − / + buttons).
  function zoom(factor) {
    const view = currentWindow();
    if (!view) return;
    const span = (view.end - view.start) * factor;
    showSpan(span, (view.start + view.end) / 2 + span / 2);
  }

  function drawZoomLevels() {
    ZOOM_LEVELS.forEach(function (level) {
      const button = document.createElement("button");
      button.type = "button";
      button.textContent = level.label;
      button.title = "Show " + level.title;
      button.dataset.spanMs = level.spanMs;
      button.addEventListener("click", function () { zoomToLevel(level.spanMs); });
      elements.zoomLevels.append(button);
    });
  }

  // markZoomLevel highlights the zoom level matching the shown span (none after a free zoom).
  function markZoomLevel() {
    const span = currentSpan();
    elements.zoomLevels.querySelectorAll("button").forEach(function (button) {
      const levelSpan = Number(button.dataset.spanMs);
      button.classList.toggle("active", Math.abs(span - levelSpan) <= levelSpan * ZOOM_MATCH_TOLERANCE);
    });
  }

  function showSelectedPeriod(period) {
    const box = elements.selectedPeriod;
    box.replaceChildren();
    if (!period) {
      box.hidden = true;
      return;
    }
    const swatch = document.createElement("span");
    swatch.className = "swatch";
    swatch.style.background = modeColor(period.modeId);
    const name = document.createElement("span");
    name.className = "info";
    name.textContent = modeName(period.modeId, period.modeName);
    const time = document.createElement("strong");
    const end = period.endedAt ? formatTime(new Date(period.endedAt)) : "now";
    time.textContent = formatDateTime(new Date(period.startedAt)) + " — " + end + " - " + formatDuration(period.durationSec);
    box.append(swatch, name, time);
    box.hidden = false;
  }

  function drawHeader() {
    const periods = modeHistory.periods || [];
    const latest = periods[periods.length - 1];
    elements.machineName.textContent = modeHistory.machineName || "Machine " + modeHistory.machineId;
    elements.programName.textContent = (latest && latest.programName) || "—";
    elements.currentMode.textContent = modeHistory.currentModeId != null ? modeName(modeHistory.currentModeId) : "—";
    showSelectedPeriod(latest);
  }

  function drawTotals() {
    elements.totals.replaceChildren();
    (modeHistory.totals || []).forEach(function (total) {
      const item = document.createElement("div");
      item.className = "total";
      const swatch = document.createElement("span");
      swatch.className = "swatch";
      swatch.style.background = modeColor(total.modeId);
      const text = document.createElement("span");
      text.textContent = modeName(total.modeId, total.modeName) + " (" + total.modeId + "): " +
        formatDuration(total.durationSec) + " · " + total.periodCount + "×";
      item.append(swatch, text);
      elements.totals.append(item);
    });
  }

  function cell(text, className) {
    const td = document.createElement("td");
    td.textContent = text;
    if (className) td.className = className;
    return td;
  }

  // drawPeriodList lists the periods newest first.
  function drawPeriodList() {
    const periods = (modeHistory.periods || []).slice().reverse();
    elements.periodRows.replaceChildren();
    elements.periodCount.textContent = "(" + periods.length + " periods)";
    periods.forEach(function (period, index) {
      const row = document.createElement("tr");
      const modeCell = document.createElement("td");
      const mode = document.createElement("span");
      mode.className = "mode-cell";
      const swatch = document.createElement("span");
      swatch.className = "swatch";
      swatch.style.background = modeColor(period.modeId);
      mode.append(swatch, document.createTextNode(modeName(period.modeId, period.modeName)));
      modeCell.append(mode);
      const end = period.endedAt ? cell(formatDateTime(new Date(period.endedAt))) : cell("Running", "running");
      row.append(
        cell(String(index + 1)),
        cell(String(period.modeId)),
        modeCell,
        cell(period.programName || "—"),
        cell(formatDateTime(new Date(period.startedAt))),
        end,
        cell(formatDuration(period.durationSec), "number")
      );
      elements.periodRows.append(row);
    });
  }

  // ---------------------------------------------------------------- cycles

  // cycleTimeParts splits a cycle's time by mode. A running cycle's time the server has not split
  // yet is the "Not split yet" part, so the parts always add up to the cycle time.
  function cycleTimeParts(cycle) {
    const parts = CYCLE_TIME_PARTS.map(function (part) {
      return { name: modeName(part.modeId), color: modeColor(part.modeId), seconds: cycle[part.field] || 0 };
    });
    const splitSeconds = parts.reduce(function (sum, part) { return sum + part.seconds; }, 0);
    const notSplit = cycle.status === "running" ? Math.max(0, cycle.durationSec - splitSeconds) : 0;
    parts.push({ name: NOT_SPLIT_NAME, color: NOT_SPLIT_COLOR, seconds: notSplit });
    return parts;
  }

  function latestCycle() {
    const cycles = cycleList();
    return cycles[cycles.length - 1];
  }

  function lastClosedCycle() {
    const cycles = cycleList();
    for (let index = cycles.length - 1; index >= 0; index--) {
      if (cycles[index].status !== "running") return cycles[index];
    }
    return null;
  }

  function setCycleStatus(message, isError) {
    elements.cycleStatus.textContent = message;
    elements.cycleStatus.classList.toggle("error", Boolean(isError));
    elements.cycleStatus.hidden = !message;
  }

  function drawCurrentCycle() {
    const latest = latestCycle();
    let text = "—";
    if (latest && latest.status === "running") {
      text = "running for " + formatDuration(latest.durationSec);
    } else if (latest) {
      text = "last one " + cycleStatus(latest).name.toLowerCase() + " in " + formatDuration(latest.durationSec);
    } else if (cycleHistory) {
      text = "none in the last " + CYCLE_HISTORY_HOURS + " h";
    }
    elements.currentCycle.textContent = text;
  }

  function card(label, value, note, color) {
    const box = document.createElement("div");
    box.className = "card";
    if (color) box.style.setProperty("--card-color", color);
    const labelText = document.createElement("div");
    labelText.className = "card-label";
    labelText.textContent = label;
    const valueText = document.createElement("div");
    valueText.className = "card-value";
    valueText.textContent = value;
    box.append(labelText, valueText);
    if (note) {
      const noteText = document.createElement("div");
      noteText.className = "card-note";
      noteText.textContent = note;
      box.append(noteText);
    }
    return box;
  }

  // drawCycleCards shows the cycles of the last 24 hours in a few numbers.
  function drawCycleCards() {
    elements.cycleCards.replaceChildren();
    if (!cycleHistory) return;
    const summary = cycleHistory.summary;
    const latest = latestCycle();
    const running = latest && latest.status === "running" ? latest : null;
    const lastClosed = lastClosedCycle();
    const planned = latest && latest.standardCycleSec;
    elements.cycleCards.append(
      card("Completed cycles", summary.completedCount, "runs that finished the part", CYCLE_STATUSES.completed.color),
      card("Unfinished cycles", summary.unfinishedCount, "runs stopped before the end", CYCLE_STATUSES.unfinished.color),
      card("Running now", running ? formatDuration(running.durationSec) : "No",
        running ? "cycle " + cycleNumber(running) + " started " + formatTime(new Date(running.startedAt)) : "no cycle is running",
        CYCLE_STATUSES.running.color),
      card("Average cycle time", summary.completedCount ? formatDuration(summary.averageCycleSec) : "—",
        planned ? "planned " + formatDuration(planned) : "of the completed cycles"),
      card("Last cycle", lastClosed ? formatDuration(lastClosed.durationSec) : "—",
        lastClosed ? cycleStatus(lastClosed).name + " · " + cycleEndReason(lastClosed) : "none closed yet"),
      card("Parts counted", summary.partsMade, "by the machine's part counter")
    );
  }

  // drawCycleSplit shows where the time of all shown cycles went, as one bar split by mode.
  function drawCycleSplit() {
    elements.cycleSplit.replaceChildren();
    const parts = [];
    cycleList().forEach(function (cycle) {
      cycleTimeParts(cycle).forEach(function (part, index) {
        if (!parts[index]) parts[index] = { name: part.name, color: part.color, seconds: 0 };
        parts[index].seconds += part.seconds;
      });
    });
    const total = parts.reduce(function (sum, part) { return sum + part.seconds; }, 0);
    if (!total) return;

    const title = document.createElement("div");
    title.className = "card-label";
    title.textContent = "Where the cycle time went";
    const bar = document.createElement("div");
    bar.className = "split-bar";
    const legend = document.createElement("div");
    legend.className = "split-legend";
    parts.forEach(function (part) {
      if (!part.seconds) return;
      const percent = Math.round(part.seconds / total * 100);
      const text = part.name + ": " + formatDuration(part.seconds) + " (" + percent + "%)";
      const piece = document.createElement("span");
      piece.style.flexGrow = part.seconds;
      piece.style.background = part.color;
      piece.title = text;
      bar.append(piece);
      const item = document.createElement("span");
      const swatch = document.createElement("span");
      swatch.className = "swatch";
      swatch.style.background = part.color;
      item.append(swatch, document.createTextNode(text));
      legend.append(item);
    });
    elements.cycleSplit.append(title, bar, legend);
  }

  // drawCycleChart shows each cycle's time as a bar split by mode, with the planned cycle time
  // as a dashed line; it opens on the newest CYCLE_BARS_SHOWN cycles.
  function drawCycleChart() {
    const cycles = cycleList();
    elements.cycleChart.hidden = !cycles.length;
    if (!cycles.length) {
      cycleChart.clear();
      return;
    }
    cycleChart.resize();
    const timeParts = cycles.map(cycleTimeParts);
    const planned = latestCycle().standardCycleSec;
    const series = timeParts[0].map(function (part, partIndex) {
      return {
        name: part.name,
        type: "bar",
        stack: "cycle time",
        barMaxWidth: 28,
        itemStyle: { color: part.color },
        data: timeParts.map(function (parts) { return parts[partIndex].seconds; }),
      };
    });
    if (planned) {
      series[0].markLine = {
        silent: true,
        symbol: "none",
        lineStyle: { color: "#f28c28", type: "dashed" },
        label: { color: "#f28c28", formatter: "Planned " + formatDuration(planned) },
        data: [{ yAxis: planned }],
      };
    }

    cycleChart.setOption({
      backgroundColor: "transparent",
      animation: false,
      grid: { left: 70, right: 90, top: 40, bottom: 70 },
      legend: { top: 4, textStyle: { color: "#d7dbe0" }, data: series.map(function (item) { return item.name; }) },
      tooltip: {
        trigger: "axis",
        axisPointer: { type: "shadow" },
        backgroundColor: "#1f242b",
        borderColor: "#2a2f37",
        textStyle: { color: "#d7dbe0" },
        formatter: function (params) { return cycleTooltip(cycles[params[0].dataIndex]); },
      },
      dataZoom: [
        { type: "inside", startValue: Math.max(0, cycles.length - CYCLE_BARS_SHOWN), endValue: cycles.length - 1 },
        {
          type: "slider", height: 18, bottom: 16, borderColor: "#2a2f37", textStyle: { color: "#8b939e" },
          startValue: Math.max(0, cycles.length - CYCLE_BARS_SHOWN), endValue: cycles.length - 1,
        },
      ],
      xAxis: {
        type: "category",
        data: cycles.map(function (cycle, index) { return "Cycle " + (index + 1) + "\n" + formatTime(new Date(cycle.startedAt)); }),
        axisLabel: { color: "#8b939e", hideOverlap: true },
        axisLine: { lineStyle: { color: "#2a2f37" } },
      },
      yAxis: {
        type: "value",
        name: "Cycle time",
        nameTextStyle: { color: "#8b939e" },
        axisLabel: { color: "#8b939e", formatter: function (value) { return formatDuration(value); } },
        splitLine: { lineStyle: { color: "#232830" } },
      },
      series: series,
    }, true);
  }

  function statusBadge(cycle) {
    const badge = document.createElement("span");
    badge.className = "badge";
    badge.style.setProperty("--badge-color", cycleStatus(cycle).color);
    badge.textContent = cycleStatus(cycle).name;
    return badge;
  }

  // secondsCell shows a duration, or "—" for none.
  function secondsCell(seconds) {
    return cell(seconds ? formatDuration(seconds) : "—", "number");
  }

  // drawCycleList lists the newest cycles first (at most CYCLE_TABLE_LIMIT); a click zooms the
  // timeline to the cycle.
  function drawCycleList() {
    const cycles = cycleList();
    elements.cycleRows.replaceChildren();
    elements.cycleCount.textContent = cycleHistory ? "(" + cycles.length + " cycles)" : "";
    const shown = cycles.slice(-CYCLE_TABLE_LIMIT).reverse();
    shown.forEach(function (cycle) {
      const row = document.createElement("tr");
      row.className = "clickable";
      row.title = "Show this cycle on the chart";
      const statusCell = document.createElement("td");
      statusCell.append(statusBadge(cycle));
      const end = cycle.endedAt ? cell(formatDateTime(new Date(cycle.endedAt))) : cell("Running", "running");
      row.append(
        cell(String(cycleNumber(cycle))),
        statusCell,
        cell(cycleEndReason(cycle) || "—"),
        cell(cycle.programName || "—"),
        cell(cycleJob(cycle) || "—"),
        cell(formatDateTime(new Date(cycle.startedAt))),
        end,
        cell(formatDuration(cycle.durationSec), "number"),
        secondsCell(cycle.productionSec),
        secondsCell(cycle.holdSec),
        secondsCell(cycle.downtimeSec),
        secondsCell(cycle.offlineSec),
        cell(cycle.loadingSec != null ? formatDuration(cycle.loadingSec) : "—", "number"),
        cell(cycle.partsMade ? String(cycle.partsMade) : "—", "number")
      );
      row.addEventListener("click", function () { focusCycle(cycle); });
      elements.cycleRows.append(row);
    });
    let note = "";
    if (cycleHistory && !cycles.length) {
      note = "No cycles in the last " + CYCLE_HISTORY_HOURS + " hours. A cycle starts when the program starts running.";
    } else if (cycles.length > shown.length) {
      note = "Showing the newest " + shown.length + " of " + cycles.length + " cycles.";
    }
    elements.cycleNote.textContent = note;
  }

  // focusCycle zooms the timeline to cycle with a little time around it and shows it above the chart.
  function focusCycle(cycle) {
    if (!modeHistory) return;
    const start = new Date(cycle.startedAt).getTime();
    const end = cycleEnd(cycle).getTime();
    const margin = Math.max((end - start) * 0.1, 5 * SECOND_MS);
    showSpan(end - start + 2 * margin, end + margin);
    showSelectedCycle(cycle);
    elements.chart.scrollIntoView({ behavior: "smooth", block: "center" });
  }

  function showSelectedCycle(cycle) {
    const box = elements.selectedPeriod;
    box.replaceChildren();
    const swatch = document.createElement("span");
    swatch.className = "swatch";
    swatch.style.background = cycleStatus(cycle).color;
    const name = document.createElement("span");
    name.className = "info";
    name.textContent = "Cycle " + cycleNumber(cycle) + " (" + cycleStatus(cycle).name + ")";
    const time = document.createElement("strong");
    const end = cycle.endedAt ? formatTime(new Date(cycle.endedAt)) : "now";
    time.textContent = formatDateTime(new Date(cycle.startedAt)) + " — " + end + " - " + formatDuration(cycle.durationSec);
    box.append(swatch, name, time);
    box.hidden = false;
  }

  function drawCycles() {
    drawCurrentCycle();
    drawCycleCards();
    drawCycleSplit();
    drawCycleChart();
    drawCycleList();
  }

  // groupKey is the machine list group of a machine: its mode id, or NO_MODE.
  function groupKey(machine) {
    return machine.modeId != null ? String(machine.modeId) : NO_MODE;
  }

  function groupName(key) {
    return key === NO_MODE ? "No data" : modeName(Number(key));
  }

  function groupColor(key) {
    return key === NO_MODE ? UNKNOWN_MODE_COLOR : modeColor(Number(key));
  }

  // groupKeys lists the groups in dashboard order: the known modes, other modes, then NO_MODE.
  function groupKeys(counts) {
    const keys = DASHBOARD_MODE_ORDER.map(String);
    Object.keys(counts).forEach(function (key) {
      if (key !== NO_MODE && keys.indexOf(key) === -1) keys.push(key);
    });
    keys.push(NO_MODE);
    return keys;
  }

  function countByGroup() {
    const counts = {};
    machines.forEach(function (machine) {
      const key = groupKey(machine);
      counts[key] = (counts[key] || 0) + 1;
    });
    return counts;
  }

  // tile is one mode count; a click shows only that mode's machines, a second click all again.
  function tile(key, name, count, color) {
    const button = document.createElement("button");
    button.type = "button";
    button.className = "tile" + (modeFilter === key ? " active" : "") + (count ? "" : " empty");
    button.style.setProperty("--tile-color", color);
    button.title = key === null || modeFilter === key ? "Show every machine" : "Show only " + name;
    const number = document.createElement("span");
    number.className = "tile-count";
    number.textContent = count;
    const label = document.createElement("span");
    label.className = "tile-name";
    label.textContent = name;
    button.append(number, label);
    button.addEventListener("click", function () {
      modeFilter = modeFilter === key ? null : key;
      drawDashboard();
    });
    return button;
  }

  // drawTiles shows the number of machines per mode: always the main modes, others while used.
  function drawTiles(counts) {
    elements.tiles.replaceChildren(tile(null, "All machines", machines.length, "var(--accent)"));
    groupKeys(counts).forEach(function (key) {
      const count = counts[key] || 0;
      if (count || DASHBOARD_TILE_MODES.indexOf(Number(key)) !== -1) {
        elements.tiles.append(tile(key, groupName(key), count, groupColor(key)));
      }
    });
  }

  function matchesSearch(machine, text) {
    if (!text) return true;
    return [machine.machineName, machine.customerName, machine.companyName, String(machine.machineId)]
      .some(function (value) { return value && value.toLowerCase().indexOf(text) !== -1; });
  }

  function machineButton(machine) {
    const button = document.createElement("button");
    button.type = "button";
    button.className = "machine" + (machine.machineId === selectedMachineId ? " selected" : "");
    const swatch = document.createElement("span");
    swatch.className = "swatch";
    swatch.style.background = groupColor(groupKey(machine));
    const text = document.createElement("span");
    text.className = "machine-text";
    const title = document.createElement("div");
    title.className = "machine-title";
    title.textContent = machine.machineName || "Machine " + machine.machineId;
    const details = document.createElement("div");
    details.className = "machine-sub";
    details.textContent = ["#" + machine.machineId, machine.customerName, machine.programName]
      .filter(Boolean).join(" · ");
    text.append(title, details);
    const since = document.createElement("span");
    since.className = "machine-since";
    since.textContent = machine.modeDurationSec != null ? formatDuration(machine.modeDurationSec) : "";
    if (machine.modeSince) since.title = "In this mode since " + formatDateTime(new Date(machine.modeSince));
    button.append(swatch, text, since);
    button.addEventListener("click", function () { selectMachine(machine.machineId); });
    return button;
  }

  // drawMachineList lists the machines that pass the tile filter and the search, grouped by mode.
  function drawMachineList(counts) {
    const text = elements.search.value.trim().toLowerCase();
    const groups = {};
    machines.forEach(function (machine) {
      const key = groupKey(machine);
      if ((modeFilter === null || modeFilter === key) && matchesSearch(machine, text)) {
        (groups[key] = groups[key] || []).push(machine);
      }
    });

    elements.machineList.replaceChildren();
    groupKeys(counts).forEach(function (key) {
      const members = groups[key];
      if (!members) return;
      const head = document.createElement("div");
      head.className = "group-head";
      const swatch = document.createElement("span");
      swatch.className = "swatch";
      swatch.style.background = groupColor(key);
      const name = document.createElement("span");
      name.textContent = groupName(key);
      const count = document.createElement("span");
      count.className = "group-count";
      count.textContent = members.length;
      head.append(swatch, name, count);
      elements.machineList.append(head);
      members.forEach(function (machine) { elements.machineList.append(machineButton(machine)); });
    });

    if (!elements.machineList.childElementCount) {
      const empty = document.createElement("div");
      empty.className = "list-empty";
      empty.textContent = machines.length ? "No machine matches." : "No machines.";
      elements.machineList.append(empty);
    }
  }

  function drawDashboard() {
    const counts = countByGroup();
    drawTiles(counts);
    drawMachineList(counts);
  }

  async function loadMachines() {
    try {
      const reply = await fetch(MACHINES_PATH);
      const body = await reply.json();
      if (!reply.ok || !body.success) {
        throw new Error(body.message || "Request failed (" + reply.status + ")");
      }
      machines = body.data.machines || [];
      elements.updated.textContent = "Updated " + formatTime(new Date()) +
        " · refreshes every " + LIST_REFRESH_MS / 1000 + " s";
    } catch (error) {
      elements.updated.textContent = "Could not load machines: " + error.message;
      return;
    }
    drawDashboard();
  }

  // selectMachine shows machineId's chart and keeps it in the URL (?machineId=).
  function selectMachine(machineId) {
    selectedMachineId = machineId;
    const url = new URL(window.location);
    url.searchParams.set("machineId", machineId);
    window.history.replaceState(null, "", url);
    drawDashboard();
    loadHistory(machineId);
  }

  // fetchData returns the data of a public API response, or throws its message.
  async function fetchData(path) {
    const reply = await fetch(path);
    const body = await reply.json();
    if (!reply.ok || !body.success) {
      throw new Error(body.message || "Request failed (" + reply.status + ")");
    }
    return body.data;
  }

  // loadHistory reads the machine's modes and cycles together. Cycles that cannot be read leave
  // the mode chart as it is and say so in the Cycles section.
  async function loadHistory(machineId) {
    setStatus("Loading…");
    const cyclesRequest = fetchData(API_PATH + machineId + CYCLES_SUFFIX).catch(function (error) { return error; });
    let modes;
    try {
      modes = await fetchData(API_PATH + machineId + "/modes");
    } catch (error) {
      if (machineId !== selectedMachineId) return;
      modeHistory = null;
      cycleHistory = null;
      chart.clear();
      setStatus(error.message, true);
      setCycleStatus("");
      drawCycles();
      return;
    }
    const cycles = await cyclesRequest;
    if (machineId !== selectedMachineId) return; // another machine was picked meanwhile
    modeHistory = modes;
    cycleHistory = cycles instanceof Error ? null : cycles;
    setCycleStatus(cycles instanceof Error ? "Could not load cycles: " + cycles.message : "", cycles instanceof Error);

    const count = (modeHistory.periods || []).length;
    setStatus(count ? "" : "No mode data in the last 10 days for this machine.");
    drawHeader();
    drawChart();
    drawTotals();
    drawPeriodList();
    drawCycles();
  }

  chart.on("click", function (params) {
    if (!modeHistory || !params.value) return;
    if (params.seriesIndex === 1) {
      showSelectedCycle(cycleList()[params.value[3]]);
    } else {
      showSelectedPeriod(modeHistory.periods[params.value[3]]);
    }
  });
  cycleChart.on("click", function (params) {
    const cycle = cycleList()[params.dataIndex];
    if (cycle) focusCycle(cycle);
  });
  elements.refresh.addEventListener("click", function () {
    loadMachines();
    if (selectedMachineId) loadHistory(selectedMachineId);
  });
  elements.search.addEventListener("input", drawDashboard);
  chart.on("datazoom", markZoomLevel);
  elements.zoomIn.addEventListener("click", function () { zoom(1 / ZOOM_STEP); });
  elements.zoomOut.addEventListener("click", function () { zoom(ZOOM_STEP); });
  window.addEventListener("resize", function () {
    chart.resize();
    cycleChart.resize();
  });

  // The list refreshes on its own; the chart reloads on Refresh or when a machine is picked, so
  // a zoomed chart is not reset under the reader.
  drawZoomLevels();
  loadMachines();
  window.setInterval(loadMachines, LIST_REFRESH_MS);
  const initialMachineId = Number(new URLSearchParams(window.location.search).get("machineId"));
  if (Number.isInteger(initialMachineId) && initialMachineId > 0) selectMachine(initialMachineId);
})();
