# LINE TechTrain - TODO アプリ 実装例

:earth_americas: [English version](README.md)

## 概要

このリポジトリは LINE の TechTrain 課題の参考資料です。以下の内容で構成されています

- /example
	- /client
		- FrontEndの実装例であり、HTML/CSS/JSで記載されています
	- /server
		- ApplicationServerの実装例であり、Go言語で記載されています
- /wireframe
	- /toppage.png
		- トップページのUIのワイヤーフレームです
	- /detail.png
		- 詳細ページのUIのワイヤーフレームです
- /architecture
	- system.png
		- 想定される簡易なシステム構成図です

なお、実行環境は提供していません。

## 本資料の目的

開発に入るまえに全体像を掴むことは重要なことです。
本リポジトリ内の内容は、TechTrain課題開始にあたり全体像の把握のための資料として、あくまで１例としてご利用ください。
TODO 課題を記載のコード/画像通りに作成してほしいという意味で提供はしておりません。

## 各ページの要件例

- トップページ
	- ToDoアイテムの一覧を表示
	- 各アイテムは、詳細ページのリンクとなっており、各アイテムの左に終了ボタンが表示される
	- 終了アイテムの一覧 - 終了済みアイテムは折り畳まれた状態で表示され、クリックすると一覧を表示する
- 詳細ページ
	- フォームを表示する
		- タイトル編集インプット
		- 説明文編集インプット
		- 登録ボタン
		- 削除ボタン

## 関連項目

- [line/line-liff-v2-starter](https://github.com/line/line-liff-v2-starter)
	- LIFFアプリの作成方法例を示すリポジトリ

## ライセンス

本リポジトリ内のコンテンツは TechTrain 課題の遂行を目的とする場合のみ利用可能です。
著作権は LINE 株式会社が保有し、転載を禁じます。