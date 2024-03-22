use serde::{Deserialize, Serialize};

pub mod context;
pub mod project;
pub mod toolchain;
pub mod lock;
pub mod dependency;

/// Represents an `espresso.toml` file
#[derive(Deserialize, Serialize, Debug)]
pub struct Config {
    pub project: Project,
    pub toolchain: Toolchain,
    /// Dependencies located on your filesystem
    pub dependencies_fs: std::collections::HashMap<String, String>,
    /// Dependencies from the Espresso Registry
    pub dependencies: std::collections::HashMap<String, String>
}

/// Represents information about the currently loaded Project
#[derive(Deserialize, Serialize, Debug)]
pub struct Project {
    /// Name of the project (ex: `My Espresso Project`)
    pub name: String,
    /// Version of the project (ex: `1.0.0`)
    pub version: String,
    /// Java base package in dot notation (ex: `com.me.project`)
    pub base_package: String,
}

/// Represents toolchain information
#[derive(Deserialize, Serialize, Debug)]
pub struct Toolchain {
    /// Path to the JDK toolchain (ex: `${JAVA_HOME}`)
    pub path: String,
}