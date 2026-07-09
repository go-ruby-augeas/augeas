// SPDX-License-Identifier: BSD-3-Clause
//
// Copyright (c) 2026, the go-ruby-augeas/augeas authors

package augeas

import engine "github.com/go-augeas/augeas"

// Flag values mirror the ruby-augeas AUG_* constants passed to Augeas.open.
// They are recorded on the handle so a Ruby binding round-trips them; the pure-Go
// engine does not perform file autoload, type checking or backups, so only NoLoad
// changes behaviour here (Open never touches the filesystem regardless).
const (
	None           = 0
	SaveBackup     = 1 << 0
	SaveNewFile    = 1 << 1
	TypeCheck      = 1 << 2
	NoStdinc       = 1 << 3
	SaveNoop       = 1 << 4
	NoLoad         = 1 << 5
	NoModlAutoload = 1 << 6
	EnableSpan     = 1 << 7
)

// Lens is the engine lens interface, re-exported for adapter callers.
type Lens = engine.Lens

// Span is the engine span record, re-exported for adapter callers.
type Span = engine.Span

// FileSystem is the engine filesystem seam, re-exported for adapter callers.
type FileSystem = engine.FileSystem

// ErrSpanUnsupported is returned by [Augeas.Span]; span tracking is a documented
// deferred feature of the engine.
var ErrSpanUnsupported = engine.ErrSpanUnsupported

// LensByName returns the built-in engine lens registered under name.
func LensByName(name string) (Lens, bool) { return engine.LensByName(name) }

// Augeas is a ruby-augeas-shaped handle over a go-augeas engine tree.
type Augeas struct {
	eng      *engine.Augeas
	root     string
	loadPath string
	flags    int
}

// New returns a handle with default settings, mirroring Augeas.new with no
// arguments.
func New() *Augeas {
	return &Augeas{eng: engine.New()}
}

// Open mirrors Augeas.open(root, loadpath, flags), recording the arguments on
// the handle. The pure-Go engine does not autoload files, so nothing is read
// until [Augeas.Load] is called explicitly.
func Open(root, loadPath string, flags int) *Augeas {
	a := New()
	a.root = root
	a.loadPath = loadPath
	a.flags = flags
	return a
}

// Root returns the root argument passed to Open ("" for New).
func (a *Augeas) Root() string { return a.root }

// LoadPath returns the loadpath argument passed to Open ("" for New).
func (a *Augeas) LoadPath() string { return a.loadPath }

// Flags returns the flags passed to Open (None for New).
func (a *Augeas) Flags() int { return a.flags }

// Engine exposes the underlying engine handle for advanced use.
func (a *Augeas) Engine() *engine.Augeas { return a.eng }

// SetFileSystem injects the filesystem seam used by Load and Save.
func (a *Augeas) SetFileSystem(fs FileSystem) { a.eng.SetFileSystem(fs) }

// Get returns the value at path; the boolean reports whether exactly one node
// matched (aug_get).
func (a *Augeas) Get(path string) (string, bool) { return a.eng.Get(path) }

// Exists reports whether at least one node matches path (aug_exists).
func (a *Augeas) Exists(path string) bool { return a.eng.Exists(path) }

// Set sets the value at path, creating the node if absent (aug_set).
func (a *Augeas) Set(path, value string) error { return a.eng.Set(path, value) }

// Setm sets value on every node matching sub under each node matching base,
// returning the count (aug_setm).
func (a *Augeas) Setm(base, sub, value string) (int, error) {
	return a.eng.SetMultiple(base, sub, value)
}

// Insert inserts a sibling labelled label next to path (aug_insert).
func (a *Augeas) Insert(path, label string, before bool) error {
	return a.eng.Insert(path, label, before)
}

// Rm removes every node matching path and returns the count (aug_rm).
func (a *Augeas) Rm(path string) int { return a.eng.Remove(path) }

// Mv moves the node matching src onto dst (aug_mv).
func (a *Augeas) Mv(src, dst string) error { return a.eng.Move(src, dst) }

// Match returns the paths of all nodes matching path (aug_match).
func (a *Augeas) Match(path string) []string { return a.eng.Match(path) }

// Label returns the label of the single node matching path (aug_label).
func (a *Augeas) Label(path string) (string, bool) { return a.eng.Label(path) }

// Defvar binds name to the node-set from expr and returns the count
// (aug_defvar).
func (a *Augeas) Defvar(name, expr string) (int, error) {
	return a.eng.DefineVariable(name, expr)
}

// Defnode binds name to expr, creating a node with value when expr matches
// nothing; it returns the node path and whether it was created (aug_defnode).
func (a *Augeas) Defnode(name, expr, value string) (string, bool) {
	return a.eng.DefineNode(name, expr, value)
}

// TextStore parses text with lens and stores the tree at path (aug_text_store).
func (a *Augeas) TextStore(lens Lens, path, text string) error {
	return a.eng.TextStore(lens, path, text)
}

// TextRetrieve serialises the subtree with lens (aug_text_retrieve). When node
// is nil the single subtree at path is used.
func (a *Augeas) TextRetrieve(lens Lens, path string, node *engine.Node) (string, error) {
	return a.eng.TextRetrieve(lens, path, node)
}

// Load reads files matching pattern, parses them with lens and mounts them under
// mount (aug_load, given an explicit lens/glob).
func (a *Augeas) Load(lens Lens, pattern, mount string) error {
	return a.eng.Load(lens, pattern, mount)
}

// Save serialises the subtree at mount with lens and writes it to filename
// (aug_save, given an explicit lens/mount/target).
func (a *Augeas) Save(lens Lens, mount, filename string) error {
	return a.eng.Save(lens, mount, filename)
}

// Span mirrors aug_span; span tracking is deferred, so it always returns
// ErrSpanUnsupported.
func (a *Augeas) Span(path string) (Span, error) { return a.eng.Span(path) }

// Error returns the last error recorded by the engine (aug_error).
func (a *Augeas) Error() error { return a.eng.Error() }
