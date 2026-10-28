import React from "react";
import { isStripedRow } from "helpers/isStripedRow";

import { Grid2 as Grid } from "@mui/material";
import { PicAndTextHorizontalSkeleton } from "components/PicAndTextHorizontal";

interface IStripedGridProps {
  items: React.ReactNode[];
  elementsPerRow?: number;
  // While loading: how many cells to show, and what each one shows. The
  // skeleton should be the loading state of the real items (UserRowSkeleton,
  // PicAndTextSkeleton, ...), so the rows keep their height when data arrives.
  length?: number;
  skeleton?: React.ReactNode;
  isLoading?: boolean;
}
export const StripedGrid: React.FC<IStripedGridProps> = ({
  items,
  elementsPerRow = 3,
  length = 8,
  skeleton = <PicAndTextHorizontalSkeleton />,
  isLoading,
}) => {
  const cells = isLoading ? Array.from({ length }, () => skeleton) : items;

  return (
    <Grid container direction={"row"} width={"100%"}>
      {cells.map((cell, idx) => (
        <Grid
          size={Math.ceil(12 / elementsPerRow)}
          key={idx}
          sx={{
            p: 1,
            backgroundColor: isStripedRow(idx, elementsPerRow)
              ? "yorgut.rows.striped"
              : "yorgut.rows.plain",
          }}
        >
          {cell}
        </Grid>
      ))}
    </Grid>
  );
};
