use std::{error, io, result};

use crate::backend::context::{AbsoltuePaths, ProjectContext};
use crate::backend::dependency::resolve::Package;
use crate::backend::toolchain::{self, compile_project, run_jar, ToolchainContext};
use crate::backend::{self, context};
use crate::frontend::terminal::{print_err, print_sameline};

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
    print_general("-- RUNNING ARTIFACT ----");
    match run_jar(&p_ctx, &tc_ctx) {
        Ok(_) => (),
        Err(e) => print_err(format!("Failed to run 'artifact.jar': {}", { e }).as_str()),
    }
    print_general("------------------------");
    Ok(())
}

/**
 * Service function for the `build` command
 */
pub fn build(
    p_ctx: ProjectContext,
    tc_ctx: ToolchainContext,
) -> result::Result<(ProjectContext, ToolchainContext), Box<dyn error::Error>> {
    // extract dependencies
    print_general("-- EXTRACTING DEPENDENCIES");
    for (name, dep) in p_ctx.state_lock_file.dependencies.iter() {
        print_general(format!("Extracting '{}'", name).as_str());
        backend::dependency::uberjar::extract(&p_ctx, &tc_ctx, dep)?;
    }
    print_general("------------------------");

    // merge the dependencies
    print_general("-- MERGING DEPENDENCIES");
    for (name, dep) in p_ctx.state_lock_file.dependencies.iter() {
        print_general(format!("Merging '{}' into classpath", name).as_str());
        backend::dependency::uberjar::copy_classes(&p_ctx, dep)?;
    }
    print_general("------------------------");

    // walk our src directory, find source files
    print_general("-- DISCOVERING");
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
    print_general("-- COMPILING");
    compile_project(java_files, &p_ctx, &tc_ctx);
    print_general("------------------------");

    // build our jar
    print_general("-- PACKAGING");
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
    mut p_ctx: ProjectContext,
    tc_ctx: ToolchainContext,
    q: String,
) -> result::Result<(ProjectContext, ToolchainContext), Box<dyn error::Error>> {
    // get absolute paths
    let ap: AbsoltuePaths = match context::get_absolute_paths(&context::get_debug_mode()) {
        Err(_) => {
            print_err("Failed to get absolute paths");
            panic!();
        }
        Ok(x) => x,
    };

    let packages = backend::dependency::resolve::query(&q).await?;

    // collect the package selection if there was more than one returned package
    let selected_package: &Package;
    if packages.len() == 1 {
        selected_package = packages.get(0).expect("At least one package was simultaneously returned & not returned. Schrodinger's package..?");
    } else if packages.len() > 1 {
        print_general(format!("Searching for '{}'", &q).as_str());
        for (elem, package) in packages.iter().enumerate() {
            print_general(
                format!("{}) {}:{}", elem + 1, package.group_id, package.artifact_id).as_str(),
            );
        }

        // get the selected package as a string
        let mut package_number_selection = String::new();
        print_sameline(format!("Select a package (1-{}): ", packages.len()).as_str());
        if let Err(_) = io::stdin().read_line(&mut package_number_selection) {
            print_err("Failed to read user package selection")
        }

        // remove any newlines
        package_number_selection = package_number_selection.replace("\n", "");

        // convert the input into a u64
        let package_number_selection_int: u64 = match package_number_selection.parse::<u64>() {
            Ok(v) => v,
            Err(e) => {
                print_err(
                    format!(
                        "Failed to parse user input as an unsigned integer: Input was '{}'",
                        package_number_selection
                    )
                    .as_str(),
                );
                panic!("{}", e);
            }
        };

        // set the selected package to its corresponding Package struct
        selected_package = match packages.get((package_number_selection_int - 1) as usize) {
            Some(v) => v,
            None => {
                print_err("Out of range");
                panic!("Out of range");
            }
        };
    } else {
        print_err("There were no packages matching that search term");
        panic!()
    }

    // add the package
    print_general(format!("Adding '{}'", selected_package.artifact_id).as_str());
    match backend::dependency::add(&mut p_ctx, &ap, selected_package).await {
        Ok(()) => {}
        Err(e) => {
            print_err(format!("Failed to add package: {}", e).as_str());
            panic!()
        }
    }

    print_general("Package added");

    // pass ownership back
    Ok((p_ctx, tc_ctx))
}
