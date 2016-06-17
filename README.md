[![Build Status](https://travis-ci.org/nigelsmall/go2neo.svg?branch=master)](https://travis-ci.org/nigelsmall/go2neo)
[![Coverage Status](https://coveralls.io/repos/github/nigelsmall/go2neo/badge.svg?branch=master)](https://coveralls.io/github/nigelsmall/go2neo?branch=master)

NOTE: This software is *not* officially supported by Neo Technology.

# Neo4j Bolt Driver for Go

This repository aims at becoming a Go database driver for Neo4j 3.0.0+. For easier picking up, tasks will be defined via GitHub issues. Please submit a pull request if you are confident that you solved an issue, ideally including some testing.

To get in touch, please contact [Nigel Small](mailto:nigel@neo4j.com) or [Sanja Bonic](mailto:sanja@cv2.me).

## How To Run The Server
```
$ wget http://dist.neo4j.org/neo4j-community-3.0.2-unix.tar.gz
$ tar xf neo4j-community-3.0.2-unix.tar.gz
$ cd neo4j-community-3.0.2
$ bin/neo4j start
Starting Neo4j.
Started neo4j (pid 3456). By default, it is available at http://localhost:7474/
There may be a short delay until the server is ready.
See /home/nigel/opt/neo4j-community-3.0.2/logs/neo4j.log for current status.
```

Note that there will generally be a delay (10-15 seconds) before the server port becomes available.

## How To Run The Tests
Assuming your `GOPATH` is set, you can navigate to the specific folder in your terminal and test the file via `go test`. A running server is required.
