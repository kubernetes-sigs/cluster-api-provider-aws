#!/bin/sh

if [ -z "$AWS_ACCESS_KEY_ID" ]; then
    echo "error: AWS_ACCESS_KEY_ID is not set in the environment" 2>&1
    exit 1
fi

if [ -z "$AWS_SECRET_ACCESS_KEY" ]; then
    echo "error: AWS_SECRET_ACCESS_KEY is not set in the environment" 2>&1
    exit 1
fi

script_dir="$(cd $(dirname "${BASH_SOURCE[0]}") && pwd -P)"

secrethash=$(cat $script_dir/bootstrap.sh | \
  sed "s/  aws_access_key_id: FILLIN/  aws_access_key_id: $(echo -n $AWS_ACCESS_KEY_ID | base64)/" | \
  sed "s/  aws_secret_access_key: FILLIN/  aws_secret_access_key: $(echo -n $AWS_SECRET_ACCESS_KEY | base64)/" | \
  base64 --w=0)

cat <<EOF > $script_dir/bootstrap.yaml
apiVersion: v1
kind: Secret
metadata:
  name: master-user-data-secret
  namespace: default
type: Opaque
data:
  userData: $secrethash
EOF
