# Gat
Gat (reverse of tag) is a tool that bumps your repository version using semantic versioning
Gat supports `v` in front of golang libraries

# Installation
```
go install github.com/B3rs/gat@latest
```

# Usage
```
gat minor/major/patch
```

# Known issues
```
ssh: handshake failed: knownhosts: key mismatch
```
your known_hosts file can be outdated or missing the repository fingerprint, add it by running the following command:
```
ssh-keyscan -H github.com >> ~/.ssh/known_hosts
```
obviously change github.com with your repo