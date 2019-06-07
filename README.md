# go-server
New server application written in Go - Work in progress

## Setup
To compile ActionPad Server, you will first need to install Golang and GCC.
On MacOS, this can be done by installing the Xcode Command Line Tools.
On Windows, you can install MinGW-64 or any other GCC.
Your computer must be capable of running a Makefile.

Once you have Golang and GCC installed, you can install the go dependencies.

```
make deps
```

Then run:
```
make
```
To compile into a Windows .exe file, run:
```
make win-exe
```

Once this is done, you will now have an ActionPad server executable.

Keep in mind this is in early stages of development and does not have all the intended functionality yet. There may also be bugs.
