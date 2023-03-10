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

Gauge values are random on each call. Counters are initialized at `0` and are incremented randomly on each call, same as Histograms and Summaries.

Random numbers are between 1 and 10 by default. The fields `min` and `max` can be set to define these for the random number generator. The generator is inclusive, meaning `min: 1` and `max: 1` can be used.

```json
[
  {
  "metric":"foobar",
  "type": "counter",
  "min": 1,
  "max": 10,
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
},
  {
  "metric":"foo_requests",
  "type": "histogram",
  "labels": {
    "job": ["foo"],
    "instance": ["fuu.1", "fuu.2"]
  },
  "le": [0.1, 0.5, 0.99]
},
  {
  "metric":"bar_seconds",
  "type": "summary",
  "labels": {
    "job": ["foo"],
    "instance": ["fuu.1", "fuu.2"]
  },
  "quantile": [0, 0.25, 0.5, 1]
}
]
```
