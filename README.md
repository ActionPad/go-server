# ActionPad Server
This is the repository for ActionPad Server, the desktop software that lets you control your computer using ActionPad on your iOS device.

## Setup
To compile ActionPad Server, you will first need to install Golang and GCC.
On MacOS, this can be done by installing the Xcode Command Line Tools, and then installing Golang.
On Windows, you should install MinGW-64, as the executable is 64-bit.
Your computer must be capable of running a Makefile.

### Mac
Then run:
```
make
```
To compile a standalone executable for macOS:
```
make apple-silicon
make apple-intel
make combine-universal
```

### Windows
To compile into a Windows .exe file, first run:
```
make win-manifest
```
Then, to compile a dev build of the server, run:
```
make win-exe
```
To compile a Windows production build (doesn't show command prompt window)
```
make win-exe-prod
```

Once this is done, you will now have an ActionPad server executable.

If you find a bug or want to contribute something, feel free to file an issue or make a pull request.
