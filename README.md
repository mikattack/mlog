# mlog

A wrapper around the excellent standard Go logging library, `mlog` satisfies my logging preferences in ways that no other package quite has.  It provides:

- Self-explanitory logging levels (under a [sensible abstraction](http://labs.ig.com/logging-level-wrong-abstraction))
- Easy control of where messages are logged to
- Easy control of what levels are logged
- No requirements for configuration or initialization (sensible defaults)
- Support for output to multiple streams, customizable per logging level
- No `fatal` or `panic` calls

This library is based heavily on [jWalterWeatherman](https://github.com/spf13/jWalterWeatherman),
which was the closest library I'd found that met my needs.

Additionally, while this library supports logging to any destination, if that destination is `syslog` then you're better off using Go's `log/syslog` package.


# Usage

No initialization or configuration is necessary.  The library works by creating a number of loggers which correspond to the following logging levels: `DEBUG `, `INFO `, `WARN `, and `ERROR`.

These loggers are default Go loggers and can be used like the following:

```
import "gitlab.com/mikattack/mlog"

// DEBUG
mlog.inTest.Printf('Noisey output, useful for development')

// INFO
mlog.inProd.Printf('Information needed to debug production issues')

// WARN
mlog.toInvestigateTomorrow.Printf('Needs investigation, but can wait until tomorrow')

// ERROR
mlog.wakeMeInTheMiddleOfTheNight.Printf("Needs attention RIGHT NOW")
```

The logger names are verbose and self descriptive. This makes it easier to decide which level to output at.


# Configuration

The library defaults to the following behavior:

- Log threshold level is `DEBUG`, meaning everything is outputt
- `WARN`, `ERROR` messages are logged to `STDOUT`
- Flags are: `DATE`, `TIME`, and `SFILE`

Each of these settings are configurable.

### Change logging threshold

The threshold can be changed at any time, but will only affect calls executed after the change was made. Anything below the configured level (exclusive) will not be logged.

```
// Exclue INFO ("inTest") messages
if verbose == false {
  mlog.SetLogThreshold(mlog.LEVEL_PRODUCTION)
}
```

### Change output destination

All log messages go to `STDOUT` by default, but can be customized on a per-level basis.  For example, to get that true 12-Factor setup, you can send error-related message to `STDERR`:

```
import (
  "os"
  "gitlab.com/mikattack/mlog"
)

mlog.SetOutput(mlog.LEVEL_TEST, os.Stderr)
mlog.SetOutput(mlog.LEVEL_PRODUCTION, os.Stderr)

mlog.SetOutput(mlog.LEVEL_TOMORROW, os.Stderr)
mlog.SetOutput(mlog.LEVEL_MIDDLE_OF_NIGHT, os.Stderr)
```

Because the output is just an `io.Writer`, it's also easy to write log streams to a file:

```
file := os.OpenFile("/var/tmp/warnings.log", os.O_RDWR|os.O_APPEND, 0660);
defer file.Close()

mlog.SetOutput(mlog.LEVEL_TOMORROW, file)
```

If you need to get extra fancy, you can log messages to multiple sources:

```
errlog := os.OpenFile("/var/tmp/critical.log", os.O_RDWR|os.O_APPEND, 0660);
defer errlog.Close()

// Output critical messages to STDERR and "/var/tmp/critical.log"
mlog.SetOutput(mlog.LEVEL_MIDDLE_OF_NIGHT, os.Stderr, errlog)
```

### Change log flags

Flags control what extra information gets added to every message:

- `NONE` - Adds nothing to the message (and ignored when used with other flags)
- `DATE` - Adds the date to a message
- `TIME` - Adds the time to a message
- `SFILE` - Adds the file the message originated from
- `LFILE` - Adds the line number the message originated from
- `MSEC` - Adds microsecond resolution to the time (if present)

Flags may be set per log stream or all at once:

```
// Strip all extra output for log streams, except critical messages
mlog.SetFlags(NONE)
mlog.SetFlags(DATE | TIME | SFILE, mlog.LEVEL_MIDDLE_OF_NIGHT)
```


# More Information

This is a convenience package designed for ease-of-use.  It doesn't do everything under the sun or anything radically different from other packages of its ilk.  The API is nice, but should not be considered stable.
