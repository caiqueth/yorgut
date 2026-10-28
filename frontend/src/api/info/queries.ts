import { useFindInstanceInfo } from "api/generated/endpoints/info/info";

// This instance's public details (the administrator's e-mail). Set when the
// API starts, so it's fetched once per visit.
export const useInstanceInfo = () =>
  useFindInstanceInfo({ query: { staleTime: Infinity } });
