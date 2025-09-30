use md_translate::markdown::parser::{MarkdownProcessor, Segment};

/// Unit test for segment parsing order (no translation needed)
/// Verifies that the parser maintains segment order
#[test]
fn test_parser_maintains_segment_order() {
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

    let segments = MarkdownProcessor::parse(test_markdown).expect("Failed to parse markdown");

    // Expected order of segments (all # are level 1)
    let expected_order = vec![
        "Heading(1)",      // # Heading 1
        "Paragraph",       // First paragraph
        "CodeBlock",       // rust code
        "Heading(1)",      // # Heading 2 (still level 1!)
        "Paragraph",       // Second paragraph
        "CodeBlock",       // python code
        "Paragraph",       // Third paragraph
        "HorizontalRule",  // ---
        "Paragraph",       // Fourth paragraph
        "CodeBlock",       // js code
        "Heading(1)",      // # Heading 3 (still level 1!)
        "Paragraph",       // Final paragraph
    ];

    assert_eq!(
        segments.len(),
        expected_order.len(),
        "Segment count should match expected"
    );

    for (i, (segment, expected)) in segments.iter().zip(expected_order.iter()).enumerate() {
        let actual = match segment {
            Segment::Heading { level, .. } => format!("Heading({})", level),
            Segment::Paragraph { .. } => "Paragraph".to_string(),
            Segment::CodeBlock { .. } => "CodeBlock".to_string(),
            Segment::HorizontalRule => "HorizontalRule".to_string(),
            _ => "Other".to_string(),
        };

        assert_eq!(
            &actual, expected,
            "Segment mismatch at index {}: expected {}, got {}",
            i, expected, actual
        );
    }

    println!("✓ Parser maintains correct segment order!");
}

/// Test that code blocks preserve their content and language
#[test]
fn test_code_blocks_not_modified() {
    let test_markdown = r#"# Test

Some text.

```rust
let x = 42;
```

More text.

```python
print("hello")
```
"#;

    let segments = MarkdownProcessor::parse(test_markdown).expect("Failed to parse markdown");

    // Find code blocks
    let code_blocks: Vec<&Segment> = segments
        .iter()
        .filter(|s| matches!(s, Segment::CodeBlock { .. }))
        .collect();

    assert_eq!(code_blocks.len(), 2, "Should have 2 code blocks");

    // Check first code block
    if let Segment::CodeBlock { language, code } = &code_blocks[0] {
        assert_eq!(
            language.as_deref(),
            Some("rust"),
            "First code block should be rust"
        );
        assert!(
            code.contains("let x = 42;"),
            "Code content should be preserved"
        );
    } else {
        panic!("Expected CodeBlock");
    }

    // Check second code block
    if let Segment::CodeBlock { language, code } = &code_blocks[1] {
        assert_eq!(
            language.as_deref(),
            Some("python"),
            "Second code block should be python"
        );
        assert!(
            code.contains("print(\"hello\")"),
            "Code content should be preserved"
        );
    } else {
        panic!("Expected CodeBlock");
    }

    println!("✓ Code blocks are correctly parsed and preserved!");
}

/// Test segment reconstruction maintains order
#[test]
fn test_segment_reconstruction_order() {
    let test_markdown = r#"# Title

Paragraph 1

```code
block
```

Paragraph 2
"#;

    // Parse
    let segments = MarkdownProcessor::parse(test_markdown).expect("Failed to parse");

    // Simulate what reconstruct_markdown does
    let mut output = String::new();
    for segment in &segments {
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
            _ => {}
        }
    }

    // Verify reconstructed markdown starts with heading
    assert!(output.starts_with("# Title"), "Should start with heading");

    // Verify code block is not at the beginning
    let code_block_pos = output.find("```code");
    let paragraph1_pos = output.find("Paragraph 1");

    assert!(paragraph1_pos.is_some(), "Paragraph 1 should exist");
    assert!(code_block_pos.is_some(), "Code block should exist");
    assert!(
        paragraph1_pos.unwrap() < code_block_pos.unwrap(),
        "Paragraph 1 should come before code block"
    );

    println!("✓ Reconstruction maintains correct order!");
}

