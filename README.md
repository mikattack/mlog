# mauth

A wrapper around the excellent standard Go logging library, `mauth` satisfies my logging preferences in ways that no other package quite has:

- Logging levels
- Allow easy control of where messages are logged to
- Allow easy control of what levels are logged
- Not require any configuration or initialization (sensible defaults)
- Support output to multiple streams (per log level, multiple writers per level)

This package is based heavily on (jWalterWeatherman)[https://github.com/spf13/jWalterWeatherman],
which was the closest library I knew of that met my needs.

# Usage

No initialization or configuration is necessary.  The library works by creating a number of loggers which correspond to the following logging levels:

- TRACE
- DEBUG
- INFO
- WARN
- ERROR
- CRITICAL
- FATAL

These loggers are based on the standard `log` library and operate much the same way:

```
import "gitlab.com/mikattack/mlog"

...

if err != nil {
  mauth.ERROR.Println(err)
}
if warn != nil {
  mauth.WARN.Println(warn)
}

mauth.INFO.Printf("the ice skates are %s", color)
```

While seven log levels is a lot, you can choose to use the ones appropriate to your application. Furthermore, only those messages falling within the range of the logging threshold will actually be output.

Additionally, you can create custom loggers which output to the `io.Writer` of your choice and is unaffected by logging thresholds:

```
mlog.NewCustomLogger("telemetry", "TELEMENTRY")

// Write statistics to the "telemetry" custom logger
if configs["telemetry"] == true {
  stats := ...
  mlog.Printf("telemetry", stats)
}
```


# Configuration

The library defaults to the following behavior:

- Log level is `WARN`
- Trace, debug and info messages are discarded
- Warn, error, critical, and fatal messages are logged to `STDOUT`

### Change logging threshold

The threshold can be changed at any time, but will only affect calls executed after the change was made.

```
if verbose == true {
  mlog.SetLogThreshold(mauth.LEVEL_TRACE)
}
```

### Change output destination

All log messages go to `STDOUT` by default, but can be customized on a per-level basis.  For example, to get that true 12-Factor setup, you can send error-related message to `STDERR`:

```
import (
  "os"
  "gitlab.com/mikattack/mlog"
)

mlog.SetOutput(mlog.LEVEL_WARN, os.Stderr)
mlog.SetOutput(mlog.LEVEL_ERROR, os.Stderr)
mlog.SetOutput(mlog.LEVEL_CRITICAL, os.Stderr)
mlog.SetOutput(mlog.LEVEL_FATAL, os.Stderr)

// Works for custom loggers too!
mlog.SetOutput("telemetry", os.Stderr)
```

Because the output is just an `io.Writer`, it's also easy to write log streams to a file:

```
file := os.OpenFile("/var/tmp/example.log", os.O_RDWR|os.O_APPEND, 0660);
defer file.Close()

mlog.SetOutput(mlog.LEVEL_WARN, file)
```

If you need to get extra fancy, you can log messages to multiple sources:

```
errlog := os.OpenFile("/var/tmp/error.log", os.O_RDWR|os.O_APPEND, 0660);
defer errlog.Close()

// Output ERROR messages to STDERR and "/var/tmp/errors.log"
mlog.SetOutput(mlog.LEVEL_ERROR, os.Stderr, errlog)
```


# More Information

This is a convenience package designed for ease-of-use.  It doesn't do everything under the sun or anything radically different from other packages of its ilk.  The API is nice, but should not be considered stable.
