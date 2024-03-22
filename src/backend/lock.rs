use std::{
    collections::{self, HashMap},
    error, fs, result,
};

use serde::{Deserialize, Serialize};

use super::context::AbsoltuePaths;

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct StateLockFile {
    /// The dependencies that should live in `.espresso/dependencies`. The key is the dependency name and the value is the dependency SHA512 checksum
    pub dependencies: collections::HashMap<String, StateLockFileDependency>,
}

/// Represents a dependency in the state lock file
#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct StateLockFileDependency {
    /// The dependency sha512 checksum
    pub checksum: String,
    /// The dependency source
    pub source: StateLockFileDependencySource
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub enum StateLockFileDependencySource {
    EspressoRegistry,
    FileSystem
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

/// Write the lock file
///
/// # Arguments
/// * `slf`: Reference to the `StateLockFile` to write
///
/// # Returns
/// Propagated errors
pub fn write_lock_file(
    slf: &StateLockFile,
    ap: &AbsoltuePaths
) -> result::Result<(), Box<dyn error::Error>> {
    // write it to a toml string, then write it to the config file
    let toml_string = toml::to_string(slf)?;
    fs::write(ap.state_lockfile.clone(), toml_string)?;
    Ok(())
}
