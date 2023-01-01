import type {Actions, RequestEvent} from '@sveltejs/kit';
import {tags} from "$lib/EntryLabels";
import {fail} from "@sveltejs/kit";
import {FindEntriesWithFilter} from "$lib/server/DB";

export const actions: Actions = {
    default: async (event: RequestEvent) => {
        const allData = await event.request.formData()
        const queryData = allData.get('query')

        if(queryData === null){
            return
        }
        let query = queryData as string
        query = query.replaceAll("("," ( ")
        query = query.replaceAll(")"," ) ")
        const operators = ["AND","OR","XOR","NOT","(",")"]
        const allowedWords = [...tags, ...operators]
        let failed = false
        const reconstructed: string[] = []
        query.split(" ").forEach(x=> {
            if(x.trim().length === 0){
                return
            }
            if (!allowedWords.includes(x)) {
                failed = true
            }
            if(!operators.includes(x)){
                reconstructed.push("entry:"+x)
            }else{
                reconstructed.push(x)
            }
        })

        if(failed){
            return fail(422,{query,incorrect: true})
        }
        const found = await FindEntriesWithFilter(reconstructed.join(" "))
        return { success: true, entries: found};
    }
};