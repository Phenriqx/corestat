<script>
    import { onMount } from "svelte";
    import { GetRAM } from "../wailsjs/go/main/App.js";

    let ram = {
        "Total memory": 0,
        "Used memory": 0,
        "Cached memory": 0,
        "Free memory": 0,
        "Percent used": 0,
    }

    let error = null;
    async function fetchRAM() {
        try {
            const result = await GetRAM()
            ram = {
                "Total memory": result["Total memory"],
                "Used memory": result["Used memory"],
                "Cached memory": result["Cached memory"],
                "Free memory": result["Free memory"],
                "Percent used": result["Percent used"],
            }
            error = null;
        }
        catch (err) {
            console.error("Error fetching RAM", err);
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
            <li>Used memory: {toGB(ram["Used memory"])} GB</li>
            <li>Cached memory: {toGB(ram["Cached memory"])} GB</li>
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