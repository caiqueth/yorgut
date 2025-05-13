import { Box, Stack, Typography } from "@mui/material";
import { useCommunityMembers } from "api/community/queries";
import { PicAndText } from "components/atoms/PicAndText";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import React from "react";
import { Link, useParams } from "react-router-dom";

interface IProps { }
export const CommunityMembers: React.FC<IProps> = () => {
  const params = useParams();
  const communityId = parseInt(params.communityId ?? "");

  const { data } = useCommunityMembers(communityId);

  const items = data?.members.map((m) => (
    <Link to={`/profile/${m.id}`}>
      <Stack direction={"row"} alignItems={"center"} gap={1}>
        <img src={`/${m.profilePicture}`} height={48}></img>
        <Typography>{m.firstName}</Typography>
      </Stack >
    </Link>
  )) ?? [];

  return <BoxWithTitle title={`membros da comunidade (${data?.memberCount})`}>
    <StripedGrid elementsPerRow={1} items={items} />
  </BoxWithTitle>;
};
