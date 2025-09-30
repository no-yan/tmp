/// Test that segments maintain their original order after translation
/// This is a critical test for the index-based reconstruction fix
#[cfg(test)]
mod test {
    use md_translate::markdown::parser::{MarkdownProcessor, Segment};
    use md_translate::translator::{Translator, TranslatorConfig};

    #[tokio::test]
    #[ignore]
    async fn test_segment_ordering_preserved() {
        // Create a mock translator with disabled cache and progress
        let config = TranslatorConfig {
            source_lang: "en".to_string(),
            target_lang: "ja".to_string(),
            use_cache: false,
            parallel_requests: 3,
            show_progress: false,
        };

        // Note: This test requires a running Ollama instance
        // In production, we should mock the OllamaClient
        let mut translator = Translator::new(
            "http://localhost:11434".to_string(),
            "qwen2.5:7b".to_string(),
            config,
        )
        .expect("Failed to create translator");

        // Test markdown with intentionally mixed segment types
        let test_markdown = r#"# Heading 1

First paragraph here.

```rust
fn main() {
    println!("Code block 1");
}
```

# Heading 2

Second paragraph here.

```python
print("Code block 2")
```

Third paragraph here.

---

Fourth paragraph after horizontal rule.

```js
console.log("Code block 3");
```

# Heading 3

Final paragraph.
"#;

        // Parse the markdown to get the expected segment order
        let original_segments =
            MarkdownProcessor::parse(test_markdown).expect("Failed to parse markdown");

        // Record the original segment types in order
        let original_order: Vec<String> = original_segments
            .iter()
            .map(|s| match s {
                Segment::Heading { level, text } => format!("Heading({}, {})", level, text),
                Segment::Paragraph { text } => {
                    format!("Paragraph({}...)", &text[..20.min(text.len())])
                }
                Segment::CodeBlock { language, .. } => format!("CodeBlock({:?})", language),
                Segment::HorizontalRule => "HorizontalRule".to_string(),
                _ => "Other".to_string(),
            })
            .collect();

        println!("Original segment order:");
        for (i, seg) in original_order.iter().enumerate() {
            println!("  {}: {}", i, seg);
        }

        // Translate the markdown
        let translated_markdown = translator
            .translate_markdown(test_markdown)
            .await
            .expect("Failed to translate markdown");

        // Parse the translated markdown
        let translated_segments = MarkdownProcessor::parse(&translated_markdown)
            .expect("Failed to parse translated markdown");

        // Record the translated segment types in order
        let translated_order: Vec<String> = translated_segments
            .iter()
            .map(|s| match s {
                Segment::Heading { level, .. } => format!("Heading({})", level),
                Segment::Paragraph { .. } => "Paragraph".to_string(),
                Segment::CodeBlock { language, .. } => format!("CodeBlock({:?})", language),
                Segment::HorizontalRule => "HorizontalRule".to_string(),
                _ => "Other".to_string(),
            })
            .collect();

        println!("\nTranslated segment order:");
        for (i, seg) in translated_order.iter().enumerate() {
            println!("  {}: {}", i, seg);
        }

        // Verify segment count matches
        assert_eq!(
            original_segments.len(),
            translated_segments.len(),
            "Segment count should not change after translation"
        );

        // Verify segment types and order are preserved
        for (i, (orig, trans)) in original_segments
            .iter()
            .zip(translated_segments.iter())
            .enumerate()
        {
            match (orig, trans) {
                (Segment::Heading { level: l1, .. }, Segment::Heading { level: l2, .. }) => {
                    assert_eq!(l1, l2, "Heading level mismatch at index {}", i);
                }
                (Segment::Paragraph { .. }, Segment::Paragraph { .. }) => {
                    // Types match, order is preserved
                }
                (
                    Segment::CodeBlock {
                        language: l1,
                        code: c1,
                    },
                    Segment::CodeBlock {
                        language: l2,
                        code: c2,
                    },
                ) => {
                    assert_eq!(l1, l2, "Code block language mismatch at index {}", i);
                    assert_eq!(
                        c1, c2,
                        "Code blocks should not be translated, mismatch at index {}",
                        i
                    );
                }
                (Segment::HorizontalRule, Segment::HorizontalRule) => {
                    // Types match, order is preserved
                }
                _ => {
                    panic!(
                        "Segment type mismatch at index {}: expected {:?} type, got {:?} type",
                        i,
                        std::mem::discriminant(orig),
                        std::mem::discriminant(trans)
                    );
                }
            }
        }

        // Specific check: verify code blocks are NOT at the beginning
        // (This was the bug - code blocks appeared first)
        let first_segment = &translated_segments[0];
        assert!(
            matches!(first_segment, Segment::Heading { .. }),
            "First segment should be a heading, not a code block"
        );

        // Verify code blocks appear in their expected positions
        let code_block_indices: Vec<usize> = translated_segments
            .iter()
            .enumerate()
            .filter_map(|(i, s)| {
                if matches!(s, Segment::CodeBlock { .. }) {
                    Some(i)
                } else {
                    None
                }
            })
            .collect();

        // In our test markdown, code blocks should be at indices 2, 5, 9
        let expected_code_indices = vec![2, 5, 9];
        assert_eq!(
            code_block_indices, expected_code_indices,
            "Code blocks should appear at their original positions"
        );

        println!("\n✓ Segment ordering test passed!");
    }

