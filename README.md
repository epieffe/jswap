# Jswap
Jswap allows you to quickly download and switch between different versions of the Java JDK via the command line. It works on Windows, Linux and macOS and does not require admin permissions.

Jswap uses the the [Adoptium API](https://api.adoptium.net/) to download the [Eclipse Temurin](https://adoptium.net) distribution of OpenJDK.

## Installation
Prebuilt binaries for Windows, Linux and macOS can be downloaded from the [Releases](https://github.com/epieffe/jswap/releases) page.

### Linux and macOS
1. Download the tar.gz archive for your OS and architecture, then extract the archive you just downloaded into a folder that is in your PATH, such as `/usr/local/bin` or `~/.local/bin`.
2. Append the following lines to the `.bashrc` file in your home directory (or `.zshrc` file if you are using zsh):
   ```bash
   export JAVA_HOME=~/.jswap/current-jdk
   export PATH=$JAVA_HOME/bin:$PATH
   ```

### Windows
Simply download and run the .exe installer. This will set (and eventually override) the `JAVA_HOME` environment variable for te current user and add `jswap`, `java` and `javac` to the current user PATH.

Jswap can then be easily uninstalled from the "Add/Remove Programs" section in the Windows Control Panel.

#### The Unix way
Alternatively, if you have Bash installed on Windows, you can install Jswap in the same way you would install it on Linux or macOS.
The only difference is that the Jswap Windows executable stores data in the LocalAppData folder rather than the home directory.

This will make Jswap available only in the Bash shell and will not touch the Windows environment variables.

1. Download the zip archive for Windows and extract it into a folder that is in your PATH, or into a new folder and then add it to you PATH.
2. Append the following lines to the `.bashrc` file in your home directory (create the file if it does not exists):
    ```bash
    export JAVA_HOME=$LOCALAPPDATA/Jswap/current-jdk
   export PATH=$JAVA_HOME/bin:$PATH
    ```

## Usage
Download and install the latest available release of a given JDK major:
```bash
jswap get 21
```

Modify PATH and JAVA_HOME to use the latest installed release of a given major:
```bash
jswap set 21
```

List all the installed JDKs:
```bash
jswap ls
```

Run `jswap --help` for more information.

## Build from sources
To build the Jswap executable from sources you need the following dependencies:
- Go
- Make
- Git

To build Jswap for your OS and architecture run this:
```bash
make build
```

You can also target different platforms:
```bash
make linux-amd64 # Linux x64
make win-amd64 # Windows x64
make mac-amd64 # macOS with Intel CPU
make mac-arm64 # macOS with ARM CPU
```

The executables will be found in the `build` directory.

### Windows installer
To build the Windows installer you need [NSIS](https://nsis.sourceforge.io/Main_Page) and the [EnVar plugin for NSIS](https://nsis.sourceforge.io/EnVar_plug-in).

NSIS is also available on Linux and can be easily installed via apt on Ubuntu and Debian:
```bash
sudo apt install nsis
```
Place the EnVar plugin dll files in `/usr/share/nsis/Plugins`.

Run the following command to build the Jswap Windows installer:
```bash
make win-installer
```

### Build using Docker
You can build Jswap and the Windows installer for any supported platform using Docker, without any other dependency.

First, build the `jswap-builder` Docker image:
```bash
docker build -t jswap-builder .
```

Run the following command to create a Docker container that builds Jswap and the Windows installer for all the supported platforms from the sources in the current directory:
```bash
docker run -v $PWD:/app --name jswap-builder jswap-builder
```

The `jswap-builder` image runs the Make target `all` by default. If you want to run a different Make target you can pass it as a command to `docker run`. For example, if you want to build only the Windows installer run this:
```bash
docker run -v $PWD:/app --name jswap-builder jswap-builder win-installer
```

For subsequent builds you can reuse the previously created container to improve build times:
```bash
docker start -a jswap-builder
```
