import * as neo4j from 'neo4j-driver'
import type {Entry} from "$lib/EntryType";

let db: neo4j.Driver;
process.on('SIGINT',CloseDB)
process.on('SIGTERM',CloseDB)

function GetSession() : neo4j.Session{
    if(db == null){
        db = neo4j.driver("neo4j://localhost:7687")
    }
    return db.session()
}

export function CloseDB(){
    db.close()
}

export async function GetEntry(id: string) : Promise<Entry | null> {
    const session = GetSession();
    let entry: Entry | null = null
    let result = await session.run('MATCH (entry:Entry {Id: $idParam}) RETURN properties(entry) as props, labels(entry) as labels', {
        idParam: id
    })

    result.records.forEach(record => {
        const e = record.get('props');
        delete e.ChildrenURLs
        delete e.Children
        e.Date = e.Date.toStandardDate()
        e.Tags = record.get('labels').filter((x: string) => x !== 'Entry')
        entry = e
    })

    result = await session.run('MATCH (entry:Entry {Id: $idParam})-[choice]->(child) RETURN choice.text as ch,child.Id as chId', {
        idParam: id
    })
    result.records.forEach(record => {
        if(entry != null) {
            if(entry.Choices == null){
                entry.Choices = {}
            }
            entry.Choices[record.get("ch")] = record.get("chId")
        }
    })

    await session.close()
    return entry
}