# Scripts Creator
This is an api to create scripts written in Go. The scripts are saved in the server.

## Running
The project includes a `docker-compose.yml` file, so in order to serve the api you just need to run:
```shell
docker-compose up
```

## Usage
The api only have one resource, namely:
```http request
POST /scripts
```
This will receive the commands and their dependencies between them, sort to accomplish the dependencies and create the script.
The expected body must have the following shape:
```json
{
  "filename": "myscript.sh",
  "tasks": [
    {
      "name": "rm",
      "command": "rm -f /tmp/test",
      "dependencies": [
        "cat"
      ]
    },
    {
      "name": "cat",
      "command": "cat /tmp/test",
      "dependencies": [
        "chown",
        "chmod"
      ]
    },
    {
      "name": "touch",
      "command": "touch /tmp/test"
    },
    {
      "name": "chmod",
      "command": "chmod 600 /tmp/test"

      "dependencies": [
        "touch"
      ]
    },
    {
      "name": "chown",
      "command": "chown root:root /tmp/test"

      "dependencies": [
        "touch"
      ]
    }
  ]
}
```
This body will create the script:
```shell
touch /tmp/test
chown root:root /tmp/test
chmod 600 /tmp/test
cat /tmp/test
rm -f /tmp/test
```

## Tests
The project includes a script to run the tests and print the coverage report. Also, there is a second docker compose which will run this script instead of serving the application.
Also, the repository includes a github action that builds the app, runs the tests and prints a coverage report. 

Test Coverage: 92.5%
