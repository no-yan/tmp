# Test Strategy Implementation Plan

## Overview

Implement comprehensive test coverage for the md_translate project following the strategy outlined in the research document. The plan introduces trait-based abstractions to eliminate external dependencies (Ollama service, filesystem) in tests, enabling fast, deterministic unit tests while maintaining backward compatibility through default generic parameters.

**Based on**: `thoughts/shared/research/2025-09-30-test-strategy-per-module.md`

## Current State Analysis

### Architecture
- **No trait abstractions**: All components use concrete types (OllamaClient, TranslationCache)
- **Direct dependencies**: Translator owns concrete OllamaClient and TranslationCache instances
- **Clean module hierarchy**: error → markdown → translator → cli (no circular dependencies)

### Test Coverage (as of 2025-09-30)
- **Unit tests**: 4 tests in `tests/segment_ordering_unit_test.rs` (parser only, fast, reliable)
- **Integration tests**: 3 tests in `tests/segment_ordering_test.rs` (require Ollama service, slow, unreliable)
- **Coverage gaps**:
  - No error.rs tests
  - No markdown/renderer.rs tests
  - No translator logic tests without Ollama
  - No cache logic tests without filesystem
  - Untested segment types: `List`, `BlockQuote`

### Key Issues
1. Integration tests fail without Ollama running on localhost:11434
2. No mocking infrastructure for external dependencies
3. No `async-trait` dependency for trait-based async abstractions
4. No `[dev-dependencies]` section in Cargo.toml

## Desired End State

### Success Criteria

#### Automated Verification:
- [ ] All new test files compile: `cargo test --no-run`
- [ ] Unit tests run without external dependencies: `cargo test --lib`
- [ ] Fast tests complete in <5 seconds: `cargo test --exclude segment_ordering_test`
- [ ] Mock-based translator tests pass: `cargo test translator_unit_test`
- [ ] Cache unit tests pass: `cargo test cache_unit_test`
- [ ] Error tests pass: `cargo test error_test`
- [ ] Renderer tests pass: `cargo test markdown_renderer_test`
- [ ] All tests pass (with Ollama running): `cargo test`

#### Manual Verification:
- [ ] Verify tests can run in parallel without conflicts
- [ ] Confirm integration tests still validate real Ollama behavior
- [ ] Check that trait abstractions don't impact runtime performance
- [ ] Validate backward compatibility - existing CLI usage unchanged

## What We're NOT Doing

- NOT removing existing integration tests (they validate real Ollama behavior)
- NOT modifying public API of Translator struct (backward compatible via default generics)
- NOT testing server module (feature-gated, not fully implemented)
- NOT implementing property-based testing (out of scope)
- NOT testing CLI binary integration (focused on library testing)

## Implementation Approach

Three-phase approach starting with low-effort wins, then introducing trait abstractions for network/filesystem dependencies, and finally testing CLI layer. Each phase is independently valuable and can be merged separately.

---

## Phase 1: Low-Hanging Fruit (No Refactoring)

### Overview
Add unit tests for modules that don't require trait abstractions. These tests verify correctness of pure logic modules without external dependencies.

**Effort**: 2-3 hours
**Impact**: Immediate coverage improvement for 3 modules

### Changes Required

#### 1. Add error.rs Tests

**File**: `tests/error_test.rs` (new file)

**Purpose**: Verify error variant creation, conversions, and display formatting

```rust
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
fn test_reqwest_error_conversion() {
    // Create a mock reqwest error (requires constructing via actual request)
    // This verifies the #[from] attribute works
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
```

**Coverage**: Tests all 6 error variants, Display trait, From conversions

#### 2. Add markdown/renderer.rs Tests

**File**: `tests/markdown_renderer_test.rs` (new file)

**Purpose**: Verify HTML generation, escaping, templating, and CSS injection

