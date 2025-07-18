import { Box } from "@mui/material";
import { useMeFriends } from "api/me/queries";
import { PicAndTextHorizontal } from "components/atoms/PicAndTextHorizontal";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import React from "react";

interface IProps { }
export const Friends: React.FC<IProps> = () => {
  const { data, isLoading } = useMeFriends();

  const items =
    data?.friends.map((f) => (
      <PicAndTextHorizontal
        imageSrc={f.profilePicture}
        title={f.firstName}
        url={`/profile/${f.id}`}
        imageSize={64}
      />
    )) ?? [];

  return (
    <Box width={"100%"}>
      <BoxWithTitle title={"meus amigos"} isLoading={isLoading}>
        <StripedGrid items={items} elementsPerRow={1} />
      </BoxWithTitle>
    </Box>
  );
};
