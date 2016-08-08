# Copyright 2016 Google Inc. All rights reserved.
# Use of this source code is governed by the Apache 2.0
# license that can be found in the LICENSE file.

""" Sample pyspark script to be uploaded to Cloud Storage and run on
Cloud Dataproc.

Note this file is not intended to be run directly, but run inside a PySpark
environment.
"""

# [START pyspark]
import pyspark

sc = pyspark.SparkContext()
rdd = sc.parallelize(['Hello,', 'world!', 'dog', 'elephant', 'panther'])
words = sorted(rdd.collect())
print words
# [END pyspark]
