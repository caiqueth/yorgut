import { useParams } from "react-router-dom";

// Reads a numeric id from the URL (":userId", ":communityId", ...). NaN when
// the route has no such param or it isn't a number.
export const useIdParam = (name: string): number => {
  const params = useParams();
  return parseInt(params[name] ?? "");
};
