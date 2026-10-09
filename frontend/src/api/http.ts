// The single HTTP entry point. The generated client (api/generated, see
// orval.config.ts) calls apiRequest for every endpoint, so the auth header,
// the 401 handling and the error messages below apply everywhere.

// What the generated client passes for each call: `url` is the spec path
// (e.g. `/users/1/profile`), relative to the API base URL.
export interface ApiRequestConfig {
  url: string;
  method: "GET" | "POST" | "PUT" | "DELETE" | "PATCH";
  params?: Record<string, unknown>;
  headers?: Record<string, string>;
  data?: unknown;
  signal?: AbortSignal;
}

// Errors thrown by apiRequest are plain Errors; their message starts with
// "Unauthorized", "Forbidden", "Not Found", ... (App.tsx's QueryCache matches on it).
export type ErrorType<_Error> = Error;
// Request bodies are sent as given.
export type BodyType<BodyData> = BodyData;

const getAuthorizationHeaders = (): Record<string, string> => {
  try {
    const token = localStorage.getItem("token");
    return token ? { Authorization: "Bearer " + token } : {};
  } catch (e) {
    return {};
  }
};

const buildUrl = (url: string, params?: Record<string, unknown>): string => {
  const search = new URLSearchParams();
  for (const [key, value] of Object.entries(params ?? {})) {
    if (value !== undefined && value !== null) search.append(key, String(value));
  }
  const query = search.toString();
  return `${__API_BASE_URL__}${url}${query ? `?${query}` : ""}`;
};

export async function apiRequest<T>(
  { url, method, params, headers = {}, data }: ApiRequestConfig,
  options?: { headers?: Record<string, string> },
): Promise<T> {
  const isFormData = data instanceof FormData;
  const extraHeaders = { ...headers, ...options?.headers };
  if (isFormData) {
    // The browser sets the multipart Content-Type itself (with the boundary).
    delete extraHeaders["Content-Type"];
  }

  try {
    const response = await fetch(buildUrl(url, params), {
      method,
      headers: {
        ...getAuthorizationHeaders(),
        ...(isFormData ? {} : { "Content-Type": "application/json" }),
        ...extraHeaders,
      },
      body: isFormData ? data : data ? JSON.stringify(data) : undefined,
    });

    if (!response.ok) {
      let errorMessage = `Error ${response.status}: ${response.statusText}`;
      switch (response.status) {
        case 401:
          // The session expired or is invalid: drop it so the next screen
          // goes to login (see the QueryCache error handler in App.tsx).
          localStorage.removeItem("token");
          errorMessage = "Unauthorized: Invalid or missing credentials.";
          break;
        case 403:
          errorMessage =
            "Forbidden: You do not have permission to access this resource.";
          break;
        case 404:
          errorMessage =
            "Not Found: The requested resource could not be found.";
          break;
        case 500:
          errorMessage =
            "Server Error: Something went wrong on the server side.";
          break;
        default:
          errorMessage = `HTTP Error ${response.status}: ${response.statusText}`;
      }
      throw new Error(errorMessage);
    }

    if (response.status === 204) {
      return undefined as T;
    }

    return (await response.json()) as T;
  } catch (error: any) {
    throw new Error(`${error.message}`);
  }
}

export default apiRequest;
