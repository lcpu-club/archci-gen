# ArchLinux CI generator

A tool for generating Drone CI pipeline for Loong ArchLinux.

## Auto Mode

Requirements:

1. Triggered by a PR
2. Must have no conflicts with the origin/testing branch

All the pacakages changed in the PR will be built in parallel unless defined in the body of the PR. In the following example, `a`, `b` and `c` will be built in parallel, and `d` will be built after `a` and `b` are built.

```
a,b,c;d;e,f
```

## Manual Mode

Must be one of the following:

1. Specify the packages to be built in the `packages` field
2. Specify the commit range to be built in the `commit_range` field, in `begin_commit:end_commit` format
3. Set the `mode` field to `rebuild_all`


## Notes:

- Drone uses `httpsignatures.DefaultSha256Signer.AuthRequest("hmac-key", key, req)` to sign the request.