use std::{borrow::Cow, env, fmt::format, fs, io, path, process::{Command, ExitStatus}, vec};
use walkdir::WalkDir;
use crate::{frontend::terminal::{print_debug, print_err, print_general}, util};

use super::{context::ProjectContext, project::{self, get_absolute_project_path, get_full_base_package_path, Project}};

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
    let base_package = project::get_full_base_package_path(&p_ctx);
    let files = util::directory::read_files_recursively(base_package);

    // begin sorting out java files
    let mut java_files: Vec<String> = vec![];
    for x in files.unwrap() {
        if x.ends_with(".java"){
            java_files.push(x);
        }
    }
    return Ok(java_files);
}

/**
 * Compile a Java file using the defined toolchain
 */
pub fn compile_java_file(path: &String, tc_ctx: &ToolchainContext) {
    Command::new("sh")
    .arg("-c")
    .arg(format!("{} {}", &tc_ctx.compiler_path.to_string_lossy(), path))
    .output()
    .expect("failed to execute java compiler");
}

/**
 * Compile all Java files under a project
 */
pub fn compile_project(java_files: Vec<String>, p_ctx: &ProjectContext, tc_ctx: &ToolchainContext) {
    let compiler_path: &Cow<'_, str> = &tc_ctx.compiler_path.to_string_lossy();

    for file in java_files{
        // build our compiler string
        let cmd = format!("{} {} -d {}/build",
            &compiler_path,
            file,
            get_absolute_project_path()
        );

        print_debug(format!("Running '{}'", cmd).as_str());

        // call the java compiler
        let output = Command::new("sh")
        .arg("-c")
        .arg(cmd)
        .output()
        .expect("failed to execute java compiler");

        if !output.status.success() {
            print_err("java compiler exited with error(s)");
            println!("{:?}", output.stdout);
        }
    }
}

/**
 * Get the manifest for the JAR
 */
fn get_manifest(p_ctx: &ProjectContext) -> String {
    format!("Main-Class: {}.Main\nManifest-Version: 1.0\n", p_ctx.config.project.base_package).to_string()
}

/**
 * Write the manifest file
 */
pub fn write_manifest(p_ctx: &ProjectContext) -> io::Result<()>{
    std::fs::write(get_absolute_project_path() + "/build/MANIFEST.MF", get_manifest(p_ctx))
}

/**
 * Build our JAR file
 */
pub fn build_jar(p_ctx: &ProjectContext, tc_ctx: &ToolchainContext) {
    // write our manifest
    write_manifest(p_ctx).unwrap();

    // convert our base package (ex: com.xyz.whatever) to its filesystem equivalent (com/xyz/whatever)
    let relative_base_package_path = p_ctx.config.project.base_package.clone().replace(".", "/");

    // build our packager command
    let cmd = format!("jar -c --file=artifact.jar --manifest=MANIFEST.MF {}", relative_base_package_path);

    // run the command
    let output = Command::new("sh")
    .current_dir(get_absolute_project_path() + "/build")
    .arg("-c")
    .arg(cmd)
    .output()
    .expect("failed to run jar packager");
    if !output.status.success() {
        println!("{}", String::from_utf8(output.stderr).unwrap());
        print_err("java packager (jar) exited with error(s)");
    }
}