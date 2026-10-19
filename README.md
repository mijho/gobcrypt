### GobCrypt

#### A Bcrypt hash password generator written in golang

[![test-lint](https://github.com/mijho/gobcrypt/actions/workflows/test-lint.yml/badge.svg?branch=main)](https://github.com/mijho/gobcrypt/actions/workflows/test-lint.yml)
[![build-release](https://github.com/mijho/gobcrypt/actions/workflows/build-release.yml/badge.svg?branch=main)](https://github.com/mijho/gobcrypt/actions/workflows/build-release.yml)
![go-report](https://goreportcard.com/badge/github.com/mijho/gobcrypt)

##### Installation Instructions

With go get:
```
go install github.com/mijho/gobcrypt@latest
```

From source:
```
git clone github.com/mijho/gobcrypt.git
make build
```

Or head to the releases page to download the prebuild binary for your system.

## Pass/Hash Generation

```
NAME:
   gobcrypt generate - generate pass/hash pairs

USAGE:
   gobcrypt generate [command options] [arguments...]

OPTIONS:
   --number value, -n value       Specify the number of hashes to create (default: 1) (default: 1)
   --password value, -p value     Password to hash
   --length value, -l value       Specify the length of password required (default: 16) (default: 16)
   --input-file value, -i value   Specify a file to read passwords from
   --output-file value, -o value  Specify a file to write out the pass/hash to
   --cost value, -c value         Specify the cost to use (Min: 4, Max: 31) (default: 14)
   --help, -h                     show help (default: false)
```

## Pass/Hash Validation

```
NAME:
   gobcrypt validate - validate pass/hash pairs

USAGE:
   gobcrypt validate [command options] [arguments...]

OPTIONS:
   --password value, -p value     Password to hash
   --hash value, --hp value       Hash to verify
   --input-file value, -i value   Specify a file to read passwords from
   --output-file value, -o value  Specify a file to write out the pass/hash to
   --help, -h                     show help (default: false)
```

## Examples

Create a password and hash (cost 14 by default):
```
$ ./gobcrypt generate
HxQhL1UfzkflKeEE $2a$14$xsMB6.ZbiK701eHhmaLjEeasbkmmS9y.GuP4bFEPpEebAh4bvN4IC
```
Create a hash of a specific password
```
$ ./gobcrypt generate -p supergoodpass
supergoodpass $2a$14$ztiAYzgVD92wuxb7213aBOYl5RfR3dCSeKj9UFuyzaedoHQM///tq
```

Check the hash is valid. `validate` exits with status 0 when every pair matches and 1 when any pair fails or is malformed, so it can be used in scripts:
```
$ ./gobcrypt validate -p supergoodpass -hp '$2a$14$.pR1K36frbmfWo1NjspjTOcDELTJo.LCCCD4Ft5erL1HLLwMbwPVO'
MATCH: PASS, password: supergoodpass, hash: $2a$14$.pR1K36frbmfWo1NjspjTOcDELTJo.LCCCD4Ft5erL1HLLwMbwPVO
$ echo $?
0
$ ./gobcrypt validate -p wrongpass -hp '$2a$14$.pR1K36frbmfWo1NjspjTOcDELTJo.LCCCD4Ft5erL1HLLwMbwPVO'
MATCH: FAIL, password: wrongpass, hash: $2a$14$.pR1K36frbmfWo1NjspjTOcDELTJo.LCCCD4Ft5erL1HLLwMbwPVO
1 of 1 pair(s) did not validate
$ echo $?
1
```

Create passwords of a specific length (16 chars by default)
```
$ ./gobcrypt generate -l 20
6u8zxKM2Q9khYliAhbjP $2a$14$k3C7nYjjYbGV/1WZXV2ZCOhz5.TqVjEPxH1T9e/ONj6z30DYTmAA2
```

Create a bulk amount of password hash combinations
```
$ ./gobcrypt generate -n 2 -c 4
A8Fb5D359PE1oofD $2a$04$yyXvZdnygMFhUruixLF5cOKj/3sx44lFeQdR36Eni8q59LFfMVlUK
5GJrZSqhDXVxPCRq $2a$04$RuphMwgGfp61Tl0132oJh.08SFFeDyQpvowBusXYlR4Ji1df8QYWq
```

Create a bulk amount and write the output to file. The file is created with mode `0600` because it contains plaintext passwords.
```
$ ./gobcrypt generate -o outfile.txt && cat outfile.txt
U7Bs6wLBqcqBuVSZ $2a$14$TYhWm3aSNl8eSb/xqSzvn.B2PhU3vxAdtiThNPiFqA8IcNBg.hT3W
$ ./gobcrypt validate -i outfile.txt
MATCH: PASS, password: U7Bs6wLBqcqBuVSZ, hash: $2a$14$TYhWm3aSNl8eSb/xqSzvn.B2PhU3vxAdtiThNPiFqA8IcNBg.hT3W
```

Read passwords from file to hash or validate. Blank lines are skipped.
```
$ cat infile.txt
password
hotdamn
12hfsklfh
forbarbaz
$ ./gobcrypt generate -i infile.txt
password $2a$14$9sdf57Pn47E0PWDfSnLB.eyQQLebdUmweSvfmOTeqdoBSkj9BXwra
hotdamn $2a$14$4qgEbl9WHX1EaBlbahbDZONeiIyForWSrJakPXRkMweDgkN13DJk.
12hfsklfh $2a$14$KSWk2l630PQt4bMF4WblAesTN7bg5q5LMvQ9oA8MsNOSI6W4IacVS
forbarbaz $2a$14$yb55LAVxeV/B9ap5hE2YG.t1jiEf6o.peK3voniAsHSHwHlJacvGC
```

Manually set the cost of the BCrypt hashing process (4 to 31)
```
$ ./gobcrypt generate -p supergoodpass -c 10
supergoodpass $2a$10$m8XlaLF.kibYmab5fbGjbOfYjMCZVbQcWyQwEcRIdoFcPl0Otw3NC
```

## Limitations

bcrypt only uses the first 72 bytes of a password, and this tool refuses longer ones rather than silently truncating them:
- `generate -p` and random generation with `-l` above 72 fail with an error.
- `generate -i` skips any line longer than 72 bytes, writes the rest, prints a warning to stderr, and exits with status 1.
