<script>
    import { onMount } from "svelte";
    import { GetDiskUsage } from "../wailsjs/go/main/App.js";

    let diskData = {
        "Total": 0,
        "Used": 0,
        "Free": 0,
        "Percent used": 0,
    }

    let error = null;
    
    async function fetchDiskUsage() {
        try {
            const result = await GetDiskUsage()
            diskData= {
                "Total": result["Total"],
                "Used": result["Used"],
                "Free": result["Free"],
                "Percent used": result["Percent used"],
            }
            error = null;
        }
        catch (err) {
            console.error("Error fetching disk usage", err);
            error = err.message;
        }
    }

    function toGB(bytes) {
        return (bytes / (1024 * 1024 * 1024)).toFixed(2);
    }
    
    onMount(() => {
        fetchDiskUsage();
        const interval = setInterval(fetchDiskUsage, 1000);
        return () => clearInterval(interval);
    })

</script>

<main>
    <h1>Disk usage</h1>
    <p>Total: {diskData.Total}</p>
    <p>Used: {diskData.Used}</p>
    <p>Free: {diskData.Free}</p>
    <p>Percent used: {diskData["Percent used"]}</p>
</main>

<style>

</style>