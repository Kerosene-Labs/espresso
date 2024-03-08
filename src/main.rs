use clap::Command;
mod frontend;
mod backend;

fn main() {
    let cmd = Command::new("Espresso")
        .bin_name("espresso")
        .version("1.0.0")
        .about("Build Java apps without the fuss of antiquated build tools. Drink some Espresso.")
        .subcommand_required(true)
        .subcommand(frontend::command::get_build_cmd())
        .subcommand(frontend::command::get_init_cmd());
    
    let matches = cmd.get_matches();
    
    match matches.subcommand_name() {
        Some("build") => {
            frontend::service::build();
        }
        _ => {
            println!("Unknown subcommand")
        }
    }
}
