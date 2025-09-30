mod mocks;

use mocks::InMemoryCache;
use md_translate::translator::CacheBackend;

#[test]
fn test_cache_hit_miss() {
    let mut cache = InMemoryCache::new();

    // Miss on first access
    let result = cache.get("hello", "model", "en-ja");
    assert!(result.is_none());

    let stats = cache.stats();
    assert_eq!(stats.cache_misses, 1);
    assert_eq!(stats.total_requests, 1);

    // Set value
    cache.set("hello", "こんにちは", "model", "en-ja").expect("Set failed");

    // Hit on second access
    let result = cache.get("hello", "model", "en-ja");
    assert_eq!(result, Some("こんにちは".to_string()));

    let stats = cache.stats();
    assert_eq!(stats.cache_hits, 1);
    assert_eq!(stats.total_requests, 2);
}

#[test]
fn test_cache_key_uniqueness() {
    let mut cache = InMemoryCache::new();

    // Different models should have different keys
    cache.set("text", "trans1", "model1", "en-ja").expect("Set 1 failed");
    cache.set("text", "trans2", "model2", "en-ja").expect("Set 2 failed");

    assert_eq!(cache.get("text", "model1", "en-ja"), Some("trans1".to_string()));
    assert_eq!(cache.get("text", "model2", "en-ja"), Some("trans2".to_string()));

    // Different language pairs should have different keys
    cache.set("text", "trans3", "model1", "en-es").expect("Set 3 failed");

    assert_eq!(cache.get("text", "model1", "en-ja"), Some("trans1".to_string()));
    assert_eq!(cache.get("text", "model1", "en-es"), Some("trans3".to_string()));
}

#[test]
fn test_cache_clear() {
    let mut cache = InMemoryCache::new();

    cache.set("key1", "value1", "model", "en-ja").expect("Set 1 failed");
    cache.set("key2", "value2", "model", "en-ja").expect("Set 2 failed");

    cache.clear().expect("Clear failed");

    assert!(cache.get("key1", "model", "en-ja").is_none());
    assert!(cache.get("key2", "model", "en-ja").is_none());
}

#[test]
fn test_cache_stats_tracking() {
    let mut cache = InMemoryCache::new();

    cache.get("key1", "model", "en-ja");  // Miss
    cache.set("key1", "value1", "model", "en-ja").expect("Set failed");
    cache.get("key1", "model", "en-ja");  // Hit
    cache.get("key2", "model", "en-ja");  // Miss

    let stats = cache.stats();
    assert_eq!(stats.total_requests, 3);
    assert_eq!(stats.cache_hits, 1);
    assert_eq!(stats.cache_misses, 2);
}

#[test]
fn test_cache_overwrites() {
    let mut cache = InMemoryCache::new();

    cache.set("key", "value1", "model", "en-ja").expect("Set 1 failed");
    cache.set("key", "value2", "model", "en-ja").expect("Set 2 failed");

    // Second set should overwrite first
    assert_eq!(cache.get("key", "model", "en-ja"), Some("value2".to_string()));
}

#[test]
fn test_cache_isolation() {
    let mut cache1 = InMemoryCache::new();
    let mut cache2 = InMemoryCache::new();

    cache1.set("key", "value1", "model", "en-ja").expect("Set 1 failed");
    cache2.set("key", "value2", "model", "en-ja").expect("Set 2 failed");

    // Caches should be independent
    assert_eq!(cache1.get("key", "model", "en-ja"), Some("value1".to_string()));
    assert_eq!(cache2.get("key", "model", "en-ja"), Some("value2".to_string()));
}