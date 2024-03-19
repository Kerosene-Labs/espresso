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

/// Ensure integirty of a file
///
/// # Arguments
/// * `path`: Reference to a `String` containing the path of the file to check
/// * `expected_sha512_hex`: SHA512 hexadecimal string to compare against
pub async fn ensure_integrity_sha512(
    path: &String,
    expected_sha512_hex: &String,
) -> result::Result<(), Box<dyn error::Error>> {
    // read our file
    let contents = fs::read(path).await?;

    // hash it to sha512
    let mut hasher = Sha512::new();
    hasher.update(contents);
    let sha512hex = hex::encode(hasher.finalize().to_vec());

    // if the expected sha512 hex doesn't equal what we calculated, throw an error
    if *expected_sha512_hex != *sha512hex {
        return Err(EspressoError::nib(format!("Downloaded file does not have the same SHA512 checksum as defined by the package: Expected={}: Actual={}", expected_sha512_hex, sha512hex).as_str()));
    }
    Ok(())
}
