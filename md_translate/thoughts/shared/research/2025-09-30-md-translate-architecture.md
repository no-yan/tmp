---
date: 2025-09-30T07:36:51+0000
researcher: noyan
git_commit: 814e683cb61249518b3b9e8fc2f29e26356bc3e4
branch: main
repository: md_translate
topic: "Rust CLI + Dev Server Architecture for Markdown Translation Tool"
tags: [research, architecture, rust, markdown, translation, ollama, llm]
status: complete
last_updated: 2025-09-30
last_updated_by: noyan
---

# Research: Rust CLI + Dev Server Architecture for Markdown Translation Tool

**Date**: 2025-09-30T07:36:51+0000
**Researcher**: noyan
**Git Commit**: 814e683cb61249518b3b9e8fc2f29e26356bc3e4
**Branch**: main
**Repository**: md_translate

## Research Question

英語で書かれたマークダウンファイルを日本語で閲覧するためのツールの設計。目的は、トークン量節約のために英語で書かれたマークダウンを人間がチェックするときに翻訳する手間を減らしたり、直接読む認知コストを下げること。

技術選定と設計を行い、Rust CLI + Dev Serverの構成でアーキテクチャを完成させる。

## Summary

Rust製のマークダウン翻訳ツール `md_translate` の完全なアーキテクチャ設計を完了。以下の主要コンポーネントで構成：

- **CLI**: `translate`, `view`, `serve`, `cache` コマンド
- **Translator Engine**: Ollama統合、セグメント単位キャッシュ、並列翻訳
- **Markdown Processor**: pulldown-cmarkによるパース、セグメント分割、HTML生成
- **Dev Server**: Axum HTTPサーバー、WebSocket、ファイル監視による自動リロード

主な特徴：
- ローカルLLM (Ollama + Qwen2.5:7b) による翻訳
- セグメント単位のキャッシュで部分変更に強い
- 最大3並列の非同期翻訳処理
- リアルタイムファイル監視とブラウザ自動リロード

## Detailed Findings

### 1. Module Structure

プロジェクト構造：

```
md_translate/
├── Cargo.toml
├── src/
│   ├── lib.rs              # ライブラリエントリーポイント
│   ├── main.rs             # CLIエントリーポイント
│   │
│   ├── translator/         # 翻訳エンジン
│   │   ├── mod.rs          # モジュール定義、Translator構造体
│   │   ├── ollama.rs       # Ollama APIクライアント
│   │   └── cache.rs        # 翻訳キャッシュ管理
│   │
│   ├── markdown/           # マークダウン処理
│   │   ├── mod.rs          # モジュール定義
│   │   ├── parser.rs       # マークダウンパース・セグメント分割
│   │   └── renderer.rs     # HTML生成
│   │
│   ├── cli/                # CLI機能
│   │   ├── mod.rs          # モジュール定義
│   │   └── commands.rs     # サブコマンド実装
│   │
│   └── server/             # Dev Server (feature gated)
│       ├── mod.rs          # サーバーエントリー
│       ├── routes.rs       # HTTPルート定義
│       ├── watcher.rs      # ファイル監視
│       └── ws.rs           # WebSocket通信
│
└── templates/
    └── preview.html        # プレビュー用HTMLテンプレート
```

**依存関係の設計**: Cargo.tomlでサーバー機能をfeature gateすることで、CLIのみのビルドも可能。

### 2. Ollama Client Design (`translator/ollama.rs`)

**OllamaClient構造体**:
```rust
pub struct OllamaClient {
    base_url: String,           // デフォルト: http://localhost:11434
    model: String,              // デフォルト: qwen2.5:7b
    client: reqwest::Client,    // HTTP client
    timeout: Duration,          // リクエストタイムアウト
}
```

**API通信**:
- エンドポイント: `POST /api/generate`
- ストリーミング非対応（シンプル性優先）
- Temperature: 0.3（一貫性のある翻訳）

**プロンプト設計**:
```
System: You are a professional translator. Translate the following English markdown text to Japanese.
Keep markdown formatting intact. Do not translate code blocks, URLs, or technical terms.

User: [English text]
```

**エラーハンドリング**:
- 接続エラー時は明確なメッセージ
- リトライ戦略: 最大3回、exponential backoff

### 3. Translation Cache Architecture (`translator/cache.rs`)

**キャッシュ保存場所**: `~/.cache/md_translate/translations/`

**キャッシュキー生成**:
```
SHA256(source_text + model_name + language_pair)
```

