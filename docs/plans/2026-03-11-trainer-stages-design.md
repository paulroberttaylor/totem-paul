# Data-Driven Trainer Stages

## Goal

Replace hardcoded trainer stages with JSON files so new content (like Apex code drills) can be added without touching Go code.

## Stage File Format

Files live in `trainer/stages/`. Each JSON file defines a category:

```json
{
  "category": "Apex",
  "stages": [
    {
      "name": "apex_keywords",
      "label": "Apex Keywords",
      "keys": ["p","u","b","l","i","c"],
      "words": ["public","private","class","void","String","Integer","List","Map"],
      "snippets": []
    },
    {
      "name": "apex_patterns",
      "label": "Apex Patterns",
      "keys": [],
      "words": [],
      "snippets": [
        "public class AccountService {",
        "List<Account> accounts = [SELECT Id, Name FROM Account];",
        "@isTest static void testInsert() {"
      ]
    }
  ]
}
```

### Field semantics

- **`keys`** — character set for random word generation (existing behavior)
- **`words`** — explicit word list for this stage (replaces filtering from common.txt)
- **`snippets`** — full lines to type verbatim (new: "type this code" mode)
- If `snippets` is non-empty, exercises pull from snippets. Otherwise, words/keys filtering.

## File structure

```
trainer/stages/
├── core.json       # Existing 7 stages migrated from Go
├── apex.json       # Apex keywords + code snippet drills
```

## Code changes

### lesson.go

- Replace `AllStages()` with `LoadStages(dir string) ([]CategoryGroup, error)`
- `CategoryGroup`: `Category string`, `Stages []Stage`
- `Stage` gains `Words []string`, `Snippets []string` fields
- `common.txt` wordlist still used as fallback when stage has `keys` but no `words`

### picker.go

- Group stages by category with headers in the selection list

### typing.go

- When stage has snippets: pick a random snippet as exercise text
- User types it verbatim, scoring unchanged (per-character accuracy, WPM)

### No changes to

- stats, results, keymap parser, app routing

## Apex content

1. **Apex Keywords** — `public`, `private`, `class`, `void`, `String`, `List`, `Map`, `Set`, `SOQL`, `@isTest`, `trigger`, `insert`, `update`, `delete`, `try`, `catch`
2. **Apex Snippets** — 15-20 real code patterns: class headers, SOQL queries, DML operations, triggers, test methods
