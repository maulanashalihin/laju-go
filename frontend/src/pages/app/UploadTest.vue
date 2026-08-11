<script setup lang="ts">
import { ref, computed } from "vue";
import { Link } from "@inertiajs/vue3";
import AppLayout from "@layouts/AppLayout.vue";
import { Toast } from "@lib/notifications/toast";
import type { User } from "@lib/types";
import * as tus from "tus-js-client";
import {
  Upload,
  FileIcon,
  Trash2,
  Download,
  CheckCircle2,
  XCircle,
  LoaderCircle,
  AlertCircle,
  ExternalLink,
  Copy,
  Clock,
} from "lucide-vue-next";

interface Props {
  user?: User;
}

defineProps<Props>();

interface UploadEntry {
  id: string;
  name: string;
  size: number;
  progress: number;
  status: "pending" | "uploading" | "paused" | "done" | "error";
  error?: string;
  url?: string;
  tusRef?: tus.Upload;
  startedAt: Date;
}

// ── State ─────────────────────────────────────────────────
const uploads = ref<UploadEntry[]>([]);
const isDragOver = ref(false);
const fileInput = ref<HTMLInputElement>();

// ── Helper: update entry through reactive proxy ──────────
function setEntry(id: string, patch: Partial<UploadEntry>) {
  const idx = uploads.value.findIndex((e) => e.id === id);
  if (idx !== -1) {
    Object.assign(uploads.value[idx], patch);
  }
}

// ── Derived ───────────────────────────────────────────────
const totalUploaded = computed(
  () => uploads.value.filter((u) => u.status === "done").length
);
const totalBytes = computed(
  () =>
    uploads.value
      .filter((u) => u.status === "done")
      .reduce((sum, u) => sum + u.size, 0)
);
const activeUploads = computed(
  () => uploads.value.filter((u) => u.status === "uploading").length
);

function formatBytes(bytes: number): string {
  if (bytes === 0) return "0 B";
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + " " + sizes[i];
}

// ── Drag & Drop ───────────────────────────────────────────
function handleDragOver(e: DragEvent) {
  e.preventDefault();
  isDragOver.value = true;
}

function handleDragLeave() {
  isDragOver.value = false;
}

function handleDrop(e: DragEvent) {
  e.preventDefault();
  isDragOver.value = false;
  const files = Array.from(e.dataTransfer?.files || []);
  if (files.length > 0) {
    addFiles(files);
  }
}

function handleFileSelect(e: Event) {
  const target = e.target as HTMLInputElement;
  const files = Array.from(target.files || []);
  if (files.length > 0) {
    addFiles(files);
  }
  target.value = "";
}

function addFiles(filesList: globalThis.File[]) {
  for (const file of filesList) {
    const entry: UploadEntry = {
      id: crypto.randomUUID(),
      name: file.name,
      size: file.size,
      progress: 0,
      status: "pending",
      startedAt: new Date(),
    };
    uploads.value = [...uploads.value, entry];
    startUpload(entry.id, file);
  }
}

// ── TUS Upload ────────────────────────────────────────────
const TUS_ENDPOINT = "/tus/files";

function startUpload(entryId: string, file: globalThis.File) {
  // Set initial status through reactive proxy
  setEntry(entryId, { status: "uploading" });

  const upload = new tus.Upload(file, {
    endpoint: TUS_ENDPOINT,
    retryDelays: [0, 1000, 3000, 5000],
    chunkSize: 5 * 1024 * 1024, // 5MB chunks for responsive progress
    metadata: {
      filename: file.name,
      filetype: file.type,
    },
    onError: (err) => {
      setEntry(entryId, {
        status: "error",
        error: err.message,
      });
      Toast(`Upload failed: ${err.message}`, "error");
    },
    onProgress: (bytesSent, bytesTotal) => {
      const progress = bytesTotal > 0 ? (bytesSent / bytesTotal) * 100 : 0;
      setEntry(entryId, { progress });
    },
    onSuccess: () => {
      setEntry(entryId, {
        status: "done",
        progress: 100,
        url: upload.url ?? undefined,
      });
      Toast(`Upload selesai: ${file.name}`, "success");
    },
    onShouldRetry: (_err, retryAttempt, _options) => {
      Toast(`Retry ${retryAttempt + 1} for ${file.name}...`, "info");
      return true;
    },
  });

  // Store tusRef through reactive proxy
  const idx = uploads.value.findIndex((e) => e.id === entryId);
  if (idx !== -1) {
    uploads.value[idx].tusRef = upload;
  }

  upload.start();
}

function pauseUpload(id: string) {
  const idx = uploads.value.findIndex((e) => e.id === id);
  if (idx !== -1) {
    uploads.value[idx].tusRef?.abort();
    uploads.value[idx].status = "paused";
  }
}

