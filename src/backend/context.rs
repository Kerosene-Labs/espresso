use crate::backend::project::Config;

use super::project::get_config_from_fs;

/**
 * Represents the currently loaded project
 */
pub struct ProjectContext {
    pub config: Config,
}

/**
 * Get the Project Context
 */
pub fn get_project_context() -> ProjectContext {
    ProjectContext {
        config: get_config_from_fs()
    }
}