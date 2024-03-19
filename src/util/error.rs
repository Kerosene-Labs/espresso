use std::{error::Error, fmt};

/// General error for use in Espresso projects
#[derive(Debug)]
pub struct EspressoError {
    msg: String
}

impl EspressoError {
    /// Creates a error
    /// 
    /// # Arguments
    /// * `msg`: The error message
    pub fn new(msg: &str) -> EspressoError {
        EspressoError {
            msg: msg.to_string()
        }
    }

    /// Creates a new-in-box error
    /// 
    /// # Arguments
    /// * `msg`: The error message
    pub fn nib(msg: &str) -> Box<EspressoError> {
        Box::new(Self::new(msg))
    }
}

impl fmt::Display for EspressoError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}", self.msg)
    }
}

impl Error for EspressoError {}