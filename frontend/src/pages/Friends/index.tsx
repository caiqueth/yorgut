import { Box, Stack, Typography } from "@mui/material";
import { useMeFriends } from "api/me/queries";
import { PicAndTextHorizontal } from "components/atoms/PicAndTextHorizontal";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import React from "react";
import { Link } from "react-router-dom";
import foreverAlone from "assets/images/foreveralone.png"
import { FriendRating } from "./components/FriendRating";

interface IProps { }
export const Friends: React.FC<IProps> = () => {
  const { data: friends, isLoading } = useMeFriends();

  const items =
    friends?.friends.map((f) => (
      <Stack
        key={f.id}
        direction={"row"}
        justifyContent={"space-between"}
        alignItems={"center"}
        gap={2}
        flexWrap={"wrap"}
      >
        <PicAndTextHorizontal
          imageSrc={f.profilePicture}
          title={f.firstName}
          url={`/profile/${f.id}`}
          imageSize={64}
        />
        <FriendRating userId={f.id} />
      </Stack>
    )) ?? [];

  return (
    <Box width={"100%"}>
      <BoxWithTitle title={"meus amigos"}>
        {!items ? (
          <Box sx={{ p: 1 }}>
            <Typography>Parece que você ainda não possui amigos.</Typography>
            <span><img src={foreverAlone} /></span>
            <Typography><Link to="/search">Clique aqui</Link> e procure por conexões.</Typography>
          </Box>
        ) :
          <StripedGrid items={items} elementsPerRow={1} isLoading={isLoading} />
        }
      </BoxWithTitle>
    </Box>
  );
};
