import { Bell, Database, RefreshCw, ShieldCheck } from "lucide-react";
import { useEffect, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { getSettings, saveSettings } from "@/features/settings/api";

export function SettingsPage() {
  const [restaurantName, setRestaurantName] = useState("The Green Table");
  const [timezone, setTimezone] = useState("Asia/Kolkata");
  const [notifications, setNotifications] = useState(true);
  const [autoAdvance, setAutoAdvance] = useState(false);
  const [saved, setSaved] = useState(false);
  const queryClient = useQueryClient();
  const settingsQuery = useQuery({ queryKey: ["settings"], queryFn: getSettings });
  const settingsMutation = useMutation({
    mutationFn: saveSettings,
    onSuccess: () => {
      setSaved(true);
      queryClient.invalidateQueries({ queryKey: ["settings"] });
    },
  });

  useEffect(() => {
    if (!settingsQuery.data) return;
    setRestaurantName(settingsQuery.data.restaurant_name);
    setTimezone(settingsQuery.data.timezone);
    setNotifications(settingsQuery.data.operational_alerts);
    setAutoAdvance(settingsQuery.data.auto_advance_tickets);
  }, [settingsQuery.data]);

  const markChanged = () => setSaved(false);
  const submitSettings = () => settingsMutation.mutate({
    restaurant_name: restaurantName.trim(),
    timezone,
    operational_alerts: notifications,
    auto_advance_tickets: autoAdvance,
  });

  return (
    <section className="settings-workspace">
      <div>
        <p className="eyebrow">Administration</p>
        <h2>Settings</h2>
      </div>

      <div className="settings-grid">
        <Panel label="Restaurant profile">
          <PanelHeading eyebrow="General" title="Restaurant profile" action={<ShieldCheck size={20} aria-hidden="true" />} />
          <div className="settings-form">
            <label>Restaurant name<input value={restaurantName} onChange={(event) => { setRestaurantName(event.target.value); markChanged(); }} disabled={settingsQuery.isPending} /></label>
            <label>Timezone<select value={timezone} onChange={(event) => { setTimezone(event.target.value); markChanged(); }} disabled={settingsQuery.isPending}><option>Asia/Kolkata</option><option>UTC</option><option>America/New_York</option></select></label>
          </div>
        </Panel>

        <Panel label="Operations preferences">
          <PanelHeading eyebrow="Preferences" title="Operations preferences" action={<Bell size={20} aria-hidden="true" />} />
          <div className="settings-toggle-list">
            <label className="settings-toggle"><span><strong>Operational alerts</strong><small>Show attention items in the live dashboard.</small></span><input type="checkbox" checked={notifications} disabled={settingsQuery.isPending} onChange={(event) => { setNotifications(event.target.checked); markChanged(); }} /></label>
            <label className="settings-toggle"><span><strong>Auto-advance tickets</strong><small>Store the preference for the automation scheduler.</small></span><input type="checkbox" checked={autoAdvance} disabled={settingsQuery.isPending} onChange={(event) => { setAutoAdvance(event.target.checked); markChanged(); }} /></label>
          </div>
        </Panel>
      </div>

      <Panel label="System configuration">
        <PanelHeading eyebrow="Foundation" title="System configuration" action={<Database size={20} aria-hidden="true" />} />
        <div className="settings-info-list">
          <div><span>Data mode</span><strong>Persisted settings</strong></div>
          <div><span>Authentication</span><strong>Supabase Auth</strong></div>
          <div><span>Automation engine</span><strong>Configured separately</strong></div>
        </div>
        {settingsQuery.isError && <div className="orders-api-error" role="alert"><span>{settingsQuery.error.message}</span><button type="button" className="text-button" onClick={() => settingsQuery.refetch()}><RefreshCw size={15} aria-hidden="true" /> Retry</button></div>}
        {settingsMutation.isError && <p className="orders-api-error" role="alert">{settingsMutation.error.message}</p>}
        <div className="settings-actions"><button type="button" className="text-button" onClick={submitSettings} disabled={settingsQuery.isPending || settingsMutation.isPending}>{settingsMutation.isPending ? "Saving..." : saved ? "Saved" : "Save settings"}</button></div>
      </Panel>
    </section>
  );
}
