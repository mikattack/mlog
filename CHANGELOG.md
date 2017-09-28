# mlog Changelog

### 1.1.0 - 2017-09-27

Rewrite library. Now a reimplementation of the standard logger rather than a wrapper around it. Several things have changed:

 * Still only four logging levels, but their function names have been updated for additional clarity: `InTesting`, `InProduction`, `ToInvestigate`, and `PageMeNow`.
 * No longer have to call formatting functions per level. The logging call _is_ the logging level.
 * No longer per-logging level output streams. All levels of output are directed to the same `io.Writer`.
 * Simplified flags controlling additional information with each log message:
     - `LEVEL`: Logging level indicator.
     - `FILE`: File name and line number (akin to `log.Lshortfile`).
     - `DATE`: Date and time (akin to `log.LUTC | log.Ltime`)


### 1.0.0 - 2017-05-02

Initial release.

Features:

  * Self-explanitory logging levels (under a sensible abstraction)
  * Easy control of where messages are logged to
  * Easy control of what levels are logged
  * No requirements for configuration or initialization (sensible defaults)
  * Support for output to multiple streams, customizable per logging level
  * No fatal or panic calls
