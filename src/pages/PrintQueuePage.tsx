import { CircleAlert, Clock3, Printer, RefreshCw, RotateCcw } from "lucide-react";
import { useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { listPrintJobs, reprintOrder, type PrintJobStatus } from "@/features/printqueue/api";

const statuses: Array<PrintJobStatus | "all"> = ["all", "queued", "printing", "printed", "failed"];

export function PrintQueuePage() {
  const [filter, setFilter] = useState<PrintJobStatus | "all">("all");
  const queryClient = useQueryClient();
  const jobsQuery = useQuery({ queryKey: ["printers", "queue"], queryFn: listPrintJobs, refetchInterval: 2_000 });
  const reprintMutation = useMutation({ mutationFn: reprintOrder, onSuccess: () => queryClient.invalidateQueries({ queryKey: ["printers", "queue"] }) });
  const jobs = useMemo(() => jobsQuery.data ?? [], [jobsQuery.data]);
  const visibleJobs = useMemo(() => filter === "all" ? jobs : jobs.filter((job) => job.status === filter), [filter, jobs]);
  const queued = jobs.filter((job) => job.status === "queued" || job.status === "printing").length;
  const failed = jobs.filter((job) => job.status === "failed").length;

  return <section className="print-queue-workspace">
    <div className="automation-toolbar"><div><p className="eyebrow">Automation</p><h2>Print Queue</h2><p className="automation-subtitle">Durable kitchen-ticket processing, retries, and reprints.</p></div><button type="button" className="icon-button" onClick={() => jobsQuery.refetch()} disabled={jobsQuery.isFetching} title="Refresh print queue" aria-label="Refresh print queue"><RefreshCw size={18} aria-hidden="true" className={jobsQuery.isFetching ? "spin" : undefined} /></button></div>
    {(jobsQuery.isError || reprintMutation.isError) && <div className="orders-api-error" role="alert"><span>{jobsQuery.error?.message ?? reprintMutation.error?.message}</span><button type="button" className="text-button" onClick={() => jobsQuery.refetch()}>Retry</button></div>}
    <section className="automation-metric-grid print-queue-metrics" aria-label="Print queue metrics"><article className="automation-metric"><Clock3 size={20} aria-hidden="true" /><span>Awaiting print</span><strong>{jobsQuery.data ? queued : "Unavailable"}</strong><small>Queued or currently printing</small></article><article className="automation-metric"><Printer size={20} aria-hidden="true" /><span>Total jobs</span><strong>{jobsQuery.data ? jobs.length : "Unavailable"}</strong><small>Durable queue history</small></article><article className="automation-metric"><CircleAlert size={20} aria-hidden="true" /><span>Failed jobs</span><strong>{jobsQuery.data ? failed : "Unavailable"}</strong><small>Can be reprinted</small></article><article className="automation-metric"><RotateCcw size={20} aria-hidden="true" /><span>Reprints</span><strong>{jobsQuery.data ? jobs.filter((job) => job.reprint).length : "Unavailable"}</strong><small>Created from prior tickets</small></article></section>
    <Panel label="Print job history"><PanelHeading eyebrow="Ticket queue" title="Print jobs" action={<label className="filter-control"><span>Status</span><select value={filter} onChange={(event) => setFilter(event.target.value as PrintJobStatus | "all")}>{statuses.map((status) => <option key={status} value={status}>{status === "all" ? "All jobs" : status}</option>)}</select></label>} /><div className="print-job-list">{jobsQuery.isPending && <EmptyState message="Loading print jobs..." />}{visibleJobs.map((job) => <article className="print-job-row" key={job.id}><div><strong>{job.order_id}</strong><span>{job.destination} · {job.lines.map((line) => `${line.quantity}× ${line.text}`).join(", ")}</span>{job.last_error && <small>{job.last_error}</small>}</div><time>{new Date(job.created_at).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit", second: "2-digit" })}</time><span className={`automation-status status-${job.status}`}>{job.status}</span><span className="print-attempts">{job.attempts} attempt{job.attempts === 1 ? "" : "s"}</span><button type="button" className="back-button" onClick={() => reprintMutation.mutate(job.order_id)} disabled={reprintMutation.isPending}><RotateCcw size={15} aria-hidden="true" /> Reprint</button></article>)}{!jobsQuery.isPending && !jobsQuery.isError && visibleJobs.length === 0 && <EmptyState message={jobs.length ? "No print jobs match this status." : "No print jobs have been submitted."} />}</div></Panel>
  </section>;
}
