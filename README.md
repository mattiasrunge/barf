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

## Installation

```bash
# Linux amd64
sudo bash -c 'curl -L https://github.com/mattiasrunge/barf/releases/latest/download/barf-linux-amd64.tar.gz | tar xvz -C /usr/local/bin'

# Linux ARMv5
sudo bash -c 'curl -L https://github.com/mattiasrunge/barf/releases/latest/download/barf-linux-arm5.tar.gz | tar xvz -C /usr/local/bin'

# Linux ARMv7
sudo bash -c 'curl -L https://github.com/mattiasrunge/barf/releases/latest/download/barf-linux-arm7.tar.gz | tar xvz -C /usr/local/bin'

```

## Usage

```bash
Usage: barf [OPTIONS] COMMAND [arg...]

A tool for doing robust file operations.

Options:
  -v, --version   Show the version and exit

Commands:
  monitor, m      monitors active operations
  copy, cp        copies files or folders

Run 'barf COMMAND --help' for more information on a command.
```

## Features

### Visibility

For longer running operations it is usually nice see what is happening. *barf* will present a nice progress visualization for each operation. If the CLI is terminated while an operation is still running it can be found again via the monitor flag.

<img src="./docs/svg/copy-monitor.svg?raw=true" />

### Background

*barf* starts a background process, if one is not already running, which runs all operations. This means that once the operation has been sent to the background process the CLI can be aborted and the operation will still continue.

### Resumable

If for some reason the background process dies all operations are stored in a journal and will be resumed when the background process starts again. Since *[rsync](https://rsync.samba.org/)* does synchronization, all operations should be resumable.

### Remote - SSH

Since [rsync](https://rsync.samba.org/) is used for the heavy lifting and [rsync](https://rsync.samba.org/) supports synchronizing with remote systems via SSH, *barf* supports that too.

### Logging

The output of all operations are stored in log files to enable debugging and verification.

## Operations

### Copy

<small>*Append files and folders from source at remote*</small>

Copy files and folders to a new location. Since [rsync](https://rsync.samba.org/) is used, what actually happens under the hood is a synchronization which might behave slightly different than an ordinary *cp* command. It can probably be seen more like an append, take everything local and put remote, overwrite if necessary but remove nothing.

### Move (Not implemented yet)

<small>*Append files and folders from source at remote and then remove source*</small>

Move is essentially the same as Copy but it removes the source files and folders after the Copy has completed successfully.

### Push (Not implemented yet)

<small>*Replace files and folders from source at remote and delete files and folders from remote that are not present at source*</small>

Similar to Copy but makes the remove exactly the same. Only different from Copy if folders are involved. Files and folders missing at the source will be removed from the remote.

### Pull (Not implemented yet)

<small>*Replace files and folders from remote at source and delete files and folders from source that are not present at remove*</small>

Exactly like push but with source and remove inverted.

### Backup (Not implemented yet)

TODO

### Restore (Not implemented yet)

TODO

## Goals

Every project needs a reason for being.

* Good progress visualization for file operations
* Same syntax for local and remote operations
* No need for screen to keep operations running if the shell is closed
* Easy backup and restore functionality

## Technical

* Written in [Go](https://golang.org/)
* Uses a domain socket for communication (CLI -> background process)
* Stores state and logs under ```~/.config/barf```
* Uses the installed version of [rsync](https://rsync.samba.org/)
* Only tested on Linux, but might work on other systems
