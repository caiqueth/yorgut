import React from "react";
import { PicAndTextHorizontal, PicAndTextHorizontalSkeleton } from "components/PicAndTextHorizontal";
import { friendCountText } from "helpers/counts";

interface IUserRowProps {
  user: {
    id: number;
    firstName: string;
    friendCount: number;
    profilePicture?: string | null;
  };
}

// A user in a list: picture, first name and friend count, linking to the profile.
export const UserRow: React.FC<IUserRowProps> = ({ user }) => (
  <PicAndTextHorizontal
    title={user.firstName}
    text={friendCountText(user.friendCount)}
    url={`/profile/${user.id}`}
    imageSrc={user.profilePicture}
    imageSize={56}
  />
);

export const UserRowSkeleton: React.FC = () => (
  <PicAndTextHorizontalSkeleton imageSize={56} lineWidths={[96, 72]} />
);
