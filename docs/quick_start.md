Quick Start
-----------

Now that you've have the `sops` binary [installed](installation.md), you should
be able to run the binary:

```bash
$ sops -v
sops 3.3.1 (latest)
```

For simplicity, we will use PGP as an encryption provider. We assume you
already have a PGP key, but if you don't, [Github's
documentation](https://help.github.com/en/articles/generating-a-new-gpg-key)
explains how you can create one. Take note of your public key's fingerprint, as
you'll need it later.

?> Tip: you can find your key's fingerprint by running `gpg --list-keys`. The
fingerprint is a long string (typically 40 characters) of hexadecimal
characters.

In my case, the PGP fingerprint is `1022470DE3F0BC54BC6AB62DE05550BC07FB1A0A`,
so I will be using that for the rest of the quick start. Replace it with your
own fingerprint whenever I use it.  This is used by the SOPS project for
testing infrastructure. If you struggle generating your own key, you can use
this fingerprint as well, as `sops` will try to download the key from the
Internet if it is not available locally.

!> If you use the SOPS testing infrastructure key, you won't be able to decrypt
files unless you also download and import the private key. It is available
[here](../pgp/sops_functional_tests_key.asc). Download it and then import it
with `gpg --import sops_functional_tests_key.asc`. 

Now that we have our PGP key set up, let's create a new SOPS encrypted file:

```bash
$ sops --pgp 1022470DE3F0BC54BC6AB62DE05550BC07FB1A0A quickstart.yaml
```

Upon running this, you should be taken to an editor and be shown an example
file:

```yaml
hello: Welcome to SOPS! Edit this file as you please!
example_key: example_value
# Example comment
example_array:
- example_value1
- example_value2
example_number: 1234.5679
example_booleans:
- true
- false
```

Make some edits to the file, for example, add some data to the end of the
file, so it looks like this:

```yaml
hello: Welcome to SOPS! Edit this file as you please!
example_key: example_value
# Example comment
example_array:
- example_value1
- example_value2
example_number: 1234.5679
example_booleans:
- true
- false
some_more_data:
  hello: world
```

Save the file. Let's see what it looks like:

```bash
$ cat quickstart.yaml
```
```yaml
hello: ENC[AES256_GCM,data:r8pcvLb9EZyZFdy/Q1S0lI6JNxsZuf1QM9ZUw7n/A8QkBf1e9B7yOcI1BbzjcQ==,iv:mI5g/blnelhbr75h/yhxyrExPDPaEProWyh9Q0OItv4=,tag:w8ixFTZkRtc+OVSyZOPlrw==,type:str]
example_key: ENC[AES256_GCM,data:87zFG19/7ooxavTFHA==,iv:E3po1Pf4vvqfOaZtfNMPA6BnUCxs1FHu95AVVv3zbrc=,tag:gM7hn6jlXn9QKYPalUkL1g==,type:str]
#ENC[AES256_GCM,data:eJZ9g+Ikh1ijDuwADlbksg==,iv:k8yO0gFoFkpWX/dTZ5IsNvh0ld/cxSkiE0kgyVegZs8=,tag:7qB7CANMs4ZndzfAqrmNiA==,type:comment]
example_array:
- ENC[AES256_GCM,data:rDGwfQzcPB7MWWj+ks8=,iv:LEKN7Yq2KXpukQJcEFG9VJAJ4vSs+2r3QlljHYbsIuA=,tag:JG4WyZOqUzzPkAlPAo/bnQ==,type:str]
- ENC[AES256_GCM,data:TdLQadOzpSTZ4dMqSG0=,iv:AUE5L4VMy+N6tw/aeyZhQcckFxjQgTI3UyI6p6o2HOE=,tag:AAh+WQpO/jSvGvbEH5/QPw==,type:str]
example_number: ENC[AES256_GCM,data:vZkNkTjce/tC,iv:8FlC3NF+AJpHJ6JGtl1oMgw+8btxsKlZqCbXlhQ60Ss=,tag:2AaHYIC7OSYU27JxQTlhXQ==,type:float]
example_booleans:
- ENC[AES256_GCM,data:4oWlWw==,iv:YxEiZ+uweb8rLZBY5nvtt/8LpCl9q+TTCJEyyu6gzhE=,tag:pcZl3mh0y+Uyo/bRePXohg==,type:bool]
- ENC[AES256_GCM,data:FBeNQA4=,iv:dQOSHy8Tgs8f2lsRmUiIZdACfyL+R1OqA+DS7vxuyeI=,tag:Dy3+kQl+vBY8GLObY/qqDw==,type:bool]
some_more_data:
    hello: ENC[AES256_GCM,data:UKybdeY=,iv:SyP3DdAFPOOXZkAkHtKiauH5t0T/lUZbGdo+SI9g80k=,tag:4mgLyrtWWdbNbXolho2v2A==,type:str]
sops:
    kms: []
    gcp_kms: []
    azure_kv: []
    lastmodified: '2019-07-17T15:52:25Z'
    mac: ENC[AES256_GCM,data:CURHSaJhLt4l9yQQqUtkrXZNDCj00SMYM6eal4U/u2XbitH1mR0nsLsRlePo2KPu7NbrJ7P46nqIwuvBbWrF2MGKa3ucTdY1pyJJCWD0+OK1+lV3P8T47/ZxN7r+aCDHVue8StRIcufvtLbtZ/01pDeg3gviSH1JwtP3VaCeihc=,iv:OQ0imfPs9JtAB2Np3V+WXbA+LwhUdQzZaIz6ZwJ3glU=,tag:HAHjEXsaPDvw2zW+d9WFCg==,type:str]
    pgp:
    -   created_at: '2019-07-17T15:51:36Z'
        enc: |-
            -----BEGIN PGP MESSAGE-----

            wYwDEEVDpnzXnMABBAB/NzaqcAN5K1JEC3RTCZcQZ/9GfZ8gnxIeVWArH1S7qN1m
            6y+m7/oYM/xEH76abBdE85lLdoFHsuKEzMcfhdYiyjjwy2zJHkzzb3NMwFqgBeeu
            ok9L3xn4LswDNxRgM+OhQqvkg2+i5QgXUcDNxEBLsNIEhyR1eZnIUafeRX3ht9Lg
            AeTRLKmYvzSNcsTLFEXGVs2w4bi14KjgzOGiNuCJ4o6QJ9DgUuW8SVR1ytG0Mj+e
            FoUTmga/cV996OasIPD/iUqOgJCFD+Cj5AtclYK8JbSoPIJYpc2CGAbi6bZ6LOGQ
            sgA=
            =O39v
            -----END PGP MESSAGE-----
        fp: 1022470DE3F0BC54BC6AB62DE05550BC07FB1A0A
    unencrypted_suffix: _unencrypted
    version: 3.3.1
```

None of the values in our plain text YAML file are visible, they've all been
encrypted! We can make some further edits to the file by running
`sops quickstart.yaml`.  We don't need to specify our PGP key fingerprint
again, as the SOPS file format stores this information already.

For more ways to use `sops`, please consult the rest of the documentation.
