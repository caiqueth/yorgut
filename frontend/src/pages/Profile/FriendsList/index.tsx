import React from "react";

interface IFriendsListProps {
  userId: number;
}
export const FriendsList: React.FC<IFriendsListProps> = ({ userId }) => {
  return <>{userId}</>;
};
