# mlog

A reimplementation of the the excellent standard Go logging library, `mlog` satisfies my logging preferences in ways that no other package quite has.  It provides:

- Self-explanitory logging levels (read [these thoughts](http://labs.ig.com/logging-level-wrong-abstraction) for the motivation behind them)
- No requirements for configuration or initialization (sensible defaults)
- No `fatal` or `panic` calls

This library was motivated by a desire for a different logging API and as a learning exercise.


# Usage

No initialization or configuration is necessary.  Although the library supports logging levels, they are highly "opinionated" in favor of code clarity over flexibility. There are only four levels available:

```
import "gitlab.com/mikattack/mlog"

def main() {
    // DEBUG
    mlog.InTesting('Noisey output, useful for development')

    // INFO
    mlog.InProduction('Information needed to debug production issues')

    // WARN
    mlog.ToInvestigate('Needs investigation, but can wait until tomorrow')

    // ERROR
    mlog.PageMeNow("Needs attention RIGHT NOW")
}
```

The logger names are verbose and self descriptive. This makes it easier to decide which level to output at.


# Configuration

The library defaults to the following behavior:

- Log threshold level is `INFO` ("In Production")
- Output is `STDOUT` for all logging levels
- Flags are: `DATE`, `FILE`, and `LEVEL`

Each of these settings are configurable.

### Change logging threshold

The threshold can be changed at any time, but will only affect calls executed after the change was made. Anything below the configured level (exclusive) will not be logged.

```
// Exclude "InTesting" messages
if verbose == false {
  mlog.SetThreshold(mlog.IN_PRODUCTION)
}
```

### Change output destination

All log messages go to `STDOUT` by default but can be redirected anywhere:

```
mlog.SetOutput(os.Stderr)
```

Because the output is just an `io.Writer`, it's also easy to write log streams to a file:

```
file := os.OpenFile("/var/tmp/warnings.log", os.O_RDWR|os.O_APPEND, 0660);
defer file.Close()

mlog.SetOutput(file)
```

If you need to get extra fancy, you can log messages to multiple sources:

```
alog := os.OpenFile("/var/tmp/activity.log", os.O_RDWR|os.O_APPEND, 0660);
defer alog.Close()

// Output to STDERR and "/var/tmp/activity.log"
mlog.SetOutput(os.Stderr, alog)
```

### Change log flags

Flags control what extra information gets added to every message:

- `DATE` - Adds the UTC date and time to a message
- `FILE` - Adds the file and line number the message originated from
- `LEVEL` - Prefixes the message with the logging level

All flags are enabled by default, which is identical to the following:

```
mlog.SetFlags(DATE | LEVEL | FILE)
```


# More Information

This is a convenience package designed for ease-of-use.  It doesn't do everything under the sun or anything radically different from other packages of its ilk.  This project was mostly a learning exercise.
