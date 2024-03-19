use backend::{
    context::{get_project_context, ProjectContext},
    toolchain::{get_toolchain_context, ToolchainContext},
};
use clap::Command;
use frontend::terminal::print_err;
mod backend;
mod frontend;
mod util;

/// Get runtime contexts required for command service functions
/// 
/// # Returns
/// A tuple containing a `ProjectContext` and a `ToolchainContext`
fn get_contexts() -> (ProjectContext, ToolchainContext) {
    let p_ctx = match get_project_context() {
        Ok(v) => v,
        Err(e) => {
            print_err(format!("Failed to get project context! \nERROR MESSAGE: \n{}", e).as_str());
            panic!("{}", e);
        }
    };
    let tc_ctx = get_toolchain_context(&p_ctx);
    (p_ctx, tc_ctx)
}

#[tokio::main]
async fn main() {
    let cmd = Command::new("Espresso")
        .bin_name("espresso")
        .version("1.0.0")
        .about("Build Java apps without the fuss of antiquated build tools. Drink some Espresso.")
        .subcommand_required(true)
        .subcommand((frontend::command::BUILD_CMD).clone())
        .subcommand((frontend::command::INIT_CMD).clone())
        .subcommand((frontend::command::RUN_CMD).clone())
        .subcommand((frontend::command::ADD_CMD).clone());

    let matches = cmd.get_matches();

    // ensure the espresso_debug directory exists if ESPRESSO_DEBUG=1
    // ensure_debug_directory_exists_if_debug();

    match matches.subcommand_name() {
        Some("build") => {
            let (p_ctx, tc_ctx) = get_contexts();
            match frontend::service::build(p_ctx, tc_ctx) {
                Ok(_) => (),
                Err(e) => print_err(format!("Error occurred running command: {}", e).as_str()),
            }
        }
        Some("init") => {
            frontend::service::init();
        }
        Some("run") => {
            let (p_ctx, tc_ctx) = get_contexts();
            match frontend::service::run(p_ctx, tc_ctx) {
                Ok(_) => (),
                Err(e) => print_err(format!("Error occurred running command: {}", e).as_str()),
            }
        }
        Some("add") => {
            let (p_ctx, tc_ctx) = get_contexts();
            match frontend::service::add(p_ctx, tc_ctx).await {
                Ok(_) => (),
                Err(e) => print_err(format!("Error occurred running command: {}", e).as_str()),
            }
        }
        _ => print_err("Unknown subcommand"),
    }
}
