import { Box, Stack } from "@mui/material";
import { Button } from "components/Button";
import { palette } from "theme";

interface IProps {
  tabs: {
    title: string,
    onClick: () => void
  }[];
  extraHeader: React.ReactNode;
  children?: React.ReactNode;
}
export const TabsWithContent: React.FC<IProps> = ({ tabs, extraHeader, children }) => {
  return (
    <Box sx={{ width: "100%" }}>
      <Stack direction={"row"} spacing={"1px"} alignItems={"end"}>
        {tabs.map(({ title, onClick }, idx) => (
          <Button sx={{
            borderRadius: "4px 4px 0 0",
            boxShadow: "none",
            borderTop: `2px solid ${palette.darkerBlue}`,
            borderLeft: `2px solid ${palette.darkerBlue}`,
            borderRight: `2px solid ${palette.darkerBlue}`,
          }}
            key={idx}
            onClick={onClick}>
            {title}
          </Button>

        ))
        }
        <Box flexGrow={1} />
        {extraHeader}
      </Stack >
      {children}
    </Box >
  )
}