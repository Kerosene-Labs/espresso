use std::error;

use tokio::{fs::File, io::AsyncWriteExt};

/// Download a file from the internet.
/// 
/// # Arguments
/// * `url`: URL of the file
/// * `path`: Filesystem path to the desired location
/// 
/// # Returns
/// Propagates any errors
pub async fn download_file(url: &String, path: &String) -> Result<(), Box<dyn error::Error>> {
    let response = reqwest::get(url).await?;

    if response.status().is_success() {
        let body = response.bytes().await?;
        let mut file = File::create(path).await?;
        file.write_all(&body).await?;
    }
    Ok(())
}