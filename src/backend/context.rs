use std::{env, error, io, result};

use super::{lock::{self, StateLockFile}, project::get_config_from_fs, Config};

/// Represents the context of the currently loaded project.
pub struct ProjectContext {
    /// Project config (espresso.toml)
    pub config: Config,
    /// Lock File
    pub state_lock_file: StateLockFile,
    /// Absolute paths (with suffixes known at compile time) that're relavent to this project (ex: path to src)
    pub absolute_paths: AbsoltuePaths,
    /// Absolute paths (with suffixes NOT known at compile time) that're relavent to this project (ex: path to base package)
    pub dynamic_absolute_paths: DynamicAbsolutePaths,
    /// If we're running in debug mode (ESPRESSO_DEBUG=1)
    pub debug_mode: bool,
}

/// Contains absolute paths to critical resources within the currently loaded project.
pub struct AbsoltuePaths {
    /// Path to the currently loaded projects directory. Should be the current working directory.
    pub project: String,
    /// Path to the src/ directory within the currently loaded project.
    pub source: String,
    /// Path to the config file within the currently loaded project.
    pub config: String,
    /// Path to the directory that contains inner working files (ex: state_lockfile, dependency jars, etc)
    pub inner_workings: String,
    /// Path to the directory containing downloaded dependencies.
    pub dependencies: String,
    /// Path to the state lockfile within the currently loaded project.
    pub state_lockfile: String,
}

/// Contains absolute paths to critical resources within the currently loaded project. Determined at runtime.
pub struct DynamicAbsolutePaths {
    /// Path to the base package. Should be {source}/package/path/here. The Main.java file will live here.
    pub base_package: String,
}

/// Get if debug mode is active. You can enable debug mode by setting the `ESPRESSO_DEBUG`
/// environment variable to `1`.
///
/// # Returns
///
/// `true` if `ESPRESSO_DEBUG=1`, `false` if `ESPRESSO_DEBUG=0` or not set
pub fn get_debug_mode() -> bool {
    match env::var("ESPRESSO_DEBUG") {
        Ok(v) => {
            if v == "1" {
                return true;
            } else if v == "0" {
                return false;
            } else {
                return false;
            }
        }
        Err(_) => return false,
    };
}

/// Get an AbsolutePaths struct
///
/// # Arguments
///
/// * `config`: Reference to a Config
/// * `debug_mode`: Reference to a bool if we're in debug mode
///
/// # Returns
///
/// AbsolutePaths
pub fn get_absolute_paths(debug_mode: &bool) -> io::Result<AbsoltuePaths> {
    let cwd = env::current_dir()?;
    let mut cwd_string: String = cwd.to_string_lossy().into_owned();

    if *debug_mode {
        cwd_string += "/espresso_debug";
    }

    Ok(AbsoltuePaths {
        project: cwd_string.clone(),
        source: cwd_string.clone() + "/src/java",
        config: cwd_string.clone() + "/espresso.toml",
        inner_workings: cwd_string.clone() + "/.espresso",
        dependencies: cwd_string.clone() + "/.espresso/dependencies",
        state_lockfile: cwd_string.clone() + "/.espresso/state.lock.toml"
    })
}

/// Get a DynamicAbsolutePaths struct.
///
/// # Arguments
/// * `ap`: Reference to an `AbsolutePaths` struct. Used to get the `src/` directory.
/// * `config`: Reference to a `Config` struct. used to get the current `base_package`.
/// * `debug_mode`: Reference to a bool if we're in debug mode
///
/// # Returns
///
/// DynamicAbsolutePaths
pub fn get_dynamic_absolute_paths(ap: &AbsoltuePaths, config: &Config) -> DynamicAbsolutePaths {
    let base_package = ap.source.clone()
            + "/" + config
                .project
                .base_package
                .clone()
                .replace(".", "/")
                .as_str();
    DynamicAbsolutePaths { base_package }
}

/// Get context about the currently loaded project.
///
/// # Returns
///
/// ProjectContext
pub fn get_project_context() -> result::Result<ProjectContext, Box<dyn error::Error>> {
    let debug_mode = get_debug_mode();
    let absolute_paths = get_absolute_paths(&debug_mode)?;
    let config = get_config_from_fs(&absolute_paths)?;
    let state_lock_file = lock::get_state_lockfile_from_fs(&absolute_paths)?;
    let dynamic_absolute_paths = get_dynamic_absolute_paths(&absolute_paths, &config);
    
    Ok(ProjectContext {
        config,
        state_lock_file,
        absolute_paths,
        debug_mode,
        dynamic_absolute_paths,
    })
}