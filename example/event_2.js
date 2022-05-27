const { client } = require("websocket")

const _client = new client()

_client.on("connectFailed", (e) => {
  console.error(`[!] Failed to connect : ${e}`)
  process.exit(1)
})

// For listening to any event being emitted
// by smart contract at `0xcb3fA413B23b12E402Cfcd8FA120f983FB70d8E8`
_client.on("connect", (c) => {
  c.on("close", (d) => {
    console.log(`[!] Closed connection : ${d}`)
    process.exit(0)
  })

  c.on("message", (d) => {
    console.log(JSON.parse(d.utf8Data))
  })

  c.send(
    JSON.stringify({
      name: "event/0xB6E3d179E941Ed21627717834C098e3e56006C85/*/*/*/*",
      type: "subscribe",
      apiKey:
        "0x39a3859caa81aa4b881220f09ad870f81675b29dc5b2b4bfb1fcdb2db62a1b2a",
    })
  )
})

_client.connect("ws://localhost:7000/v1/ws", null)
