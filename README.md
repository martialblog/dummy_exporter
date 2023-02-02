# Dummy exporter

The dummy exporter is a customizable metric exporter for testing, demonstrations and learning.

## Running this software

### From binaries

Download the most suitable binary from the releases tab.

Then:

    ./dummy_exporter <flags>

## Building the software

### Local Build

    make build

## Configuration

Dummy exporter is configured via a configuration file and command-line flags (such as what configuration file to load, what port to listen on).

To view all available command-line flags, run `./dummy_exporter -h`.

To specify which configuration file to load, use the `--config.file` flag.

## Configuration File

The configuration file is written in JSON (to rely solely on the Golang standard library).

Dummy exporter will render metrics for each `labels` permutation.

Gauge values are random on each call.

Counter are initialized at `0` and are incremented randomly on each call.

```json
[
  {
  "metric":"foobar",
  "type": "counter",
  "labels": {
    "job": ["foo"],
    "instance": ["foo.1", "foo.2"]
  }
},
  {
  "metric":"barfoo",
  "type": "gauge",
  "labels": {
    "job": ["foo"],
    "instance": ["bar.1", "bar.2"],
    "code": ["200"]
  }
}
]
```
