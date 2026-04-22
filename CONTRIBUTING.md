
## Project philosophy

A set is easy to write until it needs to be good.

Anyone can throw together `map[T]struct{}`, and for many cases that is enough. The point of `set` is not to make the trivial possible. It is to make the common useful cases ergonomic, consistent, and cheap.

Go has plenty of developers who do not care about allocations until they suddenly have to. This project is for the moment after that. The package aims to provide familiar set operations and practical collection semantics with straightforward Go-style APIs, while keeping unnecessary heap traffic and overhead out of hot paths.

Writing a set is simple. Writing one that behaves well around the compiler, allocator, garbage collector, iteration costs, concurrency, and real workloads is less simple.

`set` exists to handle that part without making users think about it.


## Comment style guide

Go documentation is intentionally unstructured, which keeps it lightweight but often results in inconsistent and harder-to-scan comments. 
This project applies a consistent, sentence-ordered style to reduce cognitive load when reading and writing documentation.

Doc comments must follow these principles:

* Documentation must be concise, ordered, and explicit about behaviour and contract.
* Comments should read as a sequence of short, self-contained sentences, each introducing one idea.
* Include only information that is not immediately obvious from the function signature, core Go semantics, or trivial inspection.
* Do not use headings or bullet points in doc comments.
* Use doc links whenever another function, method, or type is referenced.
* Use `Note:` only for secondary clarifications that do not belong in the primary flow.

Refer to the template below as a guide for ordering information:

```go
// <Name> <primary purpose sentence>.
//
// <Core behaviour description under valid usage>.
//
// <Termination or completion conditions>.
//
// <Ordering or determinism guarantees>.
//
// <Side effects or non-obvious performance characteristics>.
//
// <Constraints or requirements on inputs or callbacks>.
//
// <Nil handling semantics>.
//
// <Panic behaviour introduced by this function>.
//
// Note: <secondary constraint or clarification, if needed>.
func Function()
```

Maintain this consistent progression of concepts so readers can reliably locate details. Each section is optional and must be included only when it provides meaningful, non-obvious information. 