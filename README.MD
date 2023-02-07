Teams Status Scrapper (TSS)

Built on Go version 1.20 and compatible with python 3.5.3

Compile:    
    `go build -ldflags -H=windowsgui .\TSS.go`

Setup config file:
    Set "host" and "port" the the values where your Raspberry Pi is (I used IP address for "host")

Run the exe:
    `.\TSS.exe`

To stop the process, you'll need to kill it from the task manager as there's no GUI. Find "TSS.exe" in background tasks.

For debugging, you can run:
    `go run .\TSS.go`
    This will not hide the command prompt and you can see the programs output.

Put the `pipyserver.py` on a Raspberry Pi with an RGB LED connected to pins 12 (R), 13 (G), and 19 (B) with appropriate resistors (I used 220Ohm). 

Note: This is obviously running on the HTTP stack, so if you're on a VPN or similar, you won't be able to access devices on your local network unless you open your firewall. (Not recommended!)