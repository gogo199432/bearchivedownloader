<script lang="ts">
    import {Entry} from "$lib/EntryType";
    export let entryId;
    async function LoadData(){
        const response = await fetch(`/api/entry/?entryId=${entryId}`);
        return await response.json()
    }

    let promise : Promise<Entry> = LoadData();
</script>
{#await promise}
    <h1>Loading...</h1>
{:then entry}
    <h1>{entry.Title}</h1>
    <h4>Tags:</h4>
    <ul>
        {#each entry.Tags as tag}
            <li>{tag}</li>
        {/each}
    </ul>
    <p>{@html entry.Text}</p>
    {#if entry.Choices != null}
    <p>Choices:</p>
    <ol style="margin-top: 1rem">
        {#each Object.entries(entry.Choices) as choice}
            <li><a href="/entry/{choice[1]}">{choice[0]}</a></li>
            {/each}
    </ol>
    {/if}
    <hr>
    <address>{entry.Author}</address>
    <p>{entry.Date}</p>
{:catch error}
    <p style="color: red">Welp...: {error.message}</p>
{/await}