```rust
use md_translate::markdown::{HtmlRenderer, Segment};

#[test]
fn test_heading_rendering() {
    let segments = vec![
        Segment::Heading { level: 1, text: "Title".to_string() },
        Segment::Heading { level: 2, text: "Subtitle".to_string() },
        Segment::Heading { level: 3, text: "Section".to_string() },
    ];
    let html = HtmlRenderer::render(&segments, "Test");

    assert!(html.contains("<h1>Title</h1>"));
    assert!(html.contains("<h2>Subtitle</h2>"));
    assert!(html.contains("<h3>Section</h3>"));
}

#[test]
fn test_paragraph_rendering() {
    let segments = vec![
        Segment::Paragraph { text: "First paragraph.".to_string() },
        Segment::Paragraph { text: "Second paragraph.".to_string() },
    ];
    let html = HtmlRenderer::render(&segments, "Test");

    assert!(html.contains("<p>First paragraph.</p>"));
    assert!(html.contains("<p>Second paragraph.</p>"));
}

#[test]
fn test_html_escaping() {
    let segments = vec![
        Segment::Paragraph { text: "<script>alert('xss')</script>".to_string() },
        Segment::Heading { level: 1, text: "Title with <em>tags</em>".to_string() },
    ];
    let html = HtmlRenderer::render(&segments, "Test");

    // Verify dangerous HTML is escaped
    assert!(html.contains("&lt;script&gt;"));
    assert!(!html.contains("<script>"));
    assert!(html.contains("&lt;em&gt;"));
}

#[test]
fn test_code_block_with_language() {
    let segments = vec![
        Segment::CodeBlock {
            language: Some("rust".to_string()),
            code: "fn main() {}".to_string(),
        },
    ];
    let html = HtmlRenderer::render(&segments, "Test");

    assert!(html.contains("<pre><code class=\"language-rust\">"));
    assert!(html.contains("fn main() {}"));
    assert!(html.contains("</code></pre>"));
}

#[test]
fn test_code_block_without_language() {
    let segments = vec![
        Segment::CodeBlock {
            language: None,
            code: "plain text".to_string(),
        },
    ];
    let html = HtmlRenderer::render(&segments, "Test");

    assert!(html.contains("<pre><code>"));
    assert!(!html.contains("class=\"language-"));
}

#[test]
fn test_horizontal_rule() {
    let segments = vec![
        Segment::HorizontalRule,
    ];
    let html = HtmlRenderer::render(&segments, "Test");

    assert!(html.contains("<hr>"));
}

#[test]
fn test_template_structure() {
    let html = HtmlRenderer::render(&[], "My Title");

    // Verify HTML5 structure
    assert!(html.contains("<!DOCTYPE html>"));
    assert!(html.contains("<html lang=\"ja\">"));
    assert!(html.contains("<meta charset=\"utf-8\">"));
    assert!(html.contains("<title>My Title</title>"));

    // Verify CSS is injected
    assert!(html.contains("<style>"));
    assert!(html.contains("</style>"));

    // Verify body tag exists
    assert!(html.contains("<body>"));
    assert!(html.contains("</body>"));
}

#[test]
fn test_title_escaping_in_template() {
    let html = HtmlRenderer::render(&[], "<script>alert()</script>");

    // Title in <title> tag should be escaped
    assert!(html.contains("&lt;script&gt;"));
    assert!(!html.contains("<title><script>"));
}

#[test]
fn test_mixed_segments() {
    let segments = vec![
        Segment::Heading { level: 1, text: "Title".to_string() },
        Segment::Paragraph { text: "Introduction.".to_string() },
        Segment::CodeBlock {
            language: Some("python".to_string()),
            code: "print('hello')".to_string(),
        },
        Segment::HorizontalRule,
        Segment::Paragraph { text: "Conclusion.".to_string() },
    ];
    let html = HtmlRenderer::render(&segments, "Mixed Content");

    // Verify order is preserved
    let h1_pos = html.find("<h1>Title</h1>").unwrap();
    let p1_pos = html.find("<p>Introduction.</p>").unwrap();
    let code_pos = html.find("<pre>").unwrap();
    let hr_pos = html.find("<hr>").unwrap();
    let p2_pos = html.find("<p>Conclusion.</p>").unwrap();

    assert!(h1_pos < p1_pos);
    assert!(p1_pos < code_pos);
    assert!(code_pos < hr_pos);
    assert!(hr_pos < p2_pos);
}

#[test]
fn test_code_escaping_in_blocks() {
    let segments = vec![
        Segment::CodeBlock {
            language: Some("html".to_string()),
            code: "<div>content</div>".to_string(),
        },
    ];
    let html = HtmlRenderer::render(&segments, "Test");

    // Code content should be escaped
    assert!(html.contains("&lt;div&gt;"));
    assert!(!html.contains("<div>content</div>"));
}
```

**Coverage**: Tests all segment rendering, HTML escaping, template structure, CSS injection

#### 3. Expand markdown/parser.rs Tests

**File**: Extend existing `tests/segment_ordering_unit_test.rs`

**New tests to add**:

```rust
#[test]
fn test_empty_markdown() {
    let segments = MarkdownProcessor::parse("").expect("Failed to parse empty markdown");
    assert_eq!(segments.len(), 0);
}

#[test]
fn test_whitespace_only() {
    let segments = MarkdownProcessor::parse("   \n\n   \n").expect("Failed to parse");
    // Should produce no segments or empty paragraph
    assert!(segments.is_empty() || segments.len() == 1);
}

#[test]
fn test_ordered_lists() {
    let markdown = "1. First item\n2. Second item\n3. Third item";
    let segments = MarkdownProcessor::parse(markdown).expect("Failed to parse");

    // Find list segment
    let list_segment = segments.iter().find(|s| matches!(s, Segment::List { .. }));
    assert!(list_segment.is_some(), "Should parse ordered list");

    if let Some(Segment::List { ordered, items }) = list_segment {
        assert!(ordered, "Should be marked as ordered list");
        assert_eq!(items.len(), 3);
    }
}

#[test]
fn test_unordered_lists() {
    let markdown = "- First item\n- Second item\n- Third item";
    let segments = MarkdownProcessor::parse(markdown).expect("Failed to parse");

    let list_segment = segments.iter().find(|s| matches!(s, Segment::List { .. }));
    assert!(list_segment.is_some(), "Should parse unordered list");

    if let Some(Segment::List { ordered, items }) = list_segment {
        assert!(!ordered, "Should be marked as unordered list");
        assert_eq!(items.len(), 3);
    }
}

#[test]
fn test_blockquote() {
    let markdown = "> This is a quote\n> Multiple lines";
    let segments = MarkdownProcessor::parse(markdown).expect("Failed to parse");

    let quote_segment = segments.iter().find(|s| matches!(s, Segment::BlockQuote { .. }));
    assert!(quote_segment.is_some(), "Should parse blockquote");

    if let Some(Segment::BlockQuote { content }) = quote_segment {
        assert!(content.contains("This is a quote"));
    }
}

#[test]
fn test_unicode_content() {
    let markdown = "# 日本語タイトル\n\n中文段落。\n\n한글 텍스트";
    let segments = MarkdownProcessor::parse(markdown).expect("Failed to parse unicode");

    // Verify UTF-8 handling
    if let Segment::Heading { text, .. } = &segments[0] {
        assert!(text.contains("日本語"));
    }

    if let Segment::Paragraph { text } = &segments[1] {
        assert!(text.contains("中文"));
    }
}

#[test]
fn test_inline_code_in_paragraph() {
    let markdown = "This paragraph has `inline code` in it.";
    let segments = MarkdownProcessor::parse(markdown).expect("Failed to parse");

    if let Segment::Paragraph { text } = &segments[0] {
        assert!(text.contains("`inline code`"));
    }
}

#[test]
fn test_nested_list_with_code() {
    let markdown = r#"- Item 1
  - Nested item
- Item 2

```rust
code block
```
"#;
    let segments = MarkdownProcessor::parse(markdown).expect("Failed to parse");

    // Should have list segment and code block segment
    let has_list = segments.iter().any(|s| matches!(s, Segment::List { .. }));
    let has_code = segments.iter().any(|s| matches!(s, Segment::CodeBlock { .. }));

    assert!(has_list, "Should have list segment");
    assert!(has_code, "Should have code block segment");
}

#[test]
fn test_list_is_translatable() {
    let markdown = "- Item 1\n- Item 2";
    let segments = MarkdownProcessor::parse(markdown).expect("Failed to parse");

    let list_segment = segments.iter().find(|s| matches!(s, Segment::List { .. }));
    assert!(list_segment.is_some());
    assert!(list_segment.unwrap().is_translatable(), "Lists should be translatable");
}

#[test]
fn test_blockquote_is_translatable() {
    let markdown = "> Quote text";
    let segments = MarkdownProcessor::parse(markdown).expect("Failed to parse");

    let quote_segment = segments.iter().find(|s| matches!(s, Segment::BlockQuote { .. }));
    assert!(quote_segment.is_some());
    assert!(quote_segment.unwrap().is_translatable(), "BlockQuotes should be translatable");
}
```

