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

    // Find the code block and verify it's not first
    let code_idx = segments.iter().position(|s| matches!(s, Segment::CodeBlock { .. }))
        .expect("Should have code block");
    assert!(code_idx > 0, "Code block should not be at index 0");

    // Verify we have expected segment types
    let has_heading = segments.iter().any(|s| matches!(s, Segment::Heading { .. }));
    let has_paragraph = segments.iter().any(|s| matches!(s, Segment::Paragraph { .. }));
    let has_code = segments.iter().any(|s| matches!(s, Segment::CodeBlock { .. }));

    assert!(has_heading, "Should have heading");
    assert!(has_paragraph, "Should have paragraph");
    assert!(has_code, "Should have code block");
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