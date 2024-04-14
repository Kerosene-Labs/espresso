use std::{error, fs, io::Error, path, result, vec};

/// Get a list of all files within a directory
///
/// # Arguments
/// * `path`: A `String`, the path to search
///
/// # Returns
/// `Result`, where `Ok` is a `Vec<String>` containing all files within the directory. Propagated errors.
#[deprecated(note = "use walk_dirs_for_dirs or walk_dir_for_files")]
pub fn read_files_recursively(path: String) -> Result<Vec<String>, Error> {
    let mut files: Vec<String> = vec![];
    for entry in fs::read_dir(path)? {
        match entry {
            Ok(f) => {
                let path = f.path().to_string_lossy().into_owned();
                if f.file_type()?.is_dir() {
                    files = [files, read_files_recursively(path)?].concat();
                } else {
                    files.push(path);
                }
            }
            Err(_) => unimplemented!("error case on reading files recursively"),
        }
    }
    return Ok(files);
}

/// Walk a directory recursively for other directories, effectively assembling a directory tree.
///
/// # Arguments
/// * `path`: A reference to a `String`, the directory to walk through.
///
/// # Returns
/// Propagated errors, a `Vec<String>` containing absolute paths to all the dirs.
pub fn walk_dir_tree(path: &path::PathBuf) -> result::Result<Vec<path::PathBuf>, Box<dyn error::Error>> {
    let mut dirs: Vec<path::PathBuf> = vec![];
    for i in fs::read_dir(path)? {
        let entry = i?;
        if entry.file_type()?.is_dir() {
            let path = entry.path().into();
            let rec_dir_tree = walk_dir_tree(&path)?;
            if rec_dir_tree.len() == 0 {
                dirs.push(path);
            } else {
                dirs = [dirs, walk_dir_tree(&path)?].concat();
            }
        }
    }
    Ok(dirs)
}

/// Walk a directory recursively for files, effectively assembling a tree of files.
///
/// # Arguments
/// * `path`: A reference to a `String`, the directory to walk through.
///
/// # Returns
/// Propagated errors, a `Vec<String>` containing absolute paths to all the files.
pub fn walk_file_tree(path: &String) -> result::Result<Vec<String>, Box<dyn error::Error>> {
    let mut files: Vec<String> = vec![];
    for i in fs::read_dir(path)? {
        let entry = i?;
        let path = entry.path().to_string_lossy().into_owned();
        if entry.file_type()?.is_dir() {
            files = [files, walk_file_tree(&path)?].concat();
        } else {
            files.push(path);
        }
    }
    Ok(files)
}