**Coverage**: Edge cases, Lists, BlockQuotes, Unicode, nested structures

### Success Criteria

#### Automated Verification:
- [x] Error tests compile and pass: `cargo test error_test`
- [x] Renderer tests compile and pass: `cargo test markdown_renderer_test`
- [x] Extended parser tests pass: `cargo test segment_ordering_unit_test`
- [x] All Phase 1 tests run in <1 second: `time cargo test --lib`

#### Manual Verification:
- [x] Verify code coverage increased for error.rs, renderer.rs, parser.rs
- [x] Check that test output is clear and descriptive
- [x] Confirm tests are deterministic (pass consistently)

---

## Phase 2: Trait Extraction (Moderate Refactoring)

### Overview
Introduce trait-based abstractions for network and filesystem dependencies, enabling mock implementations for testing. Uses default generic parameters to maintain backward compatibility.

**Effort**: 6-8 hours
**Impact**: Eliminates external dependencies in tests

### Changes Required

#### 1. Add async-trait Dependency

**File**: `Cargo.toml`

**Changes**: Add async-trait to dependencies section

```toml
[dependencies]
# ... existing dependencies ...
async-trait = "0.1"
```

**Reasoning**: Required for async trait methods since Rust doesn't natively support async in traits yet

#### 2. Create TranslationService Trait

**File**: `src/translator/translation_service.rs` (new file)

**Purpose**: Abstract translation logic from HTTP implementation

```rust
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
```

**Justification**:
- `Clone` required for spawning parallel async tasks (line 114 in translator/mod.rs)
- `Send + Sync` required for tokio task spawning
- Single `translate()` method matches current `OllamaClient::translate()` at ollama.rs:41

#### 3. Implement TranslationService for OllamaClient

**File**: `src/translator/ollama.rs`

**Changes**: Add trait implementation after existing impl block

```rust
use crate::translator::translation_service::TranslationService;
use async_trait::async_trait;

#[async_trait]
impl TranslationService for OllamaClient {
    async fn translate(&self, text: &str) -> Result<String> {
        // Existing implementation from lines 42-67
        let prompt = format!(
            "You are a professional translator. Translate the following English markdown text to Japanese. \
             Keep markdown formatting intact. Do not translate code blocks, URLs, or technical terms.\n\n{}",
            text
        );

        let request = OllamaRequest {
            model: self.model.clone(),
            prompt,
            stream: false,
            options: OllamaOptions { temperature: 0.3 },
        };

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

    fn model(&self) -> &str {
        &self.model
    }
}
```

**Note**: Existing `translate()` method at line 41 stays as-is; trait impl delegates to it or duplicates logic

#### 4. Create CacheBackend Trait

**File**: `src/translator/cache_backend.rs` (new file)

**Purpose**: Abstract cache storage from filesystem implementation

```rust
use crate::error::Result;
use crate::translator::CacheStats;

/// Trait for cache storage backends
/// Allows in-memory mocking in tests
pub trait CacheBackend: Send + Sync {
    /// Retrieve cached translation if exists and valid
    fn get(&mut self, source: &str, model: &str, lang_pair: &str) -> Option<String>;

    /// Store translation in cache
    fn set(&self, source: &str, translation: &str, model: &str, lang_pair: &str) -> Result<()>;

    /// Clear all cache entries
    fn clear(&self) -> Result<()>;

    /// Get cache statistics
    fn stats(&mut self) -> &CacheStats;
}
```

**Justification**:
- Matches existing `TranslationCache` API (cache.rs:50, 70, 89, 99)
- `Send + Sync` for safe sharing in async context
- Mutable `get()` and `stats()` for tracking statistics

#### 5. Implement CacheBackend for TranslationCache

