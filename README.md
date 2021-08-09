![Logo of the project](https://upload.wikimedia.org/wikipedia/commons/thumb/6/68/Eye_open_font_awesome.svg/200px-Eye_open_font_awesome.svg.png)

# Eye
> Simple file watcher

Eye watches file pattern regex and executes a command on changes.

## Synopsis

With an `.eyerc` file present.

```shell
eye
```

With command line arguments.

```shell
eye [-R] -p <regex_pattern> -c <command>
```

Eye will search for an `.eyerc` file in current directory or in the home
directory. If present it will load the settings.  Any command line
parameters will override the resource file settings.

### Example Usage

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
Flags:

* `-h` | `--help` - print usage information.
* `-p <regex_pattern>` - matching pattern for files to watch, mandatory.
* `-c <command>` - the command to execute on changes, mandatory.
* `-R` - do not recurse sub directories, optional.
* `-i <milliseconds>` - number of milli seconds to wait between checks,
  default value 1000 ms (1 s).

### Resource file

Resource file format is simple, just a `:` separated list of key and
value.  Lines starting with `!` is considered to be comments and empty
lines are ignored.

Example resource file with all possible settings:

```
! Eye resource file
pattern: \.go$
command: go test ./...
recursive: true
interval-millis: 500
```


## Installing

```shell
git clone https://github.com/callerobertsson/eye.git
cd eye
go install
```

## Licensing

The code in this project is licensed under GPLv3 license.
