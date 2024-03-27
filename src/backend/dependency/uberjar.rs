use std::{error, fs, result};

use crate::{backend::{context::ProjectContext, dependency::manifest, lock::StateLockFileDependency, toolchain::{self, ToolchainContext}}, util};

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

/// Copy classes from 
/// 
/// # Arguments
/// * `p_ctx`: Reference to a `ProjectContext` struct
/// * `tc_ctx`: Reference to a `ToolChainContext` struct
/// * `dependency`: Reference to a `StateLockFileDependency`, representing the local dependency to extract.
///
/// # Returns
/// Propagated errors
pub fn copy_classes(
    p_ctx: &ProjectContext,
    dependency: &StateLockFileDependency
) -> result::Result<(), Box<dyn error::Error>>{
    let source = p_ctx.absolute_paths.dependencies_extracted.clone() + "/" + dependency.checksum.as_str();

    // get our directory tree
    let dir_tree = util::directory::walk_dir_tree(&source)?;

    let manifest = manifest::parse(&(source + "/META-INF/MANIFEST.MF"));
    
    // copy
    // println!("{:?}", class_files);
    Ok(())
}