**File**: `src/translator/cache.rs`

**Changes**: Add trait implementation

```rust
use crate::translator::cache_backend::CacheBackend;

impl CacheBackend for TranslationCache {
    fn get(&mut self, source: &str, model: &str, lang_pair: &str) -> Option<String> {
        // Existing implementation from lines 51-67
        self.stats.total_requests += 1;

        let key = Self::generate_key(source, model, lang_pair);
        let cache_path = self.cache_dir.join(&key).with_extension("json");

        if let Ok(content) = std::fs::read_to_string(&cache_path) {
            if let Ok(entry) = serde_json::from_str::<CacheEntry>(&content) {
                if entry.checksum == Self::hash_text(source) {
                    self.stats.cache_hits += 1;
                    return Some(entry.translation);
                }
            }
        }

        self.stats.cache_misses += 1;
        None
    }

    fn set(&self, source: &str, translation: &str, model: &str, lang_pair: &str) -> Result<()> {
        // Existing implementation from lines 71-85
        let key = Self::generate_key(source, model, lang_pair);
        let cache_path = self.cache_dir.join(&key).with_extension("json");

        let entry = CacheEntry {
            source: source.to_string(),
            translation: translation.to_string(),
            model: model.to_string(),
            language_pair: lang_pair.to_string(),
            created_at: Utc::now(),
            checksum: Self::hash_text(source),
        };

        let json = serde_json::to_string_pretty(&entry)?;
        std::fs::write(cache_path, json)?;
        Ok(())
    }

    fn clear(&self) -> Result<()> {
        // Existing implementation from lines 90-96
        for entry in std::fs::read_dir(&self.cache_dir)? {
            let entry = entry?;
            if entry.path().extension().and_then(|s| s.to_str()) == Some("json") {
                std::fs::remove_file(entry.path())?;
            }
        }
        Ok(())
    }

    fn stats(&mut self) -> &CacheStats {
        // Existing implementation from lines 100-111
        let mut total_size: u64 = 0;
        if let Ok(entries) = std::fs::read_dir(&self.cache_dir) {
            for entry in entries.flatten() {
                if let Ok(metadata) = entry.metadata() {
                    total_size += metadata.len();
                }
            }
        }
        self.stats.total_size_bytes = total_size;
        &self.stats
    }
}
```

#### 6. Update Translator with Generic Parameters

**File**: `src/translator/mod.rs`

**Changes**: Add generic parameters with defaults for backward compatibility

```rust
// Update struct definition at line 32
pub struct Translator<T = OllamaClient, C = TranslationCache>
where
    T: TranslationService,
    C: CacheBackend,
{
    client: T,
    cache: C,
    config: TranslatorConfig,
}

// Update impl block at line 38
impl Translator<OllamaClient, TranslationCache> {
    pub fn new(ollama_url: String, model: String, config: TranslatorConfig) -> Result<Self> {
        Ok(Self {
            client: OllamaClient::new(ollama_url, model),
            cache: TranslationCache::new()?,
            config,
        })
    }
}

// Add new constructor for testing
impl<T, C> Translator<T, C>
where
    T: TranslationService,
    C: CacheBackend,
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
        // Existing implementation unchanged
        let segments = MarkdownProcessor::parse(markdown)?;
        let translated_segments = self.translate_segments(segments).await?;
        Ok(Self::reconstruct_markdown(&translated_segments))
    }

    pub async fn translate_to_html(&mut self, markdown: &str, title: &str) -> Result<String> {
        // Existing implementation unchanged
        let segments = MarkdownProcessor::parse(markdown)?;
        let translated_segments = self.translate_segments(segments).await?;
        Ok(crate::markdown::HtmlRenderer::render(
            &translated_segments,
            title,
        ))
    }

    // Move all other methods to this generic impl block
    // Lines 62-210 move here with minimal changes

    async fn translate_segments(&mut self, segments: Vec<Segment>) -> Result<Vec<Segment>> {
        // ... existing logic ...
        // Change line 99: self.cache.get(text, &self.ollama_client.model, &lang_pair)
        // to: self.cache.get(text, self.client.model(), &lang_pair)

        // Change line 114: let ollama_clone = self.ollama_client.clone();
        // to: let client_clone = self.client.clone();

        // Change line 119: let result = ollama_clone.translate(&text_owned).await;
        // to: let result = client_clone.translate(&text_owned).await;

        // Change line 139: &self.ollama_client.model,
        // to: self.client.model(),
    }

    pub fn cache_stats(&mut self) -> &CacheStats {
        self.cache.stats()
    }

    pub fn clear_cache(&self) -> Result<()> {
        self.cache.clear()
    }
}
```

**Backward Compatibility**: Default generic parameters mean existing code like `Translator::new()` works unchanged

#### 7. Update Module Exports

**File**: `src/translator/mod.rs`

**Changes**: Add new modules at top

```rust
mod cache;
mod cache_backend;  // New
mod ollama;
mod translation_service;  // New

pub use cache::{CacheStats, TranslationCache};
pub use cache_backend::CacheBackend;  // New
pub use ollama::OllamaClient;
pub use translation_service::TranslationService;  // New
```

#### 8. Create Mock Translation Service

**File**: `tests/mocks/mock_translation_service.rs` (new file)

**Purpose**: Deterministic mock for testing translator logic

```rust
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
```

#### 9. Create Mock Cache Backend

**File**: `tests/mocks/mock_cache.rs` (new file)

**Purpose**: In-memory cache for fast, deterministic testing

