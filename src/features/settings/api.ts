import { apiRequest } from "@/lib/api";

export interface RestaurantSettings {
  restaurant_name: string;
  timezone: string;
  operational_alerts: boolean;
  auto_advance_tickets: boolean;
  updated_at: string;
}

export type SaveSettingsInput = Omit<RestaurantSettings, "updated_at">;

export const getSettings = () => apiRequest<RestaurantSettings>("/api/v1/settings");

export const saveSettings = (input: SaveSettingsInput) => apiRequest<RestaurantSettings>("/api/v1/settings", {
  method: "PUT",
  body: JSON.stringify(input),
});
