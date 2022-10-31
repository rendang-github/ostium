<script>
import { fly } from 'svelte/transition';
import { getObjects } from "../js/ApiOstium";
import { createEventDispatcher } from 'svelte';

const dispatch = createEventDispatcher();
export let collection;
export let mode;
let open = false;
let loading = false;
let objs = [];

const toggle = () => {
    open = !open;
    if (open == true) {
        console.log("Loading");
        if (objs.length == 0) {
            loading = true;
        }
        getObjects(collection).then(res => {
            if (res != null) {
                objs = res;
            } else {
                objs = [];
            }
            console.log("Got objects", objs);

            // Convert timestamps
            for (let index = 0; index < objs.length; index++) {
                objs[index].modified = new Date(objs[index].modified);
            }
            objs.sort(function(a,b){
                return b.modified - a.modified;
            });
            loading = false;
        });
    }
}

const selector = (obj) => {
    console.log("Selector triggered", obj);
    // Issue event
    dispatch('navigate', {
        mode: mode,
        collection: collection,
        obj: obj
    });

    // Close list
    open = false;
}

</script>

<section>
    <button title="{collection}" on:click={toggle}><slot name="label">Default Label</slot></button>
    {#if open == true}
    <div transition:fly class="content">
        <button class="menuitem" title="New {collection}" on:click={e => selector(null)}><slot prop={null}></slot></button>
        {#if loading == true}
            <span class="menuitem">Loading {collection}...</span>
        {:else}
            {#each objs as obj}
                <button class="menuitem" title="{obj.name}" on:click={e => selector(obj)}><slot obj={obj}></slot></button>
            {/each}
        {/if}
    </div>
    {/if}
</section>

<style>
    section {
        position: relative;
        display: inline-block;
    }

    .content {
        position: absolute;
        background-color: #f6f6f6;
        min-width: 230px;
        border: 1px solid #ddd;
        z-index: 1;
    }

    button {
        color: black;
        background-color: lightgrey;
        padding: 12px 16px;
        margin: 0px;
        text-decoration: none;
        display: block;
    }

    button:hover {
        background-color: grey;
    }

    .menuitem {
        width: 100%;
    }
</style>
