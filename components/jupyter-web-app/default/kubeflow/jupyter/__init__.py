from flask import Flask, url_for, jsonify

app = Flask(__name__)

from kubeflow.jupyter import routes
