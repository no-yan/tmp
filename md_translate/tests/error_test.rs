use md_translate::error::MdTranslateError;
use std::io;

#[test]
fn test_io_error_conversion() {
    let io_err = io::Error::new(io::ErrorKind::NotFound, "file not found");
    let md_err: MdTranslateError = io_err.into();
    assert!(matches!(md_err, MdTranslateError::IoError(_)));
}

#[test]
fn test_ollama_error_display() {
    let err = MdTranslateError::OllamaError("timeout".to_string());
    assert_eq!(err.to_string(), "Ollama API error: timeout");
}

#[test]
fn test_serde_error_conversion() {
    let json_err = serde_json::from_str::<i32>("invalid").unwrap_err();
    let md_err: MdTranslateError = json_err.into();
    assert!(matches!(md_err, MdTranslateError::SerdeError(_)));
}

#[test]
fn test_cache_error_display() {
    let err = MdTranslateError::CacheError("invalid checksum".to_string());
    assert!(err.to_string().contains("Cache error"));
}

#[test]
fn test_markdown_error_display() {
    let err = MdTranslateError::MarkdownError("parse failure".to_string());
    assert!(err.to_string().contains("Markdown parsing error"));
}

#[test]
fn test_network_error_display() {
    let err = MdTranslateError::OllamaError("Network timeout".to_string());
    let display = format!("{}", err);
    assert!(display.contains("Ollama API error"));
    assert!(display.contains("Network timeout"));
}