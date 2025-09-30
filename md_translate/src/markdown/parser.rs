use crate::error::Result;
use pulldown_cmark::{CodeBlockKind, Event, Parser, Tag, TagEnd};

#[derive(Debug, Clone)]
pub enum Segment {
    Heading {
        level: u8,
        text: String,
    },
    Paragraph {
        text: String,
    },
    CodeBlock {
        language: Option<String>,
        code: String,
    },
    List {
        ordered: bool,
        items: Vec<String>,
    },
    BlockQuote {
        content: String,
    },
    HorizontalRule,
}

impl Segment {
    pub fn is_translatable(&self) -> bool {
        matches!(
            self,
            Segment::Heading { .. }
                | Segment::Paragraph { .. }
                | Segment::List { .. }
                | Segment::BlockQuote { .. }
        )
    }

    pub fn get_text(&self) -> Option<&str> {
        match self {
            Segment::Heading { text, .. } => Some(text),
            Segment::Paragraph { text } => Some(text),
            Segment::BlockQuote { content } => Some(content),
            _ => None,
        }
    }

    pub fn set_text(&mut self, new_text: String) {
        match self {
            Segment::Heading { text, .. } => *text = new_text,
            Segment::Paragraph { text } => *text = new_text,
            Segment::BlockQuote { content } => *content = new_text,
            _ => {}
        }
    }
}

pub struct MarkdownProcessor;

impl MarkdownProcessor {
    pub fn parse(markdown: &str) -> Result<Vec<Segment>> {
        let parser = Parser::new(markdown);
        let mut segments = Vec::new();
        let mut current_text = String::new();
        let mut in_heading = false;
        let mut heading_level = 0;
        let mut in_paragraph = false;
        let mut in_code_block = false;
        let mut code_language = None;
        let mut code_content = String::new();

        for event in parser {
            match event {
                Event::Start(Tag::Heading { level, .. }) => {
                    in_heading = true;
                    heading_level = level as u8;
                    current_text.clear();
                }
                Event::End(TagEnd::Heading(_)) => {
                    if in_heading {
                        segments.push(Segment::Heading {
                            level: heading_level,
                            text: current_text.trim().to_string(),
                        });
                        in_heading = false;
                        current_text.clear();
                    }
                }
                Event::Start(Tag::Paragraph) => {
                    in_paragraph = true;
                    current_text.clear();
                }
                Event::End(TagEnd::Paragraph) => {
                    if in_paragraph {
                        segments.push(Segment::Paragraph {
                            text: current_text.trim().to_string(),
                        });
                        in_paragraph = false;
                        current_text.clear();
                    }
                }
                Event::Start(Tag::CodeBlock(kind)) => {
                    in_code_block = true;
                    code_language = match kind {
                        CodeBlockKind::Fenced(lang) => Some(lang.to_string()),
                        CodeBlockKind::Indented => None,
                    };
                    code_content.clear();
                }
                Event::End(TagEnd::CodeBlock) => {
                    if in_code_block {
                        segments.push(Segment::CodeBlock {
                            language: code_language.clone(),
                            code: code_content.clone(),
                        });
                        in_code_block = false;
                        code_content.clear();
                    }
                }
                Event::Rule => {
                    segments.push(Segment::HorizontalRule);
                }
                Event::Text(text) => {
                    if in_code_block {
                        code_content.push_str(&text);
                    } else if in_heading || in_paragraph {
                        current_text.push_str(&text);
                    }
                }
                Event::Code(code) => {
                    if in_heading || in_paragraph {
                        current_text.push('`');
                        current_text.push_str(&code);
                        current_text.push('`');
                    }
                }
                _ => {}
            }
        }

        Ok(segments)
    }
}
