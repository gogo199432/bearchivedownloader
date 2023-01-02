<script lang="ts">
    export let data;
</script>

<h1>Welcome to BEArchive Rescue!</h1>

<p>This project is still very early, the UI is in a "just enough to be working" state.</p>
<p>Any contributions are welcome <a href="https://github.com/gogo199432/bearchivedownloader">over here</a></p>

{#if data.root !== null}
    <h2>To go to the beginning of the adventure click here: </h2>
    <a href="/entrylist/{data.root}"><b>Adventure starts here</b></a>
{:else}
    <p>There is no root configured. You can find the root in the Neo4J DB and use its ID to add this to the end of the url: /entrylist/whateveridyoufound</p>
{/if}
<h2>Or you can find entries based on tags here: </h2>
<a href="/find">FIND</a>