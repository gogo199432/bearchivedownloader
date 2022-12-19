import {GetEntry} from "$lib/DB";
import type { LoadEvent } from '@sveltejs/kit';
import {error} from "@sveltejs/kit";
export async function load({params}: LoadEvent) {
    if (params.entryId == null) {
        throw error(404, "There is no entry with this ID in our database.")
    }
    const entry = await GetEntry(params.entryId)
    if(entry == null){
        throw error(404, "There is no entry with this ID in our database.")
    }
    return {
        entry
    };
}