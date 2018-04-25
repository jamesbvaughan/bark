# bark [![build status](https://travis-ci.org/jamesbvaughan/bark.svg)](https://travis-ci.org/jamesbvaughan/bark)

bark is a bookmarking tool, similar to
[Pocket](https://getpocket.com/).
It has a simple command line interface as well as a web interface.

## Installation

The best way to install bark right now is to

1. Download a precompiled binary from [the releases page](https://github.com/jamesbvaughan/bark/releases).
2. Place it in your `PATH`.
3. Make it executable.

Once you've done that you're ready to go!
Just run `bark add https://github.com/jamesbvaughan/bark`
to add your first bookmark!

Alternatively, you can use bark from the web interface, which looks like this:
![web interface screenshot](https://raw.githubusercontent.com/jamesbvaughan/bark/master/web-ui.png "bark web interface")
and can be served locally by running `bark serve`.

## Usage

### Add a bookmark

```sh
bark add <URL>
```

### List bookmarks

```sh
bark list
```

### Open a bookmark

```sh
bark open <ID>
```

### Archive a bookmark

```sh
bark archive <ID>
```

### Start the webserver for the web interface

```sh
bark serve
```

### Permanently delete a bookmark

```sh
bark delete <ID>
```

### Get help

```sh
bark help
```

### Get help for a specific command

```sh
bark help <COMMAND>
```

## [Licence](LICENSE)

The MIT License (MIT)

Copyright (c) 2018 James Vaughan
