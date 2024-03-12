use backend::project::ensure_debug_directory_exists_if_debug;
use clap::Command;
use frontend::terminal::print_err;
mod backend;
mod frontend;
mod util;

fn main() {
    let cmd = Command::new("Espresso")
        .bin_name("espresso")
        .version("1.0.0")
        .about("Build Java apps without the fuss of antiquated build tools. Drink some Espresso.")
        .subcommand_required(true)
        .subcommand((&*frontend::command::BUILD_CMD).clone())
        .subcommand((&*frontend::command::INIT_CMD).clone())
        .subcommand((&*frontend::command::RUN_CMD).clone());

    let matches = cmd.get_matches();

    // ensure the espresso_debug directory exists if ESPRESSO_DEBUG=1
    ensure_debug_directory_exists_if_debug();

    match matches.subcommand_name() {
        Some("build") => {
            frontend::service::build(None, None);
        }
        Some("init") => {
            frontend::service::init();
        }
        Some("run") => {
            frontend::service::run(None, None);
        }
        _ => print_err("Unknown subcommand"),
    }
}
