import pandas as pd
import numpy as np
import json

# CHANGE THESE BASED ON YOUR REPLAY SCENARIO
filename = 'replay-out.json'
models = ['original', 'fpn']
mentor = 'original'

all_results = {'top5': [], 'top1': [], 'latency_ms': [], 'model': [], 'request_idx': []}

# UPDATE THIS BASED ON YOUR MODEL OUTPUT FORMAT
def parse_results(input):#
  body = json.loads(input['resp_body'])['predictions'][0]
  latency_ms = input['latency_ms']
  num_detections = int(body['num_detections'])
  detection_scores = body['detection_scores'][0:num_detections]
  detection_classes = body['detection_classes'][0:num_detections]
  return {
      'latency_ms': latency_ms,
      'num_detections': num_detections,
      'top1': None if num_detections == 0 else detection_classes[0],
      'detection_scores': set(detection_scores),
      'detection_classes': set(detection_classes)
  }

with open(filename) as f:
    for line in f:
        content = json.loads(line)
        mentor_results = content['results'][mentor]
        request_idx = content['request_idx']
        mentor_results = parse_results(mentor_results)

        for model in models:
          model_results = parse_results(content['results'][model])
          all_results['request_idx'].append(request_idx)
          all_results['model'].append(model)
          all_results['top5'].append(mentor_results['top1'] in model_results['detection_classes'])
          all_results['top1'].append(mentor_results['top1'] == model_results['top1'])
          all_results['latency_ms'].append(model_results['latency_ms'])


df = pd.DataFrame(data=all_results)
df = df.groupby('model').agg({'top5': 'mean', 'top1': 'mean', 'latency_ms': 'mean'})
print(df)
