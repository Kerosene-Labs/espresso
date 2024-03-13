use std::{fs, io::Error, vec};

/**
 * Get a list of all files in a directory. Will recurse into subdirectories.
 */
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
