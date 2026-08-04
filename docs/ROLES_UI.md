# Roles UI

The Roles workspace reads and creates durable role templates through authenticated Administration endpoints. It shows seeded system roles and supports custom roles using allow-listed permissions.

The workspace does not list, create, or administer Supabase Auth users. The Users workspace assigns these templates only to local records observed from verified JWTs, so no privileged Supabase credential is exposed to the browser.
