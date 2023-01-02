import {GetRootId} from "$lib/server/DB";

export async function load() {
    const root = await GetRootId();
    if(root == ""){
        console.warn("No root found!")
        return {
            root: null
        }
    }
    console.log("Root: "+root)
    return {
        root
    };
}