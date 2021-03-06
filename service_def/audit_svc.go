// Copyright 2013-Present Couchbase, Inc.
//
// Use of this software is governed by the Business Source License included in
// the file licenses/BSL-Couchbase.txt.  As of the Change Date specified in that
// file, in accordance with the Business Source License, use of this software
// will be governed by the Apache License, Version 2.0, included in the file
// licenses/APL2.txt.

package service_def

import ()

type AuditEventIface interface {
	Redact() AuditEventIface
	Clone() AuditEventIface
}

type AuditSvc interface {
	//	Write(eventId uint32, event interface{}) error
	Write(eventId uint32, event AuditEventIface) error
}
