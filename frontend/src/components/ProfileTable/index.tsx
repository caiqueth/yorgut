import { Box, Paper, Table, TableBody } from "@mui/material";
import React from "react";
import { StyledTableCell, StyledTableRow } from "./styles";

interface IProfileTableData {
  label: string;
  value?: string;
}

interface IProfileTableProps {
  data: IProfileTableData[];
}
export const ProfileTable: React.FC<IProfileTableProps> = ({ data }) => {
  return (
    <Paper>
      <Box sx={{ width: "100%" }}>
        <Table
          size="small"
          sx={{
            width: "100%",
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
