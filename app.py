from flask import Flask, request, jsonify
from flask_cors import CORS
from product_script import get_response

app = Flask(__name__)
CORS(app)

@app.route("/", methods=["POST"])
def handle_product():
    product = request.args.get("product")
    message = request.get_json().get("message")
    history = request.get_json().get("history")
    
    
    response = get_response(message, product, history)
    return jsonify(response)


app.run(host="0.0.0.0", port=8080)