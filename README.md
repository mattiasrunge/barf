# BAckground Robust File transfer (BARF)

![Latest release](https://github.com/mattiasrunge/barf/workflows/Build%20and%20release%20binaries/badge.svg)

BARF is a CLI tool written for doing robust file operations. Under the hood it uses [rsync](https://rsync.samba.org/) to do the heavy lifting.

<img src="./docs/svg/copy-normal.svg?raw=true" />

## Table of contents

- [Project status](#project-status)
- [Installation](#installation)
- [Usage](#usage)
- [Features](#features)
- [Operations](#operations)
- [Goals](#goals)
- [Technical](#technical)
- [Releases](https://github.com/mattiasrunge/barf/releases)

## Project status

This project is still in the early phases, try if you wish but not all of the features and operations are implemented yet.

### TODO

- [ ] Implement `push` operation
- [ ] Implement `pull` operation
- [ ] Implement `backup` operation
- [ ] Implement `restore` operation
- [ ] Restart daemon if version mismatch
- [ ] Option to output JSON when used in scripts
- [ ] Possibility to output live logging from rsync, like tail
- [ ] Create a better CI workflow
- [ ] Write tests
- [ ] Handle non-color terminals gracefully
- [ ] Improve scp/rm stuff in `move` operation
- [ ] Bash auto complete of CLI arguments
- [ ] Remove EventBus dep, implement own for the limited functionality used
- [ ] Find a good bytesize dep which supports parse and humanize as we need it
- [ ] Improve daemon logging to output time etc, via the log.Logger
- [ ] Improve daemon logging to print operation creation and statuses
- [ ] Investigate if we should be a download tools as well... use curl or wget and present a nice progressbar?

## Installation

```bash
# Using the install script, installing in /usr/local/bin
curl -sL https://raw.githubusercontent.com/mattiasrunge/barf/main/scripts/install.sh | sudo -E bash -

# For manual installation

# Linux amd64
sudo bash -c 'curl -L https://github.com/mattiasrunge/barf/releases/latest/download/barf-linux-amd64.tar.gz | tar xvz -C /usr/local/bin'

# Linux ARMv5
sudo bash -c 'curl -L https://github.com/mattiasrunge/barf/releases/latest/download/barf-linux-arm5.tar.gz | tar xvz -C /usr/local/bin'

# Linux ARMv7
sudo bash -c 'curl -L https://github.com/mattiasrunge/barf/releases/latest/download/barf-linux-arm7.tar.gz | tar xvz -C /usr/local/bin'

```

## Usage

```plain_text
Usage: barf [OPTIONS] COMMAND [arg...]

A tool for doing robust file operations.

Options:
  -v, --version   Show the version and exit
  -w, --width,    terminal width to use, if not set (or zero) it will be auto detected and failing that set to 132 (default 0)

Commands:
  list, l         list active operations
  monitor, m      monitors active operations
  abort, a        aborts an active operation
  copy, cp        copies files or directories
  move, mv        moves files or directories
  update, u       check for updates

Run 'barf COMMAND --help' for more information on a command.

```

## Features

### Background

*barf* starts a background process, if one is not already running, which runs all operations. This means that once the operation has been sent to the background process the CLI can be aborted and the operation will still continue.

<img src="./docs/svg/copy-monitor.svg?raw=true" />

### Visibility

For longer running operations it is usually nice see what is happening. *barf* will present a nice progress visualization for each operation. If the CLI is terminated while an operation is still running it can be found again via the monitor flag.

<img src="./docs/svg/copy-monitor-many.svg?raw=true" />

### Remote - SSH

Since [rsync](https://rsync.samba.org/) is used for the heavy lifting and [rsync](https://rsync.samba.org/) supports synchronizing with remote systems via SSH, *barf* supports that too.

<img src="./docs/svg/copy-remote.svg?raw=true" />

### Resumable

If for some reason the background process dies all operations are stored in a journal and will be restarted when the background process starts again. Since *[rsync](https://rsync.samba.org/)* does synchronization, all operations should be resumable. The output of all operations are stored in log files to enable debugging and verification.

<img src="./docs/svg/daemon-journal.svg?raw=true" />

## Operations

*barf* supports different kinds of operations which are described here.

| Operation | Description |
| --- | --- |
| `copy` | <small>*Short: Append files and directories from the source at the remote*<br><br>Copies files and directories to the remote. Since [rsync](https://rsync.samba.org/) is used, what actually happens under the hood is a synchronization which might behave slightly different than an ordinary `cp` command. It can probably be seen more like appending; take everything at the source and put it at the remote, overwrite if necessary but remove nothing.</small> |
| `move` | <small>*Short: Append files and directories from the source at the remote and remove at source*<br><br>Move is essentially the same as `copy` but it removes the source files and directories after the copy has completed successfully. Thus if copying locally there will be two copies of the files taking up space until the end when the source is removed.</small> |
| `push` | <small>*Short: Make the remote exactly like the source*<br><br>Similar to `copy` but makes the remote exactly the same as the source. Files and directories found remote but not at the source will be deleted.</small> |
| `pull` | <small>*Short: Make the source exactly like the remote*<br><br>Exactly like `push` but with source and remote inverted.</small> |
| `backup` | <small>*Short: Create a new copy of the source at the remote*<br><br>Backup will create a new directory, named as the current date and time. It will then create hard links from the previous backup if there is one. After that it will do a `push` operation to the new directory.</small> |
| `restore` | <small>*Short: Will restore the specified remote backup at the source*<br><br>Restore will take the specified backup at the remote and `pull` it to the source, overwriting everything at the source in the process.</small> |

## Goals

Every project needs a reason for being.

- Good progress visualization for file operations
- Same syntax for local and remote operations
- No need for screen to keep operations running if the shell is closed
- Easy backup and restore functionality

## Technical

- Written in [Go](https://golang.org/)
- Uses a domain socket for communication (CLI -> background process)
- Stores state and logs under ```~/.config/barf```
- Uses the installed version of [rsync](https://rsync.samba.org/), make sure there is one
- Well defined socket protocol, allowing for other types of clients
- Only tested on Linux, but might work on other systems
