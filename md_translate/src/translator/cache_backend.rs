use crate::error::Result;
use crate::translator::CacheStats;

/// Trait for cache storage backends
/// Allows in-memory mocking in tests
pub trait CacheBackend: Send + Sync {
    /// Retrieve cached translation if exists and valid
    fn get(&mut self, source: &str, model: &str, lang_pair: &str) -> Option<String>;

    /// Store translation in cache
    fn set(&self, source: &str, translation: &str, model: &str, lang_pair: &str) -> Result<()>;

    /// Clear all cache entries
    fn clear(&self) -> Result<()>;

    /// Get cache statistics (returned by value to avoid lifetime issues)
    fn stats(&mut self) -> CacheStats;
}