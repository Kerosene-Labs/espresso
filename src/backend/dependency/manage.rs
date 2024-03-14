use std::{collections, result};

use std::error;

use crate::backend::context::ProjectContext;

/// Add a .jar from the filesystem to the dependencies. This will update the `espreso.toml` file.
/// 
/// # Arguments
/// * `path`: Reference to a string of the filesystem path to the .jar. For example, `/home/user/Downloads/artifact.jar`.
/// * `name`: Reference to a string of the name of the dependency. For example, `lombok`.
/// * `version`: Reference to a string of the version of the dependency. For example, `1.0.0`.
/// * `p_ctx`: Reference to a `ProjectContext` struct
/// 
/// # Returns
/// Propagated `error:Error`(s)
pub fn add(path: &String, name: &String, version: &String, p_ctx: &ProjectContext) -> result::Result<(), Box<dyn error::Error>> {
    todo!()
}
