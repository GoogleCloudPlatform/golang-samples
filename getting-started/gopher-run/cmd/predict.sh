PROJECT_ID="maralder-start"
BUCKET_NAME="maralder-start-ml"
REGION="us-central1"
gcloud config set project $PROJECT_ID
gcloud config set compute/region $REGION

TRAINING_DATA_PATH="gs://${BUCKET_NAME}/pldata.csv"

# Specify the Docker container URI specific to the algorithm.
IMAGE_URI="gcr.io/cloud-ml-algos/linear_learner_cpu:latest"

DATASET_NAME="playerdata"
ALGORITHM="linear"
MODEL_TYPE="classification"
MODEL_NAME="${DATASET_NAME}_${ALGORITHM}_${MODEL_TYPE}"
# Not unique, so newer versions override older ones
JOB_DIR="gs://${BUCKET_NAME}/algorithms_training/${MODEL_NAME}/"
FRAMEWORK="TENSORFLOW"

# ID must be unique
DATE="$(date '+%Y%m%d_%H%M%S')"
JOB_ID="${MODEL_NAME}_${DATE}"
# gcloud beta ai-platform jobs submit training $JOB_ID \
# --master-image-uri=$IMAGE_URI --scale-tier=BASIC --job-dir=$JOB_DIR \
# -- \
# --preprocess --model_type=$MODEL_TYPE --batch_size=250 \
# --learning_rate=0.1 --max_steps=1000 --training_data_path=$TRAINING_DATA_PATH

VERSION_NAME="V_${DATE}"
gcloud ai-platform versions create $VERSION_NAME \
--model $MODEL_NAME \
--origin "${JOB_DIR}model/" \
--runtime-version=1.13 \
--framework $FRAMEWORK \
--python-version=3.5

ACT = ${1}
G = ${2}
G = ${3}
G = ${4}
G = ${5}
G = ${6}
G = ${7}
G = ${8}