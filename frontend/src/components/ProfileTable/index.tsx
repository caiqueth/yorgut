import { Box, Paper, SxProps, Table, TableBody, Theme, useTheme } from "@mui/material";
import React from "react";
import { StyledTableCell, StyledTableRow } from "components/ProfileTable/styles";

interface IProfileTableData {
  label: string;
  value?: string | React.ReactNode | null;
}

interface IProfileTableProps {
  data: IProfileTableData[];
  // For the box around the table, e.g. to join it to the box above.
  sx?: SxProps<Theme>;
}
export const ProfileTable: React.FC<IProfileTableProps> = ({ data, sx }) => {
  // Themes with white lines between the rows also keep them off the box's edges.
  const spacing = useTheme().yorgut.tableRowSpacing;
  return (
    <Paper sx={sx}>
      <Box sx={{ width: "100%", boxSizing: "border-box", px: `${spacing}px` }}>
        <Table
          size="small"
          sx={{
            width: "100%",
            ...(spacing > 0 && { borderCollapse: "separate", borderSpacing: `0 ${spacing}px` }),
          }}
        >
          <TableBody>
            {data.map(({ label, value }, idx) => (
              <StyledTableRow key={idx}>
                <StyledTableCell sx={{ width: 140 }}>{label}</StyledTableCell>
                <StyledTableCell>{value ?? "-"}</StyledTableCell>
              </StyledTableRow>
            ))}
          </TableBody>
        </Table>
      </Box>
    </Paper>
  );
};
