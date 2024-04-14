# Jswap
Jswap allows you to quickly download and use different versions of the Java JDK via the command line.

It works on Windows, Linux and macOS and does not require admin permissions.

Jswap uses the the official [Adoptium API](https://api.adoptium.net/) to download the [Eclipse Temurin](https://adoptium.net) distribution of OpenJDK.

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
To download and install the latest available release of a given JDK major do this:
```bash
jswap get 21
```
To download and install a specific JDK release:
```bash
jswap get jdk-21.0.2+13
```
To list all the JDK releases available for download:
```bash
jswap releases
```
To use the latest installed release of a given major:
```bash
jswap use 21
```
To use a specific installed release:
```bash
jswap use jdk-21.0.2+13
```
To list all the installed JDKs:
```bash
jswap ls
```
