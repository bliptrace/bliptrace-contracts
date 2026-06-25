<?php

declare(strict_types=1);

namespace G8Scope\Contracts;

use InvalidArgumentException;

/**
 * Locator + decoder for the versioned g8scope error-envelope JSON Schema.
 *
 * This is the Composer consumption surface (spec §10, §20.1): g8scope-core and
 * g8scope-app derive their envelope DTOs from the schema returned here, instead of
 * redefining the envelope. The Go services consume the same files as a Go module.
 */
final class Schema
{
    /** Schema majors shipped by this release. */
    public const SUPPORTED_MAJORS = [1];

    private static function root(): string
    {
        return \dirname(__DIR__);
    }

    /** Absolute path to the envelope schema for a major version. */
    public static function path(int $major = 1): string
    {
        $path = self::root()."/schema/v{$major}/envelope.schema.json";
        if (! is_file($path)) {
            throw new InvalidArgumentException("Unsupported envelope schema major: v{$major}");
        }

        return $path;
    }

    /** Absolute path to the conforming example envelope for a major version. */
    public static function examplePath(int $major = 1): string
    {
        return self::root()."/examples/v{$major}/envelope.example.json";
    }

    /** Raw JSON Schema string. */
    public static function contents(int $major = 1): string
    {
        return (string) file_get_contents(self::path($major));
    }

    /** Decoded JSON Schema as an associative array. */
    public static function decode(int $major = 1): array
    {
        return (array) json_decode(self::contents($major), true, 512, JSON_THROW_ON_ERROR);
    }

    /** Major version extracted from an envelope's schema_version (e.g. "1.0" => 1). */
    public static function majorFromVersion(string $schemaVersion): int
    {
        $head = explode('.', trim($schemaVersion))[0] ?? '';
        if ($head === '' || ! ctype_digit($head)) {
            throw new InvalidArgumentException("Invalid schema_version: {$schemaVersion}");
        }

        return (int) $head;
    }
}
