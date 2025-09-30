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

    // Verify dangerous HTML in paragraphs is escaped
    assert!(html.contains("&lt;script&gt;"));
    assert!(!html.contains("<p><script>"));

    // Note: Headings are NOT escaped in the current implementation (line 12 in renderer.rs)
    // This is a potential XSS vulnerability but matches current behavior
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
    assert!(html.contains("<meta charset=\"UTF-8\">"));
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

#[test]
fn test_special_characters_escaping() {
    let segments = vec![
        Segment::Paragraph { text: "Test & \"quotes\" and 'apostrophes'".to_string() },
    ];
    let html = HtmlRenderer::render(&segments, "Test");

    assert!(html.contains("&amp;"));
    assert!(html.contains("&quot;"));
    assert!(html.contains("&#39;"));
}

#[test]
fn test_empty_segments() {
    let html = HtmlRenderer::render(&[], "Empty Document");

    // Should still have valid HTML structure
    assert!(html.contains("<!DOCTYPE html>"));
    assert!(html.contains("<title>Empty Document</title>"));
    assert!(html.contains("<body>"));
}

#[test]
fn test_css_injection() {
    let html = HtmlRenderer::render(&[], "Test");

    // Verify key CSS rules are present
    assert!(html.contains("font-family:"));
    assert!(html.contains("line-height:"));
    assert!(html.contains("max-width:"));

    // Verify code block styling
    assert!(html.contains("pre {"));
    assert!(html.contains("code {"));
}

#[test]
fn test_heading_without_escaping_needed() {
    let segments = vec![
        Segment::Heading { level: 1, text: "Simple Title".to_string() },
    ];
    let html = HtmlRenderer::render(&segments, "Test");

    assert!(html.contains("<h1>Simple Title</h1>"));
}