from http.server import BaseHTTPRequestHandler, HTTPServer
import RPi.GPIO as GPIO
import socket

serverPort = 8080
pins = (12, 13, 19)  # (R, G, B)

colors = {"red":    (GPIO.HIGH, GPIO.LOW, GPIO.LOW),
          "green":  (GPIO.LOW, GPIO.HIGH, GPIO.LOW),
          "yellow": (GPIO.HIGH, GPIO.HIGH, GPIO.LOW),
          "orange": (GPIO.HIGH, GPIO.HIGH, GPIO.LOW),
          "gray":   (GPIO.HIGH, GPIO.HIGH, GPIO.HIGH)}


def gpioSetup():
    GPIO.setwarnings(False)
    GPIO.setmode(GPIO.BCM)
    GPIO.setup(pins[0], GPIO.OUT)
    GPIO.setup(pins[1], GPIO.OUT)
    GPIO.setup(pins[2], GPIO.OUT)


def getIP():
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    s.settimeout(0)
    try:
        s.connect(('1.1.1.1', 1))
        ip = s.getsockname()[0]
    except Exception:
        ip = '127.0.0.1'
    return ip


class RGBLEDServer(BaseHTTPRequestHandler):
    def do_GET(self):
        leaf = self.path.strip("/").lower()
        if leaf in colors.keys():
            self.setRGB(leaf)
            self.sendMyResponse(200, "LED Changer",
                                "Request: {} found!".format(self.path))
        else:
            self.sendMyResponse(404, "LED Changer - Not found",
                                "Request: {} not found!".format(self.path))

    def sendMyResponse(self, resp, title, body):
        self.send_response(resp)
        self.send_header("Content-type", "text/html")
        self.end_headers()
        self.wfile.write(bytes("<html>", "utf-8"))
        self.wfile.write(
            bytes("<head><title>{}</title></head>".format(title), "utf-8"))
        self.wfile.write(
            bytes("<body><p>{}</p></body>".format(body), "utf-8"))
        self.wfile.write(bytes("</html>", "utf-8"))

    def setRGB(self, color):
        print("RGB: {} value: {}".format(color, colors[color]))
        GPIO.output(pins, colors[color])


def main():
    gpioSetup()
    ip = getIP()
    webServer = HTTPServer((ip, serverPort), RGBLEDServer)
    print("Server started http://{}:{}".format(ip, serverPort))

    try:
        webServer.serve_forever()
    except KeyboardInterrupt:
        webServer.server_close()
        GPIO.cleanup(pins)
        print("Server stopped.")


if __name__ == "__main__":
    main()
