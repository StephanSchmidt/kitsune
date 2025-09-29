#!/bin/bash

# Generate coverage report
go test -coverprofile=coverage.out ./...
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')

# Remove the old coverage section if it exists (everything after the coverage marker)
sed -i '/<!-- COVERAGE:START -->/,/<!-- COVERAGE:END -->/d' README.md

# If no coverage section exists, just remove any trailing coverage info
sed -i '/^## Code Coverage$/,$d' README.md

# Append the new coverage section
cat >> README.md << EOF

## Code Coverage

<!-- COVERAGE:START -->
**Current Coverage: ${COVERAGE}**

Last updated: $(date)
<!-- COVERAGE:END -->
EOF

echo "Coverage updated in README.md: ${COVERAGE}"

# Clean up
rm -f coverage.out