# Copyright 2019 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

PROJECT_ID="go-gopher-run"
BUCKET_NAME="go-gopher-run-ml"
REGION="us-central1"
gcloud config set project $PROJECT_ID
gcloud config set compute/region $REGION

TRAINING_DATA_PATH="gs://${BUCKET_NAME}/pldata.csv"

# Specify the Docker container URI specific to the algorithm.
IMAGE_URI="gcr.io/cloud-ml-algos/boosted_trees:latest"

DATASET_NAME="playerdata"
ALGORITHM="linear"
MODEL_TYPE="classification"
MODEL_NAME="${DATASET_NAME}_${ALGORITHM}_${MODEL_TYPE}"
# Constant, so newer models override older ones and the deployed model won't need different directory
JOB_DIR="gs://${BUCKET_NAME}/algorithms_training/${MODEL_NAME}/"
FRAMEWORK="XGBOOST"

while true; do
  DATE="$(date '+%Y%m%d_%H%M%S')"
  # Include date so JOB_ID is unique
  JOB_ID="${MODEL_NAME}_${DATE}"
  gcloud beta ai-platform jobs submit training $JOB_ID \
  --master-image-uri=$IMAGE_URI --scale-tier=BASIC --job-dir=$JOB_DIR \
  -- \
  --preprocess --training_data_path=$TRAINING_DATA_PATH --objective=multi:softmax --num_class=4

  VERSION_NAME="V_${DATE}"
  gcloud ai-platform versions create $VERSION_NAME \
  --model ${MODEL_NAME} \
  --origin "${JOB_DIR}model/" \
  --runtime-version=1.11 \
  --framework ${FRAMEWORK} \
  --python-version=2.7
  gcloud ai-platform versions set-default ${VERSION_NAME} --model=${MODEL_NAME}

  # Repeate every 30 minutes
  sleep 1800
done

# bash training.sh &> /dev/null &
# --model_type=$MODEL_TYPE --batch_size=250 --learning_rate=0.1 --max_steps=1000
