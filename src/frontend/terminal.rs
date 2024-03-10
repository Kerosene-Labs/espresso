use colored::*;

pub fn print_err(msg: &str) {
    eprintln!("{} {}", "[ERROR]".red().bold(), msg.white())
}

pub fn print_general(msg: &str) {
    println!("{} {}", "[ . ]".bright_white().bold(), msg.white())
}