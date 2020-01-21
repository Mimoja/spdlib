# spdlib
A golang library and tool to extract SPD (4) binaries from firmware blobs

# Running
```
go build cmd/spdutil.go
mkdir test
./spdutil -o test firmware.bin
```

You will need an up-to-date libyara installed.
