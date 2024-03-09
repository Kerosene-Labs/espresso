use std::process::exit;

use crate::backend::project::{self, load};

/**
 * Service function for the `build` command
 */
pub fn build() {
    let loaded_project = project::load();

}

/**
 * Service function for the `init` command
 */
pub fn init() {
    // check if the project exists
    if project::does_exist() {
        eprintln!("Unable to create project: Espresso project already exists");
        exit(1);
    }

    // if the project doesn't exist, create it
    project::initialize();
    println!("Project created: Edit espresso.toml to check it out!");
}