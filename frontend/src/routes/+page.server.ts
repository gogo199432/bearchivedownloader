import {GetRootId} from "$lib/server/DB";

export async function load() {
    const root = await GetRootId();
    if(root == null){
        console.warn("No root found!")
        return {}
    }
    console.log("Root: "+root)
    return {
        root
    };
}