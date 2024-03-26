use std::{error, result};

use crate::backend::{context::ProjectContext, lock::StateLockFileDependency, toolchain::{self, ToolchainContext}};

/// Helper function to extract the specified dependency
///
/// # Arguments
/// * `p_ctx`: Reference to a `ProjectContext` struct
/// * `tc_ctx`: Reference to a `ToolChainContext` struct
/// * `dependency`: Reference to a `StateLockFileDependency`, representing the local dependency to extract.
///
/// # Returns
/// Propagated errors
pub fn extract(
    p_ctx: &ProjectContext,
    tc_ctx: &ToolchainContext,
    dependency: &StateLockFileDependency,
) -> result::Result<(), Box<dyn error::Error>> {
    toolchain::extract_jar(p_ctx, tc_ctx, &(dependency.checksum.clone() + ".jar"))
}

/// Sync a dependency into the projcets class
/// 
/// # Arguments
/// * `p_ctx`: Reference to a `ProjectContext` struct
/// * `tc_ctx`: Reference to a `ToolChainContext` struct
/// * `dependency`: Reference to a `StateLockFileDependency`, representing the local dependency to extract.
///
/// # Returns
/// Propagated errors
pub fn sync(
    p_ctx: &ProjectContext,
    tc_ctx: &ToolchainContext,
    dependency: &StateLockFileDependency
) -> result::Result<(), Box<dyn error::Error>>{
    Ok(())
}