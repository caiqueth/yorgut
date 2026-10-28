import { useEffect, useState } from "react";

// The current page of a paged list, starting at 1. It goes back to page 1 when
// resetKey changes, e.g. when navigating to another profile or community.
export const usePage = (resetKey?: unknown) => {
  const [page, setPage] = useState(1);
  useEffect(() => setPage(1), [resetKey]);
  return [page, setPage] as const;
};
