const getAuthorizationHeaders = (): Record<string, string> => {
  try {
    const token = localStorage.getItem("token");
    return token ? { Authorization: "Bearer " + token } : {};
  } catch (e) {
    return {};
  }
};

export async function doGet<T>({
  queryKey,
}: {
  queryKey: readonly unknown[];
}): Promise<T> {
  const [url, headers = {}] = queryKey as [string, any];
  return doRequest<T>(url, "GET", headers);
}

export async function doPost<T>({
  url,
  body,
  headers = {},
}: {
  url: string;
  body: any;
  headers?: Record<string, string>;
}): Promise<T> {
  return doRequest<T>(url, "POST", headers, body);
}

export async function doPut<T>({
  url,
  body,
  headers = {},
}: {
  url: string;
  body: any;
  headers?: Record<string, string>;
}): Promise<T> {
  return doRequest<T>(url, "PUT", headers, body);
}

export async function doDelete<T>({
  url,
  headers = {},
}: {
  url: string;
  headers?: Record<string, string>;
}): Promise<T> {
  return doRequest<T>(url, "DELETE", headers);
}

async function doRequest<T>(
  url: string,
  method: string,
  headers: Record<string, string>,
  body?: any,
): Promise<T> {
  try {
    const response = await fetch(url, {
      method,
      headers: {
        ...getAuthorizationHeaders(),
        // FormData sets its own multipart Content-Type (with the boundary).
        ...(body instanceof FormData ? {} : { "Content-Type": "application/json" }),
        ...headers,
      },
      body: body instanceof FormData ? body : body ? JSON.stringify(body) : undefined,
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
