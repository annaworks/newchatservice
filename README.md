# Suru - slack bot 

## Description

This slack bot service is built for a way to save questions and answers to a datastore of choice, for building a self-hosted knowledge base and to add improved search features within a slack channel. 

The service interfaces with the Slack API using an API built with Golang.

## Local Development
#### Running Go API from project location
Copy .env.example and configure ENV variables 

```
$ cp .env.example .env
```

Compile and run the go api
```
go run cmd/surubot.go
```
#### Setup slack bot
Go to [https://api.slack.com](https://api.slack.com) for instructions on setting up your own slack app.

#### Expose local server to the public
Use [https://ngrok.com](https://ngrok.com) for quick setup to connect your local go api to the slack api.

## Tests
### A health point test has been implemented.
```
go test ./...
```