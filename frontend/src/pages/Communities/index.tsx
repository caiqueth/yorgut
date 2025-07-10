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

  const tabs = [
    { title: "participando", onClick: handleGoToMemberships },
    { title: "minhas comunidades", onClick: handleGoToOwned, }
  ];

  return (
    <Box width={"100%"}>
      <Stack direction={"row"} width={"100%"} justifyContent={"end"}>
        <TabsWithContent tabs={tabs} extraHeader={
          <Button onClick={handleGoToNew}>+ criar comunidade</Button>
        }>
          <Outlet />
        </TabsWithContent>
      </Stack>
    </Box>
  );
};
