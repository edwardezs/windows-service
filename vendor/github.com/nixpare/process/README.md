# Process

This package provides a method to spawn child processes,
like os/exec package, but with one huge difference: on Windows
you have the capabilty of **sending an Interrupt signal**, that is
you can simulate a `CTRL+C` event to the child process, which
can then decide whether to handle it or not.

This is **HUGE** because Windows does not support sending signals, but with
this package at least you can send an `os.Interrupt` at least.

The underlying structure used to spawn and manage the process is
the `exec.Cmd`: for managing input, output and lifecycle of the process,
its highly recommended to use the functions provided by the package.

The package also gives for Windows users some usefull functions
to recreate the behaviour via the various functions from the `kernel32.dll`.

# Behaviour

By default, each child process is spawned in a **new console/tty**. This is required
for Windows to work so it's also the default behaviour in Linux for a consistent
implementation.

When starting a process, these are the options you can provide ...

For the `stdin` argument:
 + if you provide `nil`, the process will listen on the standard input provided by the new
   environment, and no data can be sent through the pipe
 + if you provide `process.DevNull()`, the pipe will be created and usable, the child
   will not use the standard input of the new console and will only receive data through
   the pipe by the parent
 + for `any other value`, the pipe will be created and usable, and any data inside the
   io.Reader provided will be passed through the pipe. This means that if you pass
   os.Stdin, you can send data from the parent via the methods of the pipe and via
   your console natively

For the `stdout` and `stderr` arguments:
 + if you provide `nil`:
    + on **Windows**, the child will use the standard output / error of the new console,
      you can't capture any data
    + on **Unix-like systems**, the child will write to the /dev/null device and you can't
      capture any data
 + if you provide `process.DevNull()`, you will be able to access any data sent by the child,
   but only through the process methods (no real standard output / error)
 + for `any other value`, any data will be captured and written to the io.Writer provided.
   This means that if you pass os.Stdout/Stderr, you will also have the output sent to the parent console

# OS Compatibility

The package is obviously compatible with all operating systems,
but there are some differences:

## Non Windows OSes (Linux, Mac OS, etc)

Considering its natively supported to send signals between processes,
there is full support without any external program

## Windows

By default, in windows, the process is spawned in a different console (so has
different stdin, out and err from the parent) and is hidden, but this behaviour
can be reverted with the method Process.ShowWindow() by calling it before the start

This process might still fail, the reasons are still not known. For example, if the child
process you want to stop is created by the "go run" command, this will not work
