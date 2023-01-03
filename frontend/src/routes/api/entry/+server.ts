import {error, json} from '@sveltejs/kit';
import {GetEntry} from "$lib/server/DB";
import type { RequestHandler } from '@sveltejs/kit';
export const GET = (async ({url}) => {
    const entryId = url.searchParams.get("entryId")
    if (entryId == null) {
        throw error(404, "There is no entryId given.")
    }
    const entry = await GetEntry(entryId)
    if (entry == null) {
        throw error(404, "There is no entry with this ID in our database.")
    }
    return json(entry);
}) satisfies RequestHandler