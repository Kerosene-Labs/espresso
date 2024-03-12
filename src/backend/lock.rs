use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct LockFile {
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

/**
 * Get '.espresso'
 */
pub fn get_lock_file_from_fs() -> LockFile {

}