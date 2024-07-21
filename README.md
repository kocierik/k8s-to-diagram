# K8s Diagram architecture generator

A tool for visualizing Kubernetes infrastructure manifest using [Mermaid](https://github.com/dreampuf/mermaid.go).
<div align="center">
    
![diagram generated](https://raw.githubusercontent.com/kocierik/k8s-to-diagram/main/images/k8s_infrastructure.png?token=GHSAT0AAAAAACSO5SDVETGPH3KZGPAYPZ4GZU457AA)

</div>

## Overview

This program reads Kubernetes YAML manifests and generates a visual diagram of the infrastructure and its connections. By parsing the `communication` annotations from the Kubernetes manifests, it creates a diagram showing how different services and deployments interact with each other.

## Features

- **Parse Kubernetes YAML manifests**: Reads and parses Kubernetes manifests from a directory.
- **Generate diagrams**: Creates a visual representation of the infrastructure based on communication annotations.
- **Output formats**: Generates diagrams in both SVG and PNG formats.

## Installation

1. **Clone the repository:**

    ```bash
    git clone git@github.com:kocierik/k8s-to-diagram.git
    cd k8s-to-diagram
    ```

2. **Build the application:**

    ```bash
    go build -o k8s-mermaid
    ```

## Configuration

Ensure that your Kubernetes YAML manifests include `communication` annotations in the following example format:

```yaml
metadata:
  annotations:
    communication: |
      {
        "name": "service-name-0",
        "inbound": [
          {"service": "service-name-1", "port": 1234}
        ],
        "outbound": [
          {"service": "service-name-2", "port": 1234}
        ]
      }
```
## Troubleshooting
- **Ensure valid YAML**: Make sure your YAML files are correctly formatted and valid.
- **Correct directory**: Verify that the manifest directory path is correctly specified.
- **Check annotations**: Ensure that communication annotations are correctly formatted and included in the manifests.

## Contributing

Feel free to contribute by submitting issues or pull requests. Please adhere to the coding standards and ensure your contributions are well-tested.


## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/kocierik/k8s-to-diagram/blob/main/LICENSE) file for details.
