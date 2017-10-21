#### GobCrypt

##### A Bcrypt hash password generator written in golang

```
Usage of ./gobcrypt:
  -c int
    	Specify the number of hashes to create (default 1)
  -s string
    	Hash the specified password only (1)
  -t	Validate the hash & pass
```
Example usage:

Create a password and hash:
```
$ gobcrypt 
JZA3IdJald2quxn 	 $2a$14$KjXJ/oMs3e1cwXp3BgPsHOHPLu5Ri061QJKnkTU2BapiADA1YIDXW
```
Create a hash of a specific password
```
$ gobcrypt -s supergoodpass
supergoodpass 	 $2a$14$oTmwiTh6VGvoOgMPbr5K3e7NSze/b/oisoUaZ9PdhUuqugwVSKI/K
```

Check the hash is valid:
```
$ gobcrypt -s supergoodpass -t
supergoodpass 	 $2a$14$xyy2eH4baTpn11s5fi8wGuQrcHVohw3jTNjFh.Oy9sLh2kjT/cire
Match:    true
```
Create a bulk amount of password hash combinations
```
$ gobcrypt -c 7
QBMqLooJXcbD3sW 	 $2a$14$O8Ayph.ZdMmG57qt7FwYPuEPM.zbMa.p8VANQqNJYD7QrlSUIcBz6
fboJQQLaMx1vkgP 	 $2a$14$PBzRngAjOglo9VrfalFZQ.MJ.A99eRub9m6wgyQvECW/mtgha5SYe
b1q166OO6rV60D9 	 $2a$14$qxQO1DpFF/5PbDZW6g4FWu2Q.KujrTc9z7hqC3hoccLUzSO4ZjFVO
eyrPsYSTn2LYuX1 	 $2a$14$pPeZKP3PzlepCFOxiO50Ce45Z3nBVFuWoL/RpAp7d1LdKHIoycVnW
MsUq5tySGQWhMlr 	 $2a$14$5zXPe7SVD0Hi2tAARkOQrOb4H.rbG9CYjVoKx9/fxX7w5.LZ2loY2
AaSs34ck3aWQThd 	 $2a$14$7yAraJ5JoR2G6aeoOBj/hOJKgEZ59cmi50G2tziO/Qo3X7g79B4cS
HYM75gTV3vsEFTA 	 $2a$14$GnC0/3UPWO70Ukb2kIhZhOtAGkN.1N7xkdGvwEoAV5zz.2kxFpcAK
```
