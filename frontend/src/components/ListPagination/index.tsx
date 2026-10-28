import { Stack, Typography } from "@mui/material";
import { Button } from "components/Button";
import React, { useEffect } from "react";

interface IProps {
  page: number;
  pageSize: number;
  // Total items in the list (the response's count); undefined while loading.
  count?: number;
  onChange: (page: number) => void;
}

// The page selector of every paginated list:
// "|< primeira  < anterior  3  próxima >  última >|". Hidden when everything
// fits on one page. If the list shrinks past the current page (e.g. the last
// item of the last page was removed), it moves back to the last page that
// exists.
export const ListPagination: React.FC<IProps> = ({ page, pageSize, count, onChange }) => {
  const pageCount = Math.ceil((count ?? 0) / pageSize);

  useEffect(() => {
    if (pageCount > 0 && page > pageCount) onChange(pageCount);
  }, [page, pageCount, onChange]);

  if (pageCount <= 1) return null;

  const isFirst = page <= 1;
  const isLast = page >= pageCount;
  const linkSx = { m: 0.5, px: 1, minWidth: 0 };

  return (
    <Stack
      direction={"row"}
      justifyContent={"center"}
      alignItems={"center"}
      flexWrap={"wrap"}
      sx={{ p: 1 }}
    >
      <Button size={"small"} sx={linkSx} disabled={isFirst} onClick={() => onChange(1)}>
        &#8676; primeira
      </Button>
      <Button size={"small"} sx={linkSx} disabled={isFirst} onClick={() => onChange(page - 1)}>
        &#8592; anterior
      </Button>
      <Typography variant={"body2"} sx={{ mx: 1 }} aria-label={`página ${page} de ${pageCount}`}>
        {`página ${page} de ${pageCount}`}
      </Typography>
      <Button size={"small"} sx={linkSx} disabled={isLast} onClick={() => onChange(page + 1)}>
        próxima &#8594;
      </Button>
      <Button size={"small"} sx={linkSx} disabled={isLast} onClick={() => onChange(pageCount)}>
        última &#8677;
      </Button>
    </Stack>
  );
};
