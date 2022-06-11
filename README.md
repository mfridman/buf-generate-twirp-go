# Generate Twirp code with buf

Goal of this repo is to showcase:

1. How to use the [`buf` CLI](https://docs.buf.build/generate/usage) to generate Twirp code.
2. Use the `--template` option to workaround certain plugins, such as Twirp, that generate *all* code to the same package and expect the same package name (which isn't always desired).

If you ran into "error:files have conflicting go_package settings", you're in the right place.

---

### Problem

The following `buf` command generates Go and Twirp code from the [buf.gen.yaml](buf.gen.yaml) template. Bonus, all plugins are hosted on [buf.build](https://buf.build/), no need to install anything except the `buf` CLI if you want to follow along.

```
$ buf generate --template buf.gen.yaml buf.build/acme/petapis
```

❌ The problem is the generated Go code is missing dependencies and thus will not compile:

```
could not import github.com/mfridman/buf-generate-twirp-go/go/payment/v1alpha1
(no required module provides package "github.com/mfridman/buf-generate-twirp-go/go/payment/v1alpha1")
```

💡 You could add the `--include-imports` option to buf generate, but this will result in a failure due to the way the [Twirp plugin works](https://github.com/twitchtv/twirp/blob/ff7d9f87d8598707f3465c80ee5dec1ba8520310/protoc-gen-twirp/generator.go#L181-L216):

```
$ buf generate --include-imports --template buf.gen.yaml buf.build/acme/petapis

error:files have conflicting go_package settings, must be the same: "money" and "paymentv1alpha1"
```

Taking a step back, the example in this repository is using 2 modules hosted on the BSR:

1. [`acme/petapis`](https://buf.build/acme/petapis)
2. [`acme/paymentapis`](https://buf.build/acme/paymentapis)

The `acme/petapis` module has a dependency on `acme/paymentapis`[^1] and `datetime.proto` (googleapis)[^2].

The `acme/paymentapis` module has a dependency on `money.proto`[^3].

If you're curious why the buf templates use `except` for googleapis, head over to the [docs tour for more details](https://docs.buf.build/tour/use-managed-mode#remove-modules-from-managed-mode).

---

### Solution

**The solution in most cases is to split up code generation into multiple steps.**

Whereby we generate the *dependent* code using a separate template. If you have a large dependency graph this may get out of hand, but is usually fine.

To get around the error above, we can use the `--template` option with a separate template file to generate the dependant Go code like so:

```
$ buf generate --template buf.gen-go.yaml buf.build/acme/paymentapis
```

To summarize, instead of generating all the code in one shot using `buf.gen.yaml`, we use buf generate with the `--template` option to split up code generation. This is often necessary for plugins such as Twirp.

### Example

This repo is a working example. Make sure to [install the buf cli](https://docs.buf.build/installation).

Generated code is purposefully not check in, you generate it :)

```
$ rm -rf go || true

# generate Go and Twirp code
$ buf generate --template buf.gen.yaml buf.build/acme/petapis

# generate dependant Go code (not using Twirp plugin)
$ buf generate --template buf.gen-go.yaml buf.build/acme/paymentapis

$ go run cmd/server/main.go
```

Output:

```
dante PET_TYPE_DOG
dante PET_TYPE_DOG
dante PET_TYPE_DOG
... every 1 second
```

[^1]: https://buf.build/acme/petapis/file/main:pet/v1/pet.proto#L5
[^2]: https://buf.build/acme/petapis/file/main:pet/v1/pet.proto#L6
[^3]: https://buf.build/acme/paymentapis/file/main:payment/v1alpha1/payment.proto#L5
