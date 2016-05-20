# Neo4j Bolt Driver for Go

This repository aims at becoming a Go database driver for Neo4j 3.0.0+. For easier picking up, tasks will be defined via GitHub issues. Please submit a pull request if you are confident that you solved an issue, ideally including some testing.

To get in touch, please contact [Sanja Bonic](mailto:sanja@go2neo.org) or [Nigel Small](mailto:nigel@go2neo.org).

## How To Run The Server
```
$ wget http://dist.neo4j.org/neo4j-community-3.0.0-unix.tar.gz
$ tar xf neo4j-community-3.0.0-unix.tar.gz
$ cd neo4j-community-3.0.0
$ bin/neo4j start
Starting Neo4j.
Started neo4j (pid 3456). By default, it is available at http://localhost:7474/
There may be a short delay until the server is ready.
See /home/nigel/opt/neo4j-community-3.0.0/logs/neo4j.log for current status.
```

Note that there will generally be a delay (10-15 seconds) before the server port becomes available.

## How To Run The Client
In general, for quick development testing, you can navigate to the specific folder in your terminal and test the file via `go run <filename>.go`. Assuming your `GOPATH` is set, you can install custom packages by executing `go install <packagename>`.

For example, to run `main.go` from a shell, you might use:

```
$ GOPATH=~/work/go2neo/ go run main.go
```

We will provide more information on the actual build process as this repository matures.

## Example Usage
Tbd.

## Work In Progress
The basic component for message serialisation needed now is [packstream.go](src/github.com/nigelsmall/go2neo/packstream/packstream.go). Once that is working, [socketclient.go](src/github.com/nigelsmall/go2neo/bolt/socketclient.go) is the next component that actually connects to Neo4j.

## Coming Next
Please refer to the [Issues](https://github.com/nigelsmall/go2neo/issues) for a clear list of what's next in line. (ETA: May 31st, 2016.)

### packstream.go
The most important feature of packstream.go is to make it as generic and reusable as possible. Please create new files for everything hard-coded that is specific, such as switch statements to recognize the type of byte sequences.

* [] accept byte stream
* [] stream into and decode from buffer
* [] create packer and unpacker functions

A more detailed list with explanations will be added by the end of May.

### socketclient.go
* [x] TCP socket for the Bolt protocol
* [] connect and do handshake
* [] Bolt mode initiated
