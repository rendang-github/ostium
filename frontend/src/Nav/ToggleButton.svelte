<script>
    export let on = false;
</script>
<button on:click>
    {#if on == true}
        <slot name="on">On</slot>
    {:else}
        <slot name="off">Off</slot>
    {/if}
</button>
