<!-- docs/_sidebar.md -->

* [SOPS](/)
* [Installation](installation.md)
* [Quick Start](quick_start.md)
* Usage
	* [Encrypting files](usage/encrypting_files.md)
	* [Decrypting files](usage/decrypting_files.md)
	* [Editing files](usage/editing_files.md)
	* [Git differ](usage/git_differ.md)
	* [Key rotation](usage/key_rotation.md)
	* [Key groups](usage/key_groups.md)
	* [Publishing files](usage/publishing_files.md)
	* [Key service](usage/key_service.md)
	* [Auditing](usage/auditing.md)
	* [Partial file encryption](usage/partial_file_encryption.md)
* [The `.sops.yaml` configuration file](sops_yaml_config_file.md)
* Encryption providers (master key types)
	* [AWS KMS](encryption_providers/aws_kms.md)
	* [GCP KMS](encryption_providers/gcp_kms.md)
	* [Azure KeyVault](encryption_providers/azure_keyvault.md)
	* [PGP](encryption_providers/pgp.md)
* Storage formats
	* [YAML](storage_formats/yaml.md)
	* [JSON](storage_formats/json.md)
	* [.env](storage_formats/dotenv.md)
	* [INI](storage_formats/ini.md)
	* [Arbitrary (binary) files](storage_formats/binary.md)
* Publication targets
	* [AWS S3](publication_targets/s3.md)
	* [Google Cloud Storage](publication_targets/gcs.md)
	* [Hashicorp Vault](publication_targets/vault.md)
* Internals
	* [Encryption protocol](internals/encryption_protocol.md)
* [Comparison with other tools](comparison_with_other_tools.md)

