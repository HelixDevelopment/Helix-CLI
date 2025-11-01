#!/bin/bash

# Generate SSH keys for test workers
set -e

KEY_DIR="test/workers/ssh-keys"

# Create directory if it doesn't exist
mkdir -p "$KEY_DIR"

# Generate SSH key pair if it doesn't exist
if [ ! -f "$KEY_DIR/id_rsa" ]; then
    echo "Generating SSH key pair for test workers..."
    ssh-keygen -t rsa -b 4096 -f "$KEY_DIR/id_rsa" -N "" -q
    echo "✅ SSH keys generated in $KEY_DIR/"
else
    echo "✅ SSH keys already exist in $KEY_DIR/"
fi

# Set proper permissions
chmod 600 "$KEY_DIR/id_rsa"
chmod 644 "$KEY_DIR/id_rsa.pub"

# Create known_hosts file with test worker entries
echo "Creating known_hosts file..."
cat > "$KEY_DIR/known_hosts" << EOF
[test-worker-1]:2222 ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDE...
[test-worker-2]:2223 ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDE...
[test-worker-3]:2224 ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDE...
EOF

echo "✅ Test SSH setup complete!"