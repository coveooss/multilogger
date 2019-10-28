// Copyright 2019 Coveo Solution inc. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
Package multilogger is a thin wrapper around logrus https://github.com/sirupsen/logrus that writes logs simultaneously to many outputs at the same time.
Each output (implemented as Hook) have their own logging level.
It exposes the same API as a regular logrus logger.

How to use it

So, you can use multilogger to write logs to standard error with a minimal logging level (i.e. WarningLevel) and have the full debug log (i.e. DebugLevel or TraceLevel) written to a file.

	import "github.com/coveooss/multilogger"
	import "github.com/sirupsen/logrus"

	func main() {
	    log := New("test")                                         // By default, logs are sent to stderr and the logging level is set to WarningLevel
	    log.AddHooks(NewFileHook("debug.log", logrus.DebugLevel))  // It is possible to add a file as a target for the logs
	    log.AddHooks(NewFileHook("verbose.log", "trace"))          // The logging level could be expressed as string, logrus.Level or int
	    log.Warning("This is a warning")
	    log.Infof("This is an information %s", message)
	}

Differences with logrus

 - The multilogger object is based on a logrus.Entry instead of a logrus.Logger.
 - The methods Print, Println and Printf are used to print information directly to stdout without log decoration.
   Within logrus, these methods are aliases to Info, Infoln and Infof.
 - It is not possible to set the global logging level with multilogger, the logging level is determined by hooks that are
   added to the logging object and the resulting logging level is the highest logging level defined for individual hooks.
 - Hooks are fired according to their own logging level while they were fired according to the global logging
   level with logrus.
 - The default formatter used by console and file hook is highly customizable.

Where is it used

This project is used in other Coveo projects to reduce visual clutter in CI systems while keeping debug logs available when errors arise:
    gotemplate https://github.com/coveooss/gotemplate
    terragrunt https://github.com/coveooss/terragrunt
    tgf        https://github.com/coveooss/tgf
*/
package multilogger
