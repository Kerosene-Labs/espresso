use clap::Command;
use once_cell::sync::Lazy;

pub static BUILD_CMD: Lazy<Command> = Lazy::new(|| {
    Command::new("build")
        .about("Build your Java project into a standalone .jar")
        .alias("b")
});

pub static INIT_CMD: Lazy<Command> = Lazy::new(|| {
    Command::new("init")
        .about("Initialize a new Espresso project")
        .alias("i")
});

pub static RUN_CMD: Lazy<Command> = Lazy::new(|| {
    Command::new("run")
        .about("Build & run your Java project")
        .alias("r")
});
