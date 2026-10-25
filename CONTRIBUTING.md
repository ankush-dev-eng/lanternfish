# Contributing to LANTERNFISH

Thank you for your interest in contributing to **LANTERNFISH**! We welcome bug fixes, documentation enhancements, and additional simulated attack scenarios.

---

## Code of Conduct & Principles

1. **Zero Unnecessary Dependencies:** Keep the backend standard-library-only. Avoid pulling in third-party frameworks unless strictly necessary.
2. **Deterministic Reliability:** Deception and detection rules must operate deterministically with zero false positives.
3. **Field-Instrument Aesthetics:** Preserve the deep-sea, nautical field-instrument visual design. No neon glows or emoji spam.

---

## Development Workflow

1. **Fork and Clone:**
   ```bash
   git clone https://github.com/ankush-dev-eng/lanternfish.git
   cd lanternfish
   ```

2. **Verify Formatting and Code Quality:**
   ```bash
   # Code formatting
   gofmt -w -s .

   # Static analysis
   go vet ./...

   # Run test suite
   go test -v -race ./...
   ```

3. **Conventional Commits:**
   Please format commit messages using Conventional Commits:
   - `feat:` new feature or capability
   - `fix:` bug fix
   - `test:` adding or updating tests
   - `docs:` documentation improvements
   - `style:` visual CSS or formatting edits
   - `refactor:` code refactoring without functional changes
   - `ci:` continuous integration workflow changes

4. **Pull Request Checklist:**
   - [ ] All unit and integration tests pass cleanly (`go test ./...`).
   - [ ] `go vet ./...` reports zero issues.
   - [ ] New features include corresponding unit tests.