**キャッシュエントリ構造**:
```json
{
  "source": "Original English text",
  "translation": "翻訳されたテキスト",
  "model": "qwen2.5:7b",
  "language_pair": "en-ja",
  "created_at": "2025-09-30T12:00:00Z",
  "checksum": "sha256_of_source"
}
```

**TranslationCache構造体**:
```rust
pub struct TranslationCache {
    cache_dir: PathBuf,
    stats: CacheStats,
}

pub struct CacheStats {
    total_requests: u64,
    cache_hits: u64,
    cache_misses: u64,
    total_size_bytes: u64,
}
```

**キャッシュ戦略**:
- TTL（Time To Live）なし（手動削除のみ）
- LRU削除なし（ディスク容量十分と仮定）
- 並行アクセス: ファイルロック不要

### 4. Markdown Parsing & Segmentation (`markdown/parser.rs`)

**セグメント型定義**:
```rust
pub enum Segment {
    Heading {
        level: u8,        // 1-6
        text: String,
    },
    Paragraph {
        text: String,
    },
    CodeBlock {
        language: Option<String>,
        code: String,
        // 翻訳しない
    },
    List {
        ordered: bool,
        items: Vec<String>,
    },
    BlockQuote {
        content: String,
    },
    Table {
        headers: Vec<String>,
        rows: Vec<Vec<String>>,
    },
    HorizontalRule,
}
```

**パース処理フロー**:
1. pulldown-cmarkでマークダウンをイベントストリームに変換
2. イベントをセグメント単位にグループ化
3. 各セグメントに一意のIDを付与（位置ベース）
4. 翻訳対象セグメントと除外セグメントを分類

**翻訳対象**:
- Heading、Paragraph、List、BlockQuote、Table

**翻訳除外**:
- CodeBlock、URL、技術用語（設定可能）

**MarkdownProcessor構造体**:
```rust
pub struct MarkdownProcessor {
    preserve_urls: bool,
    preserve_code: bool,
    technical_terms: Vec<String>,
}
```

### 5. Translator Core Logic (`translator/mod.rs`)

**Translator構造体**:
```rust
pub struct Translator {
    ollama_client: OllamaClient,
    cache: TranslationCache,
    markdown_processor: MarkdownProcessor,
    config: TranslatorConfig,
}

pub struct TranslatorConfig {
    source_lang: String,          // "en"
    target_lang: String,          // "ja"
    use_cache: bool,
    parallel_requests: usize,     // デフォルト: 3
    show_progress: bool,
}
```

**翻訳処理フロー**:
```
Markdown File
    ↓
Parse & Segment (MarkdownProcessor)
    ↓
Filter Segments (翻訳対象のみ)
    ↓
Check Cache
    ├─ Hit → 既存翻訳使用
    └─ Miss → Ollama API呼び出し
              ↓
         Translate (LLM)
              ↓
         Save to Cache
    ↓
Reconstruct Markdown/HTML
    ↓
Output
```

**並列翻訳設計**:
- Tokioの`JoinSet`を使用
- 最大3並列リクエスト
- 各セグメントは独立して翻訳可能
- 進捗表示: `indicatif`でプログレスバー

**エラーハンドリング**:
- 個別セグメント翻訳失敗 → 元のテキストを保持
- Ollama接続エラー → 即座に中断

### 6. CLI Interface (`cli/commands.rs`)

**CLIコマンド構造**:
```rust
#[derive(Parser)]
#[command(name = "md-translate")]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

enum Commands {
    Translate(TranslateArgs),
    View(ViewArgs),
    Serve(ServeArgs),
    Cache(CacheArgs),
}
```

**サブコマンド詳細**:

**1. `translate` - ファイル翻訳**
```bash
md-translate translate <FILE> [OPTIONS]

Options:
  -o, --output <FILE>      出力ファイル
  -m, --model <MODEL>      Ollamaモデル名
  --no-cache               キャッシュを使わない
  --format <FORMAT>        markdown|html
  --ollama-url <URL>       Ollama APIのURL

Example:
  md-translate translate README.md -o README-ja.md
```

**2. `view` - ブラウザでプレビュー**
```bash
md-translate view <FILE> [OPTIONS]

動作:
  1. ファイルを翻訳
  2. 一時HTTPサーバー起動
  3. デフォルトブラウザで開く

Example:
  md-translate view README.md
```

**3. `serve` - 開発サーバー起動**
```bash
md-translate serve [DIR] [OPTIONS]

Options:
  -p, --port <PORT>        ポート番号（デフォルト: 3000）
  -w, --watch              ファイル監視・自動リロード
  -m, --model <MODEL>      Ollamaモデル名

動作:
  1. ディレクトリ内のマークダウンファイルを監視
  2. HTTPサーバー起動
  3. ファイル変更時に自動翻訳・リロード

Example:
  md-translate serve ./docs --watch
```

