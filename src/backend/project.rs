use crate::util::pathutil;

use super::context::{AbsoltuePaths, ProjectContext};
use serde::{Deserialize, Serialize};
use std::{collections::HashMap, fs, io, error};

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
 * Load the project at the current working directory
 */
pub fn get_config_from_fs(ap: &AbsoltuePaths) -> Result<Config, Box<dyn error::Error>>{
    let contents = fs::read_to_string(ap.config.clone())?;
    let x: Config = toml::from_str(&contents)?;
    Ok(x)
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
pub fn initialize_source_tree(p_ctx: &ProjectContext) -> io::Result<()>{
    // get the base backage (dot notation) and the base package path on the fs
    let base_package_path = p_ctx.dynamic_absolute_paths.base_package.clone();
    let base_package = p_ctx.config.project.base_package.clone();

    // ensure the base package path exists
    std::fs::create_dir_all(&base_package_path)?;

    // create the Main.java file (textwrap doesn't work????)
    let base_java_content = r#"package ${BASE_PACKAGE};
import java.lang.System;

public class Main {
    public static void main(String[] args) {
        System.out.println("Hello, world!");
    }
}"#
    .replace("${BASE_PACKAGE}", &base_package);

    // write an example java file
    std::fs::write(
        base_package_path.clone() + "/Main.java",
        base_java_content,
    )?;

    Ok(())
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
pub fn initialize_config(name: String, base_package: String, ap: &AbsoltuePaths) -> io::Result<()> {

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
    fs::write(ap.config.clone(), toml_string)?;
    Ok(())
}


/// Ensure the project environment is properly setup
/// 
/// # Arguments
/// * `ap`: Reference to an `AbsolutePaths` struct
/// * `debug_mode`: Reference to a bool that defines if we're in debug mode or not
/// 
/// # Returns
/// `io::Result`, propagated from `fs::create_dir`
pub fn ensure_environment(ap: &AbsoltuePaths, debug_mode: &bool) -> io::Result<()>{
    if *debug_mode {
        fs::create_dir(&ap.project)?
    }
    Ok(())
}