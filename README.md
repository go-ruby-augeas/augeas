# go-ruby-augeas/augeas

Pure-Go (CGO=0) adapter presenting the
[`ruby-augeas`](https://github.com/hercules-team/ruby-augeas) gem API over the
configuration-editing engine [`github.com/go-augeas/augeas`](https://github.com/go-augeas/augeas).

It mirrors the engine+adapter split used by
`go-facter`/`go-ruby-facter` and `go-hiera`/`go-ruby-hiera`: all Augeas
semantics live in the engine; this package is a thin, faithful, importable
rename of that surface to the gem's method names, with no Ruby-runtime
dependency.

```go
a := augeas.Open("/", "", augeas.None)
lens, _ := augeas.LensByName("Hosts")
_ = a.TextStore(lens, "/files/etc/hosts", "127.0.0.1 localhost\n")

v, _ := a.Get("/files/etc/hosts/1/canonical") // "localhost"
_ = a.Set("/files/etc/hosts/1/alias", "loopback")
n := a.Rm("/files/etc/hosts/1/alias")         // 1
```

## API

Constructors `New()` and `Open(root, loadpath string, flags int)` plus the
`AUG_*` flag constants (`None`, `SaveBackup`, `SaveNewFile`, `TypeCheck`,
`NoStdinc`, `SaveNoop`, `NoLoad`, `NoModlAutoload`, `EnableSpan`).

Methods, each delegating to the engine: `Get`, `Exists`, `Set`, `Setm`,
`Insert`, `Rm`, `Mv`, `Match`, `Label`, `Defvar`, `Defnode`, `TextStore`,
`TextRetrieve`, `Load`, `Save`, `Span`, `Error`, plus `Root`, `LoadPath`,
`Flags`, `Engine` and `SetFileSystem`.

The `Lens`, `Span` and `FileSystem` types and `LensByName` are re-exported from
the engine so callers need not import it directly.

## Scope

This adapter adds no new Augeas behaviour; the supported path constructs, the
built-in lenses (Hosts, Fstab, Shellvars/Simplevars, Ini/Keyvalue) and the
explicit list of **deferred** features (the ~200-lens catalogue, the `.aug` lens
DSL, span tracking, extended path functions, `aug_load` autodetection) are
documented in the [engine README](https://github.com/go-augeas/augeas#deferred-not-yet-implemented--honest-scope).

Because the pure-Go engine performs no file autoload, `Open` records its
arguments but reads nothing until `Load` is called explicitly.

## License

BSD-3-Clause. See [LICENSE](LICENSE).
