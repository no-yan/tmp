use crate::error::{MdTranslateError, Result};
use crate::translator::translation_service::TranslationService;
use async_trait::async_trait;
use serde::{Deserialize, Serialize};
use std::time::Duration;

#[derive(Serialize)]
struct OllamaRequest {
    model: String,
    prompt: String,
    stream: bool,
    options: OllamaOptions,
}

#[derive(Serialize)]
struct OllamaOptions {
    temperature: f32,
}

#[derive(Deserialize)]
struct OllamaResponse {
    response: String,
}

#[derive(Clone)]
pub struct OllamaClient {
    base_url: String,
    pub model: String,
    client: reqwest::Client,
    timeout: Duration,
}

impl OllamaClient {
    pub fn new(base_url: String, model: String) -> Self {
        Self {
            base_url,
            model,
            client: reqwest::Client::new(),
            timeout: Duration::from_secs(30),
        }
    }

    pub async fn translate(&self, text: &str) -> Result<String> {
        let prompt = format!(
            r#"You are a professional translator. Translate ONLY the text content from English to Japanese.

CRITICAL RULES - Follow these strictly:
- Do NOT add any markdown syntax (no ```, ---, #, *, _, etc.)
- Do NOT add extra paragraphs or line breaks beyond what exists in the original
- Do NOT add code blocks, horizontal rules, or heading markers
- Do NOT add any formatting that wasn't in the original text
- Translate ONLY the plain text content, nothing else
- Keep the exact same structure as the input

INPUT TEXT:
{}

OUTPUT (translated text only, no explanations):"#,
            text
        );

        let request = OllamaRequest {
            model: self.model.clone(),
            prompt,
            stream: false,
            options: OllamaOptions { temperature: 0.3 },
        };

        // Retry logic: 3 attempts with exponential backoff
        for attempt in 1..=3 {
            match self.try_translate(&request).await {
                Ok(response) => return Ok(response),
                Err(_e) if attempt < 3 => {
                    let backoff = Duration::from_secs(2_u64.pow(attempt - 1));
                    tokio::time::sleep(backoff).await;
                }
                Err(e) => return Err(e),
            }
        }

        unreachable!()
    }

    async fn try_translate(&self, request: &OllamaRequest) -> Result<String> {
        let url = format!("{}/api/generate", self.base_url);

        let response = self
            .client
            .post(&url)
            .timeout(self.timeout)
            .json(request)
            .send()
            .await
            .map_err(|e| MdTranslateError::OllamaError(format!("Failed to connect: {}", e)))?;

        if !response.status().is_success() {
            return Err(MdTranslateError::OllamaError(format!(
                "API returned status {}",
                response.status()
            )));
        }

        let ollama_response: OllamaResponse = response.json().await?;
        Ok(ollama_response.response.trim().to_string())
    }
}

#[async_trait]
impl TranslationService for OllamaClient {
    async fn translate(&self, text: &str) -> Result<String> {
        // Delegate to existing implementation
        self.translate(text).await
    }

    fn model(&self) -> &str {
        &self.model
    }
}
