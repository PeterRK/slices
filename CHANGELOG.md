# Changelog

## Unreleased

- Clarify that ordered floating-point APIs do not specify results for inputs
  containing NaNs. Mutating operations must still terminate safely and preserve
  the input element multiset.
- Correct the documented v2 API signatures and include previously omitted
  partial-sort, `Order`, and general slice operations.
- Document the general slice API as a Go 1.21 semantic snapshot.
- Move the `Insert` allocation check from a functional test to a benchmark so
  instrumented test builds can run independently of performance thresholds.
- Use the language `<` operator instead of `cmp.Less` in ordered algorithm
  paths, matching the documented unspecified handling of NaNs and improving
  ordinary floating-point sorting performance.
- Add ordered floating-point stable-sort benchmarks.
