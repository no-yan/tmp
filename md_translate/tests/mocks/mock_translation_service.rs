use md_translate::error::{MdTranslateError, Result};
use md_translate::translator::TranslationService;
use async_trait::async_trait;
use std::collections::VecDeque;
use std::sync::{Arc, Mutex};
use std::time::Duration;

/// Mock translation service for testing
/// Returns predefined responses or generates deterministic output
#[derive(Clone)]
pub struct MockTranslationService {
    /// Model name to report
    pub model_name: String,
    /// Delay to simulate network latency
    pub delay_ms: u64,
    /// Predefined responses (FIFO queue)
    /// If empty, generates default response
    pub responses: Arc<Mutex<VecDeque<Result<String>>>>,
}

impl MockTranslationService {
    pub fn new() -> Self {
        Self {
            model_name: "mock-model".to_string(),
            delay_ms: 0,
            responses: Arc::new(Mutex::new(VecDeque::new())),
        }
    }

    /// Create mock with specific responses
    pub fn with_responses(responses: Vec<Result<String>>) -> Self {
        Self {
            model_name: "mock-model".to_string(),
            delay_ms: 0,
            responses: Arc::new(Mutex::new(responses.into())),
        }
    }

    /// Create mock that simulates network delay
    pub fn with_delay(delay_ms: u64) -> Self {
        Self {
            model_name: "mock-model".to_string(),
            delay_ms,
            responses: Arc::new(Mutex::new(VecDeque::new())),
        }
    }
}

impl Default for MockTranslationService {
    fn default() -> Self {
        Self::new()
    }
}

#[async_trait]
impl TranslationService for MockTranslationService {
    async fn translate(&self, text: &str) -> Result<String> {
        // Simulate network delay
        if self.delay_ms > 0 {
            tokio::time::sleep(Duration::from_millis(self.delay_ms)).await;
        }

        // Check for predefined response
        let mut responses = self.responses.lock().unwrap();
        if let Some(response) = responses.pop_front() {
            return response;
        }

        // Default: wrap text in marker for easy identification
        Ok(format!("[TRANSLATED: {}]", text))
    }

    fn model(&self) -> &str {
        &self.model_name
    }
}