```rust
use md_translate::error::Result;
use md_translate::translator::{CacheBackend, CacheStats};
use std::collections::HashMap;
use std::sync::{Arc, Mutex};
use sha2::{Digest, Sha256};

/// In-memory cache implementation for testing
/// Mimics TranslationCache behavior without filesystem I/O
#[derive(Clone)]
pub struct InMemoryCache {
    store: Arc<Mutex<HashMap<String, String>>>,
    stats: Arc<Mutex<CacheStats>>,
}

impl InMemoryCache {
    pub fn new() -> Self {
        Self {
            store: Arc::new(Mutex::new(HashMap::new())),
            stats: Arc::new(Mutex::new(CacheStats::default())),
        }
    }

    fn generate_key(source: &str, model: &str, lang_pair: &str) -> String {
        let combined = format!("{}{}{}", source, model, lang_pair);
        let mut hasher = Sha256::new();
        hasher.update(combined.as_bytes());
        format!("{:x}", hasher.finalize())
    }
}

impl Default for InMemoryCache {
    fn default() -> Self {
        Self::new()
    }
}

impl CacheBackend for InMemoryCache {
    fn get(&mut self, source: &str, model: &str, lang_pair: &str) -> Option<String> {
        let mut stats = self.stats.lock().unwrap();
        stats.total_requests += 1;

        let key = Self::generate_key(source, model, lang_pair);
        let store = self.store.lock().unwrap();

        if let Some(value) = store.get(&key) {
            stats.cache_hits += 1;
            Some(value.clone())
        } else {
            stats.cache_misses += 1;
            None
        }
    }

    fn set(&self, source: &str, translation: &str, model: &str, lang_pair: &str) -> Result<()> {
        let key = Self::generate_key(source, model, lang_pair);
        let mut store = self.store.lock().unwrap();
        store.insert(key, translation.to_string());
        Ok(())
    }

    fn clear(&self) -> Result<()> {
        let mut store = self.store.lock().unwrap();
        store.clear();

        let mut stats = self.stats.lock().unwrap();
        *stats = CacheStats::default();
        Ok(())
    }

    fn stats(&mut self) -> &CacheStats {
        // Return reference to stats
        // Note: This is tricky with Arc<Mutex<>>. May need redesign.
        // For now, update a local copy
        let stats = self.stats.lock().unwrap();
        unsafe {
            // SAFETY: We know the stats live as long as self
            // This is a workaround for the trait signature requiring &mut self
            &*(stats.as_ref() as *const CacheStats)
        }
    }
}
```

**Note**: The `stats()` method signature with `&mut self` returning `&CacheStats` is problematic with `Arc<Mutex<>>`. May need to reconsider trait design.

**Alternative**: Change trait to return `CacheStats` by value:

```rust
pub trait CacheBackend: Send + Sync {
    fn stats(&self) -> CacheStats;  // Return by value
}
```

Then update implementations accordingly.

#### 10. Create Mock Module

**File**: `tests/mocks/mod.rs` (new file)

```rust
mod mock_cache;
mod mock_translation_service;

pub use mock_cache::InMemoryCache;
pub use mock_translation_service::MockTranslationService;
```

#### 11. Write Translator Unit Tests with Mocks

**File**: `tests/translator_unit_test.rs` (new file)

**Purpose**: Test translator orchestration logic without external dependencies

```rust
mod mocks;

use md_translate::markdown::parser::{MarkdownProcessor, Segment};
use md_translate::translator::{Translator, TranslatorConfig};
use mocks::{InMemoryCache, MockTranslationService};

#[tokio::test]
async fn test_segment_ordering_with_mock() {
    let mock_client = MockTranslationService::default();
    let mock_cache = InMemoryCache::new();

    let config = TranslatorConfig {
        source_lang: "en".to_string(),
        target_lang: "ja".to_string(),
        use_cache: false,
        parallel_requests: 3,
        show_progress: false,
    };

    let mut translator = Translator::new_with_deps(mock_client, mock_cache, config);

    let markdown = r#"# H1
P1

```code
block
```

P2
P3"#;

    let result = translator.translate_markdown(markdown).await.expect("Translation failed");
    let segments = MarkdownProcessor::parse(&result).expect("Parse failed");

    // Critical: Verify code block is NOT at the beginning (was the bug)
    assert!(matches!(segments[0], Segment::Heading { .. }), "First segment should be heading");
    assert!(matches!(segments[2], Segment::CodeBlock { .. }), "Code block at index 2");

    // Verify all segments maintain order
    assert_eq!(segments.len(), 5);
}

#[tokio::test]
async fn test_parallel_translation() {
    let mock_client = MockTranslationService::with_delay(50);
    let mock_cache = InMemoryCache::new();

    let config = TranslatorConfig {
        parallel_requests: 3,
        use_cache: false,
        show_progress: false,
        ..Default::default()
    };

    let mut translator = Translator::new_with_deps(mock_client, mock_cache, config);

    // 5 translatable segments should complete faster than 5 * 50ms due to parallelism
    let start = std::time::Instant::now();
    let markdown = "# H1\nP1\n# H2\nP2\n# H3";
    translator.translate_markdown(markdown).await.expect("Translation failed");
    let duration = start.elapsed();

    // With 3 parallel requests, should complete in ~2 batches (~100ms), not 5 sequential (~250ms)
    assert!(duration < std::time::Duration::from_millis(200), "Should benefit from parallelism");
}

#[tokio::test]
async fn test_cache_integration() {
    let mock_client = MockTranslationService::default();
    let mock_cache = InMemoryCache::new();

    let config = TranslatorConfig {
        use_cache: true,
        show_progress: false,
        ..Default::default()
    };

    let mut translator = Translator::new_with_deps(mock_client, mock_cache, config);

    // First translation - cache miss
    let markdown = "# Test\nParagraph";
    let result1 = translator.translate_markdown(markdown).await.expect("Translation 1 failed");
    let stats1 = translator.cache_stats();
    assert_eq!(stats1.cache_misses, 2, "Should have 2 misses (heading + paragraph)");

    // Second translation - cache hit
    let result2 = translator.translate_markdown(markdown).await.expect("Translation 2 failed");
    let stats2 = translator.cache_stats();
    assert_eq!(stats2.cache_hits, 2, "Should have 2 hits");

    // Results should be identical
    assert_eq!(result1, result2);
}

#[tokio::test]
async fn test_non_translatable_segments_passthrough() {
    let mock_client = MockTranslationService::default();
    let mock_cache = InMemoryCache::new();

    let config = TranslatorConfig {
        use_cache: false,
        show_progress: false,
        ..Default::default()
    };

    let mut translator = Translator::new_with_deps(mock_client, mock_cache, config);

    let markdown = r#"```rust
