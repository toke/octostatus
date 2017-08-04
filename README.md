# Octostatus

A very basic Octoprint status display

## Installation

```
go get github.com/toke/octostatus/octostatus
go install github.com/toke/octostatus/octostatus
```

## Usage

Usage: Just create a config (see example) and run `octostatus`.
Currently no command line switches are available.

The output can be configured by changing the template in the configuration.

## Examples

Example Config $HOME/.config/octoclient/config.yml
```yaml
---

name: "Octoprint client config"
version: 1
printer:
  default:
    baseUrl: http://octoprint.local
    apiKey: SECRETSTUFF
output:
  default:
    template: "{{.State}} {{.JobDetail.File.Name}} {{printf \"%.01f\" .Progress.Completion }}% ({{ .Progress.PrintTimeLeftString }})"
```


Example output:
```
> /home/toke/bin/octostatus                                          
Printing 1_Z_Motor_Mount_Right_5mm.gcode 37.2% (3h39m8s)
```
