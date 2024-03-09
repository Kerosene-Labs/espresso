use std::{fs, path::Path};

use serde::{Deserialize, Serialize};
const ESPRESSO_SOURCE_PATH: &str = "src";
const ESPRESSO_CONFIG_PATH: &str = "espresso.toml";

#[derive(Deserialize, Serialize)]
pub struct Config {
    pub project: Project,
    pub toolchain: Toolchain,
    pub dependencies: Vec<String>
}

#[derive(Deserialize, Serialize)]
pub struct Project {
    pub name: String,
    pub version: String,
    pub artifact: String
}

#[derive(Deserialize, Serialize)]
pub struct Toolchain{
    pub path: String
}


/**
 * Load the project at the current working directory
 */
pub fn load() -> Config{
    let contents = fs::read_to_string(ESPRESSO_CONFIG_PATH).expect("Unable to read conig file");
    toml::from_str(&contents).unwrap()
}

/**
 * If a project exists. A project is deemed existing if it has a source directory
 * and a config file.
 */
pub fn does_exist() -> bool{
    let source_exists = does_source_exist();
    let config_exists = does_config_exist();

    // Return false if either source or config does not exist
    if !source_exists || !config_exists {
        return false;
    }

    // Return true if both source and config exist
    true
}

/**
 * If the source path (ex: src) exists
 */
pub fn does_source_exist() -> bool {
    Path::exists(Path::new(ESPRESSO_SOURCE_PATH))
}

/**
 * Checks if the config exists
 */
pub fn does_config_exist() -> bool {
    Path::exists(Path::new(ESPRESSO_CONFIG_PATH))
}

/**
 * Initialize a new project
 */
pub fn initialize() {
    // populate a base_config struct
    let base_config = Config {
        project: Project {
            name: "My Espresso Project".to_string(),
            version: "1.0.0".to_string(),
            artifact: "artifact.jar".to_string()
        },
        toolchain: Toolchain {
            path: "${JAVA_HOME}".to_string()
        },
        dependencies: vec![]      
    };

    // write it to a toml string, then write it to the config file
    let toml_string = toml::to_string(&base_config).expect("Failed to serialize struct");
    fs::write(ESPRESSO_CONFIG_PATH, toml_string)
        .expect("Failed to write config file")
}
