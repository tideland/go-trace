// Tideland Go Trace - Crumbs
//
// Copyright (C) 2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package crumbs provides a flexible but still convenient helper
// for logging. The backend is based on an interface and so can be
// exchanged. A Crumbs instance returns a leveled CrumbWriter which
// can write Info and Error messages. They messages also may be
// complemented with key/value pairs.
//
//     c := crumbs.New()
//     c.L(1).Info("a message", "id", 1337, "name", aName)
//     c.L(9).Error(err, "something happened")
//
// The additional function Crumble helps to log errors in defer
// expressions.
//
//     cw := c.L(8)
//     defer crumbs.Crumble(cw, myFile.Close, "closing file failed", "name", filename)
//
// Configuration options of Crumbs are the backend and the lowest
// level for reporting.
package crumbs // import "tideland.dev/go/trace/crumbs"

// EOF
