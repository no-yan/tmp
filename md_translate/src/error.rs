use thiserror::Error;

#[derive(Error, Debug)]
pub enum MdTranslateError {
    #[error("Ollama API error: {0}")]
    OllamaError(String),

    #[error("Cache error: {0}")]
    CacheError(String),

    #[error("Markdown parsing error: {0}")]
    MarkdownError(String),

    #[error("IO error: {0}")]
    IoError(#[from] std::io::Error),

    #[error("Network error: {0}")]
    NetworkError(#[from] reqwest::Error),

    #[error("Serialization error: {0}")]
    SerdeError(#[from] serde_json::Error),
}

pub type Result<T> = std::result::Result<T, MdTranslateError>;
