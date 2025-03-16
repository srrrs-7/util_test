<?php
 
 declare(strict_types=1);
 
 use Rector\Config\RectorConfig;
 use Rector\Set\ValueObject\SetList;
 use Rector\ValueObject\PhpVersion;
 
 return static function (RectorConfig $rectorConfig): void {
     // プロジェクトのソースコードのパスを指定
     $rectorConfig->paths([
         __DIR__ . '',
         __DIR__ . '/tests', // テストコードも含める場合
     ]);
 
     // 使用する PHP バージョンを指定 (プロジェクトの要件に合わせて)
     $rectorConfig->phpVersion(PhpVersion::PHP_83);
 
     // 適用するルールセットを指定
     $rectorConfig->sets([
         SetList::CODE_QUALITY,
         SetList::DEAD_CODE,
         SetList::PHP_83,
         SetList::TYPE_DECLARATION,
         // 必要に応じて他のセットも追加
         // SetList::PHP_80, // PHP 8.0 へのアップグレード
         // SetList::PHP_81,
         // SetList::PHP_82,
     ]);
 
     // 個別のルールを設定 (オプション)
     // $rectorConfig->rule(...);
 
     // 特定のルールやディレクトリを除外 (オプション)
     // $rectorConfig->skip(...);
 };