use crate::{
    frontend::terminal::{print_debug, print_err},
    util,
};
use std::error;
use std::{borrow::Cow, env, fs, io, path, process::Command, result, vec};

use super::context::ProjectContext;

/**
 * Represents the context of the current Java toolchain
 */
pub struct ToolchainContext {
    pub compiler_path: path::PathBuf,
    pub runtime_path: path::PathBuf,
    pub packager_path: path::PathBuf,
    pub toolchain_path: path::PathBuf,
}

/**
 * Get the toolchain path (expanded if it contains ${JAVA_HOME})
 */
pub fn get_expanded_toolchain_path<'a>(toolchain_path: &'a String) -> Cow<'a, String> {
    if toolchain_path.contains("${JAVA_HOME}") {
        match env::var("JAVA_HOME") {
            Ok(val) => return Cow::Owned(val),
            Err(_) => return Cow::Borrowed(toolchain_path),
        }
    }
    Cow::Borrowed(toolchain_path)
}

/**
 * Get the toolchain context
 */
pub fn get_toolchain_context(p_ctx: &ProjectContext) -> ToolchainContext {
    let toolchain_path = get_expanded_toolchain_path(&p_ctx.config.toolchain.path).into_owned();
    let compiler_path = path::PathBuf::from(toolchain_path.clone() + "/bin/javac");
    let runtime_path = path::PathBuf::from(toolchain_path.clone() + "/bin/java");
    let packager_path = path::PathBuf::from(toolchain_path.clone() + "/bin/jar");
    ToolchainContext {
        compiler_path,
        runtime_path,
        packager_path,
        toolchain_path: path::PathBuf::from(&toolchain_path),
    }
}

/**
 * Get a list of source files
 */
pub fn get_java_source_files(p_ctx: &ProjectContext) -> Result<Vec<String>, std::io::Error> {
    let base_package = p_ctx.dynamic_absolute_paths.base_package.clone();

    let files = util::directory::read_files_recursively(base_package);

    // begin sorting out java files
    let mut java_files: Vec<String> = vec![];
    for x in files.unwrap() {
        if x.ends_with(".java") {
            java_files.push(x);
        }
    }
    return Ok(java_files);
}

/**
 * Compile all Java files under a project
 */
pub fn compile_project(java_files: Vec<String>, p_ctx: &ProjectContext, tc_ctx: &ToolchainContext) {
    let compiler_path: &Cow<'_, str> = &tc_ctx.compiler_path.to_string_lossy();

    for file in java_files {
        // build our compiler string
        let cmd = format!(
            "{} {} -d {}/build -cp {}/build",
            &compiler_path, file, &p_ctx.absolute_paths.project, &p_ctx.absolute_paths.project
        );

        print_debug(format!("Running '{}'", cmd).as_str());

        // call the java compiler
        let output = Command::new("sh")
            .arg("-c")
            .arg(cmd)
            .output()
            .expect("failed to execute java compiler");

        if !output.status.success() {
            println!("\n\n{}\n", String::from_utf8(output.stderr).unwrap());
            print_err("java compiler exited with error(s)");
        }
    }
}

/**
 * Get the manifest for the JAR
 */
fn get_manifest(p_ctx: &ProjectContext) -> String {
    format!(
        "Main-Class: {}.Main\nManifest-Version: 1.0\n",
        p_ctx.config.project.base_package
    )
    .to_string()
}

/**
 * Write the manifest file
 */
pub fn write_manifest(p_ctx: &ProjectContext) -> io::Result<()> {
    std::fs::write(
        p_ctx.absolute_paths.project.clone() + "/build/MANIFEST.MF",
        get_manifest(p_ctx),
    )
}

/**
 * Build our JAR file
 *
 * TODO: ensure we're using the proper toolchain
 */
pub fn build_jar(p_ctx: &ProjectContext, tc_ctx: &ToolchainContext) -> result::Result<(), Box<dyn error::Error>>{
    // write our manifest
    write_manifest(p_ctx).unwrap();

    // convert our base package (ex: com.xyz.whatever) to its filesystem equivalent (com/xyz/whatever)
    let relative_base_package_path = p_ctx.config.project.base_package.clone().replace(".", "/");

    // remove the old jar
    let _ = fs::remove_file(p_ctx.absolute_paths.project.clone() + "/build/artifact.jar");

    // build our packager command
    let cmd = format!(
        "{} -c --file=artifact.jar --manifest=MANIFEST.MF {}",
        tc_ctx.packager_path.to_string_lossy().clone(),
        relative_base_package_path
    );
    print_debug(format!("Running '{}'", cmd).as_str());

    // run the command
    let output = Command::new("sh")
        .current_dir(p_ctx.absolute_paths.project.clone() + "/build")
        .arg("-c")
        .arg(cmd)
        .output()?;
    if !output.status.success() {
        let process_err = String::from_utf8(output.stderr)?;
        Err(format!("Failed to package jar: Packager output was: {}", process_err))?
    }

    Ok(())
    
}

/**
 * Run our JAR file
 */
pub fn run_jar(p_ctx: &ProjectContext, tc_ctx: &ToolchainContext) -> result::Result<(), Box<dyn error::Error>>{
    // build our packager command
    let cmd = format!(
        "{} -jar {}",
        tc_ctx.runtime_path.to_string_lossy().clone(),
        p_ctx.absolute_paths.project.clone() + "/build/artifact.jar"
    );

    // run the command
    let status = Command::new("sh")
        .current_dir(p_ctx.absolute_paths.project.clone() + "/build")
        .arg("-c")
        .arg(cmd)
        .status()?;

    // dela with our output. We wa
    if !status.success() {
        Err(format!("Java process exited with non-0 status code"))?
    }

    Ok(())
}
