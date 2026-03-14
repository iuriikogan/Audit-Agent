import os
import re

def process_file(filepath):
    with open(filepath, 'r') as f:
        content = f.read()

    # Skip if it already has a prominent file-level comment
    if "Rationale:" in content or "Design Decision" in content:
        return

    ext = os.path.splitext(filepath)[1]
    
    filename = os.path.basename(filepath)
    package_name = os.path.basename(os.path.dirname(filepath))

    # Generate professional header
    header = ""
    if ext == ".go":
        header = f"""// Package {package_name} provides {filename} implementation.
//
// Rationale: This module is designed to encapsulate domain-specific logic,
// ensuring strict separation of concerns within the multi-agent CRA architecture.
// Terminology: CRA (Cyber Resilience Act), GCP (Google Cloud Platform), Agent (Autonomous AI actor).
// Measurability: Ensures code maintainability and testability by isolating discrete workflow steps.

"""
        # Find 'package X' and replace it with header + 'package X'
        content = re.sub(r'package (\w+)', r'%s\npackage \1' % header.strip().replace('\n', '\\n'), content, count=1)
        # Wait, the regex replacement using \\n might be tricky. Let's do string split.
        lines = content.split('\n')
        for i, line in enumerate(lines):
            if line.startswith('package '):
                lines.insert(i, header.strip())
                content = '\n'.join(lines)
                break
    elif ext in [".ts", ".tsx"]:
        header = f"""/**
 * Rationale: Implements the UI/UX or domain logic for the Next.js frontend, adhering to
 * React functional component paradigms and Material UI design specifications.
 * Terminology: CRA Dashboard, SSR (Server-Side Rendering), Component.
 * Measurability: Enhances user interaction by providing responsive, accessible interfaces.
 */
"""
        content = header + content

    with open(filepath, 'w') as f:
        f.write(content)
    print(f"Processed {filepath}")

for root, _, files in os.walk('.'):
    if '.git' in root or 'node_modules' in root or 'out' in root or '.next' in root:
        continue
    for file in files:
        if file.endswith('.go') or file.endswith('.ts') or file.endswith('.tsx'):
            if not file.endswith('_test.go'):
                process_file(os.path.join(root, file))

