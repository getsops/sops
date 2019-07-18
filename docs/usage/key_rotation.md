# Key rotation

It is considered group practice to rotate cryptographic keys every once in a
while. `sops` supports rotation of its data key:

```bash
$ sops --rotate my_file.enc
```

!> Rotating the data key requires encryption access to all master keys used in
a file
