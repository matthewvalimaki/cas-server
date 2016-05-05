# cas-server
Implementation of [JASIG CAS protocol] in Go lang. Supports all protocol versions (v1, v2 and v3).

## Design goals
* Easier to develop for and easier to work with than [JASIG CAS].

## Features


## Configuration
You can configure certain aspects of `cas-server` with command line arguments but majority of the configuration will
require a [TOML] formatted configuration file.

[JASIG CAS protocol]: https://jasig.github.io/cas/4.2.x/protocol/CAS-Protocol-Specification.html
[JASIG CAS]: https://github.com/Jasig/cas
[TOML]: https://github.com/toml-lang/toml