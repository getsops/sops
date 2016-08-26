use std::rc::Rc;
use std::process::{Stdio, Command};
use std::io::{Write, Read};
use std::fs::{OpenOptions, File};
use std::str;
use std::env;
extern crate serde_json;

use serde_json::Value;

extern crate yaml_rust;
use yaml_rust::{YamlLoader, YamlEmitter};

fn main() {}

fn run_sops_and_return_output(command: &mut Command, filename: &str) -> String {
    let mut child = command.stdout(Stdio::piped())
        .arg(filename)
        .spawn()
        .expect("Could not start sops python process");
    let output = child.wait_with_output().expect("Could not retrieve sops's output");
    if !output.status.success() {
        panic!("sops did not exit successfully!");
    }
    return String::from_utf8(output.stdout).expect("Could not decode sops's output as utf-8");
}

fn get_sops_python() -> Command {
    let sops_python_path = env::var("SOPS_PYTHON_PATH")
        .expect("SOPS_PYTHON_PATH environment variable missing");
    let mut cmd = Command::new("python");
    cmd.arg(sops_python_path);
    cmd
}

fn encrypt_with_sops_python(plaintext: &str) -> String {
    let mut child = get_sops_python();
    let child = child.arg("-e");
    return run_sops_and_return_output(child, plaintext);
}

fn decrypt_with_sops_python(ciphertext: &str) -> String {
    let mut child = get_sops_python();
    let child = child.arg("-d");
    return run_sops_and_return_output(child, ciphertext);
}

fn validate_json_file(input_file_name: &str,
                      encrypt: fn(&str) -> String,
                      decrypt: fn(&str) -> String) {
    let output_file_name = "temp.json";
    let mut input = String::new();
    File::open(input_file_name).unwrap().read_to_string(&mut input);
    let input_value: Value = serde_json::from_str(&input).expect("Could not decode input json");
    let encrypted_output = encrypt(input_file_name);
    let mut output_file = OpenOptions::new()
        .write(true)
        .create(true)
        .open(output_file_name)
        .expect("Could not open output file");
    output_file.write_all(encrypted_output.as_bytes()).expect("Could not write to output file");
    let decryption = decrypt(output_file_name);
    let output_value: Value = serde_json::from_str(&decryption).unwrap();
    std::fs::remove_file(output_file_name).expect("Could not remove output file");
    assert_eq!(input_value, output_value);
}

fn validate_yaml_file(input_file_name: &str,
                      encrypt: fn(&str) -> String,
                      decrypt: fn(&str) -> String) {
    let output_file_name = "temp.yaml";
    let mut input = String::new();
    File::open(input_file_name).unwrap().read_to_string(&mut input);
    let input_value = YamlLoader::load_from_str(&input).expect("Could not decode input yaml");
    let encrypted_output = encrypt(input_file_name);
    let mut output_file = OpenOptions::new()
        .write(true)
        .create(true)
        .open(output_file_name)
        .expect("Could not open output file");
    output_file.write_all(encrypted_output.as_bytes()).expect("Could not write to output file");
    let decryption = decrypt(output_file_name);
    let output_value = YamlLoader::load_from_str(&decryption)
        .expect("Could not decode output yaml");
    std::fs::remove_file(output_file_name).expect("Could not remove output file");
    assert_eq!(input_value, output_value);
}

#[test]
fn validate_python_json() {
    validate_json_file("example.json",
                       encrypt_with_sops_python,
                       decrypt_with_sops_python);
}


#[test]
fn validate_python_yaml() {
    validate_yaml_file("example.yaml",
                       encrypt_with_sops_python,
                       decrypt_with_sops_python);
}
