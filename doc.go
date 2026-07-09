// SPDX-License-Identifier: BSD-3-Clause
//
// Copyright (c) 2026, the go-ruby-augeas/augeas authors

// Package augeas is a pure-Go (no cgo) adapter that presents the ruby-augeas
// gem API over the configuration-editing engine
// github.com/go-augeas/augeas.
//
// The go-augeas engine already implements Augeas semantics: the tree model, the
// XPath-like path language, the get/set/insert/move/remove editing operations,
// the lens framework and file load/save behind a filesystem seam. This package
// does not reimplement any of that; it wraps the engine and renames the surface
// to match the Ruby gem so a consumer such as go-embedded-ruby (rbgo) can expose
// a Ruby "Augeas" class:
//
//   - a constructor pair, [New] and [Open](root, loadpath, flags), mirroring
//     Augeas.new / Augeas.open, with the AUG_* flag constants;
//   - the gem method surface: [Augeas.Get], [Augeas.Exists], [Augeas.Set],
//     [Augeas.Setm], [Augeas.Insert], [Augeas.Rm], [Augeas.Mv], [Augeas.Match],
//     [Augeas.Save], [Augeas.Load], [Augeas.Defvar], [Augeas.Defnode],
//     [Augeas.Label], [Augeas.TextStore], [Augeas.TextRetrieve], [Augeas.Span]
//     and [Augeas.Error], each delegating to the engine;
//   - type aliases ([Lens], [Span], [FileSystem]) and [LensByName] re-exported
//     so callers need not import the engine directly;
//   - an injectable filesystem seam ([Augeas.SetFileSystem]) where the Ruby
//     interpreter's I/O would otherwise plug in.
//
// The package has no dependency on any Ruby runtime: the surface is Go-typed and
// a Ruby binding layer marshals Ruby values onto these Go types.
package augeas
