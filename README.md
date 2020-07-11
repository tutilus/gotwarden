# gotwarden

Gotwarden is an unofficial Bitwarden minimalist API server written in Go based on the awesome work made by [jcs](https://github.com/jcs) on this project [rubywarden](https://github.com/jcs/rubywarden).

It provides a private backend for the open-source password management solutions [Bitwarden](https://github.com/bitwarden).

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

Go-lang.

Needed modules can be provide using `go mod`.

```sh
go mod download  
```

### Database

So far, only `sqlite3` is supported as database. 

* `sqlite3` 
* `postgres` [WIP]

### Installing

To create the `gotwarden` binary file nothing more than using make.

```sh
# Build source
make build
```

```sh
# run using gin debug mode
make run-debug
```

At the first running a new database is created in `fixtures/test.db` except `DB_PATH` variable is set.

## Running the tests

No test set up but it would be provided soon.

```sh
make test
```

## Deployment

### Local

```sh
# Production mode
make run 

# Or choose debug mode
make run-debug
```

By default, gotwarden serves on port 3000. It can be change into Makefile file  by overiding the PORT env variable.

### Docker

```sh

# Create docker container
make container

# Launch the container
make serve

```

## Parameters

Some variables can be set to adapt the behavior of `gotwarden`.

It can be set as environment variables or into a `.env` file.

| Variables | Description | Default |
|-----------|-------------|---------|
| DB_TYPE   | Database type (so far only 'sqlite' is supported)  | sqlite  |
| DB_FILEPATH | Sqlite database path | ./fixtures/test.db |
| DB_USER   | Database user* | |
| DB_PASSWORD | Database password* | |
| DB_HOST   | Database host*  | localhost |
| DB_NAME   | Database name*  | |
| DB_PORT   | Database port* | 5432 |
| PORT | Web server port | 3000 |
| WARDEN_IDENTITY_URL|| /identity |
| WARDEN_ATTACHMENT_URL || /attachments |
| WARDEN_ICONS_URL || /icons |
| WARDEN_SECRET_PHRASE || This a secret ... sshhhshh" |
| WARDEN_STATIC_PATH || ./fixtures/assets |

> No needed for sqlite database

## Built With

Mainly with this components:

* As web-framework [gin](https://github.com/gin-gonic/gin)
* As ORM like [gorp](https://github.com/go-gorp/gorp)

And others ... 

## Contributing

Please read [CONTRIBUTING.md](https://gist.github.com/PurpleBooth/b24679402957c63ec426) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags). 

## Authors

* **Tutilus**  [Tutilus](https://github.com/tutilus)

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Acknowledgments
