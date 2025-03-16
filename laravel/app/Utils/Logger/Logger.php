<?php

namespace App\Utils\Log;

use \Illuminate\Support\Facades\Log;

class Logger
{
    public static function info(string $message, array $context = []): void
    {
        Log::info($message, $context);
    }

    public static function error(string $message, array $context = []): void
    {
        Log::error($message, $context);
    }

    public static function warning(string $message, array $context = []): void
    {
        Log::warning($message, $context);
    }

    public static function debug(string $message, array $context = []): void
    {
        Log::debug($message, $context);
    }

}