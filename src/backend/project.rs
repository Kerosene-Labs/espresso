use std::{collections::HashMap, env, fs::{self, create_dir_all}, path::Path};
use serde::{Deserialize, Serialize};
use super::context::ProjectContext;

#[derive(Deserialize, Serialize, Debug)]
pub struct Config {
    pub project: Project,
    pub toolchain: Toolchain,
    pub dependencies: std::collections::HashMap<String, String>,
}

#[derive(Deserialize, Serialize, Debug)]
pub struct Project {
    pub name: String,
    pub version: String,
    pub artifact: String,
    pub base_package: String,
}

#[derive(Deserialize, Serialize, Debug)]
pub struct Toolchain {
    pub path: String,
}

/**
 * Get if ESPRESSO_DEBUG=1
 */
pub fn get_debug_mode() -> bool {
    match env::var("ESPRESSO_DEBUG") {
        Ok(v) => {
            if v == "1" {
                return true
            } else if v == "0" {
                return false
            } else {
                return false
            }
        }
        Err(_) => return false,
    };
}

/**
 * Get the config path. Note, this value changes if ESPRESSO_DEBUG=1
 */
pub fn get_config_path() -> String {
    let debug_mode = get_debug_mode();
    if debug_mode {
        "espresso_debug/espresso.toml".to_string()
    } else {
        "espresso.toml".to_string()
    }
}

pub fn get_absolute_project_path() -> String {
    let debug_mode = get_debug_mode();
    let current_dir = env::current_dir().unwrap().to_string_lossy().into_owned();
    if debug_mode {
        current_dir + "/espresso_debug"
    } else {
        current_dir
    }
}

/**
 * Get the source path. Note, this value is prefixed with `espresso_debug` if ESPRESSO_DEBUG=1
 */
pub fn get_source_path() -> String {
    (get_absolute_project_path() + "/src/java").to_string()
}

/**
 * Ensure development directory exists
 */
pub fn ensure_debug_directory_exists_if_debug(){
    if !get_debug_mode() {
        return;
    }
    if !Path::exists(Path::new("espresso_debug")) {
        create_dir_all("espresso_debug").expect("Failed to ensure debug directory exists");
    }
}

/**
 * Load the project at the current working directory
 */
pub fn get_config_from_fs() -> Config {
    let contents = fs::read_to_string(get_config_path()).expect("Unable to read conig file");
    toml::from_str(&contents).unwrap()
}

/**
 * If a project exists. A project is deemed existing if it has a source directory
 * and a config file.
 */
pub fn does_exist() -> bool {
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
    Path::exists(Path::new(get_source_path().as_str()))
}

/**
 * Checks if the config exists
 */
pub fn does_config_exist() -> bool {
    Path::exists(Path::new(get_config_path().as_str()))
}

/**
 * Get the base package path. This value is the `location of src + base_package`
 */
pub fn get_full_base_package_path(p_ctx: &ProjectContext) -> String{
    format!("{}/{}", get_source_path(), p_ctx.config.project.base_package.replace(".", "/"))
}

/**
 * Initialize the source tree
 */
pub fn initialize_source_tree(p_ctx: &ProjectContext) {
    std::fs::create_dir_all(get_full_base_package_path(p_ctx)).expect("failed to create main package directories in file system");

    // create the Main.java file (textwrap doesn't work????)
    let base_java_content = r#"package ${BASE_PACKAGE};
import java.lang.System;

public class Main {
    public static void main(String[] args) {
        System.out.println("Hello, world!");
    }
}"#.replace("${BASE_PACKAGE}", &p_ctx.config.project.base_package);
    
    std::fs::write(get_full_base_package_path(p_ctx) + "/Main.java", base_java_content);
}

fn process_input(x: String, default: String) -> String {
    let new = x.replace("\n", "");
    if new.is_empty() {
        return default;
    }
    return new;
}

/**
 * Initialize a config
 */
pub fn initialize_config(name: String, base_package: String) {
    // process the name
    // populate a base_config struct
    let base_config = Config {
        project: Project {
            name: process_input(name, "My Espresso Project".to_string()),
            version: "1.0.0".to_string(),
            artifact: "artifact.jar".to_string(),
            base_package: process_input(base_package, "com.me.myespressoproject".to_string()),
        },
        toolchain: Toolchain {
            path: "${JAVA_HOME}".to_string(),
        },
        dependencies: HashMap::new(),
    };

    // write it to a toml string, then write it to the config file
    let toml_string = toml::to_string(&base_config).expect("Failed to serialize struct");
    fs::write(get_config_path(), toml_string).expect("Failed to write config file")
}
