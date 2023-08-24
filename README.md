# Developer Technical Test

This partial application was created to present to a technical team for review.</p>
Some features of the software;

- Filter and sort the list of jobs
- View a jobs details and add/edit notes for selected job
- Change the status of a job

The application is a standalone demo requiring no additional WAMP/LAMP dependencies. It has been built and tested on Windows, standalone Ubuntu Linux and WSL:Ubuntu. If the application requires rebuilding, there are two batch files provided. The rebuild assumes there is an existing Go installation on the host computer. The provided batch files will download the required 3rd party packages for the build process.

## Building
This application uses the Go programming language - where the latest was [Go 1.21](https://go.dev/dl/) as of writing this application. If you do not have Go installed on your system, you can acquire a copy from [Go.dev](https://go.dev/dl/). The go1.21.0.windows-amd64.msi was used to build this application.

To run the server on your Windows system:

1. Run `buildpkg.cmd` in the jobs.tradie.pg folder to build the binary (`jobs.tradie.pg`) using non vendored packages
1. Run `buildvendor.cmd` in the jobs.tradie.pg folder to build the binary (`jobs.tradie.pg`) with the vendor
1. Run the binary `jobs.tradie.pg.exe`
1. Browse to [http://localhost](http://localhost) (the port can be change in the code of 80 is not working) to test the application out.

### Non Windows
Testing has been performed on WSL & Linux but not MacOS. However, the commands in buildpkg.cmd and buildvendor.cmd can be run manually to build and run this demo.

#### Build by pkg

``` bash
export GO111MODULE="on"
export GOFLAGS="-mod=mod"
go mod download
:: strip debug info during build
go build -ldflags="-s -w" .

``` 
#### Build by vendor

``` bash
export GO111MODULE="on"
export GOFLAGS="-mod=vendor"
go mod vendor
:: strip debug info during build
go build -ldflags="-s -w" 
```

### Dependencies
The application uses the following Go packages to build;

- [JSON filed data validator](github.com/gookit/validate)
- [Datastore: PostgreSQL driver](https://github.com/jackc/pgx/)
- [HTTP router: Gorilla mux](https://github.com/gorilla/mux)

## Datastore

This version application requires a separate database to function - PostgreSQL. A JSON file is imported from the local folder. This will be imported when the application is run for the first time. Thereafter the application will use the database each time it is executed.

## Sample screens

![Job updates page](/jobslisting.jpg)
![Job updates page](/noteedit.jpg)
