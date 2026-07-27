# Kitchen Queue API Connection

The Kitchen workspace reads the durable Orders API and only displays orders that require kitchen work.

| Order status | Kitchen state | Available kitchen action |
| --- | --- | --- |
| `received` | Queued | Start preparing |
| `preparing` | Cooking | Mark ready |
| `ready` | Ready for handoff | None |

Delivered and cancelled orders are excluded from the kitchen queue. The screen uses the same `VITE_API_URL` and Supabase access-token behavior described in [ORDERS_UI.md](ORDERS_UI.md).

Station routing is intentionally not presented yet: the durable order model has no station assignment. That will be added together with recipe and kitchen-routing data rather than inferred from order channels.
