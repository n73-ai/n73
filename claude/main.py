from flask import Flask, request, jsonify
import threading
import asyncio
from claude import NewProject, ResumeProject
import os


app = Flask(__name__)


def run_async_task(coroutine_fn, *args):
    """Run an async function in a new event loop."""
    asyncio.run(coroutine_fn(*args))


def start_background_task(coroutine_fn, *args):
    thread = threading.Thread(target=run_async_task, args=(coroutine_fn, *args), daemon=True)
    thread.start()
    return thread


def validate_fields(data, required_fields):
    """Validate that all required fields exist in the request."""
    missing = [field for field in required_fields if field not in data]
    if missing:
        return False, f"Missing fields: {', '.join(missing)}"
    return True, ""


@app.route("/claude/new", methods=["POST"])
def new_endpoint():
    data = request.get_json()
    valid, error = validate_fields(data, ["prompt", "model", "work_dir", "webhook_url", "jwt"])
    if not valid:
        return jsonify({"error": error}), 400

    start_background_task(NewProject, data["prompt"], data["model"], data["work_dir"], data["webhook_url"], data["jwt"])
    return jsonify({"status": "processing"}), 200


@app.route("/claude/resume", methods=["POST"])
def resume_endpoint():
    data = request.get_json()
    valid, error = validate_fields(data, ["prompt", "model", "work_dir", "webhook_url", "session_id", "jwt"])
    if not valid:
        return jsonify({"error": error}), 400

    start_background_task(ResumeProject, data["prompt"], data["model"], data["work_dir"], data["webhook_url"], data["session_id"], data["jwt"])
    return jsonify({"status": "processing"}), 200


@app.route('/health', methods=['GET'])
def health():
    return '', 200


port = int(os.getenv('PORT', 5000))


if __name__ == "__main__":
    app.run(host="::", port=port, debug=True)
