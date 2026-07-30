import { apiRequest } from "@/lib/api";

export type StaffStatus = "on_shift" | "on_break" | "off_shift";
export interface StaffMember { id: string; name: string; role: string; station: string; status: StaffStatus; handoff: string; created_at: string; updated_at: string }
export interface ShiftTask { id: string; title: string; owner: string; due_label: string; completed: boolean; created_at: string; updated_at: string }
export interface ShiftHandoff { note: string; updated_at: string }
export interface CreateStaffInput { name: string; role: string; station: string; status?: StaffStatus; handoff?: string }
export interface CreateTaskInput { title: string; owner: string; due_label?: string }

export const listStaff = () => apiRequest<StaffMember[]>("/api/v1/staff");
export const createStaffMember = (input: CreateStaffInput) => apiRequest<StaffMember>("/api/v1/staff", { method: "POST", body: JSON.stringify(input) });
export const updateStaffStatus = (id: string, status: StaffStatus) => apiRequest<StaffMember>(`/api/v1/staff/${encodeURIComponent(id)}/status`, { method: "PATCH", body: JSON.stringify({ status }) });
export const updateStaffHandoff = (id: string, handoff: string) => apiRequest<StaffMember>(`/api/v1/staff/${encodeURIComponent(id)}/handoff`, { method: "PATCH", body: JSON.stringify({ handoff }) });
export const listShiftTasks = () => apiRequest<ShiftTask[]>("/api/v1/staff/tasks");
export const createShiftTask = (input: CreateTaskInput) => apiRequest<ShiftTask>("/api/v1/staff/tasks", { method: "POST", body: JSON.stringify(input) });
export const completeShiftTask = (id: string) => apiRequest<ShiftTask>(`/api/v1/staff/tasks/${encodeURIComponent(id)}/complete`, { method: "PATCH" });
export const getShiftHandoff = () => apiRequest<ShiftHandoff>("/api/v1/staff/handoff");
export const saveShiftHandoff = (note: string) => apiRequest<ShiftHandoff>("/api/v1/staff/handoff", { method: "PUT", body: JSON.stringify({ note }) });
