pub mod commands;

use clap::{Parser, Subcommand};

#[derive(Parser)]
#[command(name = "md-translate")]
#[command(about = "Translate markdown files using local Ollama LLM")]
#[command(version)]
pub struct Cli {
    #[command(subcommand)]
    pub command: Commands,
}

#[derive(Subcommand)]
pub enum Commands {
    /// Translate a markdown file
    Translate(TranslateArgs),

    /// View translated markdown in browser
    View(ViewArgs),

    /// Start development server with live reload
    #[cfg(feature = "server")]
    Serve(ServeArgs),

    /// Manage translation cache
    Cache(CacheArgs),
}

#[derive(Parser)]
pub struct TranslateArgs {
    /// Input markdown file
    pub file: String,

    /// Output file (defaults to stdout)
    #[arg(short, long)]
    pub output: Option<String>,

    /// Ollama model name
    #[arg(short, long, default_value = "qwen2.5:7b")]
    pub model: String,

    /// Ollama API URL
    #[arg(long, default_value = "http://localhost:11434")]
    pub ollama_url: String,

    /// Disable cache
    #[arg(long)]
    pub no_cache: bool,

    /// Output format (markdown or html)
    #[arg(long, default_value = "markdown")]
    pub format: String,
}

#[derive(Parser)]
pub struct ViewArgs {
    /// Input markdown file
    pub file: String,

    /// Ollama model name
    #[arg(short, long, default_value = "qwen2.5:7b")]
    pub model: String,

    /// Ollama API URL
    #[arg(long, default_value = "http://localhost:11434")]
    pub ollama_url: String,
}

#[cfg(feature = "server")]
#[derive(Parser)]
pub struct ServeArgs {
    /// Directory to serve (defaults to current directory)
    pub dir: Option<String>,

    /// Port number
    #[arg(short, long, default_value = "3000")]
    pub port: u16,

    /// Enable file watching
    #[arg(short, long)]
    pub watch: bool,

    /// Ollama model name
    #[arg(short, long, default_value = "qwen2.5:7b")]
    pub model: String,

    /// Ollama API URL
    #[arg(long, default_value = "http://localhost:11434")]
    pub ollama_url: String,
}

#[derive(Parser)]
pub struct CacheArgs {
    #[command(subcommand)]
    pub command: CacheCommands,
}

#[derive(Subcommand)]
pub enum CacheCommands {
    /// Show cache statistics
    Stats,

    /// Clear all cache entries
    Clear,
}
