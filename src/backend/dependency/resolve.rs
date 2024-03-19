use std::{collections::HashMap, error, fmt::format, result};
use serde::{Serialize, Deserialize};

use crate::{backend::context::ProjectContext, frontend::terminal::print_err, util::{self, error::EspressoError, net::download_file}};

/// Represents a resolved dependency
#[derive(Serialize, Deserialize)]
pub struct QueryPackagesResponse {
    pub packages: Vec<Package>
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Package {
    pub metadata: PackageMetadata,
    pub group_id: String,
    pub artifact_id: String,
    #[serde(rename="ref")]
    pub ref_: String,
}

/// Represents a package's metadata.
#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct PackageMetadata {
    pub source_repository: String,
    pub versions: Vec<PackageVersion>,
}

/// Represents a specific release/version of the package.
#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct PackageVersion {
    pub version: String,
    pub flags: Vec<Flags>,
    pub vulnerabilities: HashMap<String, String>,
    pub artifact_url: String,
    pub sha512sum: String,
}

/// Represents the supported package types. This will dictate how they're applied at compile time.
#[derive(Serialize, Deserialize, Debug, Clone)]
pub enum Flags {
    #[serde(rename = "annotation_processor")]
    AnnotationProcessor,
}

/// Query for packages from the Espresso Registry
/// 
/// # Arguments
/// * `q`: The query string
/// 
/// # Returns
/// Propagated errors, returns a `Vec` of `Package` struct(s). The `Vec` will be empty if no packages were returned in the query.
pub async fn query(q: &String) -> result::Result<Vec<Package>, Box<dyn error::Error>> {
    let client = reqwest::Client::new();

    // make a request to the registry
    let response = client.get("https://registry.espresso.hlafaille.xyz/v1/search")
        .query(&[("q", q)])
        .send()
        .await?;

    // handle our response from the registry
    let response_text = response.text().await?;
    let query_packages_response: QueryPackagesResponse = match serde_json::from_str(&response_text) {
        Ok(v) => v,
        Err(_) => {
            return Err(
                EspressoError::nib(format!("Failed to deserialize response: Response content was: {}", response_text).as_str())
            )
        }
    };
    Ok(query_packages_response.packages)
}

/// Download the latest version of a package
async fn download(p_ctx: &ProjectContext, package: &Package) -> result::Result<(), Box<dyn error::Error>> {
    // get the latest version of this project
    let version = match package.metadata.versions.get(0) {
        Some(v) => v,
        None => {
            return Err(EspressoError::nib("Failed to get the latest version of the package (there are no versions)"))
        }
    };

    // establish our full path
    let download_path = p_ctx.absolute_paths.dependencies.clone() + format!("/{}.jar", version.sha512sum).as_str();

    // TODO 
    // download the file
    download_file(&version.artifact_url, &download_path).await?;

    // ensure integrity
    util::pathutil::ensure_integrity_sha512(&download_path, &version.sha512sum).await?;
    
    Ok(())
}

/// Download the latest version of the package, adding it to the state.lock.toml & cargoespresso.toml file(s)
/// 
/// # Arguments
/// * `p_ctx` Reference to a `ProjectContext` struct
/// * `package` Reference to the `Package` to be added
/// 
/// # Returns
/// Propagated `error::Error`
pub async fn add(p_ctx: &ProjectContext, package: &Package) -> result::Result<(), Box<dyn error::Error>> {
    // download the package
    download(p_ctx, package).await?;
    
    // perform 

    Ok(())
}