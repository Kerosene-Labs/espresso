use std::{error, fs, result};

use serde::{Deserialize, Serialize};

use super::context::{AbsoltuePaths, DynamicAbsolutePaths};

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct StateLockFile {
    pub dependencies: Vec<Dependency>,
}

/**
 * Represents a dependency in '.espresso/dependencies'.
 */
#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Dependency {
    pub name: String,
    pub fs_name: String,
    pub checksum: String,
}

/// Get the `.espresso/state.lock` file from the filesystem.
/// 
/// # Arguments
/// * `ap`: Reference to an `AbsolutePaths` struct
/// 
/// # Returns
/// Returns a `LockFile` struct on success, propagating any errors in the process
pub fn get_state_lockfile_from_fs(ap: &AbsoltuePaths) -> result::Result<StateLockFile, Box<dyn error::Error>> {
    let contents = fs::read_to_string(&ap.state_lockfile)?;
    let x: StateLockFile = toml::from_str(&contents)?;
    Ok(x)
}

/// Initialize a new state lockfile
pub fn initialize_state_lockfile(ap: &AbsoltuePaths) -> result::Result<(), Box<dyn error::Error>> {
    let base = StateLockFile {
        dependencies: vec![]
    };

    let toml_string = toml::to_string(&base)?;
    fs::write(ap.state_lockfile.clone(), toml_string)?;
    Ok(())
}