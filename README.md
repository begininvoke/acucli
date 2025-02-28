# Acucli

A powerful command-line interface tool for Acunetix that streamlines web application security scanning workflows. Acucli allows you to manage targets, scans, reports, and automate entire security assessment processes directly from your terminal.

## Features

- **Target Management**
  - Create, list, and delete targets
  - Configure target settings
  - Manage target groups
- **Scan Management**
  - Configure and run scans
  - Monitor scan progress
  - Manage scan profiles
- **Report Generation**
  - Generate HTML and CSV reports
  - Manage report templates
  - Download report files
- **Automation**
  - Automate entire scanning workflow
  - Automatic resource cleanup
  - Pipeline integration support
- **Configuration**
  - YAML-based configuration
  - Environment variable support
  - Flexible API settings

## Installation

1. Ensure you have Go installed on your system
2. Install Acucli using Go:
   ```bash
   go install github.com/tosbaa/acucli@latest
   ```
3. Download the configuration template:
   ```bash
   # Download .acucli.yaml and place it in your home directory
   curl -O https://raw.githubusercontent.com/tosbaa/acucli/main/.acucli.yaml
   mv .acucli.yaml ~/
   ```
4. Configure your API settings in `~/.acucli.yaml`

## Usage

### Global Flags

- `--config, -c`: Specify config file (default: $HOME/.acucli.yaml)
- `--version, -v`: Show version information
- `--help, -h`: Show help information

### Target Management

```bash
# List all targets
acucli target list

# Add a target
echo "https://target.com" | acucli target add

# Add multiple targets to a group
cat targets.txt | acucli target add --gid=<TARGETGROUP-ID>

# Get target information
acucli target --id <TARGET-ID>

# Set target configuration
echo "<TARGET-ID>" | acucli target setConfig

# Remove a target
echo "<TARGET-ID>" | acucli target remove
```

### Target Group Management

```bash
# Create a target group
echo "TargetGroupName" | acucli targetGroup add

# List target groups
acucli targetGroup list

# Get group information
acucli targetGroup --id <TARGETGROUP-ID>

# Remove a target group
echo "<TARGETGROUP-ID>" | acucli targetGroup remove
```

### Scan Profile Management

```bash
# List scan profiles
acucli scanProfile list

# Get profile information
acucli scanProfile --id=<SCANPROFILE-ID>

# Export a scan profile
acucli scanProfile --id=<SCANPROFILE-ID> -e

# Import a scan profile
cat profile.json | acucli scanProfile add

# Remove a scan profile
echo "<SCANPROFILE-ID>" | acucli scanProfile remove
```

### Scan Management

```bash
# Start a scan for single target
acucli scan --targetID=<TARGET-ID> --scanProfileID=<SCANPROFILE-ID>

# Start scans for multiple targets
cat targets.txt | acucli scan --scanProfileID=<SCANPROFILE-ID>

# Start scans for a target group
acucli targetGroup --id=<TARGETGROUP-ID> | cut -f2 | acucli scan --scanProfileID=<SCANPROFILE-ID>
```

### Report Management

```bash
# List all reports
acucli report list

# Generate a report
echo "<SCAN-ID>" | acucli report generate --templateID=<TEMPLATE-ID>

# Get report details
echo "<REPORT-ID>" | acucli report get

# Remove a report
echo "<REPORT-ID>" | acucli report remove
```

### Automated Workflow

The `auto` command automates the entire scanning process in one command:

```bash
# Basic usage
acucli auto --target=https://example.com

# Advanced usage
acucli auto \
  --target=https://example.com \
  --scanProfileID=<SCAN-PROFILE-ID> \
  --reportTemplateID=<REPORT-TEMPLATE-ID> \
  --format=html \
  --output=/path/to/output/report.html \
  --timeout=600
```

#### Auto Command Workflow

1. Adds target with specified URL
2. Verifies target creation
3. Starts scan with specified profile
4. Monitors scan progress
5. Generates report/export
6. Downloads report files
7. Cleans up resources automatically

#### Auto Command Options

- `--target, -u`: Target URL to scan (required)
- `--format, -f`: Output format (html or csv, default: html)
- `--output, -o`: Output path for report files
- `--timeout, -i`: Timeout in seconds (default: 800)
- `--scanProfileID, -s`: Custom scan profile ID
- `--reportTemplateID, -r`: Custom report template ID

## Advanced Usage

### Pipeline Integration

```bash
# Scan all targets in a group and remove them afterward
acucli targetGroup --id=<TARGETGROUP-ID> | cut -f2 | tee >(acucli scan --scanProfileID=<SCANPROFILE-ID>) | acucli target remove

# Generate reports for multiple scans
acucli scan list | cut -f1 | acucli report generate --templateID=<TEMPLATE-ID>
```

### Configuration File (.acucli.yaml)

```yaml
API: "your-api-key-here"
URL: "https://your-acunetix-instance"
ScanConfig:
  target_timeout: 0
  max_scan_time: 0
  scanning_mode: "sequential"
  user_agent: "Mozilla/5.0..."
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For issues and feature requests, please use the [GitHub issue tracker](https://github.com/tosbaa/acucli/issues).
