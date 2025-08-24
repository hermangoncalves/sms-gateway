import { Hono } from "hono"

type Bindings = {
  DB: D1Database
}

const app = new Hono<{ Bindings: Bindings }>()

app.post("/messages", async (c) => {
  const { number, text } = await c.req.json<{ number: string; text: string }>()
  await c.env.DB.prepare(
    "INSERT INTO messages (number, text, status) VALUES (?, ?, ?)"
  ).bind(number, text, "pending").run()

  return c.json({ status: "queued", number, text })
})

app.get("/pending-sms", async (c) => {
  const limit = Number(c.req.query("limit") ?? 5)

  const { results } = await c.env.DB.prepare(
    "SELECT id, number, text FROM messages WHERE status = ? ORDER BY id LIMIT ?"
  ).bind("pending", limit).all()

  return c.json(results)
})

app.post("/confirm-sms", async (c) => {
  const { number, status } = await c.req.json<{ number: string; status: string }>()
  await c.env.DB.prepare(
    "UPDATE messages SET status = ? WHERE number = ? AND status = 'pending'"
  ).bind(status, number).run()

  return c.json({ number, deliveryStatus: status })
})

app.get("/sent", async (c) => {
  const { results } = await c.env.DB.prepare(
    "SELECT id, number, text, status, created_at FROM messages WHERE status = ? ORDER BY id DESC"
  ).bind("sent").all()
  return c.json(results)
})

export default app