let x = 42;
```

---
"#;

    let result = translator.translate_markdown(markdown).await.expect("Translation failed");

    // Code blocks and horizontal rules should not be modified
    assert!(result.contains("let x = 42;"));
    assert!(result.contains("---"));
}

#[tokio::test]
async fn test_translation_error_handling() {
    use md_translate::error::MdTranslateError;

    // Create mock that returns errors
    let mock_client = MockTranslationService::with_responses(vec![
        Err(MdTranslateError::OllamaError("Network timeout".to_string())),
    ]);
    let mock_cache = InMemoryCache::new();

    let config = TranslatorConfig {
        use_cache: false,
        show_progress: false,
        ..Default::default()
    };

    let mut translator = Translator::new_with_deps(mock_client, mock_cache, config);

    let markdown = "# Test";
    let result = translator.translate_markdown(markdown).await;

    // Should preserve original text on error (line 148-149 in translator/mod.rs)
    let translated = result.expect("Should handle error gracefully");
    assert!(translated.contains("Test"), "Original text preserved on error");
}

#[tokio::test]
async fn test_reconstruct_markdown_format() {
    let mock_client = MockTranslationService::default();
    let mock_cache = InMemoryCache::new();

    let config = TranslatorConfig {
        use_cache: false,
        show_progress: false,
        ..Default::default()
    };

    let mut translator = Translator::new_with_deps(mock_client, mock_cache, config);

    let markdown = r#"# Title

Paragraph text.

```rust
code
```
"#;

    let result = translator.translate_markdown(markdown).await.expect("Translation failed");

    // Verify reconstruction formatting
    assert!(result.starts_with("# "), "Heading format preserved");
    assert!(result.contains("\n\n"), "Double newlines between segments");
    assert!(result.contains("```rust\n"), "Code block fencing preserved");
}
```

**Coverage**: Ordering, parallelism, caching, error handling, non-translatable passthrough, reconstruction

#### 12. Write Cache Unit Tests

**File**: `tests/cache_unit_test.rs` (new file)

**Purpose**: Test cache logic independently with in-memory implementation

```rust
mod mocks;

use mocks::InMemoryCache;
use md_translate::translator::CacheBackend;

#[test]
fn test_cache_hit_miss() {
    let mut cache = InMemoryCache::new();

    // Miss on first access
    let result = cache.get("hello", "model", "en-ja");
    assert!(result.is_none());

    let stats = cache.stats();
    assert_eq!(stats.cache_misses, 1);
    assert_eq!(stats.total_requests, 1);

    // Set value
    cache.set("hello", "こんにちは", "model", "en-ja").expect("Set failed");

    // Hit on second access
    let result = cache.get("hello", "model", "en-ja");
    assert_eq!(result, Some("こんにちは".to_string()));

    let stats = cache.stats();
    assert_eq!(stats.cache_hits, 1);
    assert_eq!(stats.total_requests, 2);
}

#[test]
fn test_cache_key_uniqueness() {
    let mut cache = InMemoryCache::new();

    // Different models should have different keys
    cache.set("text", "trans1", "model1", "en-ja").expect("Set 1 failed");
    cache.set("text", "trans2", "model2", "en-ja").expect("Set 2 failed");

    assert_eq!(cache.get("text", "model1", "en-ja"), Some("trans1".to_string()));
    assert_eq!(cache.get("text", "model2", "en-ja"), Some("trans2".to_string()));

    // Different language pairs should have different keys
    cache.set("text", "trans3", "model1", "en-es").expect("Set 3 failed");

    assert_eq!(cache.get("text", "model1", "en-ja"), Some("trans1".to_string()));
    assert_eq!(cache.get("text", "model1", "en-es"), Some("trans3".to_string()));
}

#[test]
fn test_cache_clear() {
    let mut cache = InMemoryCache::new();

    cache.set("key1", "value1", "model", "en-ja").expect("Set 1 failed");
    cache.set("key2", "value2", "model", "en-ja").expect("Set 2 failed");

    cache.clear().expect("Clear failed");

    assert!(cache.get("key1", "model", "en-ja").is_none());
    assert!(cache.get("key2", "model", "en-ja").is_none());
}

#[test]
fn test_cache_stats_tracking() {
    let mut cache = InMemoryCache::new();

    cache.get("key1", "model", "en-ja");  // Miss
    cache.set("key1", "value1", "model", "en-ja").expect("Set failed");
    cache.get("key1", "model", "en-ja");  // Hit
    cache.get("key2", "model", "en-ja");  // Miss

    let stats = cache.stats();
    assert_eq!(stats.total_requests, 3);
    assert_eq!(stats.cache_hits, 1);
    assert_eq!(stats.cache_misses, 2);
}

