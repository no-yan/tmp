use crate::markdown::parser::Segment;

pub struct HtmlRenderer;

impl HtmlRenderer {
    pub fn render(segments: &[Segment], title: &str) -> String {
        let mut body = String::new();

        for segment in segments {
            match segment {
                Segment::Heading { level, text } => {
                    body.push_str(&format!("<h{level}>{text}</h{level}>\n"));
                }
                Segment::Paragraph { text } => {
                    body.push_str(&format!("<p>{}</p>\n", Self::escape_html(text)));
                }
                Segment::CodeBlock { language, code } => {
                    let lang_class = language
                        .as_ref()
                        .map(|l| format!(" class=\"language-{}\"", l))
                        .unwrap_or_default();
                    body.push_str(&format!(
                        "<pre><code{}>{}</code></pre>\n",
                        lang_class,
                        Self::escape_html(code)
                    ));
                }
                Segment::HorizontalRule => {
                    body.push_str("<hr>\n");
                }
                _ => {}
            }
        }

        Self::wrap_with_template(&body, title)
    }

    fn escape_html(text: &str) -> String {
        text.replace('&', "&amp;")
            .replace('<', "&lt;")
            .replace('>', "&gt;")
            .replace('"', "&quot;")
            .replace('\'', "&#39;")
    }

    fn wrap_with_template(body: &str, title: &str) -> String {
        format!(
            r#"<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{title}</title>
    <style>
        body {{
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Noto Sans", Helvetica, Arial, sans-serif;
            line-height: 1.6;
            color: #24292f;
            background-color: #ffffff;
            max-width: 980px;
            margin: 0 auto;
            padding: 45px;
        }}

        h1, h2, h3, h4, h5, h6 {{
            margin-top: 24px;
            margin-bottom: 16px;
            font-weight: 600;
            line-height: 1.25;
            border-bottom: 1px solid #d8dee4;
            padding-bottom: 0.3em;
        }}

        h1 {{ font-size: 2em; }}
        h2 {{ font-size: 1.5em; }}
        h3 {{ font-size: 1.25em; }}

        p {{
            margin-top: 0;
            margin-bottom: 16px;
        }}

        pre {{
            background-color: #f6f8fa;
            border-radius: 6px;
            padding: 16px;
            overflow: auto;
            font-size: 85%;
            line-height: 1.45;
        }}

        code {{
            background-color: rgba(175,184,193,0.2);
            padding: 0.2em 0.4em;
            margin: 0;
            font-size: 85%;
            border-radius: 6px;
            font-family: ui-monospace, SFMono-Regular, "SF Mono", Menlo, Consolas, monospace;
        }}

        pre code {{
            background-color: transparent;
            padding: 0;
        }}

        hr {{
            height: 0.25em;
            padding: 0;
            margin: 24px 0;
            background-color: #d8dee4;
            border: 0;
        }}
    </style>
</head>
<body>
{body}
</body>
</html>"#,
            title = Self::escape_html(title),
            body = body
        )
    }
}
