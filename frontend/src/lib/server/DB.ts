import * as neo4j from 'neo4j-driver'
import type {Entry} from "$lib/EntryType";

let db: neo4j.Driver;
process.on('SIGINT',CloseDB)
process.on('SIGTERM',CloseDB)

function GetSession() : neo4j.Session{
    if(db == null){
        if(process.env.DBPASSWORD){
            db = neo4j.driver(process.env.NEO4J || "neo4j://localhost:7687", neo4j.auth.basic(process.env.DBUSERNAME ?? "neo4j",process.env.DBPASSWORD))
        }else{
            db = neo4j.driver(process.env.NEO4J || "neo4j://localhost:7687")
        }
    }
    return db.session()
}

export function CloseDB(){
    db?.close()
}

export async function GetRootId() : Promise<string>{
    // const id = 'some value'; // To make es-lint happy
    if (process.env.ROOT) {
        return process.env.ROOT;
    }else{
        return ""
    }
    // Doesn't actually work properly, dunno why. Probably some issue with the scraper
    // the DB seems to have multiple nodes without incoming edges
    /*const session = GetSession();
    const result = await session.run(
        'MATCH (entry) WHERE NOT (entry)<-[]-() RETURN entry.Id as Id'
    );

    await session.close();
    if(result.records.length == 0){
        throw error(404,"No root found")
    }
    result.records.forEach(x=>{
        id = x.get("Id")
    })
    return id*/
}

export async function GetEntry(id: string) : Promise<Entry | null> {
    const session = GetSession();
    const result = await session.run('MATCH (entry {Id: $idParam}) RETURN properties(entry) as props, labels(entry) as labels , entry.Id as Id', {
        idParam: id
    })
    const entry = await ParseToEntry(session,result)

    await session.close()
    return entry
}

export async function FindEntriesWithFilter(filter:string) : Promise<{id: string,title: string, labels: string}[]>{
    const session = GetSession();
    const result = await session.run('MATCH (entry) WHERE '+filter+' RETURN entry.Id as Id, entry.Title as Title, labels(entry) as Labels')
    await session.close()
    const entries :{id: string,title: string, labels:string}[] = []
    result.records.forEach(record =>{
        const labels = record.get("Labels").filter((x:string) => x !== "Entry").join(" ")
        entries.push({id: record.get("Id"), title: record.get("Title"), labels})
    })
    return entries
}

async function ParseToEntry(session: neo4j.Session, result: neo4j.QueryResult) : Promise<Entry | null>{
    let entry: Entry | null = null

    let id = "some value" // To make es-lint happy

    result.records.forEach(record => {
        const e = record.get('props');
        delete e.ChildrenURLs
        delete e.Children
        e.Date = e.Date.toStandardDate()
        e.Tags = record.get('labels').filter((x: string) => x !== 'Entry')
        id = record.get("Id")
        entry = e
    })

    result = await session.run('MATCH (entry {Id: $idParam})-[choice]->(child) RETURN choice.text as ch,child.Id as chId', {
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
    return entry
}