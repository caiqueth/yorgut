import React from "react";
import { Button } from "components/Button";
import { useNavigate } from "react-router-dom";
import { Box, Stack } from "@mui/material";
import { BoxWithTitle } from "components/BoxWithTitle";
import { StripedGrid } from "components/StripedGrid";
import { PicAndText } from "components/PicAndText";
import { useMe, useMeFriends } from "api/me/queries";
import { useUser, useUserFriends } from "api/user/queries";
import { FriendItem } from "api/generated/model";

interface IFriendsListProps {
  // Whose friends to show; omit for the logged user ("meus amigos").
  userId?: number;
}

// The friends box of the home page (logged user) and of a profile page.
export const FriendsList: React.FC<IFriendsListProps> = ({ userId }) =>
  userId === undefined ? <MyFriendsList /> : <UserFriendsList userId={userId} />;

const MyFriendsList: React.FC = () => {
  const { data: user } = useMe();
  const { data: friends, isLoading } = useMeFriends();

  return (
    <FriendsGrid
      title={`meus amigos (${user?.friendCount})`}
      friends={friends?.friends}
      isLoading={isLoading}
    />
  );
};

const UserFriendsList: React.FC<{ userId: number }> = ({ userId }) => {
  const { data: user } = useUser(userId);
  const { data: friends, isLoading } = useUserFriends(userId);

  return (
    <FriendsGrid
      title={`amigos (${user?.friendCount})`}
      friends={friends?.friends}
      isLoading={isLoading}
    />
  );
};

interface IFriendsGridProps {
  title: string;
  friends?: FriendItem[];
  isLoading: boolean;
}

const FriendsGrid: React.FC<IFriendsGridProps> = ({ title, friends, isLoading }) => {
  const navigate = useNavigate();

  const handleGoToFriends = () => {
    navigate("friends");
  };

  const items =
    friends?.map((f) => (
      <PicAndText
        imageSrc={f.profilePicture}
        url={`/profile/${f.id}`}
        text={`${f.firstName} (${f.friendCount})`}
      />
    )) ?? [];

  return (
    <BoxWithTitle title={title}>
      <Stack direction={"column"} alignItems={"center"}>
        <StripedGrid items={items} elementsPerRow={4} isLoading={isLoading} />
        <Box alignSelf={"center"}>
          <Button onClick={handleGoToFriends}>ver amigos</Button>
        </Box>
      </Stack>
    </BoxWithTitle>
  );
};