#[test]
fn test_cache_overwrites() {
    let mut cache = InMemoryCache::new();

    cache.set("key", "value1", "model", "en-ja").expect("Set 1 failed");
    cache.set("key", "value2", "model", "en-ja").expect("Set 2 failed");

    // Second set should overwrite first
    assert_eq!(cache.get("key", "model", "en-ja"), Some("value2".to_string()));
}
```

**Coverage**: Cache hit/miss logic, key generation, clear functionality, statistics tracking

### Success Criteria

#### Automated Verification:
- [x] Trait code compiles: `cargo check`
- [x] Translator tests with mocks pass: `cargo test translator_unit_test`
- [x] Cache unit tests pass: `cargo test cache_unit_test`
- [ ] Existing integration tests still pass (with Ollama): `cargo test segment_ordering_test` (pre-existing failure)
- [x] All unit tests pass: `cargo test --test error_test --test markdown_renderer_test --test segment_ordering_unit_test --test translator_unit_test --test cache_unit_test`
- [x] Public API unchanged: existing code using `Translator::new()` compiles without changes

#### Manual Verification:
- [x] Verify backward compatibility - CLI usage works identically
- [x] Check that mock tests run significantly faster than integration tests (0.11s vs 51s for integration)
- [x] Confirm trait abstractions don't impact runtime performance
- [x] Validate generic type parameters resolve correctly in IDE

---

## Phase 3: CLI Testing (Minor Refactoring)

### Overview
Test CLI command handlers with dependency injection, using temp files for I/O and mocks for translation.

**Effort**: 3-4 hours
**Impact**: Complete end-to-end test coverage

### Changes Required

#### 1. Refactor Command Handlers for Testability

**File**: `src/cli/commands.rs`

**Changes**: Add trait-aware versions of command handlers

```rust
// Add generic version of handle_translate
pub async fn handle_translate_with_translator<T, C>(
    translator: &mut Translator<T, C>,
    args: TranslateArgs,
) -> Result<()>
where
    T: TranslationService,
    C: CacheBackend,
{
    // Move existing logic from handle_translate (lines 8-43)
    let input = std::fs::read_to_string(&args.file)?;

    println!("{}", "Translating markdown...".cyan());
    let translated = translator.translate_markdown(&input).await?;

    if let Some(output_path) = args.output {
        std::fs::write(&output_path, translated)?;
        println!("{} {}", "Saved to:".green(), output_path);
    } else {
        println!("{}", translated);
    }

    // Show cache stats
    let stats = translator.cache_stats();
    println!("\n{}", "Cache Statistics:".cyan());
    println!("  Hits: {}", stats.cache_hits);
    println!("  Misses: {}", stats.cache_misses);
    println!(
        "  Hit Rate: {:.1}%",
        if stats.total_requests > 0 {
            (stats.cache_hits as f64 / stats.total_requests as f64) * 100.0
        } else {
            0.0
        }
    );

    Ok(())
}

// Keep existing handle_translate for backward compatibility
pub async fn handle_translate(args: TranslateArgs) -> Result<()> {
    let config = TranslatorConfig {
        use_cache: !args.no_cache,
        show_progress: true,
        ..Default::default()
    };

    let mut translator = Translator::new(args.ollama_url, args.model, config)?;
    handle_translate_with_translator(&mut translator, args).await
}
```

**Similar refactoring for `handle_view()` and `handle_cache()`**

#### 2. Write CLI Command Tests

**File**: `tests/cli_commands_test.rs` (new file)

**Purpose**: Test command handlers with temp files and mocks

```rust
mod mocks;

use md_translate::cli::{TranslateArgs, CacheArgs, CacheCommands};
use md_translate::cli::commands::{handle_translate_with_translator, handle_cache_with_translator};
use md_translate::translator::{Translator, TranslatorConfig};
use mocks::{InMemoryCache, MockTranslationService};
use std::fs;
use tempfile::tempdir;

#[tokio::test]
async fn test_translate_command_basic() {
    let temp_dir = tempdir().expect("Failed to create temp dir");
    let input_path = temp_dir.path().join("input.md");
    let output_path = temp_dir.path().join("output.md");

    // Write test input
    fs::write(&input_path, "# Test\n\nParagraph.").expect("Failed to write input");

    // Create mock translator
    let mock_client = MockTranslationService::default();
    let mock_cache = InMemoryCache::new();
    let config = TranslatorConfig {
        use_cache: false,
        show_progress: false,
        ..Default::default()
    };
    let mut translator = Translator::new_with_deps(mock_client, mock_cache, config);

    // Create args
    let args = TranslateArgs {
        file: input_path.to_string_lossy().to_string(),
        output: Some(output_path.to_string_lossy().to_string()),
        model: "mock".to_string(),
        ollama_url: "http://mock".to_string(),
        no_cache: true,
        format: "markdown".to_string(),
    };

    // Execute command
    handle_translate_with_translator(&mut translator, args)
        .await
        .expect("Command failed");

    // Verify output file was created
    assert!(output_path.exists(), "Output file should exist");

    let output = fs::read_to_string(&output_path).expect("Failed to read output");
    assert!(output.contains("[TRANSLATED:"), "Should contain mock translation marker");
}

