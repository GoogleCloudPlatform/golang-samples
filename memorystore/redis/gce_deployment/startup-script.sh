#! /bin/bash
# [START memorystore_startup_script_sh]
set -ex

# Talk to the metadata server to get the project id and location of application binary.
PROJECTID=$(curl -s "http://metadata.google.internal/computeMetadata/v1/project/project-id" -H "Metadata-Flavor: Google")
export PROJECTID
GCS_BUCKET_NAME=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/gcs-bucket" -H "Metadata-Flavor: Google")
REDISHOST=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/redis-host" -H "Metadata-Flavor: Google")
REDISPORT=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/redis-port" -H "Metadata-Flavor: Google")

# Install dependencies from apt
apt-get update
apt-get install -yq ca-certificates supervisor

# Install logging monitor. The monitor will automatically pickup logs send to
# syslog.
curl "https://storage.googleapis.com/signals-agents/logging/google-fluentd-install.sh" --output google-fluentd-install.sh
checksum=$(sha256sum google-fluentd-install.sh | awk '{print $1;}')
if [ "$checksum" != "ec78e9067f45f6653a6749cf922dbc9d79f80027d098c90da02f71532b5cc967" ]; then
    echo "Checksum does not match"
    exit 1
fi
chmod +x google-fluentd-install.sh && ./google-fluentd-install.sh
service google-fluentd restart &

gcloud storage cp gs://"$GCS_BUCKET_NAME"/gce/app.tar /app.tar
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
