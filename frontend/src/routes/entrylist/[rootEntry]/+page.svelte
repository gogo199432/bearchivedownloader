<script lang="ts">
    import Entry from "$lib/Entry.svelte"
    import { page } from '$app/stores';
    import {onMount} from "svelte";

    let entries : string[] = []
    onMount(function(){
        console.log($page.params["rootEntry"])
        entries = [$page.params["rootEntry"],...entries]
    })

    function choiceHandler(event){
        entries = [...entries,event.detail.choice]
    }

    function backHandler(event){
        const border = entries.indexOf(event.detail.choice)
        entries.length = border
    }
</script>

{#each entries as entry}
    <Entry entryId={entry} on:chosen={choiceHandler} on:back={backHandler}></Entry>
    <hr>
    <hr>
{/each}