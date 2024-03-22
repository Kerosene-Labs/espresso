use backend::{
    context::{get_project_context, AbsoltuePaths, ProjectContext},
    toolchain::{get_toolchain_context, ToolchainContext},
};
use clap::{Parser, Subcommand, command};
use frontend::terminal::print_err;
mod backend;
mod frontend;
mod util;

#[derive(Parser, Debug)]
#[command(name="espresso")]
#[command(about="Build Java apps without the fuss of antiquated build tools. Drink some Espresso.")]
struct EspressoCli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand, Debug)]
enum Commands {
    #[command(about="Build your project")]
    Build {},

    #[command(about="Initialize a new project")]
    Init {},

    #[command(about="Run your project")]
    Run {},

    #[command(arg_required_else_help = true)]
    #[command(about="Add a package to your project from the Espresso Registry")]
    Add {
        #[arg(required = true)]
        search_term: String
    },

    #[command(arg_required_else_help = true)]
    #[command(about="Add a package to your project from the filesystem")]
    AddFs {
        #[arg(required = true, long)]
        path: String,

        #[arg(required = true, long)]
        name: String,
    }
}


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
    let args = EspressoCli::parse();
    
    match args.command {
        Commands::Build {  } => {
            let (p_ctx, tc_ctx) = get_contexts();
            match frontend::service::build(p_ctx, tc_ctx) {
                Ok(_) => (),
                Err(e) => print_err(format!("Error occurred running command: {}", e).as_str()),
            }
        },
        Commands::Init { } => {
            frontend::service::init()
        },
        Commands::Run {  } => {
            let (p_ctx, tc_ctx) = get_contexts();
            match frontend::service::run(p_ctx, tc_ctx) {
                Ok(_) => (),
                Err(e) => print_err(format!("Error occurred running command: {}", e).as_str()),
            }
        }
        Commands::Add { search_term } => {
            let (p_ctx, tc_ctx) = get_contexts();
            match frontend::service::add(p_ctx, tc_ctx, search_term).await {
                Ok(_) => (),
                Err(e) => print_err(format!("Error occurred running command: {}", e).as_str()),
            }
        },
        Commands::AddFs { path, name } => {
            todo!()
        }
    }
}
