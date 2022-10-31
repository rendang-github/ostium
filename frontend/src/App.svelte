<script>
import Mode from "./js/Mode";
import TopLevelNav from './Nav/TopLevelNav.svelte';
export let user = null;
export let mode = Mode.Campaign;
export let currentCampaign = null;
export let currentTheme = null;
export let currentSnippet = null;
</script>

<main>
    <TopLevelNav bind:user bind:mode bind:currentCampaign bind:currentTheme bind:currentSnippet />
    {#if mode == Mode.Campaign}
        <h1>Campaign Mode</h1>
    {:else if mode == Mode.Theme}
        <h1>Theme Mode</h1>
    {:else if mode == Mode.Snippet}
        <h1>Snippet Mode</h1>
    {/if}
</main>

<style>
</style>
