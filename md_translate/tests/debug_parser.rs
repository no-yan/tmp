use md_translate::markdown::parser::{MarkdownProcessor, Segment};

#[test]
fn debug_parser_output() {
    let test_markdown = r#"# Heading 1

First paragraph here.

```rust
fn main() {
    println!("Code block 1");
}
```

# Heading 2

Second paragraph here.
"#;

    let segments = MarkdownProcessor::parse(test_markdown).expect("Failed to parse markdown");

    println!("Total segments: {}", segments.len());
    for (i, segment) in segments.iter().enumerate() {
        match segment {
            Segment::Heading { level, text } => {
                println!("{}: Heading({}) - '{}'", i, level, text);
            }
            Segment::Paragraph { text } => {
                println!("{}: Paragraph - '{}'", i, &text[..text.len().min(50)]);
            }
            Segment::CodeBlock { language, code } => {
                println!("{}: CodeBlock({:?}) - {} bytes", i, language, code.len());
            }
            Segment::HorizontalRule => {
                println!("{}: HorizontalRule", i);
            }
            _ => {
                println!("{}: Other", i);
            }
        }
    }
}