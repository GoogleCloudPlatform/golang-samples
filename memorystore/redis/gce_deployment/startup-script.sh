#! /bin/bash
# [START memorystore_startup_script_sh]
set -ex

# Talk to the metadata server to get the project id and location of application binary.
PROJECTID=$(curl -s "http://metadata.google.internal/computeMetadata/v1/project/project-id" -H "Metadata-Flavor: Google")
GCS_BUCKET_NAME=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/gcs-bucket" -H "Metadata-Flavor: Google")
REDISHOST=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/redis-host" -H "Metadata-Flavor: Google")
REDISPORT=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/redis-port" -H "Metadata-Flavor: Google")

# Install dependencies from apt
apt-get update
apt-get install -yq ca-certificates supervisor

# Install logging monitor. The monitor will automatically pickup logs send to
# syslog.
curl -s "https://storage.googleapis.com/signals-agents/logging/google-fluentd-install.sh" | bash
service google-fluentd restart &

gsutil cp gs://"$GCS_BUCKET_NAME"/gce/app.tar /app.tar
mkdir -p /app
tar -x -f /app.tar -C /app
chmod +x /app/app

# Create a goapp user. The application will run as this user.
getent passwd goapp || useradd -m -d /home/goapp goapp
chown -R goapp:goapp /app

# Configure supervisor to run the Go app.
cat >/etc/supervisor/conf.d/goapp.conf << EOF
[program:goapp]
directory=/app
environment=HOME="/home/goapp",USER="goapp",REDISHOST=$REDISHOST,REDISPORT=$REDISPORT
command=/app/app
autostart=true
autorestart=true
user=goapp
stdout_logfile=syslog
stderr_logfile=syslog
EOF

supervisorctl reread
supervisorctl update
# [END memorystore_startup_script_sh]
