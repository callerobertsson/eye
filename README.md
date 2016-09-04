![Logo of the project](https://upload.wikimedia.org/wikipedia/commons/thumb/6/68/Eye_open_font_awesome.svg/200px-Eye_open_font_awesome.svg.png)

# Eye
> Simple file watcher

Eye watches file pattern regex and executes a command on changes.

## Synopsis

```shell
eye [-r] -p <regex_pattern> -c <command>
```

Flags:

* `-r` - recurse sub directories.
* `-p <regex_pattern>` - matching pattern for files to watch.
* `-c <command>` - the command to execute on changes.

### Examples

```shell
eye -r -p '\.go$' -c 'go build'
```

Watch files in current directory with extension `.go` and execute `go build` on changes:

To use for Golang development the `eye-go.sh` can be used as command. 

Copy `eye-go.sh` to your Golang project directory and change it to suit your needs. Then issue the command:

```shell
eye -r -p '\.go$' -c './eye-go.sh'
```

Eye will watch all Go files in current directory and below for changes and issue the `eye-go.sh` script if anything changes.  The script will build, lint, and vet the project.

## Installing

```shell
git clone https://github.com/callerobertsson/eye.git
cd eye
go install
```

## Licensing

The code in this project is licensed under GPLv3 license.
