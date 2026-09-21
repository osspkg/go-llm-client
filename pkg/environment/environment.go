/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package environment describes environment variables supported by provider
// servers.
package environment

// Variable describes one environment variable accepted by a provider server.
//
// An empty AllowedValues slice means that the variable accepts a documented
// scalar format instead of a finite enumeration. Default is the value used
// when the variable is not set; an empty Default means that the provider does
// not set a value by default.
type Variable struct {
	// Name is the exact case-sensitive environment variable name.
	Name string
	// AllowedValues lists finite values documented by the provider. It is empty
	// when the provider accepts a scalar format such as a path, duration, or
	// integer range.
	AllowedValues []string
	// Default is the provider's documented value when Name is not set. An empty
	// value means that the provider uses an empty or unset default; the
	// Description distinguishes those cases when they matter.
	Default string
	// Description explains the effect, format, and operational purpose of the
	// variable.
	Description string
}

// Scheme is a provider-independent description of server environment
// variables. The returned value is metadata only; it does not read or modify
// the process environment.
type Scheme struct {
	// Variables contains the provider's environment variable descriptions in a
	// stable order suitable for CLI and documentation rendering.
	Variables []Variable
}
