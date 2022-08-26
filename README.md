# Statements
The present project is a long living process that parses csv files to generate an account summary and sends it via email.

## Assumptions
This worker is implemented assuming that:
* all files are named after the email address we will be sending the account statement
* All files have the extension `.csv`. E.g.: `user@mail.com.csv`
* The files are readable and available in the configured folder
* The email is sent using [sendgrid](http://sendgrid.com) as email broker, so an API Key and a dynamic template must be created beforehand

### Email details
The email is sent using [sendgrid](http://sendgrid.com) as email broker. 
A compatible dynamic template is [included](emailTemplate.html) in this repo.

## Local installation

1. Install [go](https://golang.org/doc/install).
2. Clone or copy the repo into your `GOPATH`

## Create config.cfg file
The service requires a `config.cfg` file under the `./config/` folder.
This configuration file has the following structure:
```
# Indicates how often the process will work. E.g.: set to "730h" to run every 30 days
interval = "24h"
# Set startAt to empty to stop the worker. If you want an immediate start set it to "now", otherwise use a AM-PM hour format, e.g., 12:00PM
startAt = "now"
# Directory containing the statement CSV files, relative to the root of the project
filesDir = "statements"
# API key to be able to send emails
sendGridAPIKey = "add your sendgrid API key"
# the sendgrid dynamic template ID to use for mailing
templateID = "dynamic template ID"
```
A sample config file can be created by starting the service with the -- sampleconfig flag enabled
```bash
$ go run cmd/app/main.go --sampleconfig=true > config/config.cfg
```
You can update configuration in `./config/config.cfg` if required.

## Installing dependencies
This can be made by running go mod commands, i.e.: `go mod tidy && go mod vendor`

## Running locally
After creating a config.cfg file with a valid configuration and having a running Elasticsearch instance, you can run the service by excecuting in your terminal
```bash
$ go run cmd/app/main.go
```

### Running with Docker
Service includes a dockerfile that can be used to build and run the backend service

```bash
$ docker build --tag=statements .
docker run statements
```

## How to navigate the code
### cmd/app/main.go
Configuration and logger initialization. Then the main function "continues" in internal/app/app.go.

### config
Configuration. First, `config.cfg` is read and used to populate the config struct in `config.go`
In order to parse the `config.cfg` file, [toml](https://github.com/BurntSushi/toml) is required

### internal/app
There is always one Run function in the `app.go` file, which "continues" the main function.
We start the worker and wait for signals in select for graceful completion.


### internal/usecase
Business logic.

### internal/usecase/interfaces.go
interfaces to implement the business logic (usecases)