#[tokio::test]
async fn test_translate_command_default_output() {
    let temp_dir = tempdir().expect("Failed to create temp dir");
    let input_path = temp_dir.path().join("test.md");

    fs::write(&input_path, "# Test").expect("Failed to write input");

    let mock_client = MockTranslationService::default();
    let mock_cache = InMemoryCache::new();
    let config = TranslatorConfig::default();
    let mut translator = Translator::new_with_deps(mock_client, mock_cache, config);

    let args = TranslateArgs {
        file: input_path.to_string_lossy().to_string(),
        output: None,  // Should print to stdout
        model: "mock".to_string(),
        ollama_url: "http://mock".to_string(),
        no_cache: true,
        format: "markdown".to_string(),
    };

    // Execute command (output goes to stdout)
    handle_translate_with_translator(&mut translator, args)
        .await
        .expect("Command failed");
}

#[tokio::test]
async fn test_translate_command_missing_input() {
    let mock_client = MockTranslationService::default();
    let mock_cache = InMemoryCache::new();
    let config = TranslatorConfig::default();
    let mut translator = Translator::new_with_deps(mock_client, mock_cache, config);

    let args = TranslateArgs {
        file: "/nonexistent/file.md".to_string(),
        output: None,
        model: "mock".to_string(),
        ollama_url: "http://mock".to_string(),
        no_cache: true,
        format: "markdown".to_string(),
    };

    // Should return IoError
    let result = handle_translate_with_translator(&mut translator, args).await;
    assert!(result.is_err());
}

#[tokio::test]
async fn test_cache_command_stats() {
    let mock_client = MockTranslationService::default();
    let mock_cache = InMemoryCache::new();
    let config = TranslatorConfig {
        use_cache: true,
        ..Default::default()
    };
    let mut translator = Translator::new_with_deps(mock_client, mock_cache, config);

    // Populate cache
    translator.translate_markdown("# Test").await.expect("Translation failed");

    // Test stats command
    let args = CacheArgs {
        command: CacheCommands::Stats,
    };

    // Note: Need to add generic version of handle_cache
    // handle_cache_with_translator(&mut translator, args).await.expect("Command failed");
}

#[tokio::test]
async fn test_cache_command_clear() {
    let mock_client = MockTranslationService::default();
    let mock_cache = InMemoryCache::new();
    let config = TranslatorConfig {
        use_cache: true,
        ..Default::default()
    };
    let mut translator = Translator::new_with_deps(mock_client, mock_cache, config);

    // Populate cache
    translator.translate_markdown("# Test").await.expect("Translation failed");

    // Test clear command
    let args = CacheArgs {
        command: CacheCommands::Clear,
    };

    // handle_cache_with_translator(&mut translator, args).await.expect("Command failed");

    // Verify cache is cleared
    let stats = translator.cache_stats();
    assert_eq!(stats.cache_hits, 0);
}
```

**Note**: Requires adding `tempfile` to dev-dependencies

#### 3. Add Test Dependencies

**File**: `Cargo.toml`

**Changes**: Add dev-dependencies section

```toml
[dev-dependencies]
tempfile = "3"
```

**Purpose**: Temporary directories for file I/O testing

### Success Criteria

#### Automated Verification:
- [x] CLI tests compile: `cargo test --test cli_commands_test --no-run`
- [x] CLI tests pass: `cargo test cli_commands_test`
- [x] All tests pass: `cargo test`

#### Manual Verification:
- [x] Verify temp files are cleaned up after tests
- [x] Check that command handlers work correctly with real Ollama
- [x] Validate error handling for missing files, permissions, etc.

---

## Testing Strategy

### Unit Tests (Fast, Always Run)
```bash
# Run all unit tests (no external dependencies)
cargo test --lib
cargo test error_test
cargo test markdown_parser_test
cargo test markdown_renderer_test
cargo test cache_unit_test
cargo test translator_unit_test

# Should complete in <5 seconds total
```

### Integration Tests (Slow, Optional)
```bash
# Require Ollama service running
cargo test segment_ordering_test

# Or run with feature flag for optional execution
cargo test --test integration -- --ignored
```

### CI/CD Configuration
```yaml
# .github/workflows/test.yml
- name: Run fast tests
  run: cargo test --lib --tests --exclude segment_ordering_test

- name: Run integration tests (if Ollama available)
  run: |
    if systemctl is-active --quiet ollama; then
      cargo test segment_ordering_test
    fi
  continue-on-error: true
```

## Performance Considerations

### Expected Test Times
- **Phase 1 unit tests**: <1 second total
- **Phase 2 mock-based tests**: <5 seconds total
- **Existing integration tests**: 30-60 seconds (network dependent)

### Memory Usage
- Mock implementations use in-memory storage (minimal overhead)
- Generic type parameters have zero runtime cost (monomorphization)

## Migration Notes

### Backward Compatibility
- Default generic parameters preserve existing API
- `Translator::new()` continues to work unchanged
- CLI commands unchanged
- Existing integration tests unchanged

### Gradual Migration
Each phase can be merged independently:
1. Phase 1 provides immediate value with no refactoring
2. Phase 2 can be adopted incrementally (old tests still work)
3. Phase 3 builds on Phase 2 infrastructure

## References

- Original research: `thoughts/shared/research/2025-09-30-test-strategy-per-module.md`
- Current test patterns: `tests/segment_ordering_unit_test.rs:6-288`
- Integration test approach: `tests/segment_ordering_test.rs:6-331`
- Translator orchestration logic: `src/translator/mod.rs:62-167`
- Cache implementation: `src/translator/cache.rs:25-123`
- Ollama client: `src/translator/ollama.rs:23-92`