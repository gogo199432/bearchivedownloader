# BEArchiveDownloader

Goal of this repo is to crawl the branching story stored at https://addventure.bearchive.com/~addventure/

## Steps to run:

- Have Neo4J running somewhere. I used docker, and deployed it with the following command:
  - `docker volume create neo4j`
  - `docker run -d --publish=7474:7474 --publish=7687:7687 --volume=neo4j:/data --env=NEO4J_AUTH=none neo4j`
- Run the app with `go run main.go` in the main folder
- Go to neo4j's browser to explore the DB. If deployed with Docker, go to http://localhost:7474/browser and log in without authentication
- Enter the following query into the input field `MATCH (n) RETURN n`
- Alternatively use the propery frontend by going into the frontend folder and running `pnpm dev`

## Configure:
You can configure the scraper in its folder using the `config.yaml` file. At the moment you can define the DB connection 
and how many children down the scraper should go in one go.

## Things missing:
- Proper css framework to make it pretty
- Search based on time, author etc
- Caching for query results

## Contributing:
Just open a Pull Request and describe exactly what your changes do.