/// Test mixed translatable and non-translatable segments
#[test]
fn test_mixed_segment_types() {
    let test_markdown = r#"# H1

P1

```
code1
```

# H2

P2

---

P3

```
code2
```

P4
"#;

    let segments = MarkdownProcessor::parse(test_markdown).expect("Failed to parse");

    // Verify translatable and non-translatable are interspersed
    let types: Vec<(usize, bool)> = segments
        .iter()
        .enumerate()
        .map(|(i, s)| (i, s.is_translatable()))
        .collect();

    println!("Segment translatability:");
    for (i, translatable) in &types {
        println!("  {}: {}", i, if *translatable { "translatable" } else { "non-translatable" });
    }

    // Code blocks should not be translatable
    let code_indices: Vec<usize> = segments
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

    for &idx in &code_indices {
        assert!(
            !segments[idx].is_translatable(),
            "Code block at index {} should not be translatable",
            idx
        );
    }

    // Headings and paragraphs should be translatable
    let text_indices: Vec<usize> = segments
        .iter()
        .enumerate()
        .filter_map(|(i, s)| {
            if matches!(s, Segment::Heading { .. } | Segment::Paragraph { .. }) {
                Some(i)
            } else {
                None
            }
        })
        .collect();

    for &idx in &text_indices {
        assert!(
            segments[idx].is_translatable(),
            "Text segment at index {} should be translatable",
            idx
        );
    }

    println!("✓ Mixed segment types correctly identified!");
}

/// Test empty markdown input
#[test]
fn test_empty_markdown() {
    let segments = MarkdownProcessor::parse("").expect("Failed to parse empty markdown");
    assert_eq!(segments.len(), 0);
}

/// Test whitespace-only input
#[test]
fn test_whitespace_only() {
    let segments = MarkdownProcessor::parse("   \n\n   \n").expect("Failed to parse");
    // Should produce no segments or empty paragraph
    assert!(segments.is_empty() || segments.len() == 1);
}

/// Test unicode content (Japanese, Chinese, Korean)
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

/// Test inline code in paragraph
#[test]
fn test_inline_code_in_paragraph() {
    let markdown = "This paragraph has `inline code` in it.";
    let segments = MarkdownProcessor::parse(markdown).expect("Failed to parse");

    if let Segment::Paragraph { text } = &segments[0] {
        assert!(text.contains("`inline code`"));
    }
}

