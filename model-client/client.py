import pathlib
import numpy as np
import requests
import json
from PIL import Image
import cv2
import tensorflow as tf
import os
import time
import random

requests.packages.urllib3.disable_warnings()
import ssl
import urllib.request

model_server_ip = os.environ['HOST_IP']
model_server_port = os.environ['SERVER_PORT']
model_image = os.environ['MODEL_IMAGE']
image_urls = json.loads(os.environ.get('IMAGE_PATHS'))
video_stream = os.environ.get('VIDEO') == 'true'
seed = os.environ['RANDOM_SEED']

random.seed(seed)
cap = None

if video_stream:
    cap = cv2.VideoCapture(0)    
    #Check if camera was opened correctly
    if not (cap.isOpened()):
        print("Could not open video device")
        exit()
    cap.set(cv2.CAP_PROP_FRAME_WIDTH, 400)
    cap.set(cv2.CAP_PROP_FRAME_HEIGHT, 300)
    print("Taking data from /dev/video0")
else:
    # Fetch images.
    image_paths = []
    for image in image_urls:
        print("Fetching image at %s" % (image), flush=True)
        filename = image.split('/')[-1]
        urllib.request.urlretrieve(image, filename)
        test_image_path = pathlib.Path(filename)
        image_paths.append(test_image_path)    
    print("Taking data from env=IMAGE_PATHS")

print("Hitting model server at %s" % (model_server_ip), flush=True)
try:
    _create_unverified_https_context = ssl._create_unverified_context
except AttributeError:
    # Legacy Python that doesn't verify HTTPS certificates by default
    pass
else:
    # Handle target environment that doesn't support HTTPS verification
    ssl._create_default_https_context = _create_unverified_https_context

def get_next_image():
    if video_stream:
        ret, frame = cap.read()
        if not ret:
            print("Couldn't capture frame")
            return None
        return frame
    else:
        rand_img = random.choice(image_paths)
        image_np = np.array(Image.open(rand_img))
        return image_np


while True:
    images = []
    image_np = get_next_image()
    if image_np is None:
        continue

    images.append(image_np.tolist())

    data = json.dumps({"instances": images})

    headers = {"content-type": "application/json"}
    json_response = requests.post('http://%s:%s/v1/models/%s:predict' % (model_server_ip, model_server_port, model_image), data=data, headers=headers)
    resp = json_response.json()
    resp['timestamp'] = time.time()
    print(json.dumps(resp), flush=True)
    time.sleep(3)
                     
