![Logo of the project](https://upload.wikimedia.org/wikipedia/commons/thumb/6/68/Eye_open_font_awesome.svg/200px-Eye_open_font_awesome.svg.png)

# Eye
> Simple file watcher

Eye watches file pattern regex and executes a command on changes.

Watch interval is hard coded to 1 second.

## Synopsis

```shell
eye [-R] -p <regex_pattern> -c <command>
```

Flags:

* `-R` - do not recurse sub directories, optional.
* `-p <regex_pattern>` - matching pattern for files to watch, mandatory.
* `-c <command>` - the command to execute on changes, mandatory.
* `-i <milliseconds> - number of milli seconds to wait between checks,
  default value 1000 ms (1 s).
* `-h` | `--help` - print usage information.

### Examples

```shell
eye -r -p '\.go$' -c 'go build'
```

Watch files in current directory with extension `.go` and execute `go build` on changes:

To use for Golang development the `eye-go.sh` can be used as command. 

Copy `eye-go.sh` to your Golang project directory and change it to suit your needs. Then issue the command:

```shell
eye -p '\.go$' -c './eye-go.sh -i 500'
```

Eye will watch all Go files in current directory and below for changes
every half second. If any matching files changes the `eye-go.sh` script
will be executed.  The script will build, lint, vet, and look for TODOs
in the project.

## Installing

```shell
git clone https://github.com/callerobertsson/eye.git
cd eye
go install
```

## Licensing

The code in this project is licensed under GPLv3 license.