**4. `cache` - キャッシュ管理**
```bash
md-translate cache <SUBCOMMAND>

Subcommands:
  stats    キャッシュ統計表示
  clear    キャッシュ全削除
  prune    古いキャッシュエントリを削除
```

### 7. Dev Server Architecture (`server/`)

**HTTPサーバー設計 (Axum)**:

**ルート構成**:
```rust
Router::new()
    .route("/", get(index_handler))                    // ファイル一覧
    .route("/preview/:path", get(preview_handler))     // プレビュー
    .route("/api/translate", post(translate_api))      // 翻訳API
    .route("/ws", get(websocket_handler))              // WebSocket
    .nest_service("/static", ServeDir::new("static"))  // 静的ファイル
```

**主要エンドポイント**:

1. `GET /` - ファイル一覧ページ（HTMLレスポンス）
2. `GET /preview/:path` - マークダウンプレビュー
3. `POST /api/translate` - 翻訳APIエンドポイント
4. `GET /ws` - WebSocket接続（ファイル変更通知）

**サーバー状態管理**:
```rust
pub struct AppState {
    translator: Arc<Translator>,
    watch_dir: PathBuf,
    ws_clients: Arc<Mutex<Vec<WsClient>>>,
}
```

**HTMLテンプレート設計** (`templates/preview.html`):
- GitHub風マークダウンスタイル
- WebSocket接続による自動リロード
- レスポンシブデザイン

### 8. File Watching & WebSocket (`server/watcher.rs`, `server/ws.rs`)

**ファイル監視設計**:
```rust
pub struct FileWatcher {
    watcher: RecommendedWatcher,
    watch_dir: PathBuf,
    debounce_duration: Duration,  // 500ms
}
```

**監視イベントフロー**:
```
File System Change
    ↓
notify::Event
    ↓
Filter: *.md files only
    ↓
Debounce (500ms)
    ↓
Notify WebSocket Clients
    ↓
Browser Auto-reload
```

**WebSocket通信設計**:

**接続管理**:
```rust
pub struct WsClient {
    id: Uuid,
    sender: mpsc::UnboundedSender<Message>,
}
```

**メッセージプロトコル**:

Server → Client:
```json
{
  "type": "reload",
  "path": "docs/README.md",
  "timestamp": "2025-09-30T12:00:00Z"
}
```

Client → Server:
```json
{
  "type": "ping"
}
```

**並行処理設計**:
```
Main Task (Axum Server)
    ├─ HTTP Request Handlers
    ├─ File Watcher Task
    └─ WebSocket Manager Task
```

## Code References

現在のコードベース:
- `Cargo.toml:1-46` - 依存関係とfeature定義
- `src/main.rs:1-3` - 現在はHello Worldのみ

## Architecture Documentation

### システム全体図

```
┌─────────────────────────────────────────────────────────┐
│                      md_translate                        │
│                                                          │
│  ┌────────────┐         ┌──────────────┐               │
│  │    CLI     │         │  Dev Server  │               │
│  │  (main.rs) │         │   (Axum)     │               │
│  └─────┬──────┘         └──────┬───────┘               │
│        │                       │                        │
│        └───────────┬───────────┘                        │
│                    │                                    │
│        ┌───────────v───────────┐                       │
│        │   Translator Engine   │                       │
│        │  ┌─────────────────┐  │                       │
│        │  │ OllamaClient    │  │                       │
│        │  └─────────────────┘  │                       │
│        │  ┌─────────────────┐  │                       │
│        │  │ TranslationCache│  │                       │
│        │  └─────────────────┘  │                       │
│        │  ┌─────────────────┐  │                       │
│        │  │ MarkdownProcessor│ │                       │
│        │  └─────────────────┘  │                       │
│        └───────────┬───────────┘                       │
│                    │                                    │
└────────────────────┼────────────────────────────────────┘
                     │
        ┌────────────┼────────────┐
        │            │            │
        v            v            v
   ┌────────┐  ┌─────────┐  ┌─────────┐
   │ Ollama │  │  Cache  │  │ .md     │
   │ (Local)│  │  Files  │  │ Files   │
   └────────┘  └─────────┘  └─────────┘
```

### データフロー

**CLI Translateコマンド**:
```
User → CLI → Translator → [Cache Check] → [Ollama API] → Output
                              │   Hit ─────────────^
                              │   Miss → LLM Translation
```

**Dev Server**:
```
Browser → HTTP Request → Axum Handler → Translator → HTML Response
    ↑                                                      │
    │                                                      │
    └──── WebSocket ←─ File Watcher ←─ File Change ──────┘
```

