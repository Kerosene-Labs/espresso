use core::panic;
use std::{error, path::Path, result};
use tokio::fs;
use sha2::{Digest, Sha512};
use super::error::EspressoError;

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

/// Get the SHA512 checksum of a path
/// 
/// # Arguments
/// * `path`: Reference to a String of the path to checksum
/// 
/// # Returns
/// SHA512 checksum as a hex string, propagated errors.
pub async fn get_sha512_of_path(
    path: &String
) -> result::Result<String, Box<dyn error::Error>> {
    let contents = fs::read(path).await?;

    // hash it to sha512
    let mut hasher = Sha512::new();
    hasher.update(contents);
    let sha512hex = hex::encode(hasher.finalize().to_vec());

    Ok(sha512hex)
}


/// Ensure integirty of a file
///
/// # Arguments
/// * `path`: Reference to a `String` containing the path of the file to check
/// * `expected_sha512_hex`: SHA512 hexadecimal string to compare against
pub async fn ensure_integrity_sha512(
    path: &String,
    expected_sha512_hex: &String,
) -> result::Result<(), Box<dyn error::Error>> {
    let sha512hex = get_sha512_of_path(path).await?;

    // if the expected sha512 hex doesn't equal what we calculated, throw an error
    if *expected_sha512_hex != *sha512hex {
        return Err(EspressoError::nib(format!("Downloaded file does not have the same SHA512 checksum as defined by the package: Expected={}: Actual={}", expected_sha512_hex, sha512hex).as_str()));
    }
    Ok(())
}

/// Get the parent directory of a file. For example, with `/home/x/EspressoProjects/project/espresso.toml` as
/// the full path, `/home/x/EspressoProjects/project` would be the parent path.
/// 
/// # Arguments
/// * `file_path`: A reference to a `String`, containing the path to the file.
/// 
/// # Returns
/// `String`, the parent path.
/// 
#[deprecated(note = "Use std::path::Path::parent() instead")]
pub fn get_parent_path_of_file(file_path: &String) -> result::Result<String, Box<dyn error::Error>> {
    let split_path: Vec<&str> = file_path.split("/").collect();
    let mut parent_path: String = String::new();
    if let Some(slice) = split_path.get(0..split_path.len() - 1) {
        for part in slice {
            parent_path.push_str(part);
            parent_path.push_str("/")
        }
    } else {
        return Err(EspressoError::nib(""));
    }
    Ok(parent_path)

}
