mod cache;
mod cache_backend;
mod ollama;
mod translation_service;

pub use cache::{CacheStats, TranslationCache};
pub use cache_backend::CacheBackend;
pub use ollama::OllamaClient;
pub use translation_service::TranslationService;

use crate::error::Result;
use crate::markdown::parser::{MarkdownProcessor, Segment};
use indicatif::{ProgressBar, ProgressStyle};
use tokio::task::JoinSet;

pub struct TranslatorConfig {
    pub source_lang: String,
    pub target_lang: String,
    pub use_cache: bool,
    pub parallel_requests: usize,
    pub show_progress: bool,
}

impl Default for TranslatorConfig {
    fn default() -> Self {
        Self {
            source_lang: "en".to_string(),
            target_lang: "ja".to_string(),
            use_cache: true,
            parallel_requests: 3,
            show_progress: true,
        }
    }
}

pub struct Translator<T = OllamaClient, C = TranslationCache>
where
    T: TranslationService + 'static,
    C: CacheBackend + 'static,
{
    client: T,
    cache: C,
    config: TranslatorConfig,
}

impl Translator<OllamaClient, TranslationCache> {
    pub fn new(ollama_url: String, model: String, config: TranslatorConfig) -> Result<Self> {
        Ok(Self {
            client: OllamaClient::new(ollama_url, model),
            cache: TranslationCache::new()?,
            config,
        })
    }
}

impl<T, C> Translator<T, C>
where
    T: TranslationService + 'static,
    C: CacheBackend + 'static,
{
    /// Create translator with custom implementations (for testing)
    pub fn new_with_deps(client: T, cache: C, config: TranslatorConfig) -> Self {
        Self {
            client,
            cache,
            config,
        }
    }

    pub async fn translate_markdown(&mut self, markdown: &str) -> Result<String> {
        let segments = MarkdownProcessor::parse(markdown)?;
        let translated_segments = self.translate_segments(segments).await?;
        Ok(Self::reconstruct_markdown(&translated_segments))
    }

    pub async fn translate_to_html(&mut self, markdown: &str, title: &str) -> Result<String> {
        let segments = MarkdownProcessor::parse(markdown)?;
        let translated_segments = self.translate_segments(segments).await?;
        Ok(crate::markdown::HtmlRenderer::render(
            &translated_segments,
            title,
        ))
    }

    async fn translate_segments(&mut self, segments: Vec<Segment>) -> Result<Vec<Segment>> {
        let total = segments.iter().filter(|s| s.is_translatable()).count();

        let progress = if self.config.show_progress {
            let pb = ProgressBar::new(total as u64);
            pb.set_style(
                ProgressStyle::default_bar()
                    .template("[{elapsed_precise}] {bar:40.cyan/blue} {pos}/{len} {msg}")
                    .unwrap()
                    .progress_chars("=>-"),
            );
            Some(pb)
        } else {
            None
        };

        // Pre-allocate result vector with None values to preserve ordering
        let segment_count = segments.len();
        let mut translated: Vec<Option<Segment>> = vec![None; segment_count];

        let mut join_set = JoinSet::new();
        let semaphore =
            std::sync::Arc::new(tokio::sync::Semaphore::new(self.config.parallel_requests));

        for (index, segment) in segments.into_iter().enumerate() {
            if !segment.is_translatable() {
                // Store non-translatable segment at correct index immediately
                translated[index] = Some(segment);
                continue;
            }

            if let Some(text) = segment.get_text() {
                // Check cache first
                let lang_pair = format!("{}-{}", self.config.source_lang, self.config.target_lang);

                if self.config.use_cache {
                    if let Some(cached) =
                        self.cache.get(text, self.client.model(), &lang_pair)
                    {
                        let mut new_segment = segment.clone();
                        new_segment.set_text(cached);
                        translated[index] = Some(new_segment); // Store at correct index
                        if let Some(pb) = &progress {
                            pb.inc(1);
                        }
                        continue;
                    }
                }

                // Translate via client (with semaphore for parallelism control)
                let text_owned = text.to_string();
                let segment_clone = segment.clone();
                let client_clone = self.client.clone();
                let permit = semaphore.clone().acquire_owned().await.unwrap();

                // Spawn async task with index to preserve ordering
                join_set.spawn(async move {
                    let result = client_clone.translate(&text_owned).await;
                    drop(permit);
                    (index, segment_clone, text_owned, result) // Include index in return
                });
            }
        }

        // Collect results and store at original indices
        while let Some(result) = join_set.join_next().await {
            let (index, mut segment, original_text, translation_result) = result.unwrap();

            match translation_result {
                Ok(translated_text) => {
                    // Save to cache
                    if self.config.use_cache {
                        let lang_pair =
                            format!("{}-{}", self.config.source_lang, self.config.target_lang);
                        let _ = self.cache.set(
                            &original_text,
                            &translated_text,
                            self.client.model(),
                            &lang_pair,
                        );
                    }

                    segment.set_text(translated_text);
                    translated[index] = Some(segment); // Store at original index
                }
                Err(_e) => {
                    // On error, keep original text
                    translated[index] = Some(segment); // Store at original index
                }
            }

            if let Some(pb) = &progress {
                pb.inc(1);
            }
        }

        if let Some(pb) = progress {
            pb.finish_with_message("Translation complete");
        }

        // Unwrap all Options (all should be Some at this point)
        Ok(translated
            .into_iter()
            .map(|opt| opt.expect("All segments should be processed"))
            .collect())
    }

    fn reconstruct_markdown(segments: &[Segment]) -> String {
        let mut output = String::new();

        for segment in segments {
            match segment {
                Segment::Heading { level, text } => {
                    output.push_str(&"#".repeat(*level as usize));
                    output.push(' ');
                    output.push_str(text);
                    output.push_str("\n\n");
                }
                Segment::Paragraph { text } => {
                    output.push_str(text);
                    output.push_str("\n\n");
                }
                Segment::CodeBlock { language, code } => {
                    output.push_str("```");
                    if let Some(lang) = language {
                        output.push_str(lang);
                    }
                    output.push('\n');
                    output.push_str(code);
                    output.push_str("\n```\n\n");
                }
                Segment::HorizontalRule => {
                    output.push_str("---\n\n");
                }
                _ => {}
            }
        }

        output
    }

    pub fn cache_stats(&mut self) -> CacheStats {
        self.cache.stats()
    }

    pub fn clear_cache(&self) -> Result<()> {
        self.cache.clear()
    }
}
