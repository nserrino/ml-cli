import io
import json                    
import base64                  
import logging             
import numpy as np
import requests
from PIL import Image

from flask import Flask, request, jsonify, abort

app = Flask(__name__)          
app.logger.setLevel(logging.DEBUG)
  
  
@app.route('/', defaults={'u_path': ''}, methods=['POST'])
@app.route('/<path:u_path>', methods=['POST'])
def catch_all(u_path):
    if not request.json or 'image' not in request.json or 'model_server' not in request.json:
        abort(400)
             
    print(len(request.json['image']))
    # get the base64 encoded string
    im_b64 = request.json['image'].encode('utf-8')

    decoded = base64.b64decode(im_b64)
    foo = Image.open(io.BytesIO(decoded))
    img_arr = np.asarray(foo)
    images = []
    images.append(img_arr.tolist())
    data = json.dumps({"instances": images})

    headers = {"content-type": "application/json"}
    json_response = requests.post(request.json['model_server'], data=data, headers=headers)
    resp = json_response.json()
    return json_response.json()
  
  
def run_server_api():
    app.run(host='0.0.0.0', port=1234)
  
  
if __name__ == "__main__":     
    run_server_api()