from flask import Flask, request, jsonify
from flask_cors import CORS
from product_script import get_response

app = Flask(__name__)
CORS(app, resources={r"/*": {"origins": "*"}})

@app.route("/", methods=["POST"])
def handle_product():
    try:
        app.logger.debug("Received request: %s", request.args)

        data = request.get_json()
        app.logger.debug("Received JSON data: %s", data)
        if not data:
            return jsonify({"error": "Request body is required"}), 400

        message = data.get("message")
        history = data.get("history")
        product = data.get("product")

        app.logger.debug("Calling get_response with message: %s, product: %s, history: %s", message, product, history)
        response = get_response(message, product, history)
        app.logger.debug("Received response: %s", response)
        print(response)
        
        return jsonify(response)
    except Exception as e:
        app.logger.error(f"Error handling request: {e}")
        return jsonify({"error": "Internal Server Error"}), 500

app.run(host="0.0.0.0", port=8000, debug=True)