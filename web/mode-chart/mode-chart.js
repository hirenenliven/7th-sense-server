// Machine mode dashboard: the left panel lists every machine grouped by its current mode
// (GET /api/v1/health/machines), with a count per mode that filters the list. Selecting a machine
// reads GET /api/v1/health/machines/{id}/modes (both public, no session) and draws its mode
// periods as a timeline, one row per mode, with the list of periods below.
(function () {
  "use strict";

  const API_PATH = "/api/v1/health/machines/";
  const MACHINES_PATH = "/api/v1/health/machines";
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
  const DEFAULT_ZOOM_MS = 10 * DAY_MS; // a newly selected machine shows its whole history
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
  };

  const chart = echarts.init(elements.chart, null, { renderer: "canvas" });
  let modeHistory = null;
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

  // Each period is drawn twice: on the "All modes" row and on its own mode's row.
  function chartData(periods, modeIds) {
    const data = [];
    periods.forEach(function (period, index) {
      const start = new Date(period.startedAt).getTime();
      const end = periodEnd(period).getTime();
      const itemStyle = { color: modeColor(period.modeId) };
      data.push({ value: [0, start, end, index], itemStyle: itemStyle });
      data.push({ value: [modeIds.indexOf(period.modeId) + 1, start, end, index], itemStyle: itemStyle });
    });
    return data;
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
    const rowNames = [ALL_MODES_ROW].concat(modeIds.map(function (modeId) { return modeName(modeId); }));
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
        formatter: function (params) { return periodTooltip(periods[params.value[3]]); },
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

  async function loadHistory(machineId) {
    setStatus("Loading…");
    try {
      const reply = await fetch(API_PATH + machineId + "/modes");
      const body = await reply.json();
      if (!reply.ok || !body.success) {
        throw new Error(body.message || "Request failed (" + reply.status + ")");
      }
      if (machineId !== selectedMachineId) return; // another machine was picked meanwhile
      modeHistory = body.data;
    } catch (error) {
      if (machineId !== selectedMachineId) return;
      modeHistory = null;
      chart.clear();
      setStatus(error.message, true);
      return;
    }

    const count = (modeHistory.periods || []).length;
    setStatus(count ? "" : "No mode data in the last 10 days for this machine.");
    drawHeader();
    drawChart();
    drawTotals();
    drawPeriodList();
  }

  chart.on("click", function (params) {
    if (modeHistory && params.value) showSelectedPeriod(modeHistory.periods[params.value[3]]);
  });
  elements.refresh.addEventListener("click", function () {
    loadMachines();
    if (selectedMachineId) loadHistory(selectedMachineId);
  });
  elements.search.addEventListener("input", drawDashboard);
  chart.on("datazoom", markZoomLevel);
  elements.zoomIn.addEventListener("click", function () { zoom(1 / ZOOM_STEP); });
  elements.zoomOut.addEventListener("click", function () { zoom(ZOOM_STEP); });
  window.addEventListener("resize", function () { chart.resize(); });

  // The list refreshes on its own; the chart reloads on Refresh or when a machine is picked, so
  // a zoomed chart is not reset under the reader.
  drawZoomLevels();
  loadMachines();
  window.setInterval(loadMachines, LIST_REFRESH_MS);
  const initialMachineId = Number(new URLSearchParams(window.location.search).get("machineId"));
  if (Number.isInteger(initialMachineId) && initialMachineId > 0) selectMachine(initialMachineId);
})();
