<script>
    import { onMount } from 'svelte';
    import { GetCPU } from '../wailsjs/go/main/App.js';

    let cpuData = {
        Percentage: [], 
        Cores: 0,
        GeneralInfo: [],
        ModelName: "",
    }
    let temperature = {}
    let avgFrequency = 0
    
    let error = null
    async function fetchCPUData() {
        try {
            const result = await GetCPU()
            cpuData = {
                Percentage: result.CPUPercent,
                Cores: result.CPUCores,
                GeneralInfo: result.CPUInfo,
                ModelName: result.CPUInfo[0].modelName
            }

            temperature = result.CPUTemperature

            const sum = result.CPUFrequency.reduce((a, b) => a + b, 0)
            avgFrequency = (sum / result.CPUFrequency.length) / 1000

            error = null
        }
        catch (err) {
            console.error('Error fetching CPU data', err)
            error = err.message
        }
    }

    onMount(() => {
        fetchCPUData()
        const interval = setInterval(fetchCPUData, 1000)
        return () => clearInterval(interval)
    })

    function getTemperatureClass(temperature) {
        if (temperature > 70) return 'temp-red'
        if (temperature > 60) return 'temp-orange'
        if (temperature > 55) return 'temp-yellow'
        return 'temp-blue'
    }
</script>

<main>
    <h1>System CPU Monitor</h1>
    {#if error}
        <p style="color: red;">Error: {error}</p>
    {:else}
        <ul>
            <li>Model: {cpuData.ModelName}</li>
            <li>Cores: {cpuData.Cores}</li>
            <li>Average Frequency: {avgFrequency.toFixed(2)} GHz</li>
            {#if cpuData.Percentage.length}
                <li>CPU Usage (Total): {(cpuData.Percentage.reduce((a, b) => a + b, 0) / cpuData.Percentage.length).toFixed(2)}%</li>
                {#each cpuData.Percentage as percent, i}
                    {#if temperature[i] > 0}
                        <li>Core {i}: {percent.toFixed(2)}% - 
                            <span class="temperature {getTemperatureClass(temperature[i])}">{temperature[i]}</span> ºC
                        </li>
                    {:else}
                        <li>Core {i}: {percent.toFixed(2)}% - 
                            <span class="temperature {getTemperatureClass(temperature[10])}">{temperature[10]}</span> ºC
                        </li>
                    {/if}
                {/each}
            {/if}
        </ul>
    {/if}
</main>
<style>
    ul {
        list-style: none;
        padding: 0;
    }
    li {
        margin: 10px 0;
    }

    .temperature { font-weight: bold; }
    .temp-blue { color: rgb(92, 92, 255); }
    .temp-yellow { color: rgb(254, 254, 70);}
    .temp-orange { color: rgb(213, 139, 0); }
    .temp-red { color: rgb(222, 23, 12); }
</style>