function removeUpload(index: number) {
  const entry = uploads.value[index];
  if (entry?.status === "uploading") {
    entry.tusRef?.abort();
  }
  uploads.value = uploads.value.filter((_, i) => i !== index);
}

function clearCompleted() {
  uploads.value = uploads.value.filter((u) => u.status !== "done");
}

// ── File actions ──────────────────────────────────────────
function getUploadUrl(uploadUrl: string | undefined): string {
  if (!uploadUrl) return "";
  try {
    const url = new URL(uploadUrl);
    return url.pathname;
  } catch {
    return uploadUrl;
  }
}

function copyDownloadLink(index: number) {
  const url = getUploadUrl(uploads.value[index]?.url);
  if (!url) return;
  navigator.clipboard.writeText(window.location.origin + url);
  Toast("Link copied!", "success");
}

function deleteUpload(index: number) {
  const entry = uploads.value[index];
  if (!entry?.url) return;
  const url = getUploadUrl(entry.url);

  fetch(url, {
    method: "DELETE",
    headers: {
      "Tus-Resumable": "1.0.0",
    },
  })
    .then((res) => {
      if (res.ok || res.status === 204 || res.status === 404) {
        uploads.value = uploads.value.filter((_, i) => i !== index);
        Toast(`Deleted: ${entry.name}`, "success");
      } else {
        Toast("Failed to delete", "error");
      }
    })
    .catch(() => {
      Toast("Failed to delete", "error");
    });
}

// ── Protocol endpoint data ───────────────────────────────
const tusEndpoints = [
  { method: "POST", path: "/tus/files", desc: "Create upload" },
  { method: "HEAD", path: "/tus/files/:id", desc: "Get upload offset/info" },
  { method: "PATCH", path: "/tus/files/:id", desc: "Upload chunk" },
  { method: "GET", path: "/tus/files/:id", desc: "Download file" },
  { method: "DELETE", path: "/tus/files/:id", desc: "Terminate upload" },
  { method: "OPTIONS", path: "/tus/files", desc: "Protocol discovery" },
];

// ── Status helpers ────────────────────────────────────────
function statusColor(status: string): string {
  switch (status) {
    case "done": return "text-green-500";
    case "error": return "text-red-500";
    case "uploading": return "text-brand-500";
    case "paused": return "text-amber-500";
    default: return "text-neutral-500";
  }
}

function statusIcon(status: string) {
  switch (status) {
    case "done": return CheckCircle2;
    case "error": return XCircle;
    case "uploading": return LoaderCircle;
    case "paused": return Clock;
    default: return AlertCircle;
  }
}

function statusLabel(status: string): string {
  switch (status) {
    case "done": return "Completed";
    case "error": return "Failed";
    case "uploading": return "Uploading...";
    case "paused": return "Paused";
    case "pending": return "Pending";
    default: return status;
  }
}
</script>

