<script>
import ObjectSelector from './ObjectSelector.svelte';
import Mode from "../js/Mode";
import Login from './Login.svelte';
export let user;
export let mode;
export let currentCampaign;
export let currentLayout;
export let currentSnippet;
export let currentResource;

function onNavigate(event) {
    console.log("Navigate", event);
    mode = event.detail.mode;
    switch (mode) {
        case Mode.Campaign:
            currentCampaign = event.detail.obj;
            break;
        case Mode.Layout:
            currentLayout = event.detail.obj;
            break;
        case Mode.Snippet:
            currentSnippet = event.detail.obj;
            break;
        case Mode.Resource:
            currentResource = event.detail.obj;
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
    <ObjectSelector collection={"layout"} on:navigate={onNavigate} mode={Mode.Layout} let:obj={obj}>
        <span slot="label">Layouts</span>
        {#if obj}
            <div><strong>{obj.name}</strong></div>
            <div>{obj.description}</div>
            <div><em>{obj.modified.toLocaleString()}</em></div>
        {:else}
            + Add new layout
        {/if}
    </ObjectSelector>
    <ObjectSelector collection={"snippet"} on:navigate={onNavigate} mode={Mode.Snippet} let:obj={obj}>
        <span slot="label">Snippets</span>
        {#if obj}
            <div><strong>{obj.name}</strong></div>
            <div>{obj.description}</div>
            <div><em>{obj.modified.toLocaleString()}</em></div>
        {:else}
            + Add new snippet
        {/if}
    </ObjectSelector>
    <ObjectSelector collection={"resource"} on:navigate={onNavigate} mode={Mode.Resource} let:obj={obj}>
        <span slot="label">Resources</span>
        {#if obj}
            <div><strong>{obj.name}</strong></div>
            <div>{obj.description}</div>
            <div><em>{obj.modified.toLocaleString()}</em></div>
        {:else}
            + Add new resource
        {/if}
    </ObjectSelector>
    <Login bind:user />
</div>
