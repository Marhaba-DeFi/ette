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
      name: "event/0xfB114de60340ec5053fEe1E2782d090EFBEe16B1/*/*/*/*",
      type: "subscribe",
      apiKey:
        "0xb7a81a56fb912fbe618e6af7bbba70cac4ca342a61fc09f6e95e4d0db0e852ad",
    })
  )
})

_client.connect("ws://localhost:7000/v1/ws", null)
