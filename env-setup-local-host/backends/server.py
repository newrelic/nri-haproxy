#!/usr/bin/env python3
"""Minimal HTTP server used as a HAProxy backend.

Usage: python3 server.py <port> [name]

Returns a simple text response identifying the server.
"""
import sys
from http.server import HTTPServer, BaseHTTPRequestHandler


class Handler(BaseHTTPRequestHandler):
    server_name = "backend"

    def do_GET(self):
        try:
            self.send_response(200)
            self.send_header("Content-Type", "text/plain")
            self.end_headers()
            self.wfile.write(f"Hello from {self.server_name} (path={self.path})\n".encode())
        except BrokenPipeError:
            # HAProxy health checks close the connection early — harmless.
            pass

    def log_message(self, format, *args):
        # Suppress request logs to keep terminal clean
        pass

    def handle(self):
        try:
            super().handle()
        except BrokenPipeError:
            pass


def main():
    port = int(sys.argv[1]) if len(sys.argv) > 1 else 9001
    name = sys.argv[2] if len(sys.argv) > 2 else f"backend-{port}"
    Handler.server_name = name
    server = HTTPServer(("127.0.0.1", port), Handler)
    print(f"  Backend '{name}' listening on 127.0.0.1:{port}")
    server.serve_forever()


if __name__ == "__main__":
    main()
