<script>
    import { onMount } from 'svelte';
    import { GetProcesses } from '../wailsjs/go/main/App.js'

    let processData = {}  
    let error = null

    async function fetchProcessData() {
        try {
            const result = await GetProcesses()
            if (!result) {
                throw new Error('No process data received from backend')
            }
            
            processData = result

            error = null
        }
        catch (err) {
            console.error('Error fetching process data', err)
            error = err.message
        }
    }

    onMount(() => {
        fetchProcessData()
        const interval = setInterval(fetchProcessData, 1000)
        return () => clearInterval(interval)
    })

    function toGB(bytes) {
        if (typeof bytes !== 'number') {
            console.error('Invalid bytes value:', bytes)
            return '0.00'
        }
        return (bytes / (1024 * 1024 * 1024)).toFixed(2)
    }
</script>

<main>
    <h1>Process Monitor</h1>
    {#if error}
        <p style="color: red;">{error}</p>
    {:else}
        {#each Object.entries(processData) as [pid, process]}
            <div class="process-info">
                <h3>{process["Name"]}</h3>
                <small>{process["Cwd"]}</small>
                <p>PID: {process["PID"]}</p>
                <p>Memory: {toGB(process["MemoryInfo"]["rss"])} GB</p>
                <p>CPU: {process.cpu}%</p>
            </div>
        {/each}
    {/if}
</main>

<style>

</style>