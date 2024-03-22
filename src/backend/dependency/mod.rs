use std::{error, result};

use self::resolve::Package;

use super::context::ProjectContext;
use super::context::AbsoltuePaths;
pub mod resolve;
pub mod uberjar;

/// Download the latest version of the package, adding it to the state.lock.toml & cargoespresso.toml file(s)
/// 
/// # Arguments
/// * `p_ctx` Reference to a `ProjectContext` struct
/// * `package` Reference to the `Package` to be added
/// 
/// # Returns
/// Propagated `error::Error`
pub async fn add(p_ctx: &ProjectContext, ap: &AbsoltuePaths, package: &Package) -> result::Result<(), Box<dyn error::Error>> {
    // download the package
    resolve::download(p_ctx, package).await?;

    Ok(())
}