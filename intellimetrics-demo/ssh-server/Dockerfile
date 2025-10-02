FROM alpine:latest

# VULNERABILITY: Weak SSH configuration for penetration testing practice
RUN apk add --no-cache openssh bash curl netcat-openbsd python3 nmap

# Configure SSH
RUN ssh-keygen -A

# VULNERABILITY: Create users with weak passwords
# Perfect for password spraying and pivoting exercises
RUN adduser -D -s /bin/bash admin && echo "admin:Admin123!" | chpasswd
RUN adduser -D -s /bin/bash developer && echo "developer:Dev2024!" | chpasswd
RUN adduser -D -s /bin/bash backup && echo "backup:Backup123" | chpasswd
RUN adduser -D -s /bin/bash testuser && echo "testuser:test" | chpasswd
RUN adduser -D -s /bin/bash dbadmin && echo "dbadmin:Database!" | chpasswd

# VULNERABILITY: Grant sudo access (for privilege escalation practice)
RUN apk add --no-cache sudo
RUN echo "admin ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers
RUN echo "developer ALL=(ALL) NOPASSWD: /usr/bin/python3" >> /etc/sudoers

# Create interesting files for post-exploitation
RUN mkdir -p /home/admin/documents /home/admin/.ssh
RUN echo "Internal Network: 172.18.0.0/16" > /home/admin/documents/network.txt
RUN echo "Database: mongo:27017" > /home/admin/documents/infrastructure.txt
RUN echo "API Keys stored in: /opt/secrets/" > /home/admin/documents/notes.txt

# Create fake sensitive data
RUN mkdir -p /opt/secrets
RUN echo "PROD_DB_PASSWORD=ProductionSecret2024!" > /opt/secrets/prod.env
RUN echo "API_MASTER_KEY=mk_live_abc123xyz789" > /opt/secrets/api_keys.txt
RUN chmod 600 /opt/secrets/*
RUN chown -R admin:admin /opt/secrets

# VULNERABILITY: SSH configuration with weak settings
RUN echo "PermitRootLogin yes" >> /etc/ssh/sshd_config
RUN echo "PasswordAuthentication yes" >> /etc/ssh/sshd_config
RUN echo "PermitEmptyPasswords no" >> /etc/ssh/sshd_config
RUN echo "ChallengeResponseAuthentication no" >> /etc/ssh/sshd_config
RUN echo "PubkeyAuthentication yes" >> /etc/ssh/sshd_config
RUN echo "AllowTcpForwarding yes" >> /etc/ssh/sshd_config
RUN echo "GatewayPorts yes" >> /etc/ssh/sshd_config
RUN echo "X11Forwarding no" >> /etc/ssh/sshd_config

# Create .bash_history with interesting commands (for post-exploitation reconnaissance)
RUN echo "mysql -u root -p" > /home/admin/.bash_history
RUN echo "ssh dbadmin@mongo" >> /home/admin/.bash_history
RUN echo "cat /opt/secrets/prod.env" >> /home/admin/.bash_history
RUN echo "curl http://api:8080/api/admin/config" >> /home/admin/.bash_history
RUN chown admin:admin /home/admin/.bash_history

EXPOSE 22

CMD ["/usr/sbin/sshd", "-D", "-e"]
