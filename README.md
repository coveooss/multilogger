# multilogger

Multilogger is a small logging wrapper based on [logrus](https://github.com/sirupsen/logrus) that writes logs simultaneously to the console,
files or any other hook. Each of those logging methods can be given a different logging level. It exposes the same API as a regular `logrus` logger.

See [multilogger godoc](https://godoc.org/github.com/coveooss/multilogger) for full documentation.

## multicolor

There is also a sub package used to handle colors and attributes. It is based on [color](https://github.com/fatih/color) which is an archived project
but still useful. We based our implementation on [ghishadow](https://github.com/ghishadow/color) fork since it is maintained and have been migrated
to go module.

See [multicolor godoc](https://godoc.org/github.com/coveooss/multilogger/color) for full documentation.
