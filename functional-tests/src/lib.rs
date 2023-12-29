extern crate serde;
extern crate serde_json;
extern crate serde_yaml;
extern crate tempdir;
#[macro_use]
extern crate lazy_static;
#[macro_use]
extern crate serde_derive;

#[cfg(test)]
mod tests {
    extern crate serde;
    extern crate serde_json;
    extern crate serde_yaml;

    use serde_yaml::Value;
    use std::env;
    use std::fs::File;
    use std::io::{Read, Write};
    use std::path::Path;
    use std::process::Command;
    use tempdir::TempDir;
    const SOPS_BINARY_PATH: &'static str = "./sops";
    const KMS_KEY: &'static str = "FUNCTIONAL_TEST_KMS_ARN";

    macro_rules! assert_encrypted {
        ($object:expr, $key:expr) => {
            assert!($object.get(&$key).is_some());
            match *$object.get(&$key).unwrap() {
                Value::String(ref s) => {
                    assert!(s.starts_with("ENC["), "Value is not encrypted");
                }
                _ => panic!("Value under key was not a string"),
            }
        };
    }

    lazy_static! {
        static ref TMP_DIR: TempDir =
            TempDir::new("sops-functional-tests").expect("Unable to create temporary directory");
    }

    fn prepare_temp_file(name: &str, contents: &[u8]) -> String {
        let file_path = TMP_DIR.path().join(name);
        let mut tmp_file =
            File::create(file_path.clone()).expect("Unable to create temporary file");
        tmp_file
            .write_all(&contents)
            .expect("Error writing to temporary file");
        file_path.to_string_lossy().into_owned()
    }

