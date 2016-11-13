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
    use std::collections::BTreeMap;
    const SOPS_BINARY_PATH: &'static str = "./sops";

    // Boring code to use a single value enum instead of one for each json and yaml
    
    #[derive(PartialOrd, Ord, Eq, PartialEq, Debug)]
    enum Value {
        String(String),
        Mapping(BTreeMap<Value, Value>),
        Sequence(Vec<Value>),
        Null,
    }

    impl From<serde_yaml::Value> for Value {
        fn from(val: serde_yaml::Value) -> Value {
            match val {
                serde_yaml::Value::String(s) => Value::String(s),
                serde_yaml::Value::Mapping(m) => {
                    let mut map: BTreeMap<Value, Value> = BTreeMap::new();
                    for (key, value) in m {
                        map.insert(key.into(), value.into());
                    }
                    Value::Mapping(map)
                }
                serde_yaml::Value::Null => Value::Null,
                serde_yaml::Value::Sequence(in_vec) => {
                    let mut vec: Vec<Value> = Vec::new();
                    for v in in_vec {
                        vec.push(v.into());
                    }
                    Value::Sequence(vec)
                }
                _ => unreachable!("{:?}", val),
            }
        }
    }

    impl From<serde_json::Value> for Value {
        fn from(val: serde_json::Value) -> Value {
            match val {
                serde_json::Value::String(s) => Value::String(s),
                serde_json::Value::Object(m) => {
                    let mut map: BTreeMap<Value, Value> = BTreeMap::new();
                    for (key, value) in m {
                        map.insert(key.as_str().into(), value.into());
                    }
                    Value::Mapping(map)
                }
                serde_json::Value::Null => Value::Null,
                serde_json::Value::Array(in_vec) => {
                    let mut vec: Vec<Value> = Vec::new();
                    for v in in_vec {
                        vec.push(v.into());
                    }
                    Value::Sequence(vec)
                }
                _ => unreachable!("{:?}", val),
            }
        }
    }

    impl<'a> From<&'a str> for Value {
        fn from(val: &'a str) -> Value {
            Value::String(val.to_owned())
        }
    }

    macro_rules! assert_encrypted {
        ($object:expr, $key:expr) => {
            let key : Value = $key.into();
            assert!($object.get(&key).is_some());
            match *$object.get(&key).unwrap() {
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
        let data: serde_json::Value = serde_json::from_str(json)
            .expect("Error parsing sops's JSON output");
        match data.into() {
            Value::Mapping(m) => {
                assert!(m.get(&"sops".into()).is_some(),
                        "sops metadata branch not found");
                assert_encrypted!(&m, "foo");
                assert_encrypted!(&m, "bar");
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
        let data: serde_yaml::Value = serde_yaml::from_str(&json)
            .expect("Error parsing sops's JSON output");
        match data.into() {
            Value::Mapping(m) => {
                assert!(m.get(&"sops".into()).is_some(),
                        "sops metadata branch not found");
                assert_encrypted!(&m, "foo");
                assert_encrypted!(&m, "bar");
            }
            _ => panic!("sops's YAML output is not a mapping"),
        }
    }
}
