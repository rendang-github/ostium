<script>
import { getLogin, postLogin, deleteLogin } from "../js/ApiOstium";

export let user;
let username = "";
let password = "";
let waiting = false;
let failed = false;

const logout = () => {
    deleteLogin().then(rest => {
        user = null;
        username = "";
        password = "";
    });
}

const handlePasswordKey = () => {
    if (event.code == 'Enter') {
        event.preventDefault();
        waiting = true;
        postLogin(username, password).then(res => {
            waiting = false;
            if (! res) {
                failed = true;
            } else {
                user = res;
            }
        });
        return false;
    }
}

// Check to see if we need to initialize the user object
if (user == null) {
    // We need to see if we're logged in
    getLogin().then(res => {
        if (res) {
            user = res
        }
    });
}
</script>

<div>
{#if user != null}
    Hello {user.name} <button on:click={logout}>Log Out</button>
{:else if waiting == true}
    spinner
{:else}
    <input type="text" bind:value={username} placeholder="Enter Login Email" />
    <input type="password" bind:value={password} placeholder="Enter Password" on:keydown={handlePasswordKey}/>
    {#if failed == true}
        You failed
    {/if}
{/if}
</div>

<style>
    div {
        display: inline-block;
    }
</style>