    #[test]
    fn encrypt_json_file() {
        let file_path = prepare_temp_file(
            "test_encrypt.json",
            b"{
    \"foo\": 2,
    \"bar\": \"baz\"
}",
        );
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("encrypt")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(output.status.success(), "sops didn't exit successfully");
        let json = &String::from_utf8_lossy(&output.stdout);
        let data: Value = serde_json::from_str(json).expect("Error parsing sops's JSON output");
        match data.into() {
            Value::Mapping(m) => {
                assert!(
                    m.get(&Value::String("sops".to_owned())).is_some(),
                    "sops metadata branch not found"
                );
                assert_encrypted!(&m, Value::String("foo".to_owned()));
                assert_encrypted!(&m, Value::String("bar".to_owned()));
            }
            _ => panic!("sops's JSON output is not an object"),
        }
    }

    #[test]
    #[ignore]
    fn publish_json_file_s3() {
        let file_path = prepare_temp_file(
            "test_encrypt_publish_s3.json",
            b"{
    \"foo\": 2,
    \"bar\": \"baz\"
}",
        );
        assert!(
            Command::new(SOPS_BINARY_PATH)
                .arg("encrypt")
                .arg("-i")
                .arg(file_path.clone())
                .output()
                .expect("Error running sops")
                .status
                .success(),
            "SOPS failed to encrypt a file"
        );
        assert!(
            Command::new(SOPS_BINARY_PATH)
                .arg("publish")
                .arg("--yes")
                .arg(file_path.clone())
                .output()
                .expect("Error running sops")
                .status
                .success(),
            "sops failed to publish a file to S3"
        );

        //TODO: Check that file exists in S3 Bucket
    }

    #[test]
    fn publish_json_file_vault() {
        let file_path = prepare_temp_file(
            "test_encrypt_publish_vault.json",
            b"{
    \"foo\": 2,
    \"bar\": \"baz\"
}",
        );
        assert!(
            Command::new(SOPS_BINARY_PATH)
                .arg("encrypt")
                .arg("-i")
                .arg(file_path.clone())
                .output()
                .expect("Error running sops")
                .status
                .success(),
            "SOPS failed to encrypt a file"
        );
        assert!(
            Command::new(SOPS_BINARY_PATH)
                .arg("publish")
                .arg("--yes")
                .arg(file_path.clone())
                .output()
                .expect("Error running sops")
                .status
                .success(),
            "sops failed to publish a file to Vault"
        );

        //TODO: Check that file exists in Vault
    }

    #[test]
    fn publish_json_file_vault_version_1() {
        let file_path = prepare_temp_file(
            "test_encrypt_publish_vault_version_1.json",
            b"{
    \"foo\": 2,
    \"bar\": \"baz\"
}",
        );
        assert!(
            Command::new(SOPS_BINARY_PATH)
                .arg("encrypt")
                .arg("-i")
                .arg(file_path.clone())
                .output()
                .expect("Error running sops")
                .status
                .success(),
            "SOPS failed to encrypt a file"
        );
        assert!(
            Command::new(SOPS_BINARY_PATH)
                .arg("publish")
                .arg("--yes")
                .arg(file_path.clone())
                .output()
                .expect("Error running sops")
                .status
                .success(),
            "sops failed to publish a file to Vault"
        );

        //TODO: Check that file exists in Vault
    }

    #[test]
    #[ignore]
    fn encrypt_json_file_kms() {
        let kms_arn =
            env::var(KMS_KEY).expect("Expected $FUNCTIONAL_TEST_KMS_ARN env var to be set");

        let file_path = prepare_temp_file(
            "test_encrypt_kms.json",
            b"{
    \"foo\": 2,
    \"bar\": \"baz\"
}",
        );

        let output = Command::new(SOPS_BINARY_PATH)
            .arg("encrypt")
            .arg("--kms")
            .arg(kms_arn)
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(output.status.success(), "sops didn't exit successfully");
        let json = &String::from_utf8_lossy(&output.stdout);
        let data: Value = serde_json::from_str(json).expect("Error parsing sops's JSON output");
        match data.into() {
            Value::Mapping(m) => {
                assert!(
                    m.get(&Value::String("sops".to_owned())).is_some(),
                    "sops metadata branch not found"
                );
                assert_encrypted!(&m, Value::String("foo".to_owned()));
                assert_encrypted!(&m, Value::String("bar".to_owned()));
            }
            _ => panic!("sops's JSON output is not an object"),
        }
    }

    #[test]
    fn encrypt_yaml_file() {
        let file_path = prepare_temp_file(
            "test_encrypt.yaml",
            b"foo: 2
bar: baz",
        );
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("encrypt")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(output.status.success(), "sops didn't exit successfully");
        let json = &String::from_utf8_lossy(&output.stdout);
        let data: Value = serde_yaml::from_str(&json).expect("Error parsing sops's JSON output");
        match data.into() {
            Value::Mapping(m) => {
                assert!(
                    m.get(&Value::String("sops".to_owned())).is_some(),
                    "sops metadata branch not found"
                );
                assert_encrypted!(&m, Value::String("foo".to_owned()));
                assert_encrypted!(&m, Value::String("bar".to_owned()));
            }
            _ => panic!("sops's YAML output is not a mapping"),
        }
    }

    #[test]
    fn set_json_file_update() {
        let file_path =
            prepare_temp_file("test_set_update.json", r#"{"a": 2, "b": "ba"}"#.as_bytes());
        assert!(
            Command::new(SOPS_BINARY_PATH)
                .arg("encrypt")
                .arg("-i")
                .arg(file_path.clone())
                .output()
                .expect("Error running sops")
                .status
                .success(),
            "sops didn't exit successfully"
        );
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("set")
            .arg(file_path.clone())
            .arg(r#"["a"]"#)
            .arg(r#"{"aa": "aaa"}"#)
            .output()
            .expect("Error running sops");
        assert!(output.status.success(), "sops didn't exit successfully");
        println!(
            "stdout: {}, stderr: {}",
            String::from_utf8_lossy(&output.stdout),
            String::from_utf8_lossy(&output.stderr)
        );
        let mut s = String::new();
        File::open(file_path)
            .unwrap()
            .read_to_string(&mut s)
            .unwrap();
        let data: Value = serde_json::from_str(&s).expect("Error parsing sops's JSON output");
        if let Value::Mapping(data) = data {
            let a = data.get(&Value::String("a".to_owned())).unwrap();
            if let &Value::Mapping(ref a) = a {
                assert_encrypted!(&a, Value::String("aa".to_owned()));
                return;
            }
        }
        panic!("Output JSON does not have the expected structure");
    }

    #[test]
    fn set_json_file_insert() {
        let file_path =
            prepare_temp_file("test_set_insert.json", r#"{"a": 2, "b": "ba"}"#.as_bytes());
        assert!(
            Command::new(SOPS_BINARY_PATH)
                .arg("encrypt")
                .arg("-i")
                .arg(file_path.clone())
                .output()
                .expect("Error running sops")
                .status
                .success(),
            "sops didn't exit successfully"
        );
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("set")
            .arg(file_path.clone())
            .arg(r#"["c"]"#)
            .arg(r#"{"cc": "ccc"}"#)
            .output()
            .expect("Error running sops");
        assert!(output.status.success(), "sops didn't exit successfully");
        println!(
            "stdout: {}, stderr: {}",
            String::from_utf8_lossy(&output.stdout),
            String::from_utf8_lossy(&output.stderr)
        );
        let mut s = String::new();
        File::open(file_path)
            .unwrap()
            .read_to_string(&mut s)
            .unwrap();
        let data: Value = serde_json::from_str(&s).expect("Error parsing sops's JSON output");
        if let Value::Mapping(data) = data {
            let a = data.get(&Value::String("c".to_owned())).unwrap();
            if let &Value::Mapping(ref a) = a {
                assert_encrypted!(&a, Value::String("cc".to_owned()));
                return;
            }
        }
        panic!("Output JSON does not have the expected structure");
    }

    #[test]
    fn set_yaml_file_update() {
        let file_path = prepare_temp_file(
            "test_set_update.yaml",
            r#"a: 2
b: ba"#
                .as_bytes(),
        );
        assert!(
            Command::new(SOPS_BINARY_PATH)
                .arg("encrypt")
                .arg("-i")
                .arg(file_path.clone())
                .output()
                .expect("Error running sops")
                .status
                .success(),
            "sops didn't exit successfully"
        );
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("set")
            .arg(file_path.clone())
            .arg(r#"["a"]"#)
            .arg(r#"{"aa": "aaa"}"#)
            .output()
            .expect("Error running sops");
        assert!(output.status.success(), "sops didn't exit successfully");
        println!(
            "stdout: {}, stderr: {}",
            String::from_utf8_lossy(&output.stdout),
            String::from_utf8_lossy(&output.stderr)
        );
        let mut s = String::new();
        File::open(file_path)
            .unwrap()
            .read_to_string(&mut s)
            .unwrap();
        let data: Value = serde_yaml::from_str(&s).expect("Error parsing sops's JSON output");
        if let Value::Mapping(data) = data {
            let a = data.get(&Value::String("a".to_owned())).unwrap();
            if let &Value::Mapping(ref a) = a {
                assert_encrypted!(&a, Value::String("aa".to_owned()));
                return;
            }
        }
        panic!("Output JSON does not have the expected structure");
    }

    #[test]
    fn set_yaml_file_insert() {
        let file_path = prepare_temp_file(
            "test_set_insert.yaml",
            r#"a: 2
b: ba"#
                .as_bytes(),
        );
        assert!(
            Command::new(SOPS_BINARY_PATH)
                .arg("encrypt")
                .arg("-i")
                .arg(file_path.clone())
                .output()
                .expect("Error running sops")
                .status
                .success(),
            "sops didn't exit successfully"
        );
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("set")
            .arg(file_path.clone())
            .arg(r#"["c"]"#)
            .arg(r#"{"cc": "ccc"}"#)
            .output()
            .expect("Error running sops");
        assert!(output.status.success(), "sops didn't exit successfully");
        println!(
            "stdout: {}, stderr: {}",
            String::from_utf8_lossy(&output.stdout),
            String::from_utf8_lossy(&output.stderr)
        );
        let mut s = String::new();
        File::open(file_path)
            .unwrap()
            .read_to_string(&mut s)
            .unwrap();
        let data: Value = serde_yaml::from_str(&s).expect("Error parsing sops's JSON output");
        if let Value::Mapping(data) = data {
            let a = data.get(&Value::String("c".to_owned())).unwrap();
            if let &Value::Mapping(ref a) = a {
                assert_encrypted!(&a, Value::String("cc".to_owned()));
                return;
            }
        }
        panic!("Output YAML does not have the expected structure");
    }

    #[test]
    fn set_yaml_file_string() {
        let file_path = prepare_temp_file(
            "test_set_string.yaml",
            r#"a: 2
b: ba"#
                .as_bytes(),
        );
        assert!(
            Command::new(SOPS_BINARY_PATH)
                .arg("encrypt")
                .arg("-i")
                .arg(file_path.clone())
                .output()
                .expect("Error running sops")
                .status
                .success(),
            "sops didn't exit successfully"
        );
        assert!(
            Command::new(SOPS_BINARY_PATH)
                .arg("set")
                .arg(file_path.clone())
                .arg(r#"["a"]"#)
                .arg(r#""aaa""#)
                .output()
                .expect("Error running sops")
                .status
                .success(),
            "sops didn't exit successfully"
        );
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("decrypt")
            .arg("-i")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(output.status.success(), "sops didn't exit successfully");
        println!(
            "stdout: {}, stderr: {}",
            String::from_utf8_lossy(&output.stdout),
            String::from_utf8_lossy(&output.stderr)
        );
        let mut s = String::new();
        File::open(file_path)
            .unwrap()
            .read_to_string(&mut s)
            .unwrap();
        let data: Value = serde_yaml::from_str(&s).expect("Error parsing sops's YAML output");
        if let Value::Mapping(data) = data {
            let a = data.get(&Value::String("a".to_owned())).unwrap();
            assert_eq!(a, &Value::String("aaa".to_owned()));
        } else {
            panic!("Output JSON does not have the expected structure");
        }
    }

    #[test]
    fn set_yaml_file_string_compat() {
        let file_path = prepare_temp_file(
            "test_set_string_compat.yaml",
            r#"a: 2
b: ba"#
                .as_bytes(),
        );
        assert!(
            Command::new(SOPS_BINARY_PATH)
                .arg("-e")
                .arg("-i")
                .arg(file_path.clone())
                .output()
                .expect("Error running sops")
                .status
                .success(),
            "sops didn't exit successfully"
        );
        assert!(
            Command::new(SOPS_BINARY_PATH)
                .arg("--set")
                .arg(r#"["a"] "aaa""#)
                .arg("-i")
                .arg(file_path.clone())
                .output()
                .expect("Error running sops")
                .status
                .success(),
            "sops didn't exit successfully"
        );
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("-d")
            .arg("-i")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(output.status.success(), "sops didn't exit successfully");
        println!(
            "stdout: {}, stderr: {}",
            String::from_utf8_lossy(&output.stdout),
            String::from_utf8_lossy(&output.stderr)
        );
        let mut s = String::new();
        File::open(file_path)
            .unwrap()
            .read_to_string(&mut s)
            .unwrap();
        let data: Value = serde_yaml::from_str(&s).expect("Error parsing sops's YAML output");
        if let Value::Mapping(data) = data {
            let a = data.get(&Value::String("a".to_owned())).unwrap();
            assert_eq!(a, &Value::String("aaa".to_owned()));
        } else {
            panic!("Output JSON does not have the expected structure");
        }
    }

    #[test]
    fn decrypt_file_no_mac() {
        let file_path = prepare_temp_file(
            "test_decrypt_no_mac.yaml",
            include_bytes!("../res/no_mac.yaml"),
        );
        assert!(
            !Command::new(SOPS_BINARY_PATH)
                .arg("decrypt")
                .arg(file_path.clone())
                .output()
                .expect("Error running sops")
                .status
                .success(),
            "SOPS allowed decrypting a file with no MAC without --ignore-mac"
        );

        assert!(
            Command::new(SOPS_BINARY_PATH)
                .arg("decrypt")
                .arg("--ignore-mac")
                .arg(file_path.clone())
                .output()
                .expect("Error running sops")
                .status
                .success(),
            "SOPS failed to decrypt a file with no MAC with --ignore-mac passed in"
        );
    }

    #[test]
    fn encrypt_comments() {
        let file_path = "res/comments.yaml";
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("encrypt")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(output.status.success(), "SOPS didn't return successfully");
        assert!(
            !String::from_utf8_lossy(&output.stdout).contains("first comment in file"),
            "Comment was not encrypted"
        );
        assert!(
            !String::from_utf8_lossy(&output.stdout).contains("this-is-a-comment"),
            "Comment was not encrypted"
        );
    }

    #[test]
    fn encrypt_comments_list() {
        let file_path = "res/comments_list.yaml";
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("encrypt")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(output.status.success(), "SOPS didn't return successfully");
        assert!(
            !String::from_utf8_lossy(&output.stdout).contains("this-is-a-comment"),
            "Comment was not encrypted"
        );
        assert!(
            !String::from_utf8_lossy(&output.stdout).contains("this-is-a-comment"),
            "Comment was not encrypted"
        );
    }

    #[test]
    fn decrypt_comments() {
        let file_path = "res/comments.enc.yaml";
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("decrypt")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(output.status.success(), "SOPS didn't return successfully");
        assert!(
            String::from_utf8_lossy(&output.stdout).contains("first comment in file"),
            "Comment was not decrypted"
        );
        assert!(
            String::from_utf8_lossy(&output.stdout).contains("this-is-a-comment"),
            "Comment was not decrypted"
        );
    }

    #[test]
    fn decrypt_comments_unencrypted_comments() {
        let file_path = "res/comments_unencrypted_comments.yaml";
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("decrypt")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(output.status.success(), "SOPS didn't return successfully");
        assert!(
            String::from_utf8_lossy(&output.stdout).contains("first comment in file"),
            "Comment was not decrypted"
        );
        assert!(
            String::from_utf8_lossy(&output.stdout).contains("this-is-a-comment"),
            "Comment was not decrypted"
        );
    }

    #[test]
    fn roundtrip_shamir() {
        // The .sops.yaml file ensures this file is encrypted with two key groups, each with one GPG key
        let file_path = prepare_temp_file("test_roundtrip_keygroups.yaml", "a: secret".as_bytes());
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("encrypt")
            .arg("-i")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(
            output.status.success(),
            "SOPS failed to encrypt a file with Shamir Secret Sharing"
        );
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("decrypt")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(
            output.status.success(),
            "SOPS failed to decrypt a file with Shamir Secret Sharing"
        );
        assert!(String::from_utf8_lossy(&output.stdout).contains("secret"));
    }

    #[test]
    fn roundtrip_shamir_missing_decryption_key() {
        // The .sops.yaml file ensures this file is encrypted with two key groups, each with one GPG key,
        // but we don't have one of the private keys
        let file_path = prepare_temp_file(
            "test_roundtrip_keygroups_missing_decryption_key.yaml",
            "a: secret".as_bytes(),
        );
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("encrypt")
            .arg("-i")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(
            output.status.success(),
            "SOPS failed to encrypt a file with Shamir Secret Sharing"
        );
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("decrypt")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(
            !output.status.success(),
            "SOPS succeeded decrypting a file with a missing decryption key"
        );
    }

    #[test]
    fn test_decrypt_file_multiple_keys() {
        let file_path = prepare_temp_file(
            "test_decrypt_file_multiple_keys.yaml",
            include_bytes!("../res/multiple_keys.yaml"),
        );
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("decrypt")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(
            output.status.success(),
            "SOPS failed to decrypt a file that uses multiple keys"
        );
    }

    #[test]
    fn extract_string() {
        let file_path = prepare_temp_file(
            "test_extract_string.yaml",
            "multiline: |\n  multi\n  line".as_bytes(),
        );
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("encrypt")
            .arg("-i")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(output.status.success(), "SOPS failed to encrypt a file");
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("decrypt")
            .arg("--extract")
            .arg("[\"multiline\"]")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(output.status.success(), "SOPS failed to extract");

        assert_eq!(output.stdout, b"multi\nline");
    }

    #[test]
    fn roundtrip_binary() {
        let data = b"\"\"{}this_is_binary_data";
        let file_path = prepare_temp_file("test.binary", data);
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("encrypt")
            .arg("-i")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(
            output.status.success(),
            "SOPS failed to encrypt a binary file"
        );
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("decrypt")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(
            output.status.success(),
            "SOPS failed to decrypt a binary file"
        );
        assert_eq!(output.stdout, data);
    }

    #[test]
    #[ignore]
    fn roundtrip_kms_encryption_context() {
        let kms_arn =
            env::var(KMS_KEY).expect("Expected $FUNCTIONAL_TEST_KMS_ARN env var to be set");

        let file_path = prepare_temp_file(
            "test_roundtrip_kms_encryption_context.json",
            b"{
    \"foo\": 2,
    \"bar\": \"baz\"
}",
        );

        let output = Command::new(SOPS_BINARY_PATH)
            .arg("encrypt")
            .arg("--kms")
            .arg(kms_arn)
            .arg("--encryption-context")
            .arg("foo:bar,one:two")
            .arg("-i")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(output.status.success(), "sops didn't exit successfully");

        let output = Command::new(SOPS_BINARY_PATH)
            .arg("decrypt")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(
            output.status.success(),
            "SOPS failed to decrypt a file with KMS Encryption Context"
        );
        assert!(String::from_utf8_lossy(&output.stdout).contains("foo"));
        assert!(String::from_utf8_lossy(&output.stdout).contains("baz"));
    }

    #[test]
    fn output_flag() {
        let input_path = prepare_temp_file("test_output_flag.binary", b"foo");
        let output_path = Path::join(TMP_DIR.path(), "output_flag.txt");
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("encrypt")
            .arg("--output")
            .arg(&output_path)
            .arg(input_path.clone())
            .output()
            .expect("Error running sops");
        assert!(
            output.status.success(),
            "SOPS failed to decrypt a binary file"
        );
        assert_eq!(output.stdout, <&[u8]>::default());
        let mut f = File::open(&output_path).expect("output file not found");

        let mut contents = String::new();
        f.read_to_string(&mut contents)
            .expect("couldn't read output file contents");
        assert_ne!(contents, "", "Output file is empty");
    }

    #[test]
    fn exec_env() {
        let file_path = prepare_temp_file(
            "test_exec_env.yaml",
            r#"foo: bar
bar: |-
  baz
  bam
"#
            .as_bytes(),
        );
        assert!(
            Command::new(SOPS_BINARY_PATH)
                .arg("encrypt")
                .arg("-i")
                .arg(file_path.clone())
                .output()
                .expect("Error running sops")
                .status
                .success(),
            "sops didn't exit successfully"
        );
        let print_foo = prepare_temp_file(
            "print_foo.sh",
            r#"#!/bin/bash
echo -E "${foo}"
"#
            .as_bytes(),
        );
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("exec-env")
            .arg(file_path.clone())
            .arg(format!("/bin/bash {}", print_foo))
            .output()
            .expect("Error running sops");
        assert!(output.status.success(), "sops didn't exit successfully");
        println!(
            "stdout: {}, stderr: {}",
            String::from_utf8_lossy(&output.stdout),
            String::from_utf8_lossy(&output.stderr)
        );
        assert_eq!(String::from_utf8_lossy(&output.stdout), "bar\n");
        let print_bar = prepare_temp_file(
            "print_bar.sh",
            r#"#!/bin/bash
echo -E "${bar}"
"#
            .as_bytes(),
        );
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("exec-env")
            .arg(file_path.clone())
            .arg(format!("/bin/bash {}", print_bar))
            .output()
            .expect("Error running sops");
        assert!(output.status.success(), "sops didn't exit successfully");
        println!(
            "stdout: {}, stderr: {}",
            String::from_utf8_lossy(&output.stdout),
            String::from_utf8_lossy(&output.stderr)
        );
        assert_eq!(String::from_utf8_lossy(&output.stdout), "baz\\nbam\n");
    }

    #[test]
    fn exec_file() {
        let file_path = prepare_temp_file(
            "test_exec_file.yaml",
            r#"foo: bar
bar: |-
  baz
  bam
"#
            .as_bytes(),
        );
        assert!(
            Command::new(SOPS_BINARY_PATH)
                .arg("-e")
                .arg("-i")
                .arg(file_path.clone())
                .output()
                .expect("Error running sops")
                .status
                .success(),
            "sops didn't exit successfully"
        );
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("exec-file")
            .arg("--output-type")
            .arg("json")
            .arg(file_path.clone())
            .arg("cat {}")
            .output()
            .expect("Error running sops");
        assert!(output.status.success(), "sops didn't exit successfully");
        println!(
            "stdout: {}, stderr: {}",
            String::from_utf8_lossy(&output.stdout),
            String::from_utf8_lossy(&output.stderr)
        );
        assert_eq!(
            String::from_utf8_lossy(&output.stdout),
            r#"{
	"foo": "bar",
	"bar": "baz\nbam"
}"#
        );
    }
}
