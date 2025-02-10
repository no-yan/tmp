Downloaderは、複数のURLからファイルを並列ダウンロードするGo製サンプルです。

ワーカー数の制限やリトライ時のバックオフ処理、Pub/Subアーキテクチャによるイベント通知など、Goの並行処理・設計パターンを学習するために作成しました。

## デモ
<video width="250" src="https://github.com/user-attachments/assets/27a7679e-8d33-4867-a266-4aa8fb2ed539">
</video>


## 主な機能

- 並列ダウンロード: goroutineとchannelでワーカー数を制御し、高速に取得
- バックオフリトライ: 上限つき指数バックオフを実装
- Pub/Subアーキテクチャ: ダウンロード進捗をサブスクライバに通知
- コンテキスト制御: context.WithTimeoutとOSシグナル処理で一括キャンセル
- プログレスバー: mpbで進捗を可視化

## 使い方

```sh
go build -o downloader
./downloader --output-dir=out --workers=4 --request-timeout=30s \
  https://example.com https://example.com/api
```



