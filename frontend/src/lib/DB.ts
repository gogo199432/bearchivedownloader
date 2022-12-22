import * as neo4j from 'neo4j-driver'
import type {Entry} from "$lib/EntryType";

let db: neo4j.Driver;
process.on('SIGINT',CloseDB)
process.on('SIGTERM',CloseDB)

function GetSession() : neo4j.Session{
    if(db == null){
        db = neo4j.driver(process.env.NEO4J || "neo4j://localhost:7687")
    }
    return db.session()
}

export function CloseDB(){
    db.close()
}

export async function GetRootId() : Promise<string>{
    let id = 'some value'; // To make es-lint happy
		if (process.env.ROOT) {
			return process.env.ROOT;
		}
		const session = GetSession();
		const result = await session.run(
			'MATCH (entry:Entry) WHERE NOT (entry)<-[]-(:Entry) RETURN entry.Id as Id'
		);

		await session.close();
    result.records.forEach(x=>{
        id = x.get("Id")
    })
    return id
}

export async function GetEntry(id: string) : Promise<Entry | null> {
    const session = GetSession();
    const result = await session.run('MATCH (entry:Entry {Id: $idParam}) RETURN properties(entry) as props, labels(entry) as labels , entry.Id as Id', {
        idParam: id
    })

    const entry = await ParseToEntry(session,result)

    await session.close()
    return entry
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
    return entry
}