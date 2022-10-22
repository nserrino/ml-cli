import pathlib
import numpy as np
import requests
import json
from PIL import Image
import cv2
import os
import time
import random
import io
import base64

requests.packages.urllib3.disable_warnings()
import ssl
import urllib.request

host_ip = os.environ.get('HOST_IP')
converter_port = os.environ.get('CONVERTER_PORT')
model_image = os.environ.get('MODEL')
video_stream = os.environ.get('VIDEO_STREAM') == 'true'
image_urls = json.loads(os.environ.get('IMAGE_PATHS'))
seed = os.environ.get('RANDOM_SEED')
image_height = os.environ.get('IMAGE_HEIGHT') or 300
image_width = os.environ.get('IMAGE_WIDTH') or 400

cap = None

if video_stream:
    cap = cv2.VideoCapture(0)    
    #Check if camera was opened correctly
    if not (cap.isOpened()):
        print("Could not open video device")
        exit()
    cap.set(cv2.CAP_PROP_FRAME_WIDTH, image_width)
    cap.set(cv2.CAP_PROP_FRAME_HEIGHT, image_height)
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
    image_np = get_next_image()
    if image_np is None:
        continue

    im = Image.fromarray(image_np)
    im = im.resize((image_width, image_height))

    with io.BytesIO() as output:
        im.save(output, format='JPEG')
        encoded = base64.b64encode(output.getvalue())
        headers = {'Content-type': 'application/json', 'Accept': 'text/plain'}
        payload = json.dumps({"image": encoded.decode('utf-8'), "model": model_image})
        converter_server = 'http://%s:%s/convert' % (host_ip, converter_port)
        response = requests.post(converter_server, data=payload, headers=headers)        

        try:
            data = response.json()     
            data['timestamp'] = time.time()
            print(json.dumps(data), flush=True)         
        except requests.exceptions.RequestException:
            print(response.text)

    time.sleep(3)
                     
