<script lang="ts">
    import {Entry} from "$lib/EntryType";
    import { createEventDispatcher } from 'svelte'


    export let entryId;
    async function LoadData(){
        const response = await fetch(`/api/entry/?entryId=${entryId}`);
        return await response.json()
    }

    let promise : Promise<Entry> = LoadData();

    const dispatcher = createEventDispatcher()
    let choiceMade = false
    let chosenId
    function choiceHandler(id: string){
        dispatcher('chosen',{choice: id})
        chosenId = id
        choiceMade = true
    }
    function backHandler(){
        dispatcher('back',{choice:chosenId})
        choiceMade = false
    }
</script>
{#await promise}
    <h1>Loading...</h1>
{:then entry}
    <h1>{entry.Title}</h1>
    <address style="display: inline">{entry.Author}</address>
    <b>-</b>
    <div style="display: inline">{entry.Date}</div>
    <h4>Tags:</h4>
    <ul>
        {#each entry.Tags as tag}
            <li>{tag}</li>
        {/each}
    </ul>
    <p>{@html entry.Text}</p>
    {#if entry.Choices != null}
        {#if !choiceMade}
            <h3>Choices:</h3>
            {#each Object.entries(entry.Choices) as choice}
                <button on:click={()=>choiceHandler(choice[1])}>{choice[0].substring(1)}</button>
            {/each}
        {:else}
            <button on:click={backHandler}>Back</button>
        {/if}
    {/if}
{:catch error}
    <p style="color: red">Welp...: {error.message}</p>
{/await}

<style>
    button {
        display: block;
        width: 100%;
        min-height: 3rem;
    }
</style>