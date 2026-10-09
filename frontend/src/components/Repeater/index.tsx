import { Box } from "@mui/material";
import React from "react"

interface IProps {
  element: React.ReactNode | string;
  times: number;
  overlap?: number;
}
export const Repeater: React.FC<IProps> = ({
  element,
  times,
  overlap = 10,
}) => {
  return (
    <Box sx={{ display: "flex" }}>
      {Array(times).fill(element).map((e, idx) => (
        <Box sx={{
          ...(overlap && ({ marginLeft: idx !== 0 ? `-${overlap}px` : 0 }))
        }}>{e}</Box>
      ))}
    </Box>
  )
}