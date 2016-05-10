# cas-server
Implementation of [JASIG CAS protocol] in Go lang. Supports all protocol versions (v1, v2 and v3).

## Configuration
You can configure certain aspects of `cas-server` with command line arguments but majority of the configuration will
require a [TOML] formatted configuration file. For an example configuration please see `config/config.toml.example`.

## Running
`cas-server -config /etc/cas-server/config.toml`

## Features
* Storage
  * In memory

## Project goals
* Not to replace [JASIG CAS] but to offer competition with reduced feature set.
* Easier to develop for and easier to work with than [JASIG CAS].
  * Minimal dependencies on 3rd party Go libraries.
  * This Readme has to be enough to get `cas-server` running correctly and securely.

[JASIG CAS protocol]: https://jasig.github.io/cas/4.2.x/protocol/CAS-Protocol-Specification.html
[JASIG CAS]: https://github.com/Jasig/cas
[TOML]: https://github.com/toml-lang/toml
