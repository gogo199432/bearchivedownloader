import {GetRootId} from "$lib/DB";
import {error} from "@sveltejs/kit";

export async function load() {
    const root = await GetRootId();
    if(root == null){
        throw error(404, "There is no root in our database???")
    }
    return {
        root
    };
}