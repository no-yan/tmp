use crate::cli::*;
use crate::error::Result;
use crate::translator::{Translator, TranslatorConfig};
use colored::Colorize;
use std::path::Path;

pub async fn handle_translate(args: TranslateArgs) -> Result<()> {
    let input = std::fs::read_to_string(&args.file)?;

    let config = TranslatorConfig {
        use_cache: !args.no_cache,
        show_progress: true,
        ..Default::default()
    };

    let mut translator = Translator::new(args.ollama_url, args.model, config)?;

    println!("{}", "Translating markdown...".cyan());
    let translated = translator.translate_markdown(&input).await?;
    dbg!(translated.clone());

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

pub async fn handle_cache(args: CacheArgs) -> Result<()> {
    match args.command {
        CacheCommands::Stats => {
            let config = TranslatorConfig::default();
            let mut translator = Translator::new(
                "http://localhost:11434".to_string(),
                "qwen2.5:7b".to_string(),
                config,
            )?;

            let stats = translator.cache_stats();
            println!("{}", "Cache Statistics:".cyan().bold());
            println!("  Total Requests: {}", stats.total_requests);
            println!("  Cache Hits: {}", stats.cache_hits);
            println!("  Cache Misses: {}", stats.cache_misses);
            println!(
                "  Total Size: {} bytes ({:.2} MB)",
                stats.total_size_bytes,
                stats.total_size_bytes as f64 / 1_048_576.0
            );
        }
        CacheCommands::Clear => {
            let config = TranslatorConfig::default();
            let translator = Translator::new(
                "http://localhost:11434".to_string(),
                "qwen2.5:7b".to_string(),
                config,
            )?;

            translator.clear_cache()?;
            println!("{}", "Cache cleared successfully".green());
        }
    }

    Ok(())
}

pub async fn handle_view(args: ViewArgs) -> Result<()> {
    let input = std::fs::read_to_string(&args.file)?;
    let title = Path::new(&args.file)
        .file_stem()
        .and_then(|s| s.to_str())
        .unwrap_or("Translated Document");

    let config = TranslatorConfig {
        show_progress: true,
        ..Default::default()
    };

    let mut translator = Translator::new(args.ollama_url, args.model, config)?;

    println!("{}", "Translating and generating HTML...".cyan());
    let html = translator.translate_to_html(&input, title).await?;

    // Write to temp file
    let temp_dir = std::env::temp_dir();
    let temp_file = temp_dir.join("md_translate_preview.html");
    std::fs::write(&temp_file, html)?;

    println!("{} {}", "Preview saved to:".green(), temp_file.display());

    // Open in browser
    if let Err(e) = open::that(&temp_file) {
        eprintln!("{} {}", "Failed to open browser:".red(), e);
        println!("{}", "Please open the file manually.".yellow());
    } else {
        println!("{}", "Opening in default browser...".cyan());
    }

    Ok(())
}
