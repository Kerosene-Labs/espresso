use crate::backend::context::{AbsoltuePaths, ProjectContext};
use crate::backend::toolchain::{compile_project, run_jar, ToolchainContext};
use crate::backend::{self, context};
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
        Err(e) => print_err(format!("Failed to run 'artifact.jar': {}", { e }).as_str()),
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
            print_err(format!("Failed to build jar: {}", { e }).as_str());
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

    // initialize the project
    let debug_mode = backend::context::get_debug_mode();
    match backend::project::initialize(name, base_package, &ap, &debug_mode) {
        Ok(_) => (),
        Err(e) => {
            print_err(format!("Unable to initialize project: {}", e).as_str());
        }
    }
    print_general("Project created: Edit espresso.toml to check it out!");
}

/// Service function for the `add` command.
pub async fn add(
    p_ctx: ProjectContext,
    tc_ctx: ToolchainContext,
    q: String
) -> result::Result<(ProjectContext, ToolchainContext), Box<dyn error::Error>> {
    print_general(format!("Searching for '{}'", q).as_str());
    let packages = backend::dependency::resolve::query(q).await?;    
    for (elem, package) in packages.iter().enumerate() {
        print_general(format!("{}) G:{} | A:{}", elem + 1, package.group_id, package.artifact_id).as_str());
    }

    // collect the package selection
    let mut package_selection = String::new();
    print_sameline(format!("Select a package (1-{}): ", packages.len()).as_str());
    if let Err(_) = io::stdin().read_line(&mut package_selection) {
        print_err("Failed to read user package selection")
    }

    // pass ownership back
    Ok((p_ctx, tc_ctx))
}
