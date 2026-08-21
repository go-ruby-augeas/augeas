// SPDX-License-Identifier: BSD-3-Clause
//
// Copyright (c) 2026, the go-ruby-augeas/augeas authors

package augeas

import (
	"errors"
	"io/fs"
	"testing"

	engine "github.com/go-augeas/augeas"
)

// memFS is an in-memory FileSystem used to exercise Load and Save.
type memFS struct {
	files  map[string]string
	writes map[string][]byte
}

func (m *memFS) Glob(pattern string) ([]string, error) {
	var out []string
	for k := range m.files {
		out = append(out, k)
	}
	return out, nil
}

func (m *memFS) ReadFile(name string) ([]byte, error) {
	v, ok := m.files[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return []byte(v), nil
}

func (m *memFS) WriteFile(name string, data []byte, _ fs.FileMode) error {
	if m.writes == nil {
		m.writes = map[string][]byte{}
	}
	m.writes[name] = data
	return nil
}

func TestConstructorsAndAccessors(t *testing.T) {
	a := New()
	if a.Root() != "" || a.LoadPath() != "" || a.Flags() != None {
		t.Fatal("New defaults")
	}
	if a.Engine() == nil {
		t.Fatal("Engine handle")
	}
	o := Open("/root", "/lenses", TypeCheck|NoLoad)
	if o.Root() != "/root" || o.LoadPath() != "/lenses" || o.Flags() != TypeCheck|NoLoad {
		t.Fatalf("Open %q %q %d", o.Root(), o.LoadPath(), o.Flags())
	}
	// reference the remaining flag constants so the surface is complete
	_ = SaveBackup | SaveNewFile | NoStdinc | SaveNoop | NoModlAutoload | EnableSpan
}

func TestEditingSurface(t *testing.T) {
	a := New()
	// Set / Get / Exists
	if err := a.Set("/files/etc/hosts/1/ipaddr", "127.0.0.1"); err != nil {
		t.Fatal(err)
	}
	if err := a.Set("/files/etc/hosts/1/canonical", "localhost"); err != nil {
		t.Fatal(err)
	}
	if v, ok := a.Get("/files/etc/hosts/1/ipaddr"); !ok || v != "127.0.0.1" {
		t.Fatalf("Get %q %v", v, ok)
	}
	if !a.Exists("/files/etc/hosts/1") {
		t.Fatal("Exists")
	}
	// Set error path
	if err := a.Set("/[1]", "x"); err == nil {
		t.Fatal("Set error")
	}
	// Setm
	a.Set("/svc/1/on", "y")
	a.Set("/svc/2/on", "y")
	if n, err := a.Setm("/svc/*", "on", "n"); err != nil || n != 2 {
		t.Fatalf("Setm %d %v", n, err)
	}
	if _, err := a.Setm("/[1]", "x", "y"); err == nil {
		t.Fatal("Setm error")
	}
	// Insert
	if err := a.Insert("/files/etc/hosts/1/ipaddr", "alias", false); err != nil {
		t.Fatal(err)
	}
	if err := a.Insert("/nope", "x", false); err == nil {
		t.Fatal("Insert error")
	}
	// Label
	if l, ok := a.Label("/files/etc/hosts/1/ipaddr"); !ok || l != "ipaddr" {
		t.Fatalf("Label %q %v", l, ok)
	}
	// Match
	if m := a.Match("/files/etc/hosts/*"); len(m) != 1 {
		t.Fatalf("Match %v", m)
	}
	// Mv
	if err := a.Mv("/files/etc/hosts/1", "/files/etc/hosts/2"); err != nil {
		t.Fatal(err)
	}
	if err := a.Mv("/[1]", "/x"); err == nil {
		t.Fatal("Mv error")
	}
	// Rm
	if n := a.Rm("/files/etc/hosts/*"); n != 1 {
		t.Fatalf("Rm %d", n)
	}
}

func TestVariablesAndSpan(t *testing.T) {
	a := New()
	a.Set("/a/b", "v")
	if n, err := a.Defvar("x", "/a/*"); err != nil || n != 1 {
		t.Fatalf("Defvar %d %v", n, err)
	}
	if _, err := a.Defvar("y", "/[1]"); err == nil {
		t.Fatal("Defvar error")
	}
	if p, created := a.Defnode("z", "/a/b", "v"); created || p != "/a/b" {
		t.Fatalf("Defnode existing %q %v", p, created)
	}
	if p, created := a.Defnode("z2", "/new/node", "nv"); !created || p != "/new/node" {
		t.Fatalf("Defnode created %q %v", p, created)
	}
	if _, err := a.Span("/a/b"); err != ErrSpanUnsupported {
		t.Fatalf("Span %v", err)
	}
	// Error surface
	a.Get("/a/*") // multi-match sets last error
	if a.Error() == nil {
		t.Fatal("Error should be set")
	}
}

func TestLensSurface(t *testing.T) {
	lens, ok := LensByName("Hosts")
	if !ok {
		t.Fatal("Hosts lens")
	}
	a := New()
	// TextStore / TextRetrieve via path
	if err := a.TextStore(lens, "/files/etc/hosts", "127.0.0.1 localhost\n"); err != nil {
		t.Fatal(err)
	}
	out, err := a.TextRetrieve(lens, "/files/etc/hosts", nil)
	if err != nil || out != "127.0.0.1 localhost\n" {
		t.Fatalf("TextRetrieve %q %v", out, err)
	}
	// TextRetrieve via explicit node
	node := a.Engine().Root()
	m := a.Engine().Match("/files/etc/hosts")
	_ = m
	var target *engine.Node
	for _, c := range node.Children { // /files
		for _, e := range c.Children { // /files/etc
			for _, h := range e.Children { // /files/etc/hosts
				target = h
			}
		}
	}
	if out, err := a.TextRetrieve(lens, "", target); err != nil || out != "127.0.0.1 localhost\n" {
		t.Fatalf("TextRetrieve node %q %v", out, err)
	}

	// Load / Save through the filesystem seam
	mfs := &memFS{files: map[string]string{"hosts": "10.0.0.1 host\n"}}
	b := New()
	b.SetFileSystem(mfs)
	if err := b.Load(lens, "*", "/files"); err != nil {
		t.Fatal(err)
	}
	if v, _ := b.Get("/files/hosts/1/canonical"); v != "host" {
		t.Fatalf("loaded %q", v)
	}
	if err := b.Save(lens, "/files/hosts", "out"); err != nil {
		t.Fatal(err)
	}
	if string(mfs.writes["out"]) != "10.0.0.1 host\n" {
		t.Fatalf("saved %q", mfs.writes["out"])
	}
}

// registeredLens is a lens registered by hand, to prove a manual registration
// still wins over the engine's .aug corpus.
type registeredLens struct{}

func (registeredLens) Parse(string) (*engine.Node, error) { return &engine.Node{}, nil }
func (registeredLens) Build(*engine.Node) (string, error) { return "manual", nil }

// TestLensByNameSources covers all three ways LensByName can end: a hand
// registration, the engine's corpus, and a name that is in neither.
func TestLensByNameSources(t *testing.T) {
	// 1. A hand registration takes precedence.
	engine.Register("HandRegistered", registeredLens{})
	l, ok := LensByName("HandRegistered")
	if !ok {
		t.Fatal("a hand-registered lens must be found")
	}
	if out, _ := l.Build(nil); out != "manual" {
		t.Errorf("the registry lens should win, got %q", out)
	}

	// 2. Resolved from the engine's embedded .aug corpus.
	if _, ok := LensByName("Hosts"); !ok {
		t.Error("Hosts must resolve through the interpreter")
	}

	// 3. In neither: not found, no error surfaced.
	if _, ok := LensByName("NoSuchLensAnywhere"); ok {
		t.Error("an unknown lens must not be found")
	}
}
