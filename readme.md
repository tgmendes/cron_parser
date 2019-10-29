# Cron Expression Parser
## Intro

This small program implements a cron expression parser. It will receive a given cron expression as an input, and will output a table with the times at which the cron job will run.

For example, the following expression:

`*/15 0 1,15 * 1-5 /usr/bin/find`

will produce the following table:
```
minute        0 15 30 45
hour          0
day of month  1 15
month         1 2 3 4 5 6 7 8 9 10 11 12
day of week   1 2 3 4 5
command       /usr/bin/find
```

To learn more about crontabs, and its particularities, the following links might be useful:

* [Beginners Guide](https://www.ostechnix.com/a-beginners-guide-to-cron-jobs/)
* [Expression Checker](https://crontab.guru)


## Go Mod

This project uses go mod, and was developed with version 1.13.

## Usage

To use the program, run the following from the `cmd` directory:

```go run main.go -e "<YOUR EXPRESSION>"```

Note the usage of quotes in the command. The example given above could be run by using the following command:

```go run main.go -e "*/15 0 1,15 * 1-5 /usr/bin/find"```

## Tests

The cron parser (`parser.go`) has a set of unit tests that describe the usage of the parser with a wide variety of use cases (see `parser_test.go`). To run them use:

```go test```

while in the same directory as the code.

## TODO

* Explore the possibility of using regexes as an alternative to parsing the input strings
* Add support for more particular use cases (e.g. `0-20/2` or `*,15`)
* Add more robust error handling
* Add support for using 7 as a day of week (sunday)

## Author
Tiago Mendes (tgmendes@gmail.com)
