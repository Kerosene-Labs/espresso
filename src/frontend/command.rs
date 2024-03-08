use clap::{Command, App};

/*
 * Get the build command
 */
pub fn get_build_cmd() -> App<'static> {
    Command::new("build")
    .about("Build your Java project into a standalone .jar")
    .alias("b")
}

pub fn get_init_cmd() -> App<'static> {
    Command::new("init")
    .about("Initialize a new Espresso project")
    .alias("i")
}