use clap::Parser;
use md_translate::cli::commands;
use md_translate::cli::{Cli, Commands};

#[tokio::main]
async fn main() {
    let cli = Cli::parse();

    let result = match cli.command {
        Commands::Translate(args) => commands::handle_translate(args).await,
        Commands::View(args) => commands::handle_view(args).await,
        Commands::Cache(args) => commands::handle_cache(args).await,
        _ => {
            eprintln!("Command not yet implemented");
            std::process::exit(1);
        }
    };

    if let Err(e) = result {
        eprintln!("Error: {}", e);
        std::process::exit(1);
    }
}
