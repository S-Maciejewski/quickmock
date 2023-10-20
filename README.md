# quickmock

Run a mock server quickly. 

Read configuration from file or wireframe it interactively live.

## How to use
### Installation
quickmock comes as a single package - just download the binary for your platform and run it.

### Running
#### Batch mode
This is setting up quickmock based on endpoint definitions in a yaml file.
```bash
quickmock -f <your-definitions-file.yml>
```

#### Interactive mode (TODO)
This is setting up quickmock interactively using a TUI. 

This way you'll be able to setup server on the fly depending on what you need.
```bash
quickmock
```


## To do
- [ ] CI/CD, building for multiple platforms
- [ ] License
- [ ] Interactive mode (bubbletea?)
  - [ ] Exporting a manually created definition
- [ ] JSON endpoint definition support
- [ ] Reading config from a swagger file
- [ ] Some basic documentation
