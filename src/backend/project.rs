use std::fs;


const ESPRESSO_JSON_PATH: &str = "espresso.toml";

/**
 * Load the project at the current working directory
 */
pub fn load_project() {
    let contents = fs::read_to_string(ESPRESSO_JSON_PATH)
    .expect("Unable to read conig file");
    println!("{}", contents);
}