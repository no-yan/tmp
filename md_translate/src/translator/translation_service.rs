use crate::error::Result;
use async_trait::async_trait;

/// Trait for translation services
/// Allows mocking Ollama API calls in tests
#[async_trait]
pub trait TranslationService: Clone + Send + Sync {
    /// Translate text from source language to target language
    /// Preserves markdown formatting and code blocks
    async fn translate(&self, text: &str) -> Result<String>;

    /// Get the model name being used
    fn model(&self) -> &str;
}