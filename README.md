# Shadow
This is an educational tool, **you** are responsible for abiding by any applicative laws in your region / country

## Description
Shadow is a command line tool and a pkg that allows obfuscation and encryption of shell scripts, encoding is currently done in base64, and encryption is handled through a substitution, the key for which must be provided at runtime

## Install

### Command line install
```
go install github.com/Joey574/shadow/v2/cmd/shadow@latest
```

### Pkg usage
```
import "github.com/Joey574/shadow/pkg"
```

## Usage
```
shadow script.sh -b -e
```
Encrypt and encode script.sh and write the results to stdout (key will always be printed to stdout)

```
shadow script.sh -b -o obfs.sh
```
Encode script.sh and write the results to obfs.sh
