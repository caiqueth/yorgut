import { palette } from "theme";

import { styled, TableCell, TableRow } from "@mui/material";

export const StyledTableRow = styled(TableRow)(({ theme }) => ({
  p: 0,
  m: 0,
  "&:nth-of-type(odd)": {
    backgroundColor: palette.lightBlue,
  },
  "&:nth-of-type(even)": {
    backgroundColor: "white",
  },
}));

export const StyledTableCell = styled(TableCell)(({ theme }) => ({
  border: "none",
  whiteSpace: "pre-line",
  lineHeight: 1.6,
  padding: "0px 6px 0px 6px",
  "&:nth-of-type(odd)": {
    textAlign: "right",
    verticalAlign: "top",
  },
  "&:nth-of-type(even)": {
    textAlign: "left",
  },
}));
