use std::{collections::HashMap, error, fmt::format, result};
use serde::{Serialize, Deserialize};

use crate::util::error::EspressoError;

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
pub async fn query(q: String) -> result::Result<Vec<Package>, Box<dyn error::Error>> {
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