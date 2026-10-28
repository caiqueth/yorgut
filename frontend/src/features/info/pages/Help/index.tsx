import { Box, Stack, Typography } from "@mui/material";
import { TabsWithContent } from "components/TabsWithContent";
import React from "react";
import { Outlet, useNavigate } from "react-router-dom";

// The help pages, one per tab: about, rules, and terms of use.
export const Help: React.FC = () => {
  const navigate = useNavigate();

  const tabs = [
    { title: "sobre o Yorgut", onClick: () => navigate("/help") },
    { title: "regras de convivência", onClick: () => navigate("/help/rules") },
    { title: "termos de uso e privacidade", onClick: () => navigate("/help/terms") },
  ];

  return (
    <Stack direction={"column"} gap={2} sx={{ width: "100%" }}>
      <Box>
        <Typography>Ajuda do Yorgut</Typography>
      </Box>
      <TabsWithContent tabs={tabs} extraHeader={null}>
        <Outlet />
      </TabsWithContent>
    </Stack>
  );
};
