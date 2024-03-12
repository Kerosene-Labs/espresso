use std::path::Path;

/// If a path on the filesystem exists
///
/// # Arguments
/// * `path`: Reference to a `String` containing the path to the file you want to check for existence.
///
/// # Returns
/// `true` if exists, `false` if not.
pub fn does_path_exist(path: &String) -> bool {
    Path::exists(Path::new(path))
}
