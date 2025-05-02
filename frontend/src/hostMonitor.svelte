<script>
    import { onMount } from 'svelte';
    import { GetHostInfo } from '../wailsjs/go/main/App.js';

    let hostData = {
        hostName: "",
        uptime: "",
        OS: "",
        kernel: "",
        platform: "",
    }
    let error = null
    
    async function fetchHostData() {
        try {
            const result = await GetHostInfo()
            if (!result) {
                throw new Error('No host data received from backend')
            }

            hostData ={
                hostName: result.MajorInfo.hostname,
                uptime: result.Uptime,
                OS: result.MajorInfo.os,
                kernel: result.MajorInfo.kernelArch,
                platform: result.MajorInfo.platform,
            }
            error = null 
        }
        catch (err) {
            console.error('Error fetching host data', err)
            error = err.message
        }
    }

    onMount(() => {
        fetchHostData()
        const interval = setInterval(fetchHostData, 1000)
        return () => clearInterval(interval)
    })
</script>

<main>
    <h1>Host Monitor</h1>
    {#if error}
        <p style="color: red;">Error: {error}</p>
    {:else}
        <p>Hostname: {hostData.hostName}</p>
        <p>Uptime: {hostData.uptime}</p>
        <p>OS: {hostData.OS}</p>
        <p>Kernel: {hostData.kernel}</p>
        <p>Platform: {hostData.platform}</p>
    {/if}
</main>

<style>

</style>