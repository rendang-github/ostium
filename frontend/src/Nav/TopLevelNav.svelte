<script>
import ObjectSelector from './ObjectSelector.svelte';
import Mode from "../js/Mode";
import Login from './Login.svelte';
export let user;
export let mode;
export let currentCampaign;
export let currentTheme;
export let currentSnippet;

function onNavigate(event) {
    console.log("Navigate", event);
    mode = event.detail.mode;
    switch (mode) {
        case Mode.Campaign:
            currentCampaign = event.detail.obj;
            break;
        case Mode.Theme:
            currentTheme = event.detail.obj;
            break;
        case Mode.Snippet:
            currentSnippet = event.detail.obj;
            break;
    }
}

</script>

<div>
    <ObjectSelector collection={"campaign"} on:navigate={onNavigate} mode={Mode.Campaign} let:obj={obj}>
        <span slot="label">Campaigns</span>
        {#if obj}
            <div><strong>{obj.name}</strong></div>
            <div>{obj.description}</div>
            <div><em>{obj.modified.toLocaleString()}</em></div>
        {:else}
            + Add new campaign
        {/if}
    </ObjectSelector>
    <ObjectSelector collection={"theme"} on:navigate={onNavigate} mode={Mode.Theme}>Froom</ObjectSelector>
    <ObjectSelector collection={"snippet"} on:navigate={onNavigate} mode={Mode.Snippet}>Froom</ObjectSelector>
    <Login bind:user />
</div>
