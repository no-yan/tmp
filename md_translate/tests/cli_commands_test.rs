mod mocks;

use md_translate::cli::commands::{handle_cache_with_translator, handle_translate_with_translator};
use md_translate::cli::{CacheArgs, CacheCommands, TranslateArgs};
use md_translate::translator::{Translator, TranslatorConfig};
use mocks::{InMemoryCache, MockTranslationService};
use std::fs;
use tempfile::tempdir;

#[tokio::test]
async fn test_translate_command_basic() {
    let temp_dir = tempdir().expect("Failed to create temp dir");
    let input_path = temp_dir.path().join("input.md");
    let output_path = temp_dir.path().join("output.md");

    // Write test input
    fs::write(&input_path, "# Test\n\nParagraph.").expect("Failed to write input");

    // Create mock translator
    let mock_client = MockTranslationService::default();
    let mock_cache = InMemoryCache::new();
    let config = TranslatorConfig {
        use_cache: false,
        show_progress: false,
        ..Default::default()
    };
    let mut translator = Translator::new_with_deps(mock_client, mock_cache, config);

    // Create args
    let args = TranslateArgs {
        file: input_path.to_string_lossy().to_string(),
        output: Some(output_path.to_string_lossy().to_string()),
        model: "mock".to_string(),
        ollama_url: "http://mock".to_string(),
        no_cache: true,
        format: "markdown".to_string(),
    };

    // Execute command
    handle_translate_with_translator(&mut translator, &args)
        .await
        .expect("Command failed");

    // Verify output file was created
    assert!(output_path.exists(), "Output file should exist");

    let output = fs::read_to_string(&output_path).expect("Failed to read output");
    assert!(
        output.contains("[TRANSLATED:"),
        "Should contain mock translation marker"
    );
}

#[tokio::test]
async fn test_translate_command_missing_input() {
    let mock_client = MockTranslationService::default();
    let mock_cache = InMemoryCache::new();
    let config = TranslatorConfig::default();
    let mut translator = Translator::new_with_deps(mock_client, mock_cache, config);

    let args = TranslateArgs {
        file: "/nonexistent/file.md".to_string(),
        output: None,
        model: "mock".to_string(),
        ollama_url: "http://mock".to_string(),
        no_cache: true,
        format: "markdown".to_string(),
    };

    // Should return IoError
    let result = handle_translate_with_translator(&mut translator, &args).await;
    assert!(result.is_err());
}

#[tokio::test]
async fn test_cache_command_stats() {
    let mock_client = MockTranslationService::default();
    let mock_cache = InMemoryCache::new();
    let config = TranslatorConfig {
        use_cache: true,
        show_progress: false,
        ..Default::default()
    };
    let mut translator = Translator::new_with_deps(mock_client, mock_cache, config);

    // Populate cache
    translator
        .translate_markdown("# Test")
        .await
        .expect("Translation failed");

    // Test stats command
    let args = CacheArgs {
        command: CacheCommands::Stats,
    };

    handle_cache_with_translator(&mut translator, &args)
        .await
        .expect("Command failed");

    // Verify stats were populated
    let stats = translator.cache_stats();
    assert!(stats.total_requests > 0, "Should have cache requests");
}

#[tokio::test]
async fn test_cache_command_clear() {
    let mock_client = MockTranslationService::default();
    let mock_cache = InMemoryCache::new();
    let config = TranslatorConfig {
        use_cache: true,
        show_progress: false,
        ..Default::default()
    };
    let mut translator = Translator::new_with_deps(mock_client, mock_cache, config);

    // Populate cache
    translator
        .translate_markdown("# Test")
        .await
        .expect("Translation failed");

    // Verify cache has data
    let stats_before = translator.cache_stats();
    assert_eq!(stats_before.cache_misses, 1, "Should have 1 miss");

    // Test clear command
    let args = CacheArgs {
        command: CacheCommands::Clear,
    };

    handle_cache_with_translator(&mut translator, &args)
        .await
        .expect("Command failed");

    // Verify cache is cleared
    let stats_after = translator.cache_stats();
    assert_eq!(
        stats_after.cache_hits, 0,
        "Cache hits should be reset after clear"
    );
    assert_eq!(
        stats_after.cache_misses, 0,
        "Cache misses should be reset after clear"
    );
}

#[tokio::test]
async fn test_translate_with_cache_enabled() {
    let temp_dir = tempdir().expect("Failed to create temp dir");
    let input_path = temp_dir.path().join("cached_input.md");
    let output_path = temp_dir.path().join("cached_output.md");

    // Write test input
    fs::write(&input_path, "# Cached Test").expect("Failed to write input");

    // Create mock translator with cache enabled
    let mock_client = MockTranslationService::default();
    let mock_cache = InMemoryCache::new();
    let config = TranslatorConfig {
        use_cache: true,
        show_progress: false,
        ..Default::default()
    };
    let mut translator = Translator::new_with_deps(mock_client, mock_cache, config);

    // Create args
    let args = TranslateArgs {
        file: input_path.to_string_lossy().to_string(),
        output: Some(output_path.to_string_lossy().to_string()),
        model: "mock".to_string(),
        ollama_url: "http://mock".to_string(),
        no_cache: false,
        format: "markdown".to_string(),
    };

    // Execute command twice
    handle_translate_with_translator(&mut translator, &args)
        .await
        .expect("First command failed");

    // Second execution should hit cache
    handle_translate_with_translator(&mut translator, &args)
        .await
        .expect("Second command failed");

    // Verify cache was used
    let stats = translator.cache_stats();
    assert!(stats.cache_hits > 0, "Should have cache hits on second run");
}