/// Test multiple code blocks with different languages
#[test]
fn test_multiple_code_blocks() {
    let markdown = r#"```rust
let x = 42;
```

```python
print("hello")
```

```javascript
console.log("test");
```"#;

    let segments = MarkdownProcessor::parse(markdown).expect("Failed to parse");

    let code_blocks: Vec<&Segment> = segments
        .iter()
        .filter(|s| matches!(s, Segment::CodeBlock { .. }))
        .collect();

    assert_eq!(code_blocks.len(), 3, "Should have 3 code blocks");

    // Verify languages
    if let Segment::CodeBlock { language, .. } = &code_blocks[0] {
        assert_eq!(language.as_deref(), Some("rust"));
    }
    if let Segment::CodeBlock { language, .. } = &code_blocks[1] {
        assert_eq!(language.as_deref(), Some("python"));
    }
    if let Segment::CodeBlock { language, .. } = &code_blocks[2] {
        assert_eq!(language.as_deref(), Some("javascript"));
    }
}

/// Test code block without language specifier
#[test]
fn test_code_block_no_language() {
    let markdown = r#"```
plain code block
```"#;

    let segments = MarkdownProcessor::parse(markdown).expect("Failed to parse");

    if let Segment::CodeBlock { language, code } = &segments[0] {
        // pulldown-cmark may assign an empty string for fenced blocks without language
        // Check that it's either None or empty string
        let is_no_language = language.is_none() || language.as_ref().map(|s| s.is_empty()).unwrap_or(false);
        assert!(is_no_language, "Expected no language, got: {:?}", language);
        assert!(code.contains("plain code block"));
    }
}

/// Test multiple horizontal rules
#[test]
fn test_multiple_horizontal_rules() {
    let markdown = "# Section 1\n\n---\n\n# Section 2\n\n---\n\n# Section 3";
    let segments = MarkdownProcessor::parse(markdown).expect("Failed to parse");

    let hr_count = segments
        .iter()
        .filter(|s| matches!(s, Segment::HorizontalRule))
        .count();

    assert_eq!(hr_count, 2, "Should have 2 horizontal rules");
}

/// Test heading levels
#[test]
fn test_heading_levels() {
    let markdown = "# H1\n## H2\n### H3\n#### H4\n##### H5\n###### H6";
    let segments = MarkdownProcessor::parse(markdown).expect("Failed to parse");

    let heading_levels: Vec<u8> = segments
        .iter()
        .filter_map(|s| {
            if let Segment::Heading { level, .. } = s {
                Some(*level)
            } else {
                None
            }
        })
        .collect();

    assert_eq!(heading_levels, vec![1, 2, 3, 4, 5, 6]);
}

/// Test get_text() method
#[test]
fn test_get_text_method() {
    let heading = Segment::Heading {
        level: 1,
        text: "Test Heading".to_string(),
    };
    let paragraph = Segment::Paragraph {
        text: "Test Paragraph".to_string(),
    };
    let code_block = Segment::CodeBlock {
        language: None,
        code: "code".to_string(),
    };

    assert_eq!(heading.get_text(), Some("Test Heading"));
    assert_eq!(paragraph.get_text(), Some("Test Paragraph"));
    assert_eq!(code_block.get_text(), None);
}

/// Test set_text() method
#[test]
fn test_set_text_method() {
    let mut heading = Segment::Heading {
        level: 1,
        text: "Original".to_string(),
    };

    heading.set_text("Modified".to_string());
    assert_eq!(heading.get_text(), Some("Modified"));

    let mut paragraph = Segment::Paragraph {
        text: "Original".to_string(),
    };

    paragraph.set_text("Modified".to_string());
    assert_eq!(paragraph.get_text(), Some("Modified"));
}

/// Test is_translatable() for all segment types
#[test]
fn test_is_translatable() {
    let heading = Segment::Heading { level: 1, text: "Test".to_string() };
    let paragraph = Segment::Paragraph { text: "Test".to_string() };
    let code_block = Segment::CodeBlock { language: None, code: "test".to_string() };
    let hr = Segment::HorizontalRule;

    assert!(heading.is_translatable());
    assert!(paragraph.is_translatable());
    assert!(!code_block.is_translatable());
    assert!(!hr.is_translatable());
}

/// Test complex nested markdown
#[test]
fn test_complex_nested_markdown() {
    let markdown = r#"# Main Title

Introduction paragraph with `inline code`.

```rust
fn example() {
    println!("code");
}
```

## Subsection

More text here.

---

Final paragraph.
"#;

    let segments = MarkdownProcessor::parse(markdown).expect("Failed to parse");

    // Verify we have all expected segment types
    let has_heading = segments.iter().any(|s| matches!(s, Segment::Heading { .. }));
    let has_paragraph = segments.iter().any(|s| matches!(s, Segment::Paragraph { .. }));
    let has_code = segments.iter().any(|s| matches!(s, Segment::CodeBlock { .. }));
    let has_hr = segments.iter().any(|s| matches!(s, Segment::HorizontalRule));

    assert!(has_heading, "Should have headings");
    assert!(has_paragraph, "Should have paragraphs");
    assert!(has_code, "Should have code blocks");
    assert!(has_hr, "Should have horizontal rule");
}