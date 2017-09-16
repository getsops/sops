extern crate tempdir;
extern crate serde;
extern crate serde_json;
extern crate serde_yaml;
#[macro_use]
extern crate lazy_static;
#[macro_use]
extern crate serde_derive;

#[cfg(test)]
mod tests {
    extern crate serde;
    extern crate serde_json;
    extern crate serde_yaml;

    use std::fs::File;
    use std::io::{Write, Read};
    use tempdir::TempDir;
    use std::process::Command;
    use serde_yaml::Value;
    const SOPS_BINARY_PATH: &'static str = "./sops";

    macro_rules! assert_encrypted {
        ($object:expr, $key:expr) => {
            assert!($object.get(&$key).is_some());
            match *$object.get(&$key).unwrap() {
                Value::String(ref s) => {
                   assert!(s.starts_with("ENC["), "Value is not encrypted");
                }
                _ => panic!("Value under key was not a string"),
            }
        }
    }

    lazy_static! {
        static ref TMP_DIR: TempDir = TempDir::new("sops-functional-tests")
            .expect("Unable to create temporary directory");
    }

    fn prepare_temp_file(name: &str, contents: &[u8]) -> String {
        let file_path = TMP_DIR.path().join(name);
        let mut tmp_file = File::create(file_path.clone()).expect("Unable to create temporary file");
        tmp_file.write_all(&contents)
            .expect("Error writing to temporary file");
        file_path.to_string_lossy().into_owned()
    }

    #[test]
    fn encrypt_json_file() {
        let file_path = prepare_temp_file("test_encrypt.json",
                                          b"{
    \"foo\": 2,
    \"bar\": \"baz\"
}");
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("-e")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(output.status.success(), "sops didn't exit successfully");
        let json = &String::from_utf8_lossy(&output.stdout);
        let data: Value = serde_json::from_str(json).expect("Error parsing sops's JSON output");
        match data.into() {
            Value::Mapping(m) => {
                assert!(m.get(&Value::String("sops".to_owned())).is_some(),
                        "sops metadata branch not found");
                assert_encrypted!(&m, Value::String("foo".to_owned()));
                assert_encrypted!(&m, Value::String("bar".to_owned()));
            }
            _ => panic!("sops's JSON output is not an object"),
        }
    }

