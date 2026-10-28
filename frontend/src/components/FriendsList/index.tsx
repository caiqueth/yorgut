import React from "react";
import { Button } from "components/Button";
import { useNavigate } from "react-router-dom";
import { Box, Stack } from "@mui/material";
import { BoxWithTitle, Count } from "components/BoxWithTitle";
import { StripedGrid } from "components/StripedGrid";
import { PicAndText, PicAndTextSkeleton } from "components/PicAndText";
import { useMeFriends } from "api/me/queries";
import { useUserFriends } from "api/user/queries";
import { FriendItem } from "api/generated/model";
import { BOX_SIZE } from "helpers/paging";

interface IFriendsListProps {
  // Whose friends to show; omit for the logged user ("meus amigos").
  userId?: number;
}

// The friends box of the home page (logged user) and of a profile page: the
// most recently seen friends.
export const FriendsList: React.FC<IFriendsListProps> = ({ userId }) =>
  userId === undefined ? <MyFriendsList /> : <UserFriendsList userId={userId} />;

const MyFriendsList: React.FC = () => {
  const { data: friends, isLoading, isFetching } = useMeFriends(1, BOX_SIZE);

  return (
    <FriendsGrid
      title={<>meus amigos <Count value={friends?.count} /></>}
      friends={friends?.friends}
      isLoading={isLoading}
      isRefreshing={isFetching && !isLoading}
    />
  );
};

const UserFriendsList: React.FC<{ userId: number }> = ({ userId }) => {
  const { data: friends, isLoading, isFetching } = useUserFriends(userId, 1, BOX_SIZE);

  return (
    <FriendsGrid
      title={<>amigos <Count value={friends?.count} /></>}
      friends={friends?.friends}
      isLoading={isLoading}
      isRefreshing={isFetching && !isLoading}
    />
  );
};

interface IFriendsGridProps {
  title: React.ReactNode;
  friends?: FriendItem[];
  isLoading: boolean;
  isRefreshing: boolean;
}

const FriendsGrid: React.FC<IFriendsGridProps> = ({ title, friends, isLoading, isRefreshing }) => {
  const navigate = useNavigate();

  const handleGoToFriends = () => {
    navigate("friends");
  };

  const items =
    friends?.map((f) => (
      <PicAndText
        imageSrc={f.profilePicture}
        url={`/profile/${f.id}`}
        title={f.firstName}
        titleSuffix={` (${f.friendCount})`}
      />
    )) ?? [];

  return (
    <BoxWithTitle title={title} showProgress={isRefreshing}>
      <Stack direction={"column"} alignItems={"center"}>
        <StripedGrid items={items} elementsPerRow={4} isLoading={isLoading} skeleton={<PicAndTextSkeleton />} />
        <Box alignSelf={"center"}>
          <Button onClick={handleGoToFriends}>ver amigos</Button>
        </Box>
      </Stack>
    </BoxWithTitle>
  );
};
