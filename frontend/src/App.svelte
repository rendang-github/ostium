<script>
import Mode from "./js/Mode";
import TopLevelNav from './Nav/TopLevelNav.svelte';
import Campaign from './Campaign/Campaign.svelte';
export let user = null;
export let mode = Mode.Campaign;
export let currentCampaign = null;
export let currentLayout = null;
export let currentSnippet = null;
export let currentResource = null;
</script>

<main>
    <TopLevelNav bind:user bind:mode bind:currentCampaign bind:currentLayout bind:currentSnippet bind:currentResource />
    {#if mode == Mode.Campaign}
        <Campaign bind:user bind:campaign={currentCampaign} />
    {:else if mode == Mode.Layout}
        <h1>Layout Mode</h1>
    {:else if mode == Mode.Snippet}
        <h1>Snippet Mode</h1>
    {:else if mode == Mode.Resource}
        <h1>Resource Mode</h1>
    {/if}
</main>

<style>
</style>
