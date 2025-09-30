use md_translate::error::Result;
use md_translate::translator::{CacheBackend, CacheStats};
use std::collections::HashMap;
use std::sync::{Arc, Mutex};
use sha2::{Digest, Sha256};

/// In-memory cache implementation for testing
/// Mimics TranslationCache behavior without filesystem I/O
#[derive(Clone)]
pub struct InMemoryCache {
    store: Arc<Mutex<HashMap<String, String>>>,
    stats: Arc<Mutex<CacheStats>>,
}

impl InMemoryCache {
    pub fn new() -> Self {
        Self {
            store: Arc::new(Mutex::new(HashMap::new())),
            stats: Arc::new(Mutex::new(CacheStats::default())),
        }
    }

    fn generate_key(source: &str, model: &str, lang_pair: &str) -> String {
        let combined = format!("{}{}{}", source, model, lang_pair);
        let mut hasher = Sha256::new();
        hasher.update(combined.as_bytes());
        format!("{:x}", hasher.finalize())
    }
}

impl Default for InMemoryCache {
    fn default() -> Self {
        Self::new()
    }
}

impl CacheBackend for InMemoryCache {
    fn get(&mut self, source: &str, model: &str, lang_pair: &str) -> Option<String> {
        let mut stats = self.stats.lock().unwrap();
        stats.total_requests += 1;

        let key = Self::generate_key(source, model, lang_pair);
        let store = self.store.lock().unwrap();

        if let Some(value) = store.get(&key) {
            stats.cache_hits += 1;
            Some(value.clone())
        } else {
            stats.cache_misses += 1;
            None
        }
    }

    fn set(&self, source: &str, translation: &str, model: &str, lang_pair: &str) -> Result<()> {
        let key = Self::generate_key(source, model, lang_pair);
        let mut store = self.store.lock().unwrap();
        store.insert(key, translation.to_string());
        Ok(())
    }

    fn clear(&self) -> Result<()> {
        let mut store = self.store.lock().unwrap();
        store.clear();

        let mut stats = self.stats.lock().unwrap();
        *stats = CacheStats::default();
        Ok(())
    }

    fn stats(&mut self) -> CacheStats {
        let stats = self.stats.lock().unwrap();
        *stats
    }
}