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
   --cost value, -c value         Specify the cost to use (Min: 4, Max: 31) (default: 14) (default: 0)
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

Create a password and hash:
```
$ ./gobcrypt generate
f9fFFh7LcAaGJ89C $2a$10$pPbKZkrpHmJg5vZVGljUV.3N/gatYJN4Iv3bAnJ7b7SKgQTfZCnhS
```
Create a hash of a specific password
```
$ ./gobcrypt generate -p supergoodpass
supergoodpass $2a$10$KGm2jtacLvCXT2uitl8aQOPiDZBX8VmXIS9dZS5aAtQsQAVYetwiG
```

Check the hash is valid:
```
$ ./gobcrypt validate -p supergoodpass -hp '$2a$10$KGm2jtacLvCXT2uitl8aQOPiDZBX8VmXIS9dZS5aAtQsQAVYetwiG'
MATCH: PASS, password: supergoodpass, hash: $2a$10$KGm2jtacLvCXT2uitl8aQOPiDZBX8VmXIS9dZS5aAtQsQAVYetwiG
```

Create passwords of a specific length (15 chars by default)
```
$ ./gobcrypt generate -l 20                                                                              
6DG4D7WZnA86kD2iQnnD $2a$10$IHrYPjzwBDG2BYTIcOBFlutNHEB2mBmpP6QG8ZF4iCW8Bn3suToMu
```

Create a bulk amount of password hash combinations
```
$ ./gobcrypt generate -n 7
EQKEBdAYXPZNwMM0 $2a$10$.iB7vbaKj2L7cYpTTa/.nuHe/BKQEL3uXgqduBOpGVhvuHkSz0yoS
y8mEjf6uIYx3QyLs $2a$10$QpSqInArVO9A5s/vpkgW1eOgtdLV2TJw1dqD3aQjYF0kbX0FpDp/e
akA95A1qiwmoPcjw $2a$10$RhKxTNn.2CPB7ertSScKFO30Fqcg9xAJH.J2X8PspWLm31D2O1wR6
J5H58PoVxR6j8n3s $2a$10$WAENUArAzRSo8zXgFbB4COUKv6vhDbPyTkFIqXVr5JOf5ksjspC2K
o1xb8hiqPrXhEmCT $2a$10$Io7hlQuc2C1Z9A4IxVDSxO6cx2SCql1Gb65/YCb4Xj.djeGa1bDSS
mkI6LG4Yc0d4drvE $2a$10$L.76ZOQNZkhbjQX7GzupzO4Ke/VEdlCwwv4q5vzHX5/ojEH95RkPi
t3LjJcaHg2nEY6mZ $2a$10$mduGQTccWkQ9dQ6n1LlBxeII1IlfD6uB51Ewy/MTVQ0q88DtMHTjO
```

Create a bulk amount and write the output to file
```
$ ./gobcrypt generate -o outfile.txt && cat outfile.txt                                                 
uLTAPhWdIqGq3vPY $2a$10$zCb9z3melGqqw8I1aancVeMBS7xzXIaAAGsaDbqpj/AO3VP1/XNEi
$ ./gobcrypt validate -p supergoodpass -hp '$2a$10$KGm2jtacLvCXT2uitl8aQOPiDZBX8VmXIS9dZS5aAtQsQAVYetwiG' -o results.txt && cat results.txt
MATCH: PASS, password: supergoodpass, hash: $2a$10$KGm2jtacLvCXT2uitl8aQOPiDZBX8VmXIS9dZS5aAtQsQAVYetwiG
```

Read passwords from file to hash or validate
```
$ cat infile.txt
password
hotdamn
12hfsklfh
forbarbaz
$ ./gobcrypt generate -i infile.txt                        
password $2a$10$7Ttcr9C/6QkBnyewt/Ljg.HRpfoil96XSjTve7EXy5XOBiZB5Zr4a
hotdamn $2a$10$072dU17Zf/M.SbjAP.zhNer39Z7fPOUWHOTm6.849OzA9f3NJFaGu
12hfsklfh $2a$10$tVyXGyJ674PVJlh.dDd.seOA/RJRbDzPAlxFNjKp1Z3mVpOtOTJIG
forbarbaz $2a$10$bboUUJP0Fdp9oWBSFGN.GuV/VfBgYuHQfpQ3JcElzMH7ZeVCes8w2
$ ./gobcrypt validate -i outfile.txt 
MATCH: PASS, password: uLTAPhWdIqGq3vPY, hash: $2a$10$zCb9z3melGqqw8I1aancVeMBS7xzXIaAAGsaDbqpj/AO3VP1/XNEi
```

Manually set the cost of the BCrypt hashing process
```
$ ./gobcrypt generate -c 17
NKOSv14NuLViUqkf $2a$17$DWJWvyIZH/CamvGAa59Naefhz5AnFNFbfkQ38ryhsH27J7Zovt8Gm
```
