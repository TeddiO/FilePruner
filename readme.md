# FilePruner

**FilePruner** is a Go-based tool that scans directories recursively and deletes files older than a specified time. It allows configuration of file types, directories to scan, and deletion criteria through a YAML configuration file.

## Usage

### Configuration

FilePruner uses a YAML configuration file to define the directories, file types, and deletion time thresholds. By default, it looks for a configuration file named `filepruner-config.yml` in the working directory. You can override this using either an environment variable or a command-line argument.

#### Sample YAML Configuration (`filepruner-config.yml`)

```yaml
directories:
  - "C:/path/to/directory1"
  - "C:/path/to/directory2"
file_types:
  - ".txt"
  - ".log"
delete_after: "48h"  # Go duration format (e.g. 24h, 30m, 15s)
```

- **directories**: List of directories to scan.
- **file_types**: List of file types (extensions) to target for deletion.
- **delete_after**: Specifies the file age threshold for deletion using Go's duration format. E.g., `24h` for 24 hours, `30m` for 30 minutes, `15s` for 15 seconds.

### Launching the Application

To run FilePruner with the default configuration file (`filepruner-config.yml`), simply execute:

```bash
./filepruner
```

To preview the files that would be deleted without actually deleting them, enable dry run mode with the following command-line flag:

```bash
./filepruner -dry-run
```

You can also enable dry run mode via the `FILEPRUNER_DRY_RUN` environment variable:

```bash
export FILEPRUNER_DRY_RUN=true
./filepruner
```

To specify a different configuration file via a command-line argument, use the `-config` flag:

```bash
./filepruner -config="/path/to/your-config.yml"
```

You can also specify the configuration file using the `FILEPRUNER_CONFIG` environment variable:

```bash
export FILEPRUNER_CONFIG="/path/to/your-config.yml"
./filepruner
```

If both the environment variable and the command-line argument are used, the command-line argument takes precedence.

### Example

To prune files in your directories based on a `custom-config.yml` configuration file, you could run:

```bash
./filepruner -config="custom-config.yml"
```

Or by setting an environment variable:

```bash
export FILEPRUNER_CONFIG="custom-config.yml"
./filepruner
```

### Dry Run
To simulate deletion without actually deleting files, you can either:

- Use the dry run command-line flag:

  ```bash
  ./filepruner -config="custom-config.yml" -dry-run
  ```

- Or set the dry run environment variable:

  ```bash
  export FILEPRUNER_DRY_RUN=true
  ./filepruner -config="custom-config.yml"
  ```