use std::{fs, path::Path};

const ESPRESSO_CONFIG_PATH: &str = "espresso.toml";
const ESPRESSO_CONFIG_BASE_CONTENT: &str = r#"[project]
name = "My Espresso Project"
version = "1.0.0"
artifact = "build.jar"

[toolchain]
path = "${JAVA_HOME}"

[dependencies]
"#;

/**
 * Load the project at the current working directory
 */
pub fn load() {
    let contents = fs::read_to_string(ESPRESSO_CONFIG_PATH).expect("Unable to read conig file");
    println!("{}", contents);
}

/**
 * Checks if the project config exists
 */
pub fn does_exist() -> bool {
    Path::exists(Path::new(ESPRESSO_CONFIG_PATH))
}

/**
 * Initialize a new project
 */
pub fn initialize() {
    fs::write(ESPRESSO_CONFIG_PATH, ESPRESSO_CONFIG_BASE_CONTENT)
        .expect("Failed to write config file")
}
