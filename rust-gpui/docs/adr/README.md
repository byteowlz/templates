# Architecture Decision Records

An ADR is a short, **hard-to-reverse** decision with a real trade-off, written
down at the moment it is made. Use it for a decision a future contributor would
otherwise re-litigate; do not use it for routine changes.

## When to write one

- A choice that is expensive to reverse later (a crate seam, a storage format,
  an external dependency, a security boundary).
- A rejected alternative worth recording so it is not re-proposed.
- A cross-cutting rule the team should treat as settled.

## Status

`proposed` -> `accepted` (or `rejected` / `superseded`). Template:
[`0000-template.md`](0000-template.md).

## Naming

`NNNN-short-slug.md`, numbered in order. A new ADR extends the highest number;
superseding an ADR updates the superseded one's status and links the new one.

## Index

| ADR | Title | Status |
| --- | ----- | ------ |