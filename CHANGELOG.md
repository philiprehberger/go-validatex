# Changelog

## 0.2.1

- Fix README structure: demote Built-in Rules to subsection

## 0.2.0

- Add new built-in rules: `uuid`, `ip`, `ipv4`, `ipv6`, `alpha`, `numeric`, `alphanum`, `contains`, `excludes`
- Add custom error messages via `msg` tag option
- Add nested struct validation with dot-notation error paths
- Add `ValidateField` function for validating single values against a rules string

## 0.1.3

- Consolidate README badges onto single line

## 0.1.1

- Add badges and Development section to README

## 0.1.0

- Initial release
- Struct validation via `validate` tag
- Built-in rules: `required`, `min`, `max`, `email`, `url`, `oneof`, `len`, `pattern`
- Custom rule registration via `Register`
- Batch error collection
