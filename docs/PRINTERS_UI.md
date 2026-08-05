# Printers UI

The Printers workspace reads the authenticated `GET /api/v1/printers` fleet-health endpoint every five seconds. It presents the printers registered at API startup, their connection state, queued tickets, successful prints, and failures since the current API process started.

Network addresses, printer-driver configuration, and routing remain server-side configuration. They are deliberately not returned by the API or exposed in the browser. The API registers separate `kitchen` and `cashier` destinations at startup; configure `PRINTER_KITCHEN_ADDRESS` and `PRINTER_CASHIER_ADDRESS` for raw-TCP ESC/POS hardware, or leave either empty for its local mock driver. Ticket submission and reprinting continue through the existing authenticated printer endpoints.
