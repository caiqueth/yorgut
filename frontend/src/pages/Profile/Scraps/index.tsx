import { Box } from "@mui/material";
import React from "react";
import { useParams } from "react-router-dom";

interface IProps { }
export const UserScrapList: React.FC<IProps> = () => {
  const params = useParams();
  const userId = parseInt(params.userId ?? "");

  return <Box>Lista de recados do usuário {userId}</Box>;
};
