use crate::util::pathutil;

use super::context::{AbsoltuePaths, DynamicAbsolutePaths, ProjectContext};
use serde::{Deserialize, Serialize};
use std::{
    collections::HashMap,
    env,
    fs::{self, create_dir_all},
    path::Path,
};

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
    pub base_package: String,
}

#[derive(Deserialize, Serialize, Debug)]
pub struct Toolchain {
    pub path: String,
}

/**
 * Ensure development directory exists
 */
pub fn ensure_debug_directory_exists_if_debug(p_ctx: &ProjectContext) {
    if p_ctx.debug_mode {
        if !Path::exists(Path::new("espresso_debug")) {
            create_dir_all("espresso_debug").expect("Failed to ensure debug directory exists");
        }
    }
}

/**
 * Load the project at the current working directory
 */
pub fn get_config_from_fs(ap: &AbsoltuePaths) -> Config {
    let contents = fs::read_to_string(ap.config).expect("Unable to read conig file");
    toml::from_str(&contents).unwrap()
}

/**
 * If a project exists. A project is deemed existing if it has a source directory
 * and a config file.
 */
pub fn does_exist(ap: &AbsoltuePaths) -> bool {
    let source_exists = pathutil::does_path_exist(&ap.source);
    let config_exists = pathutil::does_path_exist(&ap.config);

    // Return false if either source or config does not exist
    if !source_exists || !config_exists {
        return false;
    }

    // Return true if both source and config exist
    true
}

/**
 * Initialize the source tree
 */
pub fn initialize_source_tree(p_ctx: &ProjectContext) {
    let base_package_path = p_ctx.dynamic_absolute_paths.base_package.take();
    std::fs::create_dir_all(base_package_path)
        .expect("failed to create main package directories in file system");

    // create the Main.java file (textwrap doesn't work????)
    let base_java_content = r#"package ${BASE_PACKAGE};
import java.lang.System;

public class Main {
    public static void main(String[] args) {
        System.out.println("Hello, world!");
    }
}"#
    .replace("${BASE_PACKAGE}", &p_ctx.config.project.base_package);

    std::fs::write(base_package_path + "/Main.java", base_java_content);
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
pub fn initialize_config(name: String, base_package: String, ap: &AbsoltuePaths) {
    // process the name
    // populate a base_config struct
    let base_config = Config {
        project: Project {
            name: process_input(name, "My Espresso Project".to_string()),
            version: "1.0.0".to_string(),
            base_package: process_input(base_package, "com.me.myespressoproject".to_string()),
        },
        toolchain: Toolchain {
            path: "${JAVA_HOME}".to_string(),
        },
        dependencies: HashMap::new(),
    };

    // write it to a toml string, then write it to the config file
    let toml_string = toml::to_string(&base_config).expect("Failed to serialize struct");
    fs::write(ap.config, toml_string).expect("Failed to write config file")
}
