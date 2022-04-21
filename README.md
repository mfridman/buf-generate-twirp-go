# Generate twirp code with buf

Running the following will generate code Go + Twirp code, but the generated code will not compile because of missing dependencies:

```
$ buf generate --template buf.gen.yaml buf.build/acme/petapis

could not import github.com/mfridman/buf-generate-twirp-go/go/payment/v1alpha1 (no required module provides package "github.com/mfridman/buf-generate-twirp-go/go/payment/v1alpha1")
```

You could add the `--include-imports` option to buf generate, but this will result in a failure:

```
$ buf generate --include-imports --template buf.gen.yaml buf.build/acme/petapis

Failure: plugin "buf.build/demolab/plugins/twirp:v8.1.1-1" exited with non-zero status 1: 2022/04/21 17:54:47 error:files have conflicting go_package settings, must be the same: "money" and "paymentv1alpha1"
```

Note, the buf generate templates use `except` for googleapis. [See the docs for more info](https://docs.buf.build/tour/use-managed-mode#remove-modules-from-managed-mode).

---

The solution in most cases is to split up code generation into multiple steps, whereby we generate the dependent code using a separate template. This can get out of hand depending on your dependency graph.

In this example [`acme/petapis`](https://buf.build/acme/petapis) module depends on the [`acme/paymentapis`](https://buf.build/acme/paymentapis) module.

https://buf.build/acme/petapis/file/main:pet/v1/pet.proto#L5

Each of those modules has a different package and even more each module depends on `money.proto` and `datetime.proto` from googleapis.

https://buf.build/acme/petapis/file/main:pet/v1/pet.proto#L6
https://buf.build/acme/paymentapis/file/main:payment/v1alpha1/payment.proto#L5

To get around the error above, we can use the `--template` option to generate the dependant Go code like so:

```
buf generate --template buf.gen-go.yaml buf.build/acme/paymentapis
```

To recap, instead of generating all the code in one go using `buf.gen.yaml`, instead use buf generate with the `--template` option to split it up. This is often necessary for plugins such as Twirp.

Here is a working example. Make sure to [install the buf cli](https://docs.buf.build/installation).

Generated code is purposefully not check in, you generate it :)

```
rm -rf go || true
buf generate --template buf.gen-go.yaml buf.build/acme/paymentapis
buf generate --template buf.gen.yaml buf.build/acme/petapis

go run cmd/server/main.go
```

Output:

```
dante PET_TYPE_DOG
dante PET_TYPE_DOG
dante PET_TYPE_DOG
... every 2 seconds
```
