<script lang="ts">
  import logo from "./assets/images/logo-universal.png";
  import { GetRAM } from "../wailsjs/go/main/App.js";
  import { onMount } from "svelte";

  interface RamStats {
    "Total Memory": number;
    "Free Memory": number;
    "Percent Used": number;
  }

  let ram: RamStats = {
    "Total Memory": 0,
    "Free Memory": 0,
    "Percent Used": 0,
  };
  let error: string | null = null;

  async function fetchRAM() {
    try {
      const result = await GetRAM();
      ram = {
        "Total Memory": result["Total Memory"],
        "Free Memory": result["Free Memory"],
        "Percent Used": result["Percent Used"],
      };
      error = null;
    } catch (err) {
      console.log("Error fetching RAM", err);
      error = err.message;
    }
  }

  onMount(() => {
    fetchRAM();
    const interval = setInterval(fetchRAM, 1000);
    return () => clearInterval(interval);
  });

  function toGB(bytes) {
    return (bytes / (1024 * 1024 * 1024)).toFixed(2);
  }
</script>

<main>
  <h1>System RAM Monitor</h1>
  {#if error}
    <p style="color: red;">Error: {error}</p>
  {:else}
    <ul>
      <li>Total Memory: {toGB(ram["Total memory"])} GB</li>
      <li>Free Memory: {toGB(ram["Free memory"])} GB</li>
      <li>Percent Used: {ram["Percent used"].toFixed(2)}%</li>
    </ul>
  {/if}
</main>

<style>
  main {
    padding: 20px;
    font-family: Arial, sans-serif;
  }
  ul {
    list-style: none;
    padding: 0;
  }
  li {
    margin: 10px 0;
  }
</style>
