use std::{collections::HashMap, error, fs, result};

/// Represents the common fields in a MANIFEST.MF file
pub enum CommonFields {
    ManifestVersion,
    PreMainClass,
    MainClass,
    LombokVersion,
    AgentClass,
    CanRedefineClasses
}

impl CommonFields {
    fn as_str(&self) -> &'static str {
        match self {
            CommonFields::ManifestVersion => "Manifest-Version",
            CommonFields::PreMainClass => "Premain-Class",
            CommonFields::MainClass => "Main-Class",
            CommonFields::LombokVersion => "Lombok-Version",
            CommonFields::AgentClass => "Agent-Class",
            CommonFields::CanRedefineClasses => "Can-Redefine-Classes"
        }
    }
}

/// Parse a MANIFEST.MF file
/// 
/// # Arguments
/// * `path`: A reference to a `String`, pointing to the `MANIFEST.MF` file on the filesystem
/// 
/// # Returns
/// A `HashMap<String, String>` representing the key/value pairs in the manifest.
pub fn parse(path: &String) -> result::Result<HashMap<String, String>, Box<dyn error::Error>> {
    let content = fs::read_to_string(path)?;
    let mut manifest_hashmap: HashMap<String, String> = HashMap::new();
    for line in content.lines() {
        if let Some(seperator_index) = line.find(": ") {
            let (key, value) = line.split_at(seperator_index);
            let value = &value[2..];
            manifest_hashmap.insert(key.to_string(), value.to_string());
        }
    }
    Ok(manifest_hashmap)
}