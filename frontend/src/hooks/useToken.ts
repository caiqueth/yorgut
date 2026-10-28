// The Yorgut session token. The backend checks its signature and expiry; on a
// 401 the HTTP client removes it, so here it's enough to know whether it exists.
export const useToken = () => localStorage.getItem("token");
