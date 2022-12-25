<script lang="ts">
    import Entry from "$lib/Entry.svelte"
    import { page } from '$app/stores';
    import {onMount} from "svelte";

    let entries : string[] = []
    onMount(function(){
        console.log($page.params["rootEntry"])
        entries = [$page.params["rootEntry"],...entries]
    })
</script>

{#each entries as entry}
    <Entry entryId={entry}></Entry>
{/each}