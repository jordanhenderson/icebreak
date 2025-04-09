# 🔊 icebreak

_A lightweight Go module that executes AWS Lambda bootstrap layers during cold starts — before your main runtime logic kicks in._

## 🔥 Why icebreak?

AWS Lambda supports [bootstrap wrapper scripts](https://docs.aws.amazon.com/lambda/latest/dg/runtimes-modify.html) via the `AWS_LAMBDA_EXEC_WRAPPER` environment variable. These wrappers can be used to initialize services, inject environment config, or run sidecar-style logic **before your actual handler starts**.

`icebreak` makes it easy to support multiple wrapper scripts as a chain — skipping self-invocation and transparently handling errors.

## ✨ Features

- Runs multiple bootstrap wrapper binaries in order
- Skips self-execution to avoid recursion
- Automatically logs all wrapper activity
- Lightweight, zero-config
- Designed for cold start phase

## 📦 Installation

```bash
go get github.com/jordanhenderson/icebreak
```

## 🚀 Usage

Import it anonymously to automatically enable `init()`:

```go
import _ "github.com/jordanhenderson/icebreak"
```

Then write your Lambda bootstrap logic as usual:

```go
package main

import (
	_ "github.com/jordanhenderson/icebreak"
	"log"
)

func main() {
	log.Println("[bootstrap] Lambda runtime starting...")
	// Your Lambda runtime logic goes here
}
```

## 🔧 How it works

At cold start, AWS sets the `AWS_LAMBDA_EXEC_WRAPPER` environment variable, usually pointing to a single binary path.

`icebreak` extends this by supporting **comma-separated lists** of wrapper paths:

```bash
export AWS_LAMBDA_EXEC_WRAPPER="/opt/wrapper1,/opt/wrapper2"
```

It will:
1. Resolve symlinks for all wrappers
2. Skip if the wrapper is the current binary (self-check)
3. Log and execute each in order
4. Continue startup

## 🧪 Example

```bash
AWS_LAMBDA_EXEC_WRAPPER="/opt/logwrap,/opt/initsecrets" ./bootstrap
```

With `icebreak` embedded, your Lambda will:
- Run `/opt/logwrap`
- Then run `/opt/initsecrets`
- Then continue with your main Go `main()` function

## 💪 Development

```bash
git clone https://github.com/jordanhenderson/icebreak
cd icebreak

```

## 📝 License

MIT © 2025 Jordan Henderson

---

**Break the ice before your Lambda gets cold.**

