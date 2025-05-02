<script>
  import { onMount } from "svelte";
  import { GetDiskUsage } from "../wailsjs/go/main/App.js";
  import { toggle_class } from "svelte/internal";

  let diskData= {};
  let error = null;

  async function fetchDiskUsage() {
    try {
        const result = await GetDiskUsage();
        if (!result) {
          throw new Error("No disk data received from backend");
        }
        diskData = result
        error = null;
    } catch (err) {
        console.error("Error fetching disk usage", err);
        error = err.message;
    }
  }

  function toGB(bytes) {
    if (typeof bytes !== "number") {
      console.error("Invalid bytes value:", bytes);
      return "0.00";
    }
    return (bytes / (1024 * 1024 * 1024)).toFixed(2);
  }

  onMount(() => {
      fetchDiskUsage();
      const interval = setInterval(fetchDiskUsage, 1000);
      return () => clearInterval(interval);
  });
</script>

<main>
  <h1>Disk Usage</h1>
  {#if error}
    <p style="color: red;">{error}</p>
  {:else}
    {#each Object.entries(diskData) as [mountpoint, usage]}
      <div class="disk-info">
        <h3>{mountpoint}</h3>
        <p>Total: {toGB(usage["total"])} GB</p>
        <p>Free: {toGB(usage["free"])} GB</p>
        <p>Used: {toGB(usage["used"])} GB</p>
        <p>Used Percent: {usage["usedPercent"].toFixed(2)}%</p>
      </div>
    {/each}
  {/if}
</main>

<style>
 
</style>
