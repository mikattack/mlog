/*
Package mlog implements a structured logger.

The easiest way to use mlog is through the package-level logger:

  package main

  import log "github.com/mikattack/mlog"

  func main() {
    log.WithFields({
      "warehouse":  "phx-4",
      "order":      51132000,
    }).ToInvestigate("Negative inventory applied to order")
  }

Output:

  time="2017-09-22T08:52:33Z" level=warning warehouse="phx-4" order=5113200 msg="Negative inventory applied to order"

The logging API is limited and intentionally simplistic. It eschews fine-grain
levels for obvious ones:

  // DEBUG
  log.InTesting("Noisey output, for development")

  // INFO
  log.InProduction("Information needed to debug production issues")

  // WARN
  log.ToInvestigate("Needs investigation, but can wait until tomorrow")

  // ERROR
  log.PageMeNow("Needs attention RIGHT NOW")
*/
package mlog
