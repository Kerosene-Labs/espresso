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
    eprintln!("{} {}", "[X]".red().bold(), msg.white());
    exit(1);
}

pub fn print_general(msg: &str) {
    println!("{} {}", "[.]".bright_white().bold(), msg.white())
}

pub fn print_debug(msg: &str) {
    println!("{} {}", "[D]".green().bold(), msg.white())
}

pub fn print_sameline(msg: &str) {
    print!("{} {}", "[?]".bold().yellow(), msg.white());
    io::stdout().flush().unwrap();
}
