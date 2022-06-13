# envdir

This is a straight port of the `envdir` C source code from the daemontools
repository (which has also been made available inside of the "testdata"
directory).

There is one external dependency (on golang.org/x/sys/unix) and there are no
additional features, versus the original tool.

### Errata

- The original envdir would only allow 256 bytes in an environment variable, this
program allows 1024 bytes.

- The original envdir program would only trim tabs and newlines from the end of
an environment variable, we trim all runes that unicode.IsSpace reports true on.
