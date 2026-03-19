# Release Guide

This document explains how to create releases for this project using `standard-version` and how to write commit messages following the **Conventional Commits** specification.

---

## How to Release

We use [`standard-version`](https://github.com/conventional-changelog/standard-version) to automate versioning and changelog generation based on your commit messages.

### Commands

- **Patch release** (default):
  Increment the patch version (e.g., `2.2.0` → `2.2.1`)

  ```bash
  npm run release
  ```

- **Minor release**:
  Increment the minor version (e.g., `2.2.1` → `2.3.0`)

  ```bash
  npm run release:minor
  ```

- **Major release**:
  Increment the major version (e.g., `2.3.0` → `3.0.0`)

  ```bash
  npm run release:major
  ```

- **Custom release version**:
  Specify the exact version number to release (e.g., `2.2.1`)

  ```bash
  npm run release -- --release-as 2.2.1
  ```

---

## How to Commit

This project follows the **Conventional Commits** specification for commit messages.

### Why?

- It allows automatic changelog generation.
- It helps determine semantic version bumps.
- It creates consistent, easy-to-understand commit history.

### Commit Message Format

A commit message must be structured like this:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types** commonly used:

| Type     | Description                                             |
| -------- | ------------------------------------------------------- |
| feat     | A new feature                                           |
| fix      | A bug fix                                               |
| docs     | Documentation only changes                              |
| style    | Code style changes (formatting, no logic)               |
| refactor | Code change that neither fixes a bug nor adds a feature |
| perf     | Performance improvements                                |
| test     | Adding or fixing tests                                  |
| chore    | Changes to build process or auxiliary tools             |

### Examples

- `feat: add user login feature`
- `fix(auth): correct token expiration bug`
- `docs: update README with new instructions`
- `chore: upgrade dependencies`

---

### Learn More

Read the full Conventional Commits specification here:
[https://www.conventionalcommits.org/en/v1.0.0/#summary](https://www.conventionalcommits.org/en/v1.0.0/#summary)

---

## Important Notes

- Commits **not following** the conventional commit format will **not be included** in the changelog.
- Version bumps are based only on commits with recognized types (e.g., `feat` triggers minor, `fix` triggers patch).
- Use a commit message linter like [`commitlint`](https://github.com/conventional-changelog/commitlint) to enforce proper commit message format.

---

### How to recreate this for another repos

```bash
npm init -y

npm install --save-dev standard-version @commitlint/cli @commitlint/config-conventional husky

# In package.json
"scripts": {
  "release": "standard-version",
  "release:minor": "standard-version --release-as minor",
  "release:major": "standard-version --release-as major",
  "prepare": "husky install"
}

npx husky install
npx husky add .husky/commit-msg 'npx --no -- commitlint --edit "$1"'

# create commitlint.config.js in root
module.exports = {
  extends: ['@commitlint/config-conventional'],
};

# copy frontend release ci/cd if needed
```
