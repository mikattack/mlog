/*
Package mlog implements a structured logger.

The easiest way to use mlog is through the package-level logger:

  package main

  import log "github.com/mikattack/mlog"

  func main() {
    log.ToInvestigate("Non-existent inventory applied to order: 51138000")
  }

Output:

  [WARNING] 2017-09-22 08:52:33 orders.go:513 Non-existent inventory applied to order: 51138000

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
