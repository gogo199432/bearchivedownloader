# BEArchiveDownloader

Goal of this repo is to crawl the branching story stored at https://addventure.bearchive.com/~addventure/

At the moment you can define the depth by setting the "MaxDepth" variable in the main.go file.

Steps to run:

- Have Neo4J running somewhere. I used docker, and deployed it with the following command:
  - `docker volume create neo4j`
  - `docker run -d --publish=7474:7474 --publish=7687:7687 --volume=neo4j:/data --env=NEO4J_AUTH=none neo4j`
- Run the app with `go run main.go` in the main folder
- Go to neo4j's browser to explore the DB. If deployed with Docker, go to http://localhost:7474/browser and log in without authentication
- Enter the following query into the input field `MATCH (n) RETURN n`

Things missing:
- App should be able to be stopped and pick up later. For this it should grab nodes that don't have outgoing connections yet and visit those instead of the starting link currently defined
- There is probably some way to multi-thread it, atm it is not async.
- Actual frontend to use the DB. While you can browse it by going to neo4j's browser window, it is not very good for exploring the story.

Pull Requests are welcome