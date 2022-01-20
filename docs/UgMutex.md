type UgMutex
===

リソースに対する複数の同時アクセスを単一のみに制限することで、リソースの変更による破損を保護します。
加えて、変更を行わない事が保証されているアクセスに限り、複数同時アクセスを許可する機構を備えています。
これらの動作に加え、変更を行わないアクセスの後、変更を行うアクセスへの移行(Upgrade)する機構を備えています。
go標準パッケージ`sync`に含まれる`Mutex`関係には同等の機能がありません。

変更を行わないアクセスによって取得されたデータへに対して、変更を行う場合、変更を行うアクセス機構を再取得する必要が生まれます。  
アクセス保護機構が独立している場合、変更を行わないアクセスの解放、変更を行うアクセスの取得という工程が発生するため、別スレッドが変更を行うアクセスの取得に成功する危険性が生まれます。  
変更を行わないアクセスの解放、変更を行うアクセスの取得をスレッドセーフに行う必要があり、Upgradeはスレッドセーフに動作する事を保証します。

### import
```
import "github.com/l4go/mutex"
```
vendoringして使うことを推奨します。

## 機能の概略

単一のみロック取得可能なlockと、  
複数がロック取得可能なread lock、  
read lockと等しい動作をしつつ、途中で取得対象をlockに移行する機能をもつupgrade lockが存在します。  

lockの取得が試みられると、新たなread lock, upgrade lockの取得を試みるスレッドはブロックされます。  
read lockが1つでも取得されている場合、全てのread lockが解放されるまで、lockの取得を試みるスレッドはブロックされます。
read lockが1つでも取得されている場合であっても、upgrade lockの取得は可能です。  

upgrade lockが取得対象をlockに移行(upgrade)できます。  
upgradeを試みた際、read lockの動作は、lock取得時と同等です。  

## 利用サンプル

値をlockしながら書き換える`writer`と、値の読み込みのみを行う`reader`、あたいの読み込みから書き込みに繊維する`upgradeWriter`が複数のスレッド（goroutine)で動作するサンプルです

[example](../examples/ex_uglock/ex_uglock.go)

## メソッド概略

### func NewUgMutex() \*UgMutex

\*UgMutex型を生成します。

### func (m \*UgMutex) Lock()

lockの取得を試みます。  
既にread lock 取得済みのスレッドが存在する場合、全てのread lock解放までスレッドをブロックします。  
lockまたは、upgrade lockは単一のみ取得が可能で、lockを取得できない場合、スレッドをブロックします。  

### func (m \*UgMutex) Unlock()

取得したlockを解放します。

### func (m \*UgMutex) RLock()

read lockの取得を試みます。すでにlockが取得されている場合、lockの解放までスレッドをブロックします。  
upgrade lockがupgrade前の場合は、read lockの取得が可能です。  
upgrade lockがupgrade後の場合は、upgrade lockの解放までスレッドをブロックします。  
複数がread lockの取得が可能です。

### func (m \*UgMutex) RUnlock()

取得したread lockを解放します。

### func (m \*UgMutex) UgLock()

upgrade lockの取得を試みます。  
read lockへは、干渉しません。  
lockまたは、upgrade lockは単一のみ取得が可能で、upgrade lockを取得できない場合、スレッドをブロックします。  

### func (m \*UgMutex) Upgrade()

upgrade lockをlock相当の動作に移行します。他スレッドでブロックされているlockに対してスレッドセーフな動きを保証します。  
read lockに対して、lockが取得された時と同様の干渉を行います。  
既にread lock 取得済みのスレッドが存在する場合、全てのread lock解放までスレッドをブロックします。  

### func (m \*UgMutex) UgUnlock()

取得したupgrade lockを解放します。
