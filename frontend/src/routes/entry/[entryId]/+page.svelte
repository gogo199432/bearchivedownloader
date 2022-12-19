<script lang="ts">
    export let data;
</script>
{#if data.entry != null}
<h1>{data.entry.Title}</h1>
<h4>Tags:</h4>
<ul>
    {#each data.entry.Tags as tag}
        <li>{tag}</li>
    {/each}
</ul>
<p>{@html data.entry.Text}</p>
{#if data.entry.Choices != null}
<p>Choices:</p>
<ol style="margin-top: 1rem">
    {#each Object.entries(data.entry.Choices) as choice}
        <li><a href="/entry/{choice[1]}">{choice[0]}</a></li>
        {/each}
</ol>
{/if}
<hr>
<address>{data.entry.Author}</address>
<p>{data.entry.Date.toDateString()}</p>
{:else}
<p>Welp. You have reached the end of the universe.</p>
{/if}