# Decrypting files

Suppose you have an encrypted SOPS file, and would like to decrypt it. As an
example, let's use the file from [Encrypting files](usage/encrypting_files.md):

```bash
$ cat my_file.enc
{
	"data": "ENC[AES256_GCM,data:2jQNm13QiRfVWf/nEk8=,iv:s7vZVoMFcQ32TvhQixVc6ZxC5EevnTBG3UxlQ8lD1iM=,tag:sa6j6Im4tQDCqt2oHCWXjg==,type:str]",
	"sops": {
		"kms": null,
		"gcp_kms": null,
		"azure_kv": null,
		"lastmodified": "2019-07-17T19:12:15Z",
		"mac": "ENC[AES256_GCM,data:5uvzX7kRrFCL5a58js1ls4ALltUtqqi5ZhoQGMbXDZdTDxYaCLkoFTg4PCOil7133g1uaye6f9AjPgvEUGpPFONWg6g9k6k7fN/AsgSiYSZUHD1yCQZhyiKIVhxrnFV+0v0fH5Zwm7bZAqrUrzaH3YGpo6ces9iBSsHMCEHDGdc=,iv:7LX+irD4XrmuHsCYKJvd4BsP5TMkWeLSc6I+heg+c0s=,tag:33WZIPkU5LOMny0L7cx6fg==,type:str]",
		"pgp": [
			{
				"created_at": "2019-07-17T19:12:15Z",
				"enc": "-----BEGIN PGP MESSAGE-----\n\nwYwDEEVDpnzXnMABBACsHmqe5BT4S4O684E39czJrmGkRTSMYX9YCnSVNUVMkwNY\n+JL6FmuLC13320weeO3xL9CCDJIKAGixTehi5JDVY1rK9bCuUTQrjN8NMEHPAZn8\nRvB/W0hKkqaOOpAjq2syp2RjTnNOn8cqkP80Jo9w3BXIJJitBJuKC850Vh0FJdLg\nAeTzXWkhERBsxwJ4ADawDWCU4XPi4M7gDOHmiOAa4shf6zjgwOU5PFnDtDG0kLCl\n49NAKBboOSyEx9sFA4on3j7GLXnroOCj5OwgcsBqIjxerkXizChBS1ziKXEME+Gc\neAA=\n=snqD\n-----END PGP MESSAGE-----",
				"fp": "1022470DE3F0BC54BC6AB62DE05550BC07FB1A0A"
			}
		],
		"unencrypted_suffix": "_unencrypted",
		"version": "3.3.1"
	}
}
```

Let's decrypt it:

```bash
$ sops --decrypt my_file.enc
hello, world!
```

As expected, the original file contents are returned and printed to standard output.
