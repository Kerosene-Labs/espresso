use crate::backend::project::{self, load_config};
use crate::frontend::terminal::print_err;
use std::process::exit;

use super::terminal::print_general;

/**
 * Service function for the `build` command
 */
pub fn build() {
    load_config();
    print_general("Building project...");
    unimplemented!("build project");
}

/**
 * Service function for the `init` command
 */
pub fn init() {
    // check if the project exists
    if project::does_exist() {
        print_err("Unable to create project: Espresso project already exists");
        exit(1);
    }

    // if the project doesn't exist, create it
    project::initialize_config();
    print_general("Project created: Edit espresso.toml to check it out!");
}
