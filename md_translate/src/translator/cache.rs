use crate::error::Result;
use crate::translator::cache_backend::CacheBackend;
use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sha2::{Digest, Sha256};
use std::path::PathBuf;

#[derive(Serialize, Deserialize)]
struct CacheEntry {
    source: String,
    translation: String,
    model: String,
    language_pair: String,
    created_at: DateTime<Utc>,
    checksum: String,
}

#[derive(Default, Debug, Clone, Copy)]
pub struct CacheStats {
    pub total_requests: u64,
    pub cache_hits: u64,
    pub cache_misses: u64,
    pub total_size_bytes: u64,
}

pub struct TranslationCache {
    cache_dir: PathBuf,
    stats: CacheStats,
}

impl TranslationCache {
    pub fn new() -> Result<Self> {
        let cache_dir = dirs::cache_dir()
            .ok_or_else(|| {
                std::io::Error::new(
                    std::io::ErrorKind::NotFound,
                    "Could not find cache directory",
                )
            })?
            .join("md_translate")
            .join("translations");

        std::fs::create_dir_all(&cache_dir)?;

        Ok(Self {
            cache_dir,
            stats: CacheStats::default(),
        })
    }

    pub fn get(&mut self, source: &str, model: &str, lang_pair: &str) -> Option<String> {
        self.stats.total_requests += 1;

        let key = Self::generate_key(source, model, lang_pair);
        let cache_path = self.cache_dir.join(&key).with_extension("json");

        if let Ok(content) = std::fs::read_to_string(&cache_path) {
            if let Ok(entry) = serde_json::from_str::<CacheEntry>(&content) {
                // Verify checksum
                if entry.checksum == Self::hash_text(source) {
                    self.stats.cache_hits += 1;
                    return Some(entry.translation);
                }
            }
        }

        self.stats.cache_misses += 1;
        None
    }

    pub fn set(&self, source: &str, translation: &str, model: &str, lang_pair: &str) -> Result<()> {
        let key = Self::generate_key(source, model, lang_pair);
        let cache_path = self.cache_dir.join(&key).with_extension("json");

        let entry = CacheEntry {
            source: source.to_string(),
            translation: translation.to_string(),
            model: model.to_string(),
            language_pair: lang_pair.to_string(),
            created_at: Utc::now(),
            checksum: Self::hash_text(source),
        };

        let json = serde_json::to_string_pretty(&entry)?;
        std::fs::write(cache_path, json)?;

        Ok(())
    }

    pub fn clear(&self) -> Result<()> {
        for entry in std::fs::read_dir(&self.cache_dir)? {
            let entry = entry?;
            if entry.path().extension().and_then(|s| s.to_str()) == Some("json") {
                std::fs::remove_file(entry.path())?;
            }
        }
        Ok(())
    }

    pub fn stats(&mut self) -> &CacheStats {
        // Calculate total size
        let mut total_size = 0u64;
        if let Ok(entries) = std::fs::read_dir(&self.cache_dir) {
            for entry in entries.flatten() {
                if let Ok(metadata) = entry.metadata() {
                    total_size += metadata.len();
                }
            }
        }
        self.stats.total_size_bytes = total_size;

        &self.stats
    }

    fn generate_key(source: &str, model: &str, lang_pair: &str) -> String {
        let combined = format!("{}{}{}", source, model, lang_pair);
        Self::hash_text(&combined)
    }

    fn hash_text(text: &str) -> String {
        let mut hasher = Sha256::new();
        hasher.update(text.as_bytes());
        format!("{:x}", hasher.finalize())
    }
}

impl CacheBackend for TranslationCache {
    fn get(&mut self, source: &str, model: &str, lang_pair: &str) -> Option<String> {
        self.stats.total_requests += 1;

        let key = Self::generate_key(source, model, lang_pair);
        let cache_path = self.cache_dir.join(&key).with_extension("json");

        if let Ok(content) = std::fs::read_to_string(&cache_path) {
            if let Ok(entry) = serde_json::from_str::<CacheEntry>(&content) {
                // Verify checksum
                if entry.checksum == Self::hash_text(source) {
                    self.stats.cache_hits += 1;
                    return Some(entry.translation);
                }
            }
        }

        self.stats.cache_misses += 1;
        None
    }

    fn set(&self, source: &str, translation: &str, model: &str, lang_pair: &str) -> Result<()> {
        let key = Self::generate_key(source, model, lang_pair);
        let cache_path = self.cache_dir.join(&key).with_extension("json");

        let entry = CacheEntry {
            source: source.to_string(),
            translation: translation.to_string(),
            model: model.to_string(),
            language_pair: lang_pair.to_string(),
            created_at: Utc::now(),
            checksum: Self::hash_text(source),
        };

        let json = serde_json::to_string_pretty(&entry)?;
        std::fs::write(cache_path, json)?;

        Ok(())
    }

    fn clear(&self) -> Result<()> {
        for entry in std::fs::read_dir(&self.cache_dir)? {
            let entry = entry?;
            if entry.path().extension().and_then(|s| s.to_str()) == Some("json") {
                std::fs::remove_file(entry.path())?;
            }
        }
        Ok(())
    }

    fn stats(&mut self) -> CacheStats {
        // Calculate total size
        let mut total_size = 0u64;
        if let Ok(entries) = std::fs::read_dir(&self.cache_dir) {
            for entry in entries.flatten() {
                if let Ok(metadata) = entry.metadata() {
                    total_size += metadata.len();
                }
            }
        }
        self.stats.total_size_bytes = total_size;

        self.stats
    }
}
