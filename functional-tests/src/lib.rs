extern crate tempdir;
extern crate serde_json;
extern crate serde_yaml;
#[macro_use]
extern crate lazy_static;

#[cfg(test)]
mod tests {
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
        let file_path = prepare_temp_file("test_decrypt_no_mac.yaml",
                                          r#"
myapp1: ENC[AES256_GCM,data:QsGJGjvQOpoVCIlrYTcOQEfQzriw,iv:ShmgdRNV6UrOJ22Rgr7habB74Nd/YFxU4lDh6jy6n+8=,tag:8GT6U8lzrI27DcFc1+icgQ==,type:str]
sops:
    pgp:
    -   fp: 1022470DE3F0BC54BC6AB62DE05550BC07FB1A0A
        created_at: '2015-11-25T00:32:57Z'
        enc: |
            -----BEGIN PGP MESSAGE-----
            Version: GnuPG v1

            hIwDEEVDpnzXnMABBACBf7lGw8B0sLbfup1Ye51FNpY6iF/4SPTdjeV4OB3uDwIJ
            FRa6z7VR+FrtWyyNYRNB2Wm5eegnEEWwui6hFw7tvlhkN8C5hWQ0B47oYMTstZDR
            TR3Eu7y70u3YLoQKZgDnPb6hQplGIoYVd/EMpDgKmKnmz5oCiIkEI68T3aXo5tJc
            AZhplIlk9eSMHIW9CmGkNp5HtZlQWzVSdGdcQcIUBG4F+Vf40max9u0Jkk1Se1do
            BJ+D4Kl5dZXBj3njvo4YdZ+FGoYPfMlX1GCw0W4caUu6tD8RjuzJA+fYo2Q=
            =Cnu4
            -----END PGP MESSAGE-----
    lastmodified: '2016-03-16T23:34:46Z'
    version: 1.7
    attention: This section contains key material that should only be modified with
        extra care. See `sops -h`.
"#
                                              .as_bytes());
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
}
