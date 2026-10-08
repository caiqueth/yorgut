// Uploaded files come from the API as paths relative to the backend
// ("storage/uploads/<uuid>"). Used as-is, the browser resolves them against the
// current page URL, so they only load on "/". Resolve them against the API
// origin instead. Absolute URLs (e.g. Google profile pictures) pass through.
const apiOrigin = new URL(__API_BASE_URL__, window.location.origin).origin;

export const uploadUrl = (path?: string | null): string | undefined => {
  if (!path) return undefined;
  if (!/^\/?storage\//.test(path)) return path;
  return `${apiOrigin}/${path.replace(/^\//, "")}`;
};
