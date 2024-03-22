use std::{error, result};

use crate::util;

use self::resolve::Package;

use super::context::AbsoltuePaths;
use super::context::ProjectContext;
use super::lock;
use super::lock::StateLockFileDependency;
use super::lock::StateLockFileDependencySource;
use super::project;
pub mod resolve;
pub mod uberjar;

/// Download the latest version of the package, adding it to the state.lock.toml & cargoespresso.toml file(s)
///
/// # Arguments
/// * `p_ctx` Temporary ownership of a `ProjectContext` struct. Must be mutable.
/// * `package` Reference to the `Package` to be added
///
/// # Returns
/// Propagated `error::Error`
pub async fn add(
    p_ctx: &mut ProjectContext,
    ap: &AbsoltuePaths,
    package: &Package,
) -> result::Result<(), Box<dyn error::Error>> {
    // TODO ensure this package does not exist in the config already

    // download the package
    let artifact_path = resolve::download(&p_ctx, package).await?;

    // get our friendly name
    let friendly_name = format!("{}:{}", package.group_id, package.artifact_id).to_string();

    // get the package version
    let latest_version = package
        .metadata
        .versions
        .get(0)
        .expect("Failed to get latest version from package: Versions metadata may be empty");

    // add the package to the config
    p_ctx
        .config
        .dependencies
        .insert(friendly_name.clone(), latest_version.version.clone());
    match project::write_config(&p_ctx.config, ap) {
        Ok(_) => (),
        Err(e) => {
            panic!("Failed to add package to config: Failed to write: {:?}", e)
        }
    }

    // add the package to the state lock file
    p_ctx.state_lock_file.dependencies.insert(
        friendly_name.clone(),
        StateLockFileDependency {
            checksum: util::pathutil::get_sha512_of_path(&artifact_path).await?,
            source: StateLockFileDependencySource::EspressoRegistry
        },
    );
    match lock::write_lock_file(&p_ctx.state_lock_file, ap) {
        Ok(_) => (),
        Err(e) => {
            panic!("Failed to add package to state lock file: Failed to write: {:?}", e)
        }
    }
    
    Ok(())
}
