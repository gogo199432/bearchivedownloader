<script lang="ts">
    import Entry from "$lib/Entry.svelte"
    import {page} from '$app/stores';
    import {onMount} from "svelte";
    import {fade, fly} from 'svelte/transition';


    let entries : string[] = []
    onMount(function(){
        entries = [$page.params["rootEntry"],...entries]
    })

    function choiceHandler(event){
        entries = [...entries,event.detail.choice]
    }

    function backHandler(event){
        entries.length = entries.lastIndexOf(event.detail.choice)
    }
</script>

{#each entries as entry}
    <div in:fly="{{ y: -50, duration: 1000 }}" out:fade>
        <Entry entryId={entry} on:chosen={choiceHandler} on:back={backHandler} ></Entry>
        <hr>
        <hr>
    </div>
{/each}