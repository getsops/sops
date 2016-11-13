extern crate tempdir;
extern crate serde_json;

#[cfg(test)]
mod tests {
    extern crate serde_json;

    use std::fs::File;
    use std::io::Write;
    use tempdir::TempDir;
    use std::process::Command;
    use serde_json::Value;
    const SOPS_BINARY_PATH: &'static str = "./sops";

    macro_rules! assert_encrypted {
        ($object:expr, $key:expr) => {
            assert!($object.get($key).is_some());
            match *$object.get($key).unwrap() {
                serde_json::Value::String(ref s) => {
                    assert!(s.starts_with("ENC["), "Value is not encrypted");
                }
                _ => panic!("Value under key {} is not a string", $key)
            }
        }
    }

    #[test]
    fn encrypt_file() {
        let tmp_dir = TempDir::new("sops-functional-tests")
            .expect("Unable to create temporary directory");
        let file_path = tmp_dir.path().join("test_encrypt.json");
        let mut tmp_file = File::create(file_path.clone())
            .expect("Unable to create temporary file");
        tmp_file.write_all(b"{
    \"foo\": 2,
    \"bar\": \"baz\"
}")
            .expect("Error writing to temporary file");
        let output = Command::new(SOPS_BINARY_PATH)
            .arg("-e")
            .arg(file_path.clone())
            .output()
            .expect("Error running sops");
        assert!(output.status.success(), "sops didn't exit successfully");
        let json = &String::from_utf8_lossy(&output.stdout);
        let data: Value = serde_json::from_str(json).expect("Error parsing sops's JSON output");
        match data {
            serde_json::Value::Object(o) => {
                assert!(o.get("sops").is_some(), "sops metadata branch not found");
                assert_encrypted!(o, "foo");
                assert_encrypted!(o, "bar");
            }
            _ => panic!("sops's JSON output is not an object"),
        }
    }
}
