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

// Every failed request throws an ApiError: status is the HTTP status (0 when
// the API couldn't be reached) and message what to tell the user: the API's
// own message (in Portuguese, e.g. "O álbum já tem 12 fotos."), or a generic
// one for the status. App.tsx redirects on 401 and 404 and shows the message
// of failed actions.
export class ApiError extends Error {
  constructor(
    public readonly status: number,
    message: string,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

export type ErrorType<_Error> = ApiError;

// When the response has no message of its own (e.g. a proxy's error page).
const statusMessages: Record<number, string> = {
  401: "Sua sessão expirou. Entre de novo.",
  403: "Você não tem permissão para fazer isso.",
  404: "Não encontramos o que você procurou.",
  413: "O arquivo é grande demais.",
  429: "Muitas ações em pouco tempo. Espere um pouco e tente de novo.",
};
const fallbackMessage = "Algo deu errado. Tente de novo mais tarde.";

// The API answers errors as {"error": "..."}.
const messageOf = async (response: Response): Promise<string | undefined> => {
  try {
    const body = await response.json();
    return typeof body?.error === "string" ? body.error : undefined;
  } catch {
    return undefined;
  }
};
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

  let response: Response;
  try {
    response = await fetch(buildUrl(url, params), {
      method,
      headers: {
        ...getAuthorizationHeaders(),
        ...(isFormData ? {} : { "Content-Type": "application/json" }),
        ...extraHeaders,
      },
      body: isFormData ? data : data ? JSON.stringify(data) : undefined,
    });
  } catch {
    throw new ApiError(0, "Não foi possível falar com o Yorgut. Confira sua conexão e tente de novo.");
  }

  if (!response.ok) {
    if (response.status === 401) {
      // The session expired or is invalid: drop it so the next screen
      // goes to login (see the QueryCache error handler in App.tsx).
      localStorage.removeItem("token");
    }
    throw new ApiError(
      response.status,
      (await messageOf(response)) ?? statusMessages[response.status] ?? fallbackMessage,
    );
  }

  // No body: null, not undefined, which TanStack Query refuses as query data
  // (e.g. /me/friends/:id on your own profile).
  if (response.status === 204) {
    return null as T;
  }

  return (await response.json()) as T;
}

export default apiRequest;
