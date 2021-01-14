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

If for some reason the background process dies all operations are stored in a journal and will be restarted when the background process starts again. Since *[rsync](https://rsync.samba.org/)* does synchronization, all operations should be resumable.

The output of all operations are stored in log files to enable debugging and verification.

## Operations

*barf* supports different kinds of operations which are described here.

### Copy <span style="color: grey; font-size: 12px; padding-left: 5px">Append files and folders from source at remote</span>

Copy files and folders to a new location. Since [rsync](https://rsync.samba.org/) is used, what actually happens under the hood is a synchronization which might behave slightly different than an ordinary *cp* command. It can probably be seen more like an append, take everything local and put remote, overwrite if necessary but remove nothing.

### Move <span style="color: grey; font-size: 12px; padding-left: 5px">Append files and folders from source at remote and then remove source</span>

*(Not implemented yet)*

Move is essentially the same as Copy but it removes the source files and folders after the Copy has completed successfully.

### Push <span style="color: grey; font-size: 12px; padding-left: 5px">Replace files and folders from source at remote and delete files and folders from remote that are not present at source</span>

*(Not implemented yet)*

Similar to Copy but makes the remove exactly the same. Only different from Copy if folders are involved. Files and folders missing at the source will be removed from the remote.

### Pull <span style="color: grey; font-size: 12px; padding-left: 5px">Replace files and folders from remote at source and delete files and folders from source that are not present at remove</span>

*(Not implemented yet)*

Exactly like push but with source and remove inverted.

### Backup <span style="color: grey; font-size: 12px; padding-left: 5px">Create a copy of local in a new backup folder at the remote</span>

*(Not implemented yet)*

Backup will create a new directory, named as the current date and time. It will then do create hard links from the previous backup if there is one. After that it will do a Push operation to that folder.

### Restore <span style="color: grey; font-size: 12px; padding-left: 5px">Will do a Pull from the specified remote backup to local</span>

*(Not implemented yet)*

Restore will take the specified backup at the remote and Pull it to the local overwriting local in the process if it is not empty.

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
- Uses the installed version of [rsync](https://rsync.samba.org/)
- Only tested on Linux, but might work on other systems
