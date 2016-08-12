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

## Installing

```shell
git clone https://github.com/callerobertsson/eye.git
cd eye
go install
```

## Licensing

The code in this project is licensed under MIT license.
