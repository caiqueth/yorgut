import { Box, Stack, Typography } from "@mui/material";
import { useCommunityMembers } from "api/community/queries";
import { PicAndText } from "components/atoms/PicAndText";
import { PicAndTextHorizontal } from "components/atoms/PicAndTextHorizontal";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import React from "react";
import { Link, useParams } from "react-router-dom";

interface IProps { }
export const CommunityMembers: React.FC<IProps> = () => {
  const params = useParams();
  const communityId = parseInt(params.communityId ?? "");

  const { data, isLoading } = useCommunityMembers(communityId);

  const items = data?.members.map((m) => (
    <PicAndTextHorizontal
      title={`${m.firstName}`}
      text={`${m.friendCount} ${m.friendCount > 1 ? "amigos" : "amigo"}`}
      url={`/profile/${m.id}`}
      imageSrc={m.profilePicture}
      imageSize={64}
    />
  )) ?? [];

  return (
    <BoxWithTitle title={`membros da comunidade (${data?.memberCount})`} isLoading={isLoading}>
      <StripedGrid elementsPerRow={1} items={items} />
    </BoxWithTitle>
  );
};