    /// Test with all translatable segments to verify async completion order doesn't matter
    #[tokio::test]
    #[ignore]
    async fn test_all_translatable_segments_ordering() {
        let config = TranslatorConfig {
            source_lang: "en".to_string(),
            target_lang: "ja".to_string(),
            use_cache: false,
            parallel_requests: 3,
            show_progress: false,
        };

        let mut translator = Translator::new(
            "http://localhost:11434".to_string(),
            "qwen2.5:7b".to_string(),
            config,
        )
        .expect("Failed to create translator");

        let test_markdown = r#"# First Heading

First paragraph with some text.

# Second Heading

Second paragraph with different text.

# Third Heading

Third paragraph with more text.

# Fourth Heading

Fourth paragraph with final text.
"#;

        let original_segments =
            MarkdownProcessor::parse(test_markdown).expect("Failed to parse markdown");

        let translated_markdown = translator
            .translate_markdown(test_markdown)
            .await
            .expect("Failed to translate markdown");

        let translated_segments = MarkdownProcessor::parse(&translated_markdown)
            .expect("Failed to parse translated markdown");

        // Verify all heading levels are in correct order
        let original_levels: Vec<u8> = original_segments
            .iter()
            .filter_map(|s| {
                if let Segment::Heading { level, .. } = s {
                    Some(*level)
                } else {
                    None
                }
            })
            .collect();

        let translated_levels: Vec<u8> = translated_segments
            .iter()
            .filter_map(|s| {
                if let Segment::Heading { level, .. } = s {
                    Some(*level)
                } else {
                    None
                }
            })
            .collect();

        assert_eq!(
            original_levels, translated_levels,
            "Heading levels should maintain their order"
        );

        println!("✓ All-translatable segments ordering test passed!");
    }

    /// Test with cache hits to verify cached segments maintain order
    #[tokio::test]
    #[ignore]
    async fn test_cached_segments_ordering() {
        let config = TranslatorConfig {
            source_lang: "en".to_string(),
            target_lang: "ja".to_string(),
            use_cache: true,
            parallel_requests: 3,
            show_progress: false,
        };

        let mut translator = Translator::new(
            "http://localhost:11434".to_string(),
            "qwen2.5:7b".to_string(),
            config,
        )
        .expect("Failed to create translator");

        // Clear cache before test
        translator.clear_cache().expect("Failed to clear cache");

        let test_markdown = r#"# Test Heading

Repeated paragraph text.

```code
block
```

Repeated paragraph text.

# Another Heading

Repeated paragraph text.
"#;

        // First translation - populate cache
        let first_translation = translator
            .translate_markdown(test_markdown)
            .await
            .expect("Failed first translation");

        // Second translation - should hit cache
        let second_translation = translator
            .translate_markdown(test_markdown)
            .await
            .expect("Failed second translation");

        // Verify both translations have same structure
        let first_segments = MarkdownProcessor::parse(&first_translation)
            .expect("Failed to parse first translation");
        let second_segments = MarkdownProcessor::parse(&second_translation)
            .expect("Failed to parse second translation");

        assert_eq!(
            first_segments.len(),
            second_segments.len(),
            "Cache hit should preserve segment count"
        );

        // Verify ordering is identical
        for (i, (s1, s2)) in first_segments
            .iter()
            .zip(second_segments.iter())
            .enumerate()
        {
            assert_eq!(
                std::mem::discriminant(s1),
                std::mem::discriminant(s2),
                "Segment types should match at index {} (cache hit)",
                i
            );
        }

        println!("✓ Cached segments ordering test passed!");
    }
}
