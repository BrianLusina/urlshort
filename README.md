# Urlshort

A simple url shortener utility written in Go.

You will require Go 1.18. to run this application. You can visit [here](https://golang.org/doc/install) for more information.

## Setup

Minimum setup is really required for this as it is a small application after all. However there is a [Makefile](./Makefile) that provides convenient scripts
to make life a little bit easier.

### Dependency installation

Simply run `make install` to install all dependencies. Based on your go installation, this will install all the dependencies required for this project, which is quite minimal.

### Running the application

You can run the application with `make run` and this will start the application on port 3000 and you can access it and add urls to shorten on http://localhost:3000/add. This will display a simple form to add a url and will respond with an id which you can then append to http://localhost:3000/<id> to get redirected to the longer url(note that the <id> here is the shortened id from the add operation)

You can optionally run the application passing in flags as you see fit. To view these flags run the following command:

``` bash
> go run app/cmd/main.go --help
-file string
    data store file name (default "store.json")
-host string
    host name and port (default "localhost")
-http string
    http listen address (default ":3000")
-primary string
    RPC address of the primary store
-rpc
    enable rpc server
```

> The output will be as above.

The default arguments(flags) are in the brackets

#### Primary and Replica applications

If you are keen, from the `--help` flag passed in to the application, you will have noticed that there are `--primary` and `--rpc` flags available. These are used to tell the application the primary server address and the `--rpc` is basically set to allow the application to run as a rpc server, allowing the primary server to handle writes, while the replicas to handle reads.

You can run multiple replicas, just ensure that the http address is different for each replica and ensure that the primary is run before the replicas are started.

```bash
go run app/cmd/main.go --http :5000
```

> This will run the primary on port 5000

While the replicas can run on other ports like below:

```bash
go run app/cmd/main.go --primary :5000 --rpc true --http :5001
go run app/cmd/main.go --primary :5000 --rpc true --http :5002
go run app/cmd/main.go --primary :5000 --rpc true --http :5003
go run app/cmd/main.go --primary :5000 --rpc true --http :5004
....(etc)
```

> of course you can either run them in the background or in separate terminal sessions.

The `--primary` flag tells the replicas where the primary is so that they can communicate with it.

### Storage

A simple note about how storage is handled.

Note that storage is in a file and is persisted across restarts as it's stored on the host file system and in this case will be stored in the root of the project directory. The `--file` flag tells the application which file to use. It defaults to `store.json`, but this can be changed to any file name you wish.

Minimal refactoring can be setup to enable a more full fledged solution to store data in a database or another storage solution. Only thing that would need to change is in the [repo package](./app/internal/repo/) specifying how the storage will be handled, either to a file system or a storage system.
