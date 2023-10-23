# quickmock

Run a mock server quickly.

Read configuration from file or wireframe it interactively live.

## How to use

### Installation

quickmock comes as a single package - just download the binary for your platform and run it.
You can find the binaries in [the Releases section.](https://github.com/S-Maciejewski/quickmock/releases)

### Running

#### Batch mode

You can run quickmock with endpoints defined in a yaml file.

```bash
quickmock -d -f <your-definitions-file.yml>
```

_Note: -d flag runs quickmock in detached mode - without TUI._

The format of the file can be seen in `example.yml` - it's pretty self-explanatory

```yaml
- method: <your method>
  path: <your path>
  response:
    code: <code>
    content: <some content, JSON or not>
```

#### Interactive mode

This is setting up quickmock interactively using a TUI.

This way you'll be able to setup server on the fly depending on what you need.

```bash
quickmock
```

_Note: you can still use -f with interactive mode to pre-load some endpoint definitions_

## To do

- [ ] Interactive mode CRUD
    - [ ] Exporting a manually created definition
- [ ] JSON endpoint definition support
- [ ] Reading config from a swagger file / export to swagger
- [ ] Some basic documentation
- [ ] Minimize the binary size (replacing gox with gccgo in the
  pipeline?) https://stackoverflow.com/questions/3861634/how-to-reduce-go-compiled-file-size
