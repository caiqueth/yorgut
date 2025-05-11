import { Box, Stack, Typography } from "@mui/material";
import { useMeFriends } from "api/me/queries";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import React from "react";
import { Link } from "react-router-dom";

interface ItemProps {
  id: number;
  firstName: string;
  profilePicture: string;
}
const Item: React.FC<ItemProps> = ({
  id,
  firstName,
  profilePicture,
}) => {
  return (
    <Stack
      alignItems={"center"}
      flexDirection={"row"}
      sx={{ textAlign: "left" }}
      gap={1}
    >
      <img height={64} src={profilePicture} alt="image not found" />
      <Typography variant={"body2"}>
        <Link to={`/profile/${id}`}>{firstName}</Link>{" "}
      </Typography>
    </Stack>
  );
};

interface IProps { }
export const Friends: React.FC<IProps> = () => {
  const { data } = useMeFriends();

  const items =
    data?.friends.map((f) => (
      <Item id={f.id} firstName={f.firstName} profilePicture={f.profilePicture}></Item>
    )) ?? [];

  return (
    <Box width={"100%"}>
      <BoxWithTitle title={"meus amigos"}>
        <StripedGrid items={items} elementsPerRow={1} />
      </BoxWithTitle>
    </Box>
  );
};
