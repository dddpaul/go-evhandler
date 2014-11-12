go-evhandler
=========

Simple key press listener and handler written in Go.

Usage:

```
go run evhandler.go [device name]
```

You can bind actions to keys in configuration file like this:

```
[params]
    device = "/dev/input/event3"

[actions]
    KEY_STOPCD = "echo stop"
    KEY_PLAYPAUSE = "echo play"
```

Sincere gratitude for authors of these open-source software:
   
* [golang-evdev](https://github.com/gvalkov/golang-evdev) - Go bindings for the linux input subsystem.
