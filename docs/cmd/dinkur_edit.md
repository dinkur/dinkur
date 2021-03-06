## dinkur edit

Edit the latest or a specific entry

### Synopsis

Applies changes to the currently active entry, or the latest entry, or
a specific entry using the --id or -i flag.

```
dinkur edit [new name of entry] [flags]
```

### Options

```
  -a, --after-id uint    sets --start time to the end time of entry with ID
  -L, --after-last       sets --start time to the end time of latest entry
  -z, --append           add name to the end of the existing name, instead of replacing it
  -b, --before-id uint   sets --end time to the start time of entry with ID
  -e, --end time         end time of entry; entry will be unmarked as active if set
  -h, --help             help for edit
  -i, --id uint          ID of entry (default is active or latest entry)
  -s, --start time       start time of entry
```

### Options inherited from parent commands

```
      --client string         Dinkur client: "db" or "grpc" (default "db")
      --color string          colored output: "auto", "always", or "never" (default "auto")
      --config string         config file (default "~/.config/dinkur/config.yaml")
      --data string           database file (default "/home/kalle/.local/share/dinkur/dinkur.db")
      --data-mkdir            create directory for data if it doesn't exist (default true)
      --grpc-address string   address of Dinkur daemon gRPC API (default "localhost:59122")
  -v, --verbose               enables debug logging
```

### SEE ALSO

* [dinkur](dinkur.md)	 - The Dinkur CLI

###### Auto generated by spf13/cobra on 20-Jan-2022
