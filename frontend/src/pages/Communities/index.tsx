import { Box, Stack } from "@mui/material";
import { Button } from "components/atoms/Button";
import React from "react";
import { Outlet, useNavigate } from "react-router-dom";

interface IProps {
  tabs: { title: string, onClick: () => void }[];
  extraHeader: React.ReactNode;
}
const TabsWithContent: React.FC<IProps> = ({ tabs, extraHeader }) => {
  return (
    <Box sx={{ width: "100%" }}>
      <Stack direction={"row"} spacing={"2px"}>
        {tabs.map(({ title, onClick }) => (
          <Button sx={{
            borderRadius: "4px 4px 0 0",
          }}
            onClick={onClick}>
            {title}
          </Button>
        ))
        }
        <Box flexGrow={1} />
        {extraHeader}
      </Stack >
      <Outlet />
    </Box >
  )
}

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
        } />
      </Stack>
    </Box>
  );
};
