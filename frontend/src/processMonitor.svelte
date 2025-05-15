<script>
  import { onMount } from "svelte";
  import {
    GetProcesses,
    SigKillProcess,
    SigTerminateProcess,
  } from "../wailsjs/go/main/App.js";

  let processData = {};
  let error = null;
  let selectedProcess = null;
  let showMenu = false;

  async function fetchProcessData() {
    try {
      const result = await GetProcesses();
      if (!result) {
        throw new Error("No process data received from backend");
      }

      processData = result;

      error = null;
    } catch (err) {
      console.error("Error fetching process data", err);
      error = err.message;
    }
  }

  function openMenu(process) {
    selectedProcess = process;
    showMenu = true;
  }

  function closeMenu() {
    showMenu = false;
    selectedProcess = null;
  }

  async function killProcess() {
    try {
      await SigKillProcess(selectedProcess.PID);
      alert(
        `Killing process ${selectedProcess.Name} (PID: ${selectedProcess.PID})`,
      );
      closeMenu();
    } catch (err) {
      console.error(err);
      throw new Error("Error terminating process: ", err);
    }
  }

  async function terminateProcess() {
    try {
      await SigTerminateProcess(selectedProcess.PID);
      alert(
        `Terminating process ${selectedProcess.Name} (PID: ${selectedProcess.PID})`,
      );
      closeMenu();
    } catch (err) {
      console.error(err);
      throw new Error("Error terminating process: ", err);
    }
  }

  onMount(() => {
    fetchProcessData();
    const interval = setInterval(fetchProcessData, 1000);
    return () => clearInterval(interval);
  });

  function toGB(bytes) {
    if (typeof bytes !== "number") {
      console.error("Invalid bytes value:", bytes);
      return "0.00";
    }
    return (bytes / (1024 * 1024 * 1024)).toFixed(2);
  }
</script>

<main>
  <h1>Process Monitor</h1>
  {#if error}
    <p style="color: red;">{error}</p>
  {:else}
    <ul>
      {#each Object.entries(processData) as [_, process]}
        <li class="process-row" on:click={() => openMenu(process)}>
          <h3>{process["Name"]}</h3>
          <small>{process["Cwd"]}</small>
          <p>PID: {process["PID"]}</p>
          <p>Memory: {toGB(process["MemoryInfo"]["rss"])} GB</p>
          <p>CPU: {process["CPUPercent"].toFixed(2)}%</p>
          <p>Threads: {process["Threads"]}</p>
        </li>
      {/each}
    </ul>
  {/if}

  {#if showMenu && selectedProcess}
    <div class="process-menu">
      <h2>{selectedProcess["Name"]} (PID: {selectedProcess["PID"]})</h2>
      <button on:click={killProcess}>Kill</button>
      <button on:click={terminateProcess}>Terminate</button>
      <button on:click={closeMenu}>Cancel</button>
    </div>
  {/if}
</main>

<style>
  main {
    padding: 20px;
  }

  .process-info {
    border: 1px solid #ccc;
    padding: 10px;
    margin-bottom: 10px;
  }

  .process-info h3 {
    margin: 0;
    color: #007bff;
  }

  .process-info small {
    display: block;
    color: #666;
  }

  ul {
    list-style: none;
    padding: 0;
  }

  .process-row {
    padding: 10px;
    border-bottom: 1px solid #444;
    cursor: pointer;
    transition: background 0.2s;
  }

  .process-row:hover {
    background: #222;
  }

  .process-menu {
    position: fixed;
    left: 0;
    right: 0;
    bottom: 0;
    background: #222;
    border-top: 2px solid #444;
    padding: 20px;
    box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.5);
    z-index: 1000;
    text-align: center;
  }

  .process-menu button {
    margin: 0 10px;
    padding: 10px 20px;
    font-size: 1em;
  }
</style>
