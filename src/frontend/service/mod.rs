use crate::backend::context::{AbsoltuePaths, ProjectContext};
use crate::backend::toolchain::{
    compile_project, run_jar, ToolchainContext,
};
use crate::backend::{self, context, project};
use crate::frontend::terminal::{print_err, print_sameline};
use std::{error, io, result};

use super::terminal::print_general;

/**
 * Service function for the `run` command
 */
pub fn run(
    mut p_ctx: ProjectContext,
    mut tc_ctx: ToolchainContext,
) -> result::Result<(), Box<dyn error::Error>> {
    // build our jar
    (p_ctx, tc_ctx) = build(p_ctx, tc_ctx)?;

    // run it
    print_general("-- RUNNING ARTIFACT -----");
    match run_jar(&p_ctx, &tc_ctx) {
        Ok(_) => (),
        Err(e) => {
            print_err(format!("Failed to run 'artifact.jar': {}", {e}).as_str())
        }
    }
    Ok(())
}

/**
 * Service function for the `build` command
 */
pub fn build(
    p_ctx: ProjectContext,
    tc_ctx: ToolchainContext,
) -> result::Result<(ProjectContext, ToolchainContext), Box<dyn error::Error>> {
    // walk our src directory, find source files
    print_general("-- DISCOVERING ----------");
    let java_files = backend::toolchain::get_java_source_files(&p_ctx).unwrap();
    print_general(
        format!(
            "Discovered {} source file(s) in base package '{}'",
            java_files.len(),
            &p_ctx.config.project.base_package,
        )
        .as_str(),
    );
    print_general("------------------------");


    // compile the project to class files
    print_general("-- COMPILING ------------");
    compile_project(java_files, &p_ctx, &tc_ctx);
    print_general("------------------------");

    // build our jar
    print_general("-- PACKAGING ------------");
    match backend::toolchain::build_jar(&p_ctx, &tc_ctx) {
        Ok(_) => (),
        Err(e) => {
            print_err(format!("Failed to build jar: {}", {e}).as_str());
        }
    }
    print_general("------------------------");

    // pass ownership back to the caller
    Ok((p_ctx, tc_ctx))
}

/**
 * Service function for the `init` command
 */
pub fn init() {
    // get absolute paths
    let ap: AbsoltuePaths = match context::get_absolute_paths(&context::get_debug_mode()) {
        Err(_) => {
            print_general("Failed to get absolute paths");
            return;
        }
        Ok(x) => x,
    };

    // check if the project exists
    if project::does_exist(&ap) {
        print_err(
            "Unable to initialize project: An Espresso project (or remnants of one) already exist",
        );
    }

    print_general("Tell us a bit about your project!");

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

    // ensure our environment is setup
    match backend::project::ensure_environment(&ap, &backend::context::get_debug_mode()) {
        Ok(_) => (),
        Err(e) => {
            print_err(format!("Failed to ensure environment: {}", {e}).as_str());
        }
    }

    // initialize the config
    match backend::project::initialize_config(name, base_package, &ap) {
        Ok(_) => (),
        Err(e) => {
            print_err(format!("Failed to run initialize config: {}", {e}).as_str());
        }
    }

    // get our project context
    let p_ctx = match backend::context::get_project_context() {
        Err(_) => {
            print_general("Failed to get project context");
            return;
        }
        Ok(x) => x,
    };

    // initialize our source tree
    match backend::project::initialize_source_tree(&p_ctx) {
        Ok(_) => (),
        Err(e) => {
            print_err(format!("Failed to initialize source tree: {}", {e}).as_str());
        }
    }
    print_general("Project created: Edit espresso.toml to check it out!");
}
