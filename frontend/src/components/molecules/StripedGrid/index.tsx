import React from "react";
import { isStripedRow } from "helpers/isStripedRow";
import { palette } from "theme";

import { Grid2 as Grid, Skeleton } from "@mui/material";

interface IStripedGridProps {
  items: React.ReactNode[];
  elementsPerRow?: number;
  isLoading?: boolean;
}
export const StripedGrid: React.FC<IStripedGridProps> = ({
  items,
  elementsPerRow = 3,
  isLoading,
}) => {
  return (
    <Grid container direction={"row"} width={"100%"}>
      {(isLoading ? Array.from({ length: 8 }, () => <></>) : items).map((item, idx) => (
        <Grid
          size={Math.ceil(12 / elementsPerRow)}
          key={idx}
          sx={{
            p: 1,
            backgroundColor: isStripedRow(idx, elementsPerRow)
              ? palette.darkBlue
              : palette.mediumBlue,
          }}
        >
          {isLoading ? <Skeleton height={82} /> : item}
        </Grid>
      ))
      }
    </Grid>
  );
};
