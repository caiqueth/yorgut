import React from "react";
import { Box, Typography } from "@mui/material";
import { BoxWithTitle, Count } from "components/BoxWithTitle";
import { ListPagination } from "components/ListPagination";
import { StripedGrid } from "components/StripedGrid";
import { LIST_PAGE_SIZE } from "helpers/paging";

interface IPagedListProps {
  // The box title; the count is added after it, "amigos de Ana (12)".
  title: React.ReactNode;
  // Total items in the list (the response's count); undefined while loading.
  count?: number;
  // One row per item, on the current page.
  items: React.ReactNode[];
  isLoading: boolean;
  isFetching: boolean;
  // One row while loading (UserRowSkeleton, CommunityRowSkeleton, ...).
  skeleton: React.ReactNode;
  // Shown instead of the rows when the list is empty: a text or a whole block.
  empty?: React.ReactNode;
  page: number;
  onPageChange: (page: number) => void;
}

// A full list page (friends, fans, members, communities): a titled box with one
// row per line and the page selector under it.
export const PagedList: React.FC<IPagedListProps> = ({
  title,
  count,
  items,
  isLoading,
  isFetching,
  skeleton,
  empty,
  page,
  onPageChange,
}) => (
  <BoxWithTitle title={<>{title} <Count value={count} /></>} showProgress={isFetching && !isLoading}>
    {empty && !items.length && !isLoading ? (
      <Box sx={{ p: 1 }}>{typeof empty === "string" ? <Typography>{empty}</Typography> : empty}</Box>
    ) : (
      <StripedGrid items={items} elementsPerRow={1} isLoading={isLoading} skeleton={skeleton} />
    )}
    <ListPagination page={page} pageSize={LIST_PAGE_SIZE} count={count} onChange={onPageChange} />
  </BoxWithTitle>
);