### 技術スタック

**依存関係** (`Cargo.toml`):
```toml
[dependencies]
clap = { version = "4", features = ["derive"] }
tokio = { version = "1", features = ["full"] }
reqwest = { version = "0.12", features = ["json"] }
pulldown-cmark = "0.11"
serde = { version = "1", features = ["derive"] }
serde_json = "1"
sha2 = "0.10"
anyhow = "1"
thiserror = "1"
axum = { version = "0.7", optional = true }
tower-http = { version = "0.5", features = ["fs", "trace"], optional = true }
tokio-tungstenite = { version = "0.21", optional = true }
notify = { version = "6", optional = true }
indicatif = "0.17"
colored = "2"

[features]
default = ["server"]
server = ["axum", "tower-http", "tokio-tungstenite", "notify"]
```

### 設定ファイル設計

**`~/.config/md_translate/config.toml`** (オプション):
```toml
[ollama]
url = "http://localhost:11434"
model = "qwen2.5:7b"
timeout_seconds = 30

[translation]
source_lang = "en"
target_lang = "ja"
use_cache = true
parallel_requests = 3

[cache]
directory = "~/.cache/md_translate"
max_size_mb = 500

[server]
default_port = 3000
watch_debounce_ms = 500

[ui]
show_progress = true
color_output = true
```

### エラーハンドリング戦略

**エラー型階層**:
```rust
#[derive(thiserror::Error, Debug)]
pub enum MdTranslateError {
    #[error("Ollama API error: {0}")]
    OllamaError(String),

    #[error("Cache error: {0}")]
    CacheError(String),

    #[error("Markdown parsing error: {0}")]
    MarkdownError(String),

    #[error("IO error: {0}")]
    IoError(#[from] std::io::Error),

    #[error("Network error: {0}")]
    NetworkError(#[from] reqwest::Error),
}
```

**リカバリー戦略**:
- Ollama接続失敗 → 明確なエラーメッセージ + 解決方法
- セグメント翻訳失敗 → 元のテキストを保持して継続
- キャッシュ読み込み失敗 → キャッシュスキップ
- ファイル監視失敗 → 監視なしモードで継続

### パフォーマンス最適化

**1. キャッシュ戦略**:
- セグメント単位キャッシュ → 部分変更に強い
- SHA256ハッシュ → O(1) 検索
- JSON形式 → 人間が読める

**2. 並列処理**:
- 最大3並列リクエスト（Ollama負荷考慮）
- Tokio非同期ランタイム
- セグメント間は完全独立

**3. メモリ効率**:
- ストリーミングAPI未使用（シンプル性優先）
- セグメント単位処理
- Arc<Translator> でサーバー内共有

### セキュリティ考慮事項

**1. ローカル実行**:
- Ollamaはローカルホスト限定
- 翻訳データは外部送信しない
- プライバシー保護

**2. Dev Server**:
- デフォルトはlocalhost バインド
- WebSocketは同一オリジンのみ

**3. ファイルアクセス**:
- 指定ディレクトリ配下のみアクセス
- パストラバーサル対策
- シンボリックリンク追跡制限

## Historical Context (from thoughts/)

このプロジェクトは新規作成であり、既存のthoughts/ディレクトリには関連ドキュメントなし。

## Related Research

関連する既存研究ドキュメントなし（新規プロジェクト）。

## Open Questions

実装フェーズで検討が必要な項目：

1. **Ollama モデル選択**:
   - Qwen2.5:7b vs Gemma2:9b vs Llama3.2:3b
   - 翻訳品質とパフォーマンスのトレードオフ

2. **キャッシュ容量管理**:
   - LRU削除ポリシーの必要性
   - キャッシュサイズ上限の設定

3. **HTMLテンプレート**:
   - テンプレートエンジンの選択（handlebars, tera, etc.）
   - CSSフレームワークの使用有無

4. **設定ファイル**:
   - TOMLパーサーの選択
   - デフォルト値の適切性

5. **エラーメッセージ**:
   - 日本語メッセージサポート
   - i18n対応の必要性

## Implementation Priority

**Phase 1 (MVP)** - 推定2-3日:
1. Ollama client実装
2. Translation cache実装
3. Markdown parser実装
4. CLI `translate` コマンド

**Phase 2 (Enhancement)** - 推定2-3日:
1. Dev server実装
2. File watcher実装
3. WebSocket実装
4. CLI `serve` コマンド

**Phase 3 (Optional)**:
- 翻訳スタイル選択
- 差分表示機能
- バッチ翻訳
- VSCode Extension