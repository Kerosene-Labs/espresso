
use clap::Command;
mod frontend;
mod backend;

fn main() {
    let cmd = Command::new("Espresso")
        .bin_name("espresso")
        .version("1.0.0")
        .about("Build Java apps without the fuss of antiquated build tools. Drink some Espresso.")
        .subcommand_required(true)
        .subcommand((&*frontend::command::BUILD_CMD).clone())
        .subcommand((&*frontend::command::INIT_CMD).clone());
    
    let matches = cmd.get_matches();
    
    match matches.subcommand_name() {
        Some("build") => {
            frontend::service::build();
        }
        Some("init") => {
            frontend::service::init();
        }
        _ => {
            println!("Unknown subcommand")
        }
    }
}
