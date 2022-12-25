import {GetRootId} from "$lib/server/DB";
import {error} from "@sveltejs/kit";

export async function load() {
    const root = await GetRootId();
    if(root == null){
        throw error(404, "There is no root in our database???")
    }
    console.log("Root: "+root)
    return {
        root
    };
}