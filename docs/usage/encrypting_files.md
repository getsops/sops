# Encrypting files

Suppose you already have a plain-text, unencrypted file, and would like `sops` to encrypt it:

```bash
$ cat my_file
hello, world!
```

You can do so as follows:

```bash
$ sops --encrypt my_file
{
	"data": "ENC[AES256_GCM,data:45IFIHXi0pogtNwIE90=,iv:Lu7aGNAHCvwi3VMyJewjhj+pxYCi05gdTpYC1DKFmg8=,tag:YlRb8za8o1lJ4XiTj/KxfA==,type:str]",
	"sops": {
		"kms": null,
		"gcp_kms": null,
		"azure_kv": null,
		"lastmodified": "2019-07-17T18:59:59Z",
		"mac": "ENC[AES256_GCM,data:bYHCF4z/wxIwUBDEEvyCwqUu2v+LWxRikYnI9fC2OVnwcIDadH/xhK1zIrRxq0HugeQYSar5GsCjgy4jhiaxfG1RlQZnkL2Iz/wT+ZdxwY4yeLJ+4/tRwjxCqhrmq2uTihskP7RFVLp6TqM9r/JAyCzqDHxScPmAMVydDTieYaU=,iv:jD4yvBNS075aCbv+OvRdFXjOD13etsux01XL2d5F2cw=,tag:vxm/paFnCtEmIgqkv4h+OQ==,type:str]",
		"pgp": [
			{
				"created_at": "2019-07-17T18:59:58Z",
				"enc": "-----BEGIN PGP MESSAGE-----\n\nwYwDEEVDpnzXnMABBAClqG2V54awfbG/6CuIlnawf5FnFizG4lCV9iXHHeth/r1W\n+1ejBKY9EcrnW/o59adGFmknDsIDAmJaQE3V9fruOdtMixr9MjkNqUfT90c/mCZ9\nV+yA4UfDoZuNzRk3Bkb/7PuT02X9eASdn4V3ShbGmbdMk37NQPel9QXcuAVlqtLg\nAeRzhNVGICGVGpLuzCnvP4DN4XZR4IfghuEe2eD14lmfIRLg7+VDb4IMf81bpomg\ndqXmTpbTsV/5gQxjdijsvh5QRDO2duDe5LTwyRWkZgxtkMee0XZKCVLi5gmRNeFO\n8gA=\n=Udaa\n-----END PGP MESSAGE-----",
				"fp": "1022470DE3F0BC54BC6AB62DE05550BC07FB1A0A"
			}
		],
		"unencrypted_suffix": "_unencrypted",
		"version": "3.3.1"
	}
}
```

!> You need at least one [master key](encryption_providers/README.md) available
to encrypt a file.

The encrypted file is written to standard output. If you want to save it
permanentely, you can either use your shell to redirect the output to a file,
or you can operate on the file [in-place](usage/operate_in_place.md).
