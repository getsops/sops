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
    use std::io::Write;
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
        let mut tmp_file = File::create(file_path.clone())
            .expect("Unable to create temporary file");
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
        let data: Value = serde_json::from_str(json)
            .expect("Error parsing sops's JSON output");
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
}