<template>
  <AppLayout :user="user" group="upload">
    <!-- Page Header -->
    <div class="pt-8 pb-10 border-b border-neutral-200/80 dark:border-white/[0.04]">
      <div class="max-w-6xl mx-auto px-6">
        <div class="flex items-center gap-2 text-sm text-neutral-500 dark:text-neutral-400 mb-4">
          <Link href="/app" class="hover:text-brand-600 dark:hover:text-brand-400 transition-colors">Dashboard</Link>
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
          <span class="text-neutral-700 dark:text-neutral-300">Upload Test</span>
        </div>
        <div class="flex items-start justify-between gap-4 flex-wrap">
          <div>
            <h1 class="text-3xl font-bold text-neutral-900 dark:text-white mb-2 tracking-tight">File Upload Test</h1>
            <p class="text-neutral-600 dark:text-neutral-400 max-w-xl">
              Test the TUS resumable upload protocol. Big files, chunked uploads, progress tracking —
              all powered by <a href="https://github.com/maulanashalihin/tusdfiber" target="_blank" rel="noopener" class="text-brand-600 hover:text-brand-700 dark:text-brand-400 dark:hover:text-brand-300 font-medium underline underline-offset-2">tusdfiber</a>.
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- Content Area -->
    <div class="relative max-w-6xl mx-auto px-6 py-8 space-y-6">
      <!-- Stats bar -->
      <div class="grid grid-cols-3 gap-4">
        <div class="rounded-2xl border border-neutral-200/80 dark:border-white/[0.06] bg-white dark:bg-neutral-925/50 p-5">
          <p class="text-sm text-neutral-500 dark:text-neutral-400 mb-1">Uploaded Files</p>
          <p class="text-3xl font-bold text-neutral-900 dark:text-white font-mono">{{ totalUploaded }}</p>
        </div>
        <div class="rounded-2xl border border-neutral-200/80 dark:border-white/[0.06] bg-white dark:bg-neutral-925/50 p-5">
          <p class="text-sm text-neutral-500 dark:text-neutral-400 mb-1">Total Size</p>
          <p class="text-3xl font-bold text-neutral-900 dark:text-white font-mono">{{ formatBytes(totalBytes) }}</p>
        </div>
        <div class="rounded-2xl border border-neutral-200/80 dark:border-white/[0.06] bg-white dark:bg-neutral-925/50 p-5">
          <p class="text-sm text-neutral-500 dark:text-neutral-400 mb-1">Active Uploads</p>
          <p class="text-3xl font-bold text-neutral-900 dark:text-white font-mono">{{ activeUploads }}</p>
        </div>
      </div>

      <!-- Upload Drop Zone -->
      <div
        :class="[
          'relative rounded-2xl border-2 border-dashed transition-all duration-300 overflow-hidden',
          isDragOver
            ? 'border-brand-400 bg-brand-400/5 scale-[1.01] shadow-xl shadow-brand-400/10'
            : 'border-neutral-300 dark:border-neutral-700 hover:border-brand-400/50 bg-white dark:bg-neutral-925/50',
        ]"
        role="button"
        tabindex="0"
        aria-label="Upload file drop zone"
        @dragover="handleDragOver"
        @dragleave="handleDragLeave"
        @drop="handleDrop"
      >
        <input type="file" multiple ref="fileInput" @change="handleFileSelect" class="hidden" />
        <button @click="fileInput?.click()" class="w-full py-16 px-8 flex flex-col items-center justify-center gap-4 cursor-pointer group">
          <div class="w-16 h-16 rounded-2xl bg-brand-400/10 flex items-center justify-center group-hover:scale-110 transition-transform">
            <Upload class="w-8 h-8 text-brand-600 dark:text-brand-400" />
          </div>
          <div class="text-center">
            <p class="text-lg font-semibold text-neutral-900 dark:text-white mb-1">{{ isDragOver ? "Drop files here" : "Drop files or click to upload" }}</p>
            <p class="text-sm text-neutral-500 dark:text-neutral-400">Any file type — up to 1GB per upload via TUS resumable protocol</p>
          </div>
        </button>
      </div>

      <!-- Upload Queue -->
      <div
        v-if="uploads.length > 0"
        class="rounded-2xl border border-neutral-200/80 dark:border-white/[0.06] bg-white dark:bg-neutral-925/50 overflow-hidden"
      >
        <div class="flex items-center justify-between px-6 py-4 border-b border-neutral-200/80 dark:border-white/[0.04]">
          <h3 class="text-base font-semibold text-neutral-900 dark:text-white">Upload Queue <span class="text-sm font-normal text-neutral-500 ml-2">({{ uploads.length }} files)</span></h3>
          <button v-if="uploads.some((u) => u.status === 'done')" @click="clearCompleted" class="text-xs font-medium text-neutral-500 hover:text-red-500 transition-colors px-3 py-1.5 rounded-lg hover:bg-red-500/10">Clear completed</button>
        </div>
        <ul class="divide-y divide-neutral-200/80 dark:divide-white/[0.04]">
          <li
            v-for="(entry, i) in uploads"
            :key="entry.id"
            class="px-6 py-4 transition-colors hover:bg-neutral-50/50 dark:hover:bg-white/[0.015]"
          >
            <div class="flex items-start gap-4">
              <div class="shrink-0 w-10 h-10 rounded-xl bg-neutral-100 dark:bg-neutral-800 flex items-center justify-center">
                <FileIcon class="w-5 h-5 text-neutral-500" />
              </div>
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2 mb-1">
                  <p class="text-sm font-medium text-neutral-900 dark:text-white truncate">{{ entry.name }}</p>
                  <span class="text-xs text-neutral-500 shrink-0">({{ formatBytes(entry.size) }})</span>
                </div>
                <div class="relative h-2 rounded-full bg-neutral-200/80 dark:bg-neutral-800 overflow-hidden mb-1.5">
                  <div :class="['h-full rounded-full transition-all duration-500 ease-out', entry.status === 'error' ? 'bg-red-500' : entry.status === 'done' ? 'bg-green-500' : 'bg-brand-500']" :style="{ width: entry.progress + '%' }"></div>
                </div>
                <div class="flex items-center gap-2 text-xs">
                  <component :is="statusIcon(entry.status)" :class="['w-3.5 h-3.5', statusColor(entry.status)]" />
                  <span :class="[statusColor(entry.status), 'font-medium']">{{ statusLabel(entry.status) }}</span>
                  <span v-if="entry.status === 'uploading'" class="text-neutral-400">· {{ Math.round(entry.progress) }}%</span>
                  <span v-if="entry.error" class="text-red-400">· {{ entry.error }}</span>
                </div>
              </div>
              <div class="flex items-center gap-1 shrink-0">
                <template v-if="entry.status === 'done' && entry.url">
                  <a :href="getUploadUrl(entry.url)" target="_blank" rel="noopener" class="p-2 rounded-lg hover:bg-neutral-100 dark:hover:bg-neutral-800 text-neutral-500 hover:text-neutral-900 dark:hover:text-white transition-colors" title="Download"><Download class="w-4 h-4" /></a>
                  <button @click="copyDownloadLink(i)" class="p-2 rounded-lg hover:bg-neutral-100 dark:hover:bg-neutral-800 text-neutral-500 hover:text-neutral-900 dark:hover:text-white transition-colors" title="Copy link"><Copy class="w-4 h-4" /></button>
                  <button @click="deleteUpload(i)" class="p-2 rounded-lg hover:bg-red-500/10 text-neutral-500 hover:text-red-500 transition-colors" title="Delete"><Trash2 class="w-4 h-4" /></button>
                </template>
                <template v-else-if="entry.status === 'uploading'">
                  <button @click="pauseUpload(entry.id)" class="p-2 rounded-lg hover:bg-amber-500/10 text-neutral-500 hover:text-amber-500 transition-colors" title="Pause"><Clock class="w-4 h-4" /></button>
                  <button @click="removeUpload(i)" class="p-2 rounded-lg hover:bg-red-500/10 text-neutral-500 hover:text-red-500 transition-colors" title="Cancel"><XCircle class="w-4 h-4" /></button>
                </template>
                <template v-else-if="entry.status === 'error'">
                  <button @click="removeUpload(i)" class="p-2 rounded-lg hover:bg-red-500/10 text-neutral-500 hover:text-red-500 transition-colors" title="Remove"><Trash2 class="w-4 h-4" /></button>
                </template>
              </div>
            </div>
          </li>
        </ul>
      </div>

      <!-- Empty state -->
      <div
        v-if="uploads.length === 0"
        class="rounded-2xl border border-neutral-200/80 dark:border-white/[0.06] bg-white dark:bg-neutral-925/50 p-12 text-center"
      >
        <div class="w-16 h-16 rounded-2xl bg-neutral-100 dark:bg-neutral-800 flex items-center justify-center mx-auto mb-4"><Upload class="w-8 h-8 text-neutral-500" /></div>
        <h3 class="text-lg font-semibold text-neutral-900 dark:text-white mb-2">No uploads yet</h3>
        <p class="text-neutral-500 dark:text-neutral-400 max-w-md mx-auto text-sm">Drag & drop files above or click the drop zone to select files. Uploads use the TUS resumable protocol — you can pause and resume big uploads.</p>
      </div>

      <!-- Protocol Info Card -->
      <div class="rounded-2xl border border-neutral-200/80 dark:border-white/[0.06] bg-white dark:bg-neutral-925/50 p-6">
        <div class="flex items-center gap-3 mb-4">
          <div class="w-10 h-10 rounded-xl bg-brand-400/10 flex items-center justify-center"><ExternalLink class="w-5 h-5 text-brand-600 dark:text-brand-400" /></div>
          <div><h3 class="text-base font-semibold text-neutral-900 dark:text-white">TUS Protocol Endpoints</h3><p class="text-sm text-neutral-500 dark:text-neutral-400">All routes are protected (auth required)</p></div>
        </div>
        <div class="grid sm:grid-cols-2 gap-3 text-sm">
          <div
            v-for="(item, i) in tusEndpoints"
            :key="i"
            class="flex items-center gap-3 p-2.5 rounded-lg bg-neutral-50 dark:bg-neutral-900/50 border border-neutral-200/80 dark:border-white/[0.04]"
          >
            <span :class="['inline-flex items-center px-2 py-0.5 rounded text-xs font-mono font-bold', item.method === 'GET' ? 'bg-green-500/10 text-green-600 dark:text-green-400' : item.method === 'POST' ? 'bg-blue-500/10 text-blue-600 dark:text-blue-400' : item.method === 'PATCH' ? 'bg-amber-500/10 text-amber-600 dark:text-amber-400' : item.method === 'DELETE' ? 'bg-red-500/10 text-red-600 dark:text-red-400' : 'bg-neutral-200/80 dark:bg-neutral-700 text-neutral-600 dark:text-neutral-400']">{{ item.method }}</span>
            <code class="text-xs text-neutral-700 dark:text-neutral-300 font-mono shrink min-w-0 truncate">{{ item.path }}</code>
            <span class="text-neutral-500 dark:text-neutral-400 ml-auto shrink-0">{{ item.desc }}</span>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>
