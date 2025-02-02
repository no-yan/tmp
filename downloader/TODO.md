## 機能
- [x] ダウンロードしたレスポンスの保存
- [x] 完了時ログ
- [x] CLIのオプション
  - e.g. 最大ダウンロード数, 出力先ディレクトリ, timeout
- [x] CTRL + Cでcontext.Cancelを実行

## 改善
- [x] ダウンロード中のプログレス表示
- [x] ダウンロード処理と保存処理の関心の分離
- [x] Abort時のエラーログ
- [ ] リトライ時にプログレスバーをdownloadingからretryingに変更

## テスト
- [ ] downloaderのテスト
- [x] backoffのテスト
- [ ] pubsubのテスト

