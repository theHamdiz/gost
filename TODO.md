- `gen -> README.md` -> change folder structure tree

- `cli->`  implement gost test

- `gen -> handlers` -> remove duplicate views declarations

- `gen -> middleware -> notifier` -> remove unused imports ( context, time ) or use them.

- `gen -> middleware -> recoverer` -> remove unused imports ( runtime/debug ) or use them.

- `CreateProject` -> Run a final go mod tidy

- `gen -> router` -> change prelude import to types/core/gost & remove unsed log import or use it.