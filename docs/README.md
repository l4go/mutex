l4go/mutex ライブラリ
===

golangの`sync`ライブラリ、`Mutex`で不足する機能を追加したモジュール群です。
	* `sync.Locker` interfaceに互換があります

* [mutex.UgMutex](./UgMutex.md)
	* 複数がロックを取得できる機構から、単一がロックを取得できる機構への変換(`Upgrade`)する機構を有するMutexです。`sync`には存在しません
