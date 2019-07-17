# Editing files

## Existing files

Suppose you have an encrypted SOPS file, and would like to edit it. As an
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

Let's edit it:

```bash
$ sops my_file.enc
```

This will open the decrypted, plain-text file in your editor. You can then make
any changes you desire, save the file, and exit your editor. The encrypted file
will be updated with the modifications you made.

## New files

`sops` lets you create new encrypted files directly as well. You will be shown
an example file for whichever storage format you are using. You can edit that
example file and save, and your changes will be encrypted and persisted in the
file.

## Advantages over manually decrypting, editing, and then reencrypting a file

Editing files directly through `sops` has a big advantage: the [data
key](internals/encryption_protocol.md) and other cryptographic values are
preserved. For formats that have some sort of key-value structure, this allows
`sops` to leave encrypted values that were not modified completely unchanged,
which in turn allows users to get an idea of *what* was modified in a file
by using version control systems and differs, even if they do not have access
to the plain-text decrypted data at all.

## What editor will `sops` use?

`sops` uses the `EDITOR` environment variable to decide what editor to use. If
the variable is not set, `sops` will then try a few common editors found on
most systems.

## GUI Editors (Sublime, Atom, et al.)

Some editors, especially GUI editors, spawn a new process and then immediately
exit when executed. `sops` relies on your editor exiting to know when you're
done editing the file. Because those editors exit immediately, `sops` will
always think you haven't modified the file at all.

To work around this, consult your editor's documentation for an option that
will make it wait until it is closed for the process to exit. This is sometimes
offered as a command line flag. In that case, you will need to wrap your editor
in a script that passes a command line flag, since the `EDITOR` environment
variable does not accept arguments. For example, for Atom:

```bash
#!/bin/bash
atom -w
```

Make that script executable and then set `EDITOR` to its location.
