# Upload file bot

Service which is engaged in sending files in telegrams

## Options

* `app-id` telegram app id. You can get it from [telegram apps](https://my.telegram.org/apps)
* `app-hash` telegram app hash. You can get it from [telegram apps](https://my.telegram.org/apps)
* `bot-token` telegram bot token. How to create it looks [here](https://core.telegram.org/bots)
* `user-id` telegram user-id which bot will send files
* `input-folder` the folder from which the service will take files to send to the user
* `output-folder` the folder in which the service will put the files after sending to the user

## Quick start

### Install

#### With docker
1. copy provided docker-compose.yml and customize
2. compile from the sources - `docker-compose build && docker-compose up -d`

#### Without docker
1. make sure that you have `go` version 1.16 or higher
2. compile binary file `make build`
3. run it
```shell
./dist/uptotg                       \
  --app-id={{app-id}}               \
  --app-hash={{app-hash}}           \
  --bot-token={{bot-token}}         \
  --user-id={{user-id}}             \
  --input-folder={{input-folder}}   \
  --output-folder={{output-folder}}
```
