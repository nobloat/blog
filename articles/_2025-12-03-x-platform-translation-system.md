# Building a cross-platform translation system

At [zeitkapsl.eu](https://zeitkapsl.eu), we're building a privacy-first, end-to-end encrypted photo storage app that runs across iOS, Android, desktop, and web. As we expand into European markets, localization became critical, but we quickly discovered that managing translations across multiple platforms with their different native formats was a nightmare.


## Existing tools did not match our workflow

As an open source project, we naturally looked first at open source translation management solutions. We evaluated several options, but each had fundamental limitations that prevented them from meeting our core requirement: 

**seamless two-way synchronization with native developer tools.**

[**Weblate**](https://weblate.org/) was our first serious consideration. It's a mature, self-hosted open source translation platform with a Git-based workflow that seemed promising. However, while Weblate can import and export various formats, it still requires maintaining separate format files for each platform. The real deal-breaker was that developers would need to either use Weblate's web interface or learn its API instead of working directly in Xcode and Android Studio. The setup and maintenance overhead was also significant—we'd need to run and maintain a web server, database, and all the associated infrastructure just to manage translations. For a small team, this felt like overkill.

[**Pootle**](https://pootle.translatehouse.org/) is another open source option, but it suffers from similar issues. The interface feels dated, and it has limited support for modern iOS formats like `.xcstrings`. More importantly, like Weblate, it doesn't solve the fundamental problem: developers still can't use their native tools. The workflow becomes "add string in Xcode, commit, wait for Pootle to import, translate in Pootle, export, commit back"—too many steps and too many opportunities for things to go out of sync.

The native tools themselves, Xcode's String Catalog Editor and Android Studio's Translations Editor, work well for developers working within a single platform, but they offer no cross-platform synchronization. There's no way to sync iOS strings to Android or web, no unified view of all translations, and translators can't use the spreadsheet tools they're comfortable with.

What we needed was clear: developers should continue using Xcode and Android Studio without changing their workflow. Translators should work in Excel, or similar familiar tools. The system should handle true two-way sync—importing from native formats, allowing edits in CSV, and exporting back to native formats—all while preserving platform-specific features like plural rules and state tracking. Everything should be stored in plain files we control, work seamlessly with Git, and require minimal setup.

None of the existing open source solutions met all these requirements. They either required developers to change their workflow, needed complex infrastructure to run, didn't support true bidirectional sync, or weren't accessible to non-technical translators. So we built our own solution.

The implementation consists of CSV files and a CLI tool. It works with any editor or translation service, stores all data in plain CSV files, integrates with Git version control, and uses direct file manipulation for basic operations (no API calls required). Developers continue using native tools, and translations sync automatically.


## No Two-Way Sync

Each platform we support uses a completely different localization format:

- **iOS**: `.xcstrings` files (JSON-based, with state tracking)
- **Android**: `values/<lang>/strings.xml` (XML-based, with plural rules)
- **Web/Backend/E-Mail templates**: `<lang>.json` files (simple key-value pairs)

Since we have full native apps (not hybrid like Flutter or React Native), our developers work directly in Android Studio and Xcode. When a developer adds a new feature, they naturally add new translation keys using the built-in tools:

- In **Xcode**, developers use the visual string editor or directly edit `.xcstrings` files
- In **Android Studio**, developers use the Translations Editor or edit `strings.xml` directly

This creates a critical workflow requirement: **translations flow in both directions**.

1. **Developer → CSV**: When developers add new strings in native tools, those changes must be imported into our central CSV so translators can work on them
2. **Translator → Native**: When translators complete translations in the CSV, those must be exported back to native formats so developers can use them
3. **Cross-platform sync**: Many strings are shared across platforms, while some are only relevant for one platform.

The native translation tools for each platform work well for single-platform development, but they don't support true two-way synchronization. You can't:
- Import new strings from iOS after a developer adds them in Xcode
- Edit those strings in a spreadsheet that translators are comfortable with
- Export back to both iOS and Android while maintaining format-specific features like plural rules and state tracking
- Keep everything synchronized across all platforms automatically

Without two-way sync, you're stuck manually copying strings between platforms, losing changes, or forcing developers to work in unfamiliar translation tools.

## CSV as the source of truth

We built a custom translation helper system in Go that uses CSV as the central storage format. CSV has several practical advantages for this use case:

- Universal compatibility: Every editor, spreadsheet application, and translation service can read/write CSV
- Human-readable: Translators can work directly in Excel, LibreOffice, or Google Sheets
- Version control friendly: CSV diffs cleanly in Git
- Simple structure: Easy to parse and manipulate programmatically

The CSV format includes columns for:
- `app`: Which platform/module the string belongs to (ios, android, web, core, server_emails)
- `key`: The translation key
- `comment`: Context for translators, also passed in to automatic translation tools
- Language columns: One column per language (e.g., `en`, `de`, `fr`, `es`)
- Regional variants: Support for regional differences (e.g., `de-AT`, `en-US`)


Here's what a sample row looks like:

```csv
app;key;comment;de;en;fr
android;welcome_message;Greeting on app start;Willkommen!;Welcome!;Bienvenue!
android;photo_count.singular;Shown when there's 1 photo;1 Foto;1 photo;1 photo
android;photo_count.plural;Shown when there are multiple photos;%d Fotos;%d photos;%d photos
```

## Import, Edit, Export

The system operates on a simple workflow:

### 1. Import from All Platforms

```bash
zeitkapsl-translations import
```

This command scans all platform directories and imports existing translations:
- Reads iOS `.xcstrings` files, extracting keys, values, and comments
- Parses Android `strings.xml` files from all `values-*` directories
- Loads JSON translation files from web, core, and server modules

All translations are normalized into the unified CSV format, preserving:
- Plural forms (`.singular` and `.plural` suffixes)
- Comments and context
- Placeholder variables (like `%d`, `%1s`)
- Platform-specific metadata

### 2. Edit in Your Preferred Tool

Once imported, the CSV file can be opened in:
- **Excel** or **LibreOffice Calc** for spreadsheet-style editing
- **Any text editor** for quick fixes
- **Translation management platforms** that support CSV import

Translators can see all languages side-by-side, which helps spot missing translations and maintain consistency.

### 3. Export Back to Native Formats

```bash
zeitkapsl-translations export --platform=all
```

The export process:
- Reconstructs iOS `.xcstrings` files with proper state tracking
- Generates Android `strings.xml` files with correct plural rules
- Creates JSON files for web/backend modules
- Handles regional fallbacks (e.g., `de-AT` falls back to `de` if missing)
- Preserves placeholder consistency across all languages

### 4. Auto-Translation Support

For initial translations, we integrated AI translation services:

```bash
zeitkapsl-translations auto-translate
```

The tool supports both DeepL and Azure Translator APIs, automatically filling in missing translations from English. It includes rate limiting and progress tracking, making it practical to translate hundreds of strings across multiple languages.

## Where Translations Are Integrated

The translation system manages strings across all parts of our application:

### Native Mobile Apps

**iOS** (`ios/Zeitkapsl/Supporting Files/Localizable.xcstrings`)
- All user-facing strings in the iOS app
- Used by Swift code via `NSLocalizedString()` and SwiftUI's `Text()` with localization
- Integrated directly into Xcode's localization workflow
- Developers can add new strings using Xcode's built-in string catalog editor

**Android** (`android/app/src/main/res/values-*/strings.xml`)
- All user-facing strings in the Android app
- Used by Kotlin code via `getString(R.string.key)` and resource references
- Integrated with Android Studio's Translations Editor
- Developers can add new strings directly in `strings.xml` or via Android Studio's UI

### Web Frontend

**Web Translations** (`web/static/translations/*.json`)
- All UI strings for the Svelte-based web application
- Loaded dynamically via `/translations/{locale}.json` endpoint
- Used throughout the frontend via an i18n store: `$t('key')`
- Supports regional variants with fallback (e.g., `de-AT` → `de` → `en`)
- Language switching updates the entire UI instantly

### Email Templates

**Server Email Translations** (`server/pkg/mail/templates/*.json`)
- All email subject lines and body text
- Used by the Go backend when sending transactional emails
- Templates include HTML and plain text versions
- Supports emails like:
  - Account welcome messages
  - Subscription change notifications
  - Password reset emails
  - Account cancellation confirmations
  - And many more transactional emails

### Core Library

**Core Translations** (`core/pkg/i18n/*.json`)
- Shared translations used across multiple platforms
- Common strings that appear in both mobile apps and web
- Ensures consistency across all user touchpoints

When a translator updates a string in the CSV, it propagates to iOS (via `.xcstrings` export), Android (via `strings.xml` export), web (via JSON export), and email templates (via server JSON export) from a single source of truth.

## Key Features

### Plural Handling

The system handles plural forms across platforms:
- **iOS**: Uses `.singular` and `.plural` key suffixes
- **Android**: Converts to `<plurals>` XML elements with `one` and `other` quantities
- **Web**: Maintains the same key structure

### Regional Variants

Support for regional language variants (like `de-AT` for Austrian German) with automatic fallback to base languages when regional translations are missing.

The regional support like `de-AT` and `de-DE` was especially important early on, since some words like `February` are different in Germany (`Februar`) and Austria (`Feber`).


### Platform-Specific Modules

The system tracks which translations belong to which module:
- `ios`: iOS app strings
- `android`: Android app strings  
- `web`: Web frontend translations
- `core`: Shared core library translations
- `server_emails`: Email template translations

This allows platform-specific strings while maintaining a unified workflow.

### Status Tracking

```bash
go run . status
Loaded 940 translations from translations.csv
Translation Status:
==================
Total languages: 3
Total translation keys: 940

Keys by platform:
  ios: 213 keys
  server_emails: 77 keys
  web: 465 keys
  android: 149 keys
  core: 36 keys

Languages: [de de_AT en]
```

Shows translation coverage:
- Total languages supported
- Number of translation keys per platform
- Missing translations per language

## The Implementation

The tool is written in Go and can be found [here](https://github.com/zeitkapsl/zeitkapsl/translations):

**Importers**: Platform-specific parsers (iOS `.xcstrings`, Android XML, JSON)

```
$> go run . import
Importing android from ../
Imported 283 translations from Android strings.xml files
Importing from ../ios/Zeitkapsl/Supporting Files/Localizable.xcstrings...
Imported 222 translations from iOS .xcstrings file
Importing server_emails from ../server/pkg/mail/templates
processing: ../server/pkg/mail/templates/de.json
Importing from JavaScript file: ../server/pkg/mail/templates/de.json...
processing: ../server/pkg/mail/templates/en.json
Importing from JavaScript file: ../server/pkg/mail/templates/en.json...
Importing core from ../core/pkg/i18n
processing: ../core/pkg/i18n/de.json
Importing from JavaScript file: ../core/pkg/i18n/de.json...
processing: ../core/pkg/i18n/de_AT.json
Importing from JavaScript file: ../core/pkg/i18n/de_AT.json...
processing: ../core/pkg/i18n/en.json
Importing from JavaScript file: ../core/pkg/i18n/en.json...
Importing web from ../web/static/translations
processing: ../web/static/translations/de.json
Importing from JavaScript file: ../web/static/translations/de.json...
processing: ../web/static/translations/de_AT.json
Importing from JavaScript file: ../web/static/translations/de_AT.json...
processing: ../web/static/translations/en.json
Importing from JavaScript file: ../web/static/translations/en.json...
Saving to CSV: translations.csv
Saved translations to translations.csv
```

**Exporters**: Format-specific generators that reconstruct native files

```
$> go run . export
Loaded 940 translations from translations.csv
Exporting android to ../
Exported Android translations to ../android/app/src/main/res/values-de/strings.xml
Exported Android translations to ../android/app/src/main/res/values/strings.xml
Exporting ios to ../ios/Zeitkapsl/Supporting Files
Successfully exported 213 entries to iOS format: ../ios/Zeitkapsl/Supporting Files/Localizable.xcstrings
Exporting server_emails to ../server/pkg/mail/templates
Create translation: ../server/pkg/mail/templates/de.json
Create translation: ../server/pkg/mail/templates/en.json
Exported web translations to separate JSON files in ../server/pkg/mail/templates
Exporting core to ../core/pkg/i18n
Create translation: ../core/pkg/i18n/de.json
Create translation: ../core/pkg/i18n/de_AT.json
Create translation: ../core/pkg/i18n/en.json
Exported web translations to separate JSON files in ../core/pkg/i18n
Exporting web to ../web/static/translations
Create translation: ../web/static/translations/de.json
Create translation: ../web/static/translations/de_AT.json
Create translation: ../web/static/translations/en.json
Exported web translations to separate JSON files in ../web/static/translations
Export completed successfully!

$> (main)> git status
On branch main
Your branch is ahead of 'origin/main' by 1 commit.
  (use "git push" to publish your local commits)

Changes not staged for commit:
  (use "git add <file>..." to update what will be committed)
  (use "git restore <file>..." to discard changes in working directory)
	modified:   ../android/app/src/main/res/values-de/strings.xml
	modified:   ../android/app/src/main/res/values/strings.xml
	modified:   ../core/pkg/i18n/de.json
	modified:   ../core/pkg/i18n/en.json
	modified:   ../ios/Zeitkapsl/Supporting Files/Localizable.xcstrings
	modified:   ../server/pkg/mail/templates/de.json
	modified:   ../server/pkg/mail/templates/en.json
	modified:   translations.csv
	modified:   ../web/static/translations/de.json
	modified:   ../web/static/translations/en.json
```

**Translation service**: Pluggable interface for AI translation (DeepL, Azure)

````
$> go run . auto-translate
Using DeepL translation service
Loaded 940 translations from translations.csv
Using DeepL for translation from English to all other languages
Target languages: de
Found 734 strings with English content
Translating [en→de]: Time travel -> Zeitreisen
Saved translations to translations.csv
Auto-translation completed. Translated 1 strings from English.
```

The codebase is straightforward—each platform has its own import/export functions, and the core `Translations` struct manages the unified data model.

## Lessons Learned

1. **Don't underestimate the overhead** that off-the-shelf solutions have if they don't fit your workflow or requirements exactly.
2. **Two-way sync is hard**: Native tools assume one-way workflows; building true bidirectional sync requires custom tooling
3. **Plural rules vary**: Each platform handles plurals differently; normalization is essential
4. **Translators prefer spreadsheets**: Most translators are more comfortable in Excel than in developer tools

## Conclusion

We now have something that **exactly matches our requirements and workflow**. The effort spent building it was very reasonable compared to the time savings we already see. Translations work across all our clients and server applications, including email templates. Everything is version controlled in Git. Developers can stick to their workflow, and so can translators. 

It might not be suitable for someone else's requirements or workflow, but it can also be an inspiration for others to avoid external dependencies and try to build something that matches exactly their requirements.

That's why I decided to write about it. Also it fits the "nobloat spirit" because it only dependes on `go` and its `stdlib` no other external dependencies are used, so it should be rather maintenance free for the upcoming years, unless google/apple decide to change their translation XML/JSON formats.

## Links / References

- [zeitkapsl.eu](https://zeitkapsl.eu)
- [zeitkapsl on codeberg](https://codeberg.org/zeitkapsl/zeitkapsl)
- [translations module](https://codeberg.org/zeitkapsl/zeitkapsl/src/branch/main/translations)
