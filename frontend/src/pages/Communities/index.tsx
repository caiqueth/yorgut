import { Box, Stack } from "@mui/material";
import { Button } from "components/atoms/Button";
import { TabsWithContent } from "components/molecules/TabsWithContent";
import React from "react";
import { Outlet, useNavigate } from "react-router-dom";

interface IProps { }
export const Communities: React.FC<IProps> = () => {
  const navigate = useNavigate();

  const handleGoToNew = () => {
    navigate("/communities/new")
  }

  const handleGoToMemberships = () => {
    navigate("/communities")
  }

  const handleGoToOwned = () => {
    navigate("/communities/owned")
  }

  const handleGoToDiscover = () => {
    navigate("/communities/discover")
  }

  const tabs = [
    { title: "participando", onClick: handleGoToMemberships },
    { title: "minhas comunidades", onClick: handleGoToOwned, },
    { title: "descobrir", onClick: handleGoToDiscover },
  ];

  return (
    <Box width={"100%"}>
      <Stack direction={"row"} width={"100%"} justifyContent={"end"}>
        <TabsWithContent tabs={tabs} extraHeader={
          <Box sx={{ p1: 1 }}>
            <Button onClick={handleGoToNew}>+ criar comunidade
            </Button>
          </Box>
        }>
          <Outlet />
        </TabsWithContent>
      </Stack>
    </Box>
  );
};
