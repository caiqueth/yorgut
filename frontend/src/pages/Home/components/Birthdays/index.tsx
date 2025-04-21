import React from "react";
import { useUserFriends } from "api/user/queries";
import { Link } from "react-router-dom";

import { Grid2 as Grid, Stack, Typography } from "@mui/material";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { dayjsFromDate } from "helpers/formatting/date";
import { StripedGrid } from "components/molecules/StripedGrid";

interface ItemProps {
  id: string;
  firstName: string;
  profilePicture: string;
  birthDate: Date;
}
const Item: React.FC<ItemProps> = ({
  id,
  firstName,
  profilePicture,
  birthDate,
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
        {`${dayjsFromDate(birthDate, "D [de] MMMM")}`}
      </Typography>
    </Stack>
  );
};

interface IBirthdaysProps {
  userId: number;
}
export const Birthdays: React.FC<IBirthdaysProps> = ({ userId }) => {
  const { data: friends } = useUserFriends(userId);

  const items =
    friends?.friends.map((f) => (
      <Item
        id={f.id}
        firstName={f.firstName}
        profilePicture={f.profilePicture}
        birthDate={f.birthDate}
      />
    )) ?? [];

  return (
    <BoxWithTitle title={`próximos aniversários`}>
      <StripedGrid items={items} elementsPerRow={1} />
    </BoxWithTitle>
  );
};
