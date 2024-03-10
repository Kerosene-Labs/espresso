use std::io::{self, Write};

use crate::backend;
use crate::frontend::terminal::{print_err, print_sameline};

use super::terminal::print_general;

/**
 * Service function for the `build` command
 */
pub fn build() {
    print_general("Building project...");
    unimplemented!("build project");
}

/**
 * Service function for the `init` command
 */
pub fn init() {
    print_general("Tell us a bit about your project...");
    // check if the project exists
    if backend::project::does_exist() {
        print_err("Unable to create project: Espresso project (or remnants) already exist");
    }

    // collect the name
    let mut name = String::new();
    print_sameline("Project name: ");
    if let Err(_) = io::stdin().read_line(&mut name) {
        print_err("Failed to read user input for project name")
    }

    // collect the base package
    let mut base_package = String::new();
    print_sameline("Base package: ");
    if let Err(_) = io::stdin().read_line(&mut base_package) {
        print_err("Failed to read user input for base package")
    }

    // initialize the config
    backend::project::initialize_config(name, base_package);

    // get our project context
    let p_ctx = backend::context::get_project_context(); 

    // initialize our source tree
    backend::project::initialize_source_tree(&p_ctx);
    print_general("Project created: Edit espresso.toml to check it out!");
}
