use std::{
    collections::{self, HashMap},
    error, fs, result,
};

use serde::{Deserialize, Serialize};

use super::context::AbsoltuePaths;

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct StateLockFile {
    /// The dependencies that should live in `.espresso/dependencies`. The key is the dependency name and the value is the dependency SHA512 checksum
    pub dependencies: collections::HashMap<String, String>,
}

/// Represents a dependency in the state lock file
#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct StateLockFileDependency {
    /// The dependency name
    pub name: String,
    /// The dependency sha512 checksum
    pub checksum: String,
}

/// Get the `.espresso/state.lock.toml` file from the filesystem.
///
/// # Arguments
/// * `ap`: Reference to an `AbsolutePaths` struct
///
/// # Returns
/// Returns a `StateLockFile` struct on success, propagating any errors in the process
pub fn get_state_lockfile_from_fs(
    ap: &AbsoltuePaths,
) -> result::Result<StateLockFile, Box<dyn error::Error>> {
    let contents = fs::read_to_string(&ap.state_lockfile)?;
    let x: StateLockFile = toml::from_str(&contents)?;
    Ok(x)
}

/// Initialize a new state lockfile
pub fn initialize_state_lockfile(ap: &AbsoltuePaths) -> result::Result<(), Box<dyn error::Error>> {
    let base = StateLockFile {
        dependencies: HashMap::new(),
    };

    let toml_string = toml::to_string(&base)?;
    fs::write(ap.state_lockfile.clone(), toml_string)?;
    Ok(())
}

/// Add a dependency to the state lock file
///
/// # Arguments
/// * `slf`: The `StateLockFile` to update
/// * `name`: The name of the dependency
/// * `sha512sum`: The SHA512 checksum of the dependency
///
/// # Returns
/// The updated `StateLockFile` & propagated errors
pub fn add_dependency(
    slf: StateLockFile,
    name: &String,
    sha512sum: &String,
) -> result::Result<StateLockFile, Box<dyn error::Error>> {
    Ok(slf)
}
