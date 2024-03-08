use crate::backend::project;

/**
 * Service function for the `build` command
 */
pub fn build() {
    project::load_project()
}