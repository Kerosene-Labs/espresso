use std::{error, result, vec};

use crate::{backend::{context::ProjectContext, lock::StateLockFileDependency}, util};

/// Represents the extracted class for a dependency
pub struct ExtractedClass {
    // The absolute path of the directory containing the extracted class (ex: /home/blah/EspressoProjects/.espresso/dependencies_extracted/19283asd/org/springframework/web)
    pub containing_directory_absolute_path: String,
    // The package-relative path of this class (ex: org/springframework/web)
    pub containing_directory_relative_path: String,
    // The class name (minus extension)
    pub name: String,
    // The absolute path to the class file, pre merge.
    pub path: String,
}

/// Get extracted classes for a particular dependency
/// 
/// # Arguments
/// * `p_ctx`: A reference to a `ProjectContext`
/// * `dependency`: A reference to a `StateLockFileDependency`
/// 
/// # Reference
/// Propagated errors, w
pub fn get(p_ctx: &ProjectContext, dependency: &StateLockFileDependency) -> result::Result<Vec<ExtractedClass>, Box<dyn error::Error>> {
    let extracted_class_source = p_ctx.absolute_paths.dependencies_extracted.clone() + "/" + dependency.checksum.as_str();
    let files = util::directory::walk_file_tree(&extracted_class_source)?;

    // sort out class files
    let mut class_files: Vec<String> = vec![];
    for file in files {
        if file.ends_with(".class") {
            class_files.push(file);
        }
    }
    
    // build the extracted class structs
    let extracted_classes: Vec<ExtractedClass> = vec![];
    for file in class_files {
        let absolute_path: Vec<&str> = file.split("/").collect();

        let containing_directory_absolute_path = "";

        extracted_classes.push(
            ExtractedClass {
                containing_directory_absolute_path
            }
        )
    }

    Ok(extracted_classes)
}