    #[test]
    fn encrypt_yaml_file() {
        let file_path = prepare_temp_file("test_encrypt.yaml",
                                          b"foo: 2
bar: baz");
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("-e")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(output.status.success(), "sops didn't exit successfully");
        let json = &String::from_utf8_lossy(&output.stdout);
        let data: Value = serde_yaml::from_str(&json).expect("Error parsing sops's JSON output");
        match data.into() {
            Value::Mapping(m) => {
                assert!(m.get(&Value::String("sops".to_owned())).is_some(),
                        "sops metadata branch not found");
                assert_encrypted!(&m, Value::String("foo".to_owned()));
                assert_encrypted!(&m, Value::String("bar".to_owned()));
            }
            _ => panic!("sops's YAML output is not a mapping"),
        }
    }

    #[test]
    fn set_json_file_update() {
        let file_path = prepare_temp_file("test_set_update.json", r#"{"a": 2, "b": "ba"}"#.as_bytes());
        Command::new(SOPS_BINARY_PATH)
            .arg("-e")
            .arg("-i")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("--set")
            .arg(r#"["a"] {"aa": "aaa"}"#)
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        println!("stdout: {}, stderr: {}",
                 String::from_utf8_lossy(&output.stdout),
                 String::from_utf8_lossy(&output.stderr));
        let mut s = String::new();
        File::open(file_path).unwrap().read_to_string(&mut s).unwrap();
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
        let file_path = prepare_temp_file("test_set_insert.json", r#"{"a": 2, "b": "ba"}"#.as_bytes());
        Command::new(SOPS_BINARY_PATH)
            .arg("-e")
            .arg("-i")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("--set")
            .arg(r#"["c"] {"cc": "ccc"}"#)
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        println!("stdout: {}, stderr: {}",
                 String::from_utf8_lossy(&output.stdout),
                 String::from_utf8_lossy(&output.stderr));
        let mut s = String::new();
        File::open(file_path).unwrap().read_to_string(&mut s).unwrap();
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
        let file_path = prepare_temp_file("test_set_update.yaml",
                                          r#"a: 2
b: ba"#
                                              .as_bytes());
        Command::new(SOPS_BINARY_PATH)
            .arg("-e")
            .arg("-i")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("--set")
            .arg(r#"["a"] {"aa": "aaa"}"#)
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        println!("stdout: {}, stderr: {}",
                 String::from_utf8_lossy(&output.stdout),
                 String::from_utf8_lossy(&output.stderr));
        let mut s = String::new();
        File::open(file_path).unwrap().read_to_string(&mut s).unwrap();
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
        let file_path = prepare_temp_file("test_set_insert.yaml",
                                          r#"a: 2
b: ba"#
                                              .as_bytes());
        Command::new(SOPS_BINARY_PATH)
            .arg("-e")
            .arg("-i")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("--set")
            .arg(r#"["c"] {"cc": "ccc"}"#)
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        println!("stdout: {}, stderr: {}",
                 String::from_utf8_lossy(&output.stdout),
                 String::from_utf8_lossy(&output.stderr));
        let mut s = String::new();
        File::open(file_path).unwrap().read_to_string(&mut s).unwrap();
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
    fn decrypt_file_no_mac() {
        let file_path = prepare_temp_file("test_decrypt_no_mac.yaml", include_bytes!("../res/no_mac.yaml"));
        assert!(!Command::new(SOPS_BINARY_PATH)
                    .arg("-d")
                    .arg(file_path.clone())
                    .output()
                    .expect("Error running sops")
                    .status
                    .success(),
                "SOPS allowed decrypting a file with no MAC without --ignore-mac");

        assert!(Command::new(SOPS_BINARY_PATH)
                    .arg("-d")
                    .arg("--ignore-mac")
                    .arg(file_path.clone())
                    .output()
                    .expect("Error running sops")
                    .status
                    .success(),
                "SOPS failed to decrypt a file with no MAC with --ignore-mac passed in");
    }

    #[test]
    fn encrypt_comments() {
        let file_contents = br#"
        lorem: ipsum
        # this-is-a-comment
        dolor: sit
        "#;
        let file_path = prepare_temp_file("test_encrypt_comments.yaml", file_contents);
        let output = Command::new(SOPS_BINARY_PATH)
                    .arg("-e")
                    .arg(file_path.clone())
                    .output()
                    .expect("Error running sops");
        assert!(output.status.success(), "SOPS didn't return successfully");
        assert!(!String::from_utf8_lossy(&output.stdout).contains("this-is-a-comment"), "Comment was not encrypted");
    }

    #[test]
    fn decrypt_comments() {
        let file_contents = br#"
        lorem: ENC[AES256_GCM,data:PhmSdTs=,iv:J5ugEWq6RfyNx+5zDXvcTdoQ18YYZkqesDED7LNzou4=,tag:0Qrom6J6aUnZMZzGz5XCxw==,type:str]
        #ENC[AES256_GCM,data:HiHCasVRzWUiFxKb3X/AcEeM,iv:bmNg+T91dqGk/CEtVH+FDC53osDCEPmWmJKpLyAU5OM=,tag:bTLDYxQSAfYDCBYccoUokQ==,type:comment]
        dolor: ENC[AES256_GCM,data:IgvT,iv:wtPNYbDTARFE810PH6ldOLzCDcAjkB/dzPsZjpgHcko=,tag:zwE8P+AwO1hrHkgF6pTbZw==,type:str]
        sops:
            kms: []
            lastmodified: '2017-08-16T03:41:16Z'
            mac: ENC[AES256_GCM,data:3ngUnY2hkK6pkDbCeAnOHsi/M6bLnGk1vkd+EeGyN/efqJZmwH0+9hUdACNnwHzofIR6NbtCGZal+cSCuTGD4eDuqNV+LbwV1/EaaVZj9RktTNXq3STSXxfzYGoHV3NOMtBhq6sYhF0U72nunreCymm3QzOTylAa2HlmRs54axM=,iv:EMXphsMa+ELK8XXX3MDfFJe3jFgXzwCSwjxNR5ah14k=,tag:gakwLdPvwyihj+FkTG/2kQ==,type:str]
            pgp:
            -   created_at: '2017-08-16T03:41:16Z'
                enc: |-
                    -----BEGIN PGP MESSAGE-----

                    wYwDEEVDpnzXnMABBAAlUcnNciv6rGJua/wmjVYBAHD95VT/M6cc8dg0bPR8XH5a
                    /GeM2RasBzX7ICuBijjesY9exsnrTkBK3/1XpAjygdiW5DciXmqRz/5nE4DLxH+w
                    nZvmnCmg8AdfPKxhr+eM+pKibiN4uEhsJggA9c2ACUQ/YMo4o04fLKZGXqGtT9Lg
                    AeRiZfM3ykiyHDbUQ3P9YAdL4fH44A3gpeHoGeBv4iBFFE7ge+XCby9rgN9Qa7NF
                    /Wahxm7U3RcwT6JSbNDHNCJtolEPeuCR5D2/Kc/2b30e6fLDnpbfSJXiRh4TbOG3
                    rAA=
                    =7P04
                    -----END PGP MESSAGE-----
                fp: 1022470DE3F0BC54BC6AB62DE05550BC07FB1A0A
            unencrypted_suffix: _unencrypted
            version: 2.0.9
        "#;
        let file_path = prepare_temp_file("test_decrypt_comments.yaml", file_contents);
        let output = Command::new(SOPS_BINARY_PATH)
                    .arg("-d")
                    .arg(file_path.clone())
                    .output()
                    .expect("Error running sops");
        assert!(output.status.success(), "SOPS didn't return successfully");
        assert!(String::from_utf8_lossy(&output.stdout).contains("this-is-a-comment"), "Comment was not decrypted");
    }

    #[test]
    fn decrypt_comments_unencrypted_comments() {
        let file_contents = br#"
        lorem: ENC[AES256_GCM,data:PhmSdTs=,iv:J5ugEWq6RfyNx+5zDXvcTdoQ18YYZkqesDED7LNzou4=,tag:0Qrom6J6aUnZMZzGz5XCxw==,type:str]
        # this-is-a-comment
        dolor: ENC[AES256_GCM,data:IgvT,iv:wtPNYbDTARFE810PH6ldOLzCDcAjkB/dzPsZjpgHcko=,tag:zwE8P+AwO1hrHkgF6pTbZw==,type:str]
        sops:
            kms: []
            lastmodified: '2017-08-16T03:41:16Z'
            mac: ENC[AES256_GCM,data:3ngUnY2hkK6pkDbCeAnOHsi/M6bLnGk1vkd+EeGyN/efqJZmwH0+9hUdACNnwHzofIR6NbtCGZal+cSCuTGD4eDuqNV+LbwV1/EaaVZj9RktTNXq3STSXxfzYGoHV3NOMtBhq6sYhF0U72nunreCymm3QzOTylAa2HlmRs54axM=,iv:EMXphsMa+ELK8XXX3MDfFJe3jFgXzwCSwjxNR5ah14k=,tag:gakwLdPvwyihj+FkTG/2kQ==,type:str]
            pgp:
            -   created_at: '2017-08-16T03:41:16Z'
                enc: |-
                    -----BEGIN PGP MESSAGE-----

                    wYwDEEVDpnzXnMABBAAlUcnNciv6rGJua/wmjVYBAHD95VT/M6cc8dg0bPR8XH5a
                    /GeM2RasBzX7ICuBijjesY9exsnrTkBK3/1XpAjygdiW5DciXmqRz/5nE4DLxH+w
                    nZvmnCmg8AdfPKxhr+eM+pKibiN4uEhsJggA9c2ACUQ/YMo4o04fLKZGXqGtT9Lg
                    AeRiZfM3ykiyHDbUQ3P9YAdL4fH44A3gpeHoGeBv4iBFFE7ge+XCby9rgN9Qa7NF
                    /Wahxm7U3RcwT6JSbNDHNCJtolEPeuCR5D2/Kc/2b30e6fLDnpbfSJXiRh4TbOG3
                    rAA=
                    =7P04
                    -----END PGP MESSAGE-----
                fp: 1022470DE3F0BC54BC6AB62DE05550BC07FB1A0A
            unencrypted_suffix: _unencrypted
            version: 2.0.9
        "#;
        let file_path = prepare_temp_file("test_decrypt_comments.yaml", file_contents);
        let output = Command::new(SOPS_BINARY_PATH)
                    .arg("-d")
                    .arg(file_path.clone())
                    .output()
                    .expect("Error running sops");
        assert!(output.status.success(), "SOPS didn't return successfully");
        assert!(String::from_utf8_lossy(&output.stdout).contains("this-is-a-comment"), "Comment was not decrypted");
    }

    fn roundtrip_shamir() {
        // The .sops.yaml file ensures this file is encrypted with two key groups, each with one GPG key
        let file_path = prepare_temp_file("test_roundtrip_keygroups.yaml", "a: secret".as_bytes());
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("-i")
            .arg("-e")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(output.status.success(),
                "SOPS failed to encrypt a file with Shamir Secret Sharing");
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("-d")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(output.status
                    .success(),
                "SOPS failed to decrypt a file with Shamir Secret Sharing");
        assert!(String::from_utf8_lossy(&output.stdout).contains("secret"));
    }

    #[test]
    fn roundtrip_shamir_missing_decryption_key() {
        // The .sops.yaml file ensures this file is encrypted with two key groups, each with one GPG key,
        // but we don't have one of the private keys
        let file_path = prepare_temp_file("test_roundtrip_keygroups_missing_decryption_key.yaml",
                                          "a: secret".as_bytes());
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("-i")
            .arg("-e")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(output.status.success(),
                "SOPS failed to encrypt a file with Shamir Secret Sharing");
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("-d")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(!output.status
                    .success(),
                "SOPS succeeded decrypting a file with a missing decrytion key");
    }

    #[test]
    fn test_decrypt_file_multiple_keys() {
        let file_path = prepare_temp_file("test_decrypt_file_multiple_keys.yaml",
                                          include_bytes!("../res/multiple_keys.yaml"));
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("-d")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(output.status
                    .success(),
                "SOPS failed to decrypt a file that uses multiple keys");
    }
}
