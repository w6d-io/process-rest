# Process REST

[![codecov](https://codecov.io/gh/w6d-io/process-rest/branch/main/graph/badge.svg?token=PZDUENZE1U)](https://codecov.io/gh/w6d-io/process-rest)

The aim of this application is to trigger script by POST data.
The data posted will be record in a file then the file path will be set as parameter for all scripts

There is 3 kinds of script

- Pre  Script
- Main Script
- Post Script

The process will be installed on alpine 3.6 along with ([Dockerfile](https://github.com/w6d-io/kubectl/blob/main/Dockerfile))

- helm (v3)
- kubectl (v1.20.2)
- jq (1.5)
- yq
- bash (4.3.48)
- python (3.6.8)
- pip (21.0.1)
- gettext (0.19.8.1)
- git (2.13.7)
- make
- curl
- gawk

## Configuration

The script has to be in folder within the container
In [kubernetes](https://k8s.io) it can be done through [configmap](https://kubernetes.io/docs/concepts/configuration/configmap) or [secret](https://kubernetes.io/docs/concepts/configuration/secret)

## Examples

- TODO
