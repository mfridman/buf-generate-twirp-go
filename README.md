# Generate twirp code with buf

Goal of this repo is to showcase:

1. Using [`buf` CLI](https://docs.buf.build/generate/usage) to generate twirp code
2. Use the `--template` option to workaround certain plugins, such as Twirp, that generate *all* code to the same package and expect the same package name (which isn't always desired).

---

### Problem

The following buf generate command generates Go and Twirp code (plugins are hosted on [buf.build](https://buf.build/), no need to install anything but the `buf` CLI). 

The problem is the generated code is missing dependencies:

```
$ buf generate --template buf.gen.yaml buf.build/acme/petapis

could not import github.com/mfridman/buf-generate-twirp-go/go/payment/v1alpha1
(no required module provides package "github.com/mfridman/buf-generate-twirp-go/go/payment/v1alpha1")
```

💡 You could add the `--include-imports` option to buf generate, but this will result in a failure due to the way the [Twirp plugin works](https://github.com/twitchtv/twirp/blob/ff7d9f87d8598707f3465c80ee5dec1ba8520310/protoc-gen-twirp/generator.go#L181-L216):

```
$ buf generate --include-imports --template buf.gen.yaml buf.build/acme/petapis

error:files have conflicting go_package settings, must be the same: "money" and "paymentv1alpha1"
```

The example in this repo uses 2 modules hosted on the BSR. The [`acme/petapis`](https://buf.build/acme/petapis) module depends on the [`acme/paymentapis`](https://buf.build/acme/paymentapis) module.

https://buf.build/acme/petapis/file/main:pet/v1/pet.proto#L5

Each of those modules has a different package and each module depends on `money.proto` and `datetime.proto` from googleapis.

https://buf.build/acme/petapis/file/main:pet/v1/pet.proto#L6
https://buf.build/acme/paymentapis/file/main:payment/v1alpha1/payment.proto#L5

Note, the buf generate templates use `except` for googleapis. [See the docs for more info](https://docs.buf.build/tour/use-managed-mode#remove-modules-from-managed-mode).

---

### Solution

The solution in most cases is to split up code generation into multiple steps, whereby we generate the dependent code using a separate template. If you have a large dependency graph this may get out of hand, but is usually fine.

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
