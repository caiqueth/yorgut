import { Box } from "@mui/material";
import React from "react";
import { useParams } from "react-router-dom";

interface IProps {}
export const CommunityMembers: React.FC<IProps> = () => {
  const params = useParams();
  const communityId = parseInt(params.communityId ?? "");
  return <Box>Membos da comunidade {communityId}</Box>;
};
