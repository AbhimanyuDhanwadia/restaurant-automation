import { Bell, Database, ShieldCheck } from "lucide-react";
import { useState } from "react";
import { Panel, PanelHeading } from "@/components/ui/Panel";

export function SettingsPage() {
  const [restaurantName, setRestaurantName] = useState("The Green Table");
  const [timezone, setTimezone] = useState("Asia/Kolkata");
  const [notifications, setNotifications] = useState(true);
  const [autoAdvance, setAutoAdvance] = useState(false);
  const [saved, setSaved] = useState(false);

  const saveSettings = () => setSaved(true);

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
            <label>Restaurant name<input value={restaurantName} onChange={(event) => { setRestaurantName(event.target.value); setSaved(false); }} /></label>
            <label>Timezone<select value={timezone} onChange={(event) => { setTimezone(event.target.value); setSaved(false); }}><option>Asia/Kolkata</option><option>UTC</option><option>America/New_York</option></select></label>
          </div>
        </Panel>

        <Panel label="Operations preferences">
          <PanelHeading eyebrow="Preferences" title="Operations preferences" action={<Bell size={20} aria-hidden="true" />} />
          <div className="settings-toggle-list">
            <label className="settings-toggle"><span><strong>Operational alerts</strong><small>Show attention items in the live dashboard.</small></span><input type="checkbox" checked={notifications} onChange={(event) => { setNotifications(event.target.checked); setSaved(false); }} /></label>
            <label className="settings-toggle"><span><strong>Auto-advance tickets</strong><small>Use mock automation to move kitchen tickets.</small></span><input type="checkbox" checked={autoAdvance} onChange={(event) => { setAutoAdvance(event.target.checked); setSaved(false); }} /></label>
          </div>
        </Panel>
      </div>

      <Panel label="System configuration">
        <PanelHeading eyebrow="Foundation" title="System configuration" action={<Database size={20} aria-hidden="true" />} />
        <div className="settings-info-list">
          <div><span>Data mode</span><strong>Mock data</strong></div>
          <div><span>Authentication</span><strong>Supabase Auth</strong></div>
          <div><span>Automation engine</span><strong>Not connected</strong></div>
        </div>
        <div className="settings-actions"><button type="button" className="text-button" onClick={saveSettings}>{saved ? "Saved" : "Save settings"}</button></div>
      </Panel>
    </section>
  );
}
