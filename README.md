# envdir

This is a straight port of the `envdir` C source code from the daemontools
repository (which has also been made available inside of the "testdata"
directory): https://cr.yp.to/daemontools/envdir.html

There is one external dependency (on golang.org/x/sys/unix) and there are no
additional features, versus the original tool.

### Usage

```
envdir dir childprog [arg1 arg2 arg3]
```

dir is a single argument. childprog is the child program, which can have as many
arguments as you like.

envdir sets various environment variables as specified by files in the directory
named dir. It then runs childprog.

If dir contains a file named s whose contents are t, envdir removes an
environment variable named s if one exists, and then adds an environment
variable named s with value t. The name s must not contain =. Spaces and tabs at
the end of t are removed. Nulls in t are changed to newlines in the environment
variable.

If the file s is completely empty (0 bytes long), envdir removes an environment
variable named s if one exists, without adding a new variable.

envdir exits 111 if it has trouble reading dir, if it runs out of memory for
environment variables, or if it cannot run childprog. Otherwise its exit code is
the same as that of childprog.

### Installation

Install from source:

```
go install github.com/kevinburke/envdir@v0.2
```

Or, on Macs you can install via Homebrew:

```
brew install kevinburke/safe/envdir
```

### Errata

- The original envdir would only allow 256 bytes in an environment variable, this
program allows 1024 bytes.

- The original envdir program would only read the first line of an environment
variable, where we allow newlines.

- The original envdir program would only trim tabs and newlines from the end of
an environment variable, we trim all runes that unicode.IsSpace reports true on.
