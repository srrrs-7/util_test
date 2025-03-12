<?php

declare(strict_types=1);

use Rector\Config\RectorConfig;
use Rector\Set\ValueObject\LevelSetList;
use Rector\Set\ValueObject\SetList;
use Rector\CodeQuality\Rector\Class_\InlineConstructorDefaultToPropertyRector;

return static function (RectorConfig $rectorConfig): void {
    // 対象とするファイル/ディレクトリ
    $rectorConfig->paths([
        __DIR__ . '/app',
        __DIR__ . '/bootstrap',
        __DIR__ . '/config',
        __DIR__ . '/database',
        __DIR__ . '/public',
        __DIR__ . '/resources',
        __DIR__ . '/routes',
        __DIR__ . '/storage',
        __DIR__ . '/tests',
    ]);

    // インポートするルールセット (レベルセット、個別セット)
    $rectorConfig->sets([
        LevelSetList::UP_TO_PHP_83, // PHP 8.1までのレベルセット
        SetList::CODE_QUALITY,       // コード品質向上ルール
        SetList::DEAD_CODE,          // 不要コード削除ルール
        SetList::TYPE_DECLARATION,   // 型宣言追加ルール
        // 他にも多数のセットがあります
    ]);

    // 個別のルールをスキップ (除外)
    $rectorConfig->skip([
        InlineConstructorDefaultToPropertyRector::class,

        // 特定のファイルやディレクトリをスキップ
        __DIR__ . '/src/SomeLegacyClass.php',
        __DIR__ . '/tests/data',

        // 特定のルールを特定のファイルに適用しない
        // Rector\CodeQuality\Rector\Identical\FlipTypeControlToUseExclusiveTypeRector::class => [ 
        // ],
    ]);

    // より詳細な設定オプション
    // $rectorConfig->disableParallel(); // 並列処理を無効化（メモリ不足の場合など）
    // $rectorConfig->importNames(); // 名前空間を自動的にインポート
    // $rectorConfig->removeUnusedImports();   //使われていないuseを削除
    // $rectorConfig->phpVersion(70400); // PHP 7.4 をターゲットとする (古い書き方)
    //$rectorConfig->phpstanConfig(__DIR__ . '/phpstan.neon'); //PHPStanの設定ファイルを使う場合
};