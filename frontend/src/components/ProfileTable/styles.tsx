import { styled, TableCell, TableRow } from "@mui/material";

export const StyledTableRow = styled(TableRow)(({ theme }) => ({
  p: 0,
  m: 0,
  "&:nth-of-type(odd)": {
    backgroundColor: theme.palette.yorgut.tableRows.striped,
  },
  "&:nth-of-type(even)": {
    backgroundColor: theme.palette.yorgut.tableRows.plain,
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
    color: theme.palette.yorgut.tableLabel,
  },
  "&:nth-of-type(even)": {
    textAlign: "left",
  },
}));
