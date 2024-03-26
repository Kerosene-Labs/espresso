use std::{
    io::{self, Write},
    process::exit,
};

use colored::*;

/**
 * Print an error to stderr, exiting with code 1. This function will end the program.
 * This differs from panics as we want a nice, readable output for the user.
 */
pub fn print_err(msg: &str) {
    let lines = msg.split("\n");
    for line in lines {
        if !line.is_empty() {
            eprintln!("{} {}", "error: ".red().bold(), line.bright_red());
        }
    }
    exit(1);
}

pub fn print_general(msg: &str) {
    println!("{} {}", "info: ".bright_white().bold(), msg.white())
}

pub fn print_debug(msg: &str) {
    println!("{} {}", "debug: ".green().bold(), msg.white())
}

pub fn print_sameline(msg: &str) {
    print!("{} {}", "input: ".bold().yellow(), msg.white());
    io::stdout().flush().unwrap();
}
