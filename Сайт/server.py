from flask import Flask, request, jsonify
from flask_cors import CORS
from info_player import get_player_data

app = Flask(__name__)
CORS(app, resources={
    r"/api/*": {
        "origins": ["*"],  # для теста
        "methods": ["GET", "OPTIONS"],
        "allow_headers": ["Content-Type", "Authorization"]
    }
})

@app.route('/api/stats', methods=['OPTIONS'])
def options():
    response = jsonify({})
    response.headers['Access-Control-Allow-Origin'] = '*'
    response.headers['Access-Control-Allow-Methods'] = 'GET, OPTIONS'
    response.headers['Access-Control-Allow-Headers'] = 'Content-Type'
    return response

@app.route('/api/stats', methods=['GET'])
def stats():
    nick = request.args.get('nick', '').strip()
    if not nick:
        return jsonify({"error": "Никнейм не указан"}), 400
    try:
        data = get_player_data(nick)  # ← без await
        return jsonify(data)
    except Exception as e:
        return jsonify({"error": f"Ошибка: {str(e)}"}), 500

if __name__ == '__main__':
    app.run(port=5